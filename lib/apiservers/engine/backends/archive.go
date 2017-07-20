// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backends

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/tchap/go-patricia/patricia"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/engine/proxy"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	vicarchive "github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/pkg/trace"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

const (
	containerStoreName = "container"
	volumeStoreName    = "volume"
)

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	defer trace.End(trace.Begin(name))

	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return nil, nil, NotFoundError(name)
	}

	stat, err = c.ContainerStatPath(name, path)
	if err != nil {
		// error is wrapped in statpath already
		return nil, nil, err
	}

	var reader io.ReadCloser

	reader, err = c.exportFromContainer(vc, path)
	if err != nil && IsResourceInUse(err) {
		log.Errorf("ContainerArchivePath failed, resource in use: %s", err.Error())
		err = fmt.Errorf("Resource in use")
	}
	if err == nil || reader != nil {
		content = reader
	}
	stat = nil

	return
}

func (c *Container) exportFromContainer(vc *viccontainer.VicContainer, path string) (io.ReadCloser, error) {
	mounts := mountsFromContainer(vc)
	mounts = append(mounts, types.MountPoint{Destination: "/"})
	readerMap := NewArchiveStreamReaderMap(mounts, path)

	readers, err := readerMap.ReadersForSourcePath(archiveProxy, vc.ContainerID, path)
	if err != nil {
		log.Errorf("Errors getting readers for export: %s", err.Error())
		return nil, err
	}

	//FIXME: We need a multi reader that can be closed.  MultiReader returns a regular reader
	log.Infof("Got %d archive readers", len(readers))
	finalTarReader := io.MultiReader(readers...)

	return ioutil.NopCloser(finalTarReader), nil
}

// ContainerCopy performs a deprecated operation of archiving the resource at
// the specified path in the container identified by the given name.
func (c *Container) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("%s does not yet implement ContainerCopy", ProductName())
}

// ContainerExport writes the contents of the container to the given
// writer. An error is returned if the container cannot be found.
func (c *Container) ContainerExport(name string, out io.Writer) error {
	return fmt.Errorf("%s does not yet implement ContainerExport", ProductName())
}

// ContainerExtractToDir extracts the given archive to the specified location
// in the filesystem of the container identified by the given name. The given
// path must be of a directory in the container. If it is not, the error will
// be ErrExtractPointNotDirectory. If noOverwriteDirNonDir is true then it will
// be an error if unpacking the given content would cause an existing directory
// to be replaced with a non-directory and vice versa.
func (c *Container) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	defer trace.End(trace.Begin(name))

	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return NotFoundError(name)
	}

	err := c.importToContainer(vc, path, content)
	if err != nil && IsResourceInUse(err) {
		log.Errorf("ContainerExtractToDir failed, resource in use: %s", err.Error())

		err = fmt.Errorf("Resouce in use")
	}

	return err
}

func (c *Container) importToContainer(vc *viccontainer.VicContainer, target string, content io.Reader) error {
	rawReader, err := archive.DecompressStream(content)
	if err != nil {
		log.Errorf("Input tar stream to ContainerExtractToDir not recognized: %s", err.Error())
		return StreamFormatNotRecognized()
	}
	tarReader := tar.NewReader(rawReader)

	mounts := mountsFromContainer(vc)
	mounts = append(mounts, types.MountPoint{Destination: "/"})
	writerMap := NewArchiveStreamWriterMap(mounts, target)
	defer writerMap.Close() // This should shutdown all the stream connections to the portlayer.

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Lookup the writer for that mount prefix
		writer, err := writerMap.WriterForAsset(archiveProxy, vc.ContainerID, target, *header)
		if err != nil {
			return err
		}

		tarWriter := tar.NewWriter(writer)
		tarWriter.WriteHeader(header)
		// do NOT call tarWriter.Close() or that will create the double nil record for end-of-archive

		_, err = io.Copy(writer, tarReader)
		if err != nil {
			log.Errorf("Error while copying tar data for %s: %s", header.Name, err.Error())
			return err
		}
	}

	return nil
}

// ContainerStatPath stats the filesystem resource at the specified path in the
// container identified by the given name.
func (c *Container) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	defer trace.End(trace.Begin(name))
	op := trace.NewOperation(context.Background(), "ContainerStatPath: %s", name)

	op.Debugf("path received by statpath %s", path)

	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return nil, NotFoundError(name)
	}

	mounts := mountsFromContainer(vc)
	mounts = append(mounts, types.MountPoint{Destination: "/"})

	store, deviceID, fs := resolvePathWithMountPoints(mounts, path, vc.ContainerID)
	// check to see if the path is a mount point, if so, return fake path
	if len(fs.Inclusions) == 1 {
		if _, ok := fs.Inclusions[""]; ok {
			stat = &types.ContainerPathStat{
				Name:       filepath.Base(fs.RebasePath),
				Size:       int64(4096),
				Mode:       os.ModeDir,
				Mtime:      time.Now(),
				LinkTarget: ""}
			op.Debugf("faking container stat path %#v", stat)
			return stat, nil
		}
	}

	stat, err = c.containerProxy.StatPath(op, store, deviceID, fs)
	if err != nil {
		op.Errorf("error getting statpath: %s", err.Error())
		switch err := err.(type) {
		case *storage.StatPathNotFound:
			return nil, ResourceNotFoundError(vc.Name, "file or directory")
		default:
			return nil, InternalServerError(err.Error())
		}
	}

	op.Debugf("container stat path %#v", stat)
	return stat, nil
}

//----------------------------------
// Docker cp utility
//----------------------------------

type ArchiveWriter struct {
	mountPoint types.MountPoint
	filterSpec vicarchive.FilterSpec
	writer     io.WriteCloser
}

type ArchiveReader struct {
	mountPoint types.MountPoint
	filterSpec vicarchive.FilterSpec
	reader     io.ReadCloser
}

// ArchiveStreamWriterMap maps mount prefix to io.WriteCloser
type ArchiveStreamWriterMap struct {
	prefixTrie *patricia.Trie
}

// ArchiveStreamReaderMap maps mount prefix to io.ReadCloser
type ArchiveStreamReaderMap struct {
	prefixTrie *patricia.Trie
}

// NewArchiveStreamWriterMap creates a new ArchiveStreamWriterMap.  The map contains all information
// needed to create a writers for every volume mounts for the container.  This includes the root
// volume of the container.
//
// mounts is the mount data from inspect
// containerDestPath is the destination path in the container
func NewArchiveStreamWriterMap(mounts []types.MountPoint, dest string) *ArchiveStreamWriterMap {
	writerMap := &ArchiveStreamWriterMap{}
	writerMap.prefixTrie = patricia.NewTrie()

	for _, m := range mounts {
		aw := ArchiveWriter{
			mountPoint: m,
			writer:     nil,
		}

		// If container destination path is part of this mount point's prefix, we must remove it and
		// add to the filterspec.  If the container destination path is "/" we do no stripping.
		//
		//	e.g. mount A at /mnt/A
		//
		//		cp /mnt cid:/mnt
		//
		// file data.txt from local /mnt/A/data.txt will come to the persona as A/data.txt.  We must
		// tell the storage portlayer to remove "A".
		//
		//	e.g. mount A at /mnt/A
		//
		//		cp / cid:/
		//
		// file data.txt from local /mnt/A/data.txt will come to the persona as mnt/A/data.txt.
		// Here, we must tell the portlayer to remove "mnt/A".  The key to determining whether to
		// strip "A" or "mnt/A" is based on the container destination path.
		isPrimary := !strings.Contains(aw.mountPoint.Destination, dest) || aw.mountPoint.Destination == dest
		aw.filterSpec = vicarchive.GenerateFilterSpec(dest, aw.mountPoint.Destination, isPrimary, vicarchive.CopyTo)

		writerMap.prefixTrie.Insert(patricia.Prefix(m.Destination), &aw)
	}

	return writerMap
}

// NewArchiveStreamReaderMap creates a new ArchiveStreamReaderMap.  After the call, it contains
// information to create readers for every volume mounts for the container
//
// mounts is the mount data from inspect
func NewArchiveStreamReaderMap(mounts []types.MountPoint, dest string) *ArchiveStreamReaderMap {
	readerMap := &ArchiveStreamReaderMap{}
	readerMap.prefixTrie = patricia.NewTrie()

	for _, m := range mounts {
		ar := ArchiveReader{
			mountPoint: m,
			reader:     nil,
		}

		// If the mount point is not the root file system, we must tell the portlayer to rebase the
		// files in the return tar stream with the mount point path since the volume does not know
		// the path it is mounted to.  It only knows it's root file system.
		//
		//	e.g. mount A at /mnt/A with a file data.txt in A
		//
		//	/mnt/A/data.txt		<-- from container point of view
		//	/data.txt			<-- from volume point of view
		//
		// Neither the volume nor the storage portlayer knows about /mnt/A.  The persona must tell
		// the portlayer to rebase all files from this volume to the /mnt/A/ in the final tar stream.
		isPrimary := !strings.Contains(ar.mountPoint.Destination, dest) || ar.mountPoint.Destination == dest
		ar.filterSpec = vicarchive.GenerateFilterSpec(dest, ar.mountPoint.Destination, isPrimary, vicarchive.CopyFrom)

		readerMap.prefixTrie.Insert(patricia.Prefix(m.Destination), &ar)
	}

	return readerMap
}

// FindArchiveWriter finds the one writer that matches the asset name.  There should only be one
// stream this asset needs to be written to.
func (wm *ArchiveStreamWriterMap) FindArchiveWriter(containerDestPath, assetName string) (*ArchiveWriter, error) {
	defer trace.End(trace.Begin(""))

	var aw *ArchiveWriter
	var err error

	// go function used later for searching
	findPrefix := func(prefix patricia.Prefix, item patricia.Item) error {
		if _, ok := item.(*ArchiveWriter); !ok {
			return fmt.Errorf("item not ArchiveWriter")
		}

		aw, _ = item.(*ArchiveWriter)
		return nil
	}

	// Find the prefix for the final destination.  Final destination is the combination of container destination path
	// and the asset's name.  For example,
	//
	//	container destination path =	/
	//	asset name =					mnt/A/file.txt
	//	mount 1	=						/mnt/A
	//	mount prefix =					/mnt/A
	//
	// In the above example, mount prefxi can only be determined by combining both the container destination path and
	// the asset name, as the final destination includes a mounted volume.
	combinedPath := path.Join(containerDestPath, assetName)
	prefix := patricia.Prefix(combinedPath)
	err = wm.prefixTrie.VisitPrefixes(prefix, findPrefix)
	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Failed to find a node for prefix %s: %s", containerDestPath, err.Error())
	}

	if aw == nil {
		return nil, fmt.Errorf("No archive writer found for container destination %s and asset name %s", containerDestPath, assetName)
	}

	return aw, nil
}

// WriterForAsset takes a destination path and subpath of the archive data and returns the
// appropriate writer for the two.  It's intention is to solve the case where there exist
// a mount point and another mount point within the first mount point.  For instance, the
// prefix map can have,
//
//		R/W -					/
//		mount 1 -				/mnt/a
//		mount 2 -				/mnt/a/b
//
//		case 1:
//				containerDestPath -			/mnt/a
//				archive header source -		b/file.txt
//
//			The correct writer would be the one corresponding to mount 2.
//
//		case 2:
//				containerDestPath -			/mnt/a
//				archive header source -		file.txt
//
//			The correct writer would be the one corresponding to mount 1.
//
//		case 3:
//				containerDestPath -			/
//				archive header source -		mnt/a/file.txt
//
//			The correct writer would be the one corresponding to mount 1
//
// As demonstrated above, the mount prefix and writer cannot be determined with just the
// container destination path.  It must be combined with the actual asset's name.
func (wm *ArchiveStreamWriterMap) WriterForAsset(proxy proxy.VicArchiveProxy, cid, containerDestPath string, assetHeader tar.Header) (io.WriteCloser, error) {
	defer trace.End(trace.Begin(assetHeader.Name))

	var err error
	var streamWriter io.WriteCloser

	aw, err := wm.FindArchiveWriter(containerDestPath, assetHeader.Name)
	if err != nil {
		return nil, err
	}

	// Perform the lazy initialization here.
	if aw.writer == nil {
		// lazy initialize.
		log.Debugf("Lazily initializing import stream for %s", aw.mountPoint.Destination)
		var deviceID string
		var store string
		if aw.mountPoint.Destination == "/" {
			// Special case. / refers to container VMDK and not a volume vmdk.
			deviceID = cid
			store = constants.ContainerStoreName
		} else {
			deviceID = aw.mountPoint.Name
			store = constants.VolumeStoreName
		}
		streamWriter, err = proxy.ArchiveImportWriter(context.Background(), store, deviceID, aw.filterSpec)
		if err != nil {
			err = fmt.Errorf("Unable to initialize import stream writer for mount prefix %s", aw.mountPoint.Destination)
			log.Errorf(err.Error())
			return nil, err
		}
		aw.writer = streamWriter
	} else {
		streamWriter = aw.writer
	}

	return streamWriter, nil
}

// Close visits all the archive writer in the trie and closes the actual io.WritCloser
func (wm *ArchiveStreamWriterMap) Close() {
	defer trace.End(trace.Begin(""))

	closeStream := func(prefix patricia.Prefix, item patricia.Item) error {
		if aw, ok := item.(*ArchiveWriter); ok && aw.writer != nil {
			aw.writer.Close()
			aw.writer = nil
		}
		return nil
	}

	wm.prefixTrie.Visit(closeStream)
}

// FindArchiveReaders finds all archive readers that are within the container source path.  For example,
//
//	mount A -			/mnt/A
//	mount B -			/mnt/B
//	mount AB -			/mnt/A/AB
//	base container -	/
//
//	container source path - /mnt/A
//
// For the above example, this function returns the readers for mount A and mount AB but not the
// readers for / or mount B.
func (rm *ArchiveStreamReaderMap) FindArchiveReaders(containerSourcePath string) ([]*ArchiveReader, error) {
	defer trace.End(trace.Begin(containerSourcePath))

	var nodes []*ArchiveReader
	var startingNode *ArchiveReader
	var err error

	findStartingPrefix := func(prefix patricia.Prefix, item patricia.Item) error {
		if _, ok := item.(*ArchiveReader); !ok {
			return fmt.Errorf("item not ArchiveReader")
		}

		startingNode = item.(*ArchiveReader)
		return nil
	}

	walkPrefixSubtree := func(prefix patricia.Prefix, item patricia.Item) error {
		if _, ok := item.(*ArchiveReader); !ok {
			return fmt.Errorf("item not ArchiveReader")
		}

		ar, _ := item.(*ArchiveReader)
		nodes = append(nodes, ar)
		return nil
	}

	// Find all mounts for the sourcepath
	prefix := patricia.Prefix(containerSourcePath)
	err = rm.prefixTrie.VisitSubtree(prefix, walkPrefixSubtree)
	if err != nil {
		msg := fmt.Sprintf("Failed to find a node for prefix %s: %s", containerSourcePath, err.Error())
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// The above subtree walking MAY NOT find the starting prefix.  For example /etc will not find /.
	// Subtree only finds prefix that starts with /etc.  VisitPrefixes will find the starting prefix.
	// If the search was for /, then it will not find the starting node.  In that case, we grab the
	// first node in the slice.
	err = rm.prefixTrie.VisitPrefixes(prefix, findStartingPrefix)
	if err != nil {
		msg := fmt.Sprintf("Failed to find starting node for prefix %s: %s", containerSourcePath, err.Error())
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	if startingNode != nil {
		found := false
		for _, node := range nodes {
			if node.mountPoint.Destination == startingNode.mountPoint.Destination {
				found = true
				break
			}
		}

		if !found {
			// prepend the starting node at the beginning
			nodes = append([]*ArchiveReader{startingNode}, nodes...)
		}
	} else if len(nodes) > 0 {
		startingNode = nodes[0]
	} else {
		msg := fmt.Sprintf("Failed to find starting node for prefix %s: %s", containerSourcePath, err.Error())
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	err = rm.buildFilterSpec(containerSourcePath, nodes, startingNode)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (rm *ArchiveStreamReaderMap) buildFilterSpec(containerSourcePath string, nodes []*ArchiveReader, startingNode *ArchiveReader) error {
	mounts := []string{}
	for _, node := range nodes {
		mounts = append(mounts, node.mountPoint.Destination)
	}

	for _, node := range nodes {
		vicarchive.AddMountInclusionsExclusions(node.mountPoint.Destination, &node.filterSpec, mounts, containerSourcePath)
	}
	return nil
}

// ReadersForSourcePath returns all an array of io.Reader for all the readers within a container source path.
//		Example:
//			Reader 1 -				/mnt/A
//			Reader 2 -				/mnt/A/B
//
//			containaerSroucePath -	/mnt/A
// In the above, both readers are within the the container source path.
func (rm *ArchiveStreamReaderMap) ReadersForSourcePath(proxy proxy.VicArchiveProxy, cid, containerSourcePath string) ([]io.Reader, error) {
	defer trace.End(trace.Begin(containerSourcePath))

	var streamReaders []io.Reader

	nodes, err := rm.FindArchiveReaders(containerSourcePath)
	if err != nil {
		return nil, err
	}

	mounts := []string{}
	for _, node := range nodes {
		mounts = append(mounts, node.mountPoint.Destination)
	}

	// Create the io.Reader for those mounts if they haven't already been initialized
	for _, node := range nodes {
		// build up the inclusions and exclusions
		vicarchive.AddMountInclusionsExclusions(node.mountPoint.Destination, &node.filterSpec, mounts, containerSourcePath)
		if node.reader == nil {
			var store, deviceID string
			if node.mountPoint.Destination == "/" {
				// Special case. / refers to container VMDK and not a volume vmdk.
				store = constants.ContainerStoreName
				deviceID = cid
			} else {
				store = constants.VolumeStoreName
				deviceID = node.mountPoint.Name
			}

			log.Infof("Lazily initializing export stream for %s [%s]", node.mountPoint.Name, node.mountPoint.Destination)
			reader, err := proxy.ArchiveExportReader(context.Background(), store, "", deviceID, "", true, node.filterSpec)
			if err != nil {
				err = fmt.Errorf("Unable to initialize export stream reader for prefix %s", node.mountPoint.Destination)
				log.Errorf(err.Error())
				return nil, err
			}
			log.Infof("Lazy initialization created reader %#v", reader)
			streamReaders = append(streamReaders, reader)
		} else {
			streamReaders = append(streamReaders, node.reader)
		}
	}

	if len(nodes) == 0 {
		log.Infof("Found no archive readers for %s", containerSourcePath)
	}

	return streamReaders, nil
}

// Close visits all the archive readers in the trie and closes the actual io.ReadCloser
func (rm *ArchiveStreamReaderMap) Close() {
	defer trace.End(trace.Begin(""))

	closeStream := func(prefix patricia.Prefix, item patricia.Item) error {
		if aw, ok := item.(*ArchiveReader); ok && aw.reader != nil {
			aw.reader.Close()
			aw.reader = nil
		}
		return nil
	}

	rm.prefixTrie.Visit(closeStream)
}

// use mountpoints to strip the target to a relative path
func resolvePathWithMountPoints(mounts []types.MountPoint, path, defaultDevice string) (string, string, vicarchive.FilterSpec) {
	var fs vicarchive.FilterSpec
	deviceID := defaultDevice
	store := containerStoreName
	mntpoint := ""

	// trim / and . off from path and then append / to ensure the format is correct
	for strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	for strings.HasSuffix(path, ".") {
		path = strings.TrimSuffix(path, ".")
	}
	for strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	path = "/" + path

	for _, mount := range mounts {
		if strings.HasPrefix(path, mount.Destination) {
			// path is exactly the mountpoint, not point in making a call to portlayer
			if len(mount.Destination) > len(mntpoint) {
				mntpoint = mount.Destination
				if mntpoint != "/" {
					deviceID = mount.Name
					store = volumeStoreName
				}
			}
		}
	}

	fs.RebasePath = mntpoint
	fs.Inclusions = make(map[string]struct{})
	fs.Exclusions = make(map[string]struct{})
	includedPath := strings.TrimPrefix(path, mntpoint)
	excludedPath := path + "/"
	fs.Inclusions[includedPath] = struct{}{}
	fs.Exclusions[excludedPath] = struct{}{}

	return store, deviceID, fs
}
