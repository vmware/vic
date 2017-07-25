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
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime"
	rc "github.com/go-openapi/runtime/client"
	httpclient "github.com/mreiferson/go-httpclient"
	"github.com/tchap/go-patricia/patricia"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	viccontainer "github.com/vmware/vic/lib/apiservers/engine/backends/container"
	"github.com/vmware/vic/lib/apiservers/portlayer/client"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	vicarchive "github.com/vmware/vic/lib/archive"
	"github.com/vmware/vic/lib/portlayer/constants"
	"github.com/vmware/vic/pkg/trace"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

type VicArchiveProxy interface {
	ArchiveExportReader(ctx context.Context, store, ancestorStore, deviceID, ancestor string, data bool, filterSpec vicarchive.FilterSpec) (io.ReadCloser, error)
	ArchiveImportWriter(ctx context.Context, store, deviceID string, filterSpec vicarchive.FilterSpec) (io.WriteCloser, error)
}

// ContainerArchivePath creates an archive of the filesystem resource at the
// specified path in the container identified by the given name. Returns a
// tar archive of the resource and whether it was a directory or a single file.
func (c *Container) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	defer trace.End(trace.Begin(name))

	vc := cache.ContainerCache().GetContainer(name)
	if vc == nil {
		return nil, nil, NotFoundError(name)
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

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(header)
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
	defer trace.End(trace.Begin(fmt.Sprintf("** statpath, name=%s, path=%s", name, path)))

	fakeStat := &types.ContainerPathStat{}
	fakeStat.Mode = os.ModeDir

	return fakeStat, nil
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

		// hickeng: the current logic is strip, then rebase for export - more efficient to collapse
		// any reducdency here than apply to every tar entry when packing but skipping for now.
		aw.filterSpec.StripPath = aw.mountPoint.Destination
		aw.filterSpec.RebasePath = strings.TrimPrefix(dest, aw.mountPoint.Destination)

		// if containerDestPath != "/" && strings.HasPrefix(aw.mountPoint.Destination, containerDestPath) {
		// 	aw.filterSpec.StripPath = strings.TrimPrefix(aw.mountPoint.Destination, containerDestPath)
		// }

		aw.filterSpec.Exclusions = make(map[string]struct{})
		aw.filterSpec.Inclusions = make(map[string]struct{})

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

		// hickeng: the current logic is rebase, then strip for export - more efficient to collapse
		// any reducdency here than apply to every tar entry when packing but skipping for now.
		ar.filterSpec.RebasePath = ar.mountPoint.Destination
		ar.filterSpec.StripPath = path.Join("/", path.Dir(dest))

		ar.filterSpec.Exclusions = make(map[string]struct{})
		ar.filterSpec.Inclusions = make(map[string]struct{})

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
func (wm *ArchiveStreamWriterMap) WriterForAsset(proxy VicArchiveProxy, cid, containerDestPath string, assetHeader tar.Header) (io.WriteCloser, error) {
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
	var err error

	// Build an exclusion filter for each writer.  For example, the reader for / should not read
	// from submounts as there are separate readers for those.
	buildExclusion := func(path string, node *ArchiveReader) error {
		childWalker := func(prefix patricia.Prefix, item patricia.Item) error {
			if _, ok := item.(*ArchiveReader); !ok {
				return fmt.Errorf("item not ArchiveReader")
			}

			ar, _ := item.(*ArchiveReader)
			dest := ar.mountPoint.Destination
			if dest != path {
				node.filterSpec.Exclusions[dest] = struct{}{}
			}
			return nil
		}

		// prefix = current node's mount path
		nodePrefix := patricia.Prefix(path)

		err = rm.prefixTrie.VisitSubtree(nodePrefix, childWalker)
		if err != nil {
			msg := fmt.Sprintf("Failed to build exclusion filter for %s: %s", path, err.Error())
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		return nil
	}

	for _, node := range nodes {
		// Clear out existing exclusions and inclusions
		node.filterSpec.Exclusions = make(map[string]struct{})
		node.filterSpec.Inclusions = make(map[string]struct{})

		err = buildExclusion(node.mountPoint.Destination, node)
		if err != nil {
			return err
		}
	}

	// Add inclusion filter.  When there is an inclusion, there should be only one node in
	// the slice that is returned.
	//
	//	Example 1:
	//		containerSourcePath -	/file.txt
	//
	//		ArchiveReader path -	/
	//		Inclusion filter -		file.txt
	//
	//	Example 2:
	//		containerSourcePath -	/mnt/A/a/file.txt
	//
	//		ArchiveReader path -	/mnt/A
	//		Inclusion filter -		a/file.txt
	inclusionPath := strings.TrimPrefix(containerSourcePath, startingNode.mountPoint.Destination)
	inclusionPath = strings.TrimPrefix(inclusionPath, "/")
	if len(nodes) == 1 {
		nodes[0].filterSpec.Inclusions[inclusionPath] = struct{}{}
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
func (rm *ArchiveStreamReaderMap) ReadersForSourcePath(proxy VicArchiveProxy, cid, containerSourcePath string) ([]io.Reader, error) {
	defer trace.End(trace.Begin(containerSourcePath))

	var streamReaders []io.Reader

	nodes, err := rm.FindArchiveReaders(containerSourcePath)
	if err != nil {
		return nil, err
	}

	// Create the io.Reader for those mounts if they haven't already been initialized
	for _, node := range nodes {
		if node.reader == nil {
			var store, deviceID string
			subpath := containerSourcePath
			if node.mountPoint.Destination == "/" {
				// Special case. / refers to container VMDK and not a volume vmdk.
				store = constants.ContainerStoreName
				deviceID = cid
			} else {
				store = constants.VolumeStoreName
				deviceID = node.mountPoint.Name
				subpath = strings.TrimPrefix(containerSourcePath, node.mountPoint.Destination)
			}

			if strings.HasPrefix(containerSourcePath, node.mountPoint.Destination) {
				// add the include path back
				if node.filterSpec.Inclusions == nil {
					node.filterSpec.Inclusions = make(map[string]struct{})
				}

				node.filterSpec.Inclusions[subpath] = struct{}{}
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

//------------------------------------
// ArchiveProxy
//------------------------------------

type ArchiveProxy struct {
}

func NewArchiveProxy() VicArchiveProxy {
	return &ArchiveProxy{}
}

// ArchiveExportReader streams a tar archive from the portlayer.  Once the stream is complete,
// an io.Reader is returned and the caller can use that reader to parse the data.
func (a *ArchiveProxy) ArchiveExportReader(ctx context.Context, store, ancestorStore, deviceID, ancestor string, data bool, filterSpec vicarchive.FilterSpec) (io.ReadCloser, error) {
	defer trace.End(trace.Begin(deviceID))

	if store == "" || deviceID == "" {
		return nil, fmt.Errorf("ArchiveExportReader called with either empty store or deviceID")
	}

	var err error

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		// make sure we get out of io.Copy if context is canceled
		select {
		case <-ctx.Done():
			// Attempt to tell the portlayer to cancel the stream.  This is one way of cancelling the
			// stream.  The other way is for the caller of this function to close the returned CloseReader.
			// Callers of this function should do one but not both.
			pipeReader.Close()
		}
	}()

	go func() {
		params := storage.NewExportArchiveParamsWithContext(ctx).
			WithStore(store).
			WithAncestorStore(&ancestorStore).
			WithDeviceID(deviceID).
			WithAncestor(&ancestor).
			WithData(data)

		// Encode the filter spec
		encodedFilter := ""
		if valueBytes, merr := json.Marshal(filterSpec); merr == nil {
			encodedFilter = base64.StdEncoding.EncodeToString(valueBytes)
			params = params.WithFilterSpec(&encodedFilter)
			log.Infof(" encodedFilter = %s", encodedFilter)
		}

		client := PortLayerClient()
		_, err = client.Storage.ExportArchive(params, pipeWriter)
		if err != nil {
			log.Errorf("Error from ExportArchive: %s", err.Error())
			switch err := err.(type) {
			case *storage.ExportArchiveInternalServerError:
				plErr := InternalServerError(fmt.Sprintf("Server error from archive reader for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeWriter.CloseWithError(plErr)
			case *storage.ImportArchiveLocked:
				plErr := ResourceLockedError(fmt.Sprintf("Resource locked for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeWriter.CloseWithError(plErr)
			default:
				//Check for EOF.  Since the connection, transport, and data handling are
				//encapsulated inside of Swagger, we can only detect EOF by checking the
				//error string
				if strings.Contains(err.Error(), swaggerSubstringEOF) {
					log.Debugf("swagger error %s", err.Error())
					pipeWriter.Close()
				} else {
					pipeWriter.CloseWithError(err)
				}
			}
		} else {
			pipeWriter.Close()
		}
	}()

	return pipeReader, nil
}

// ArchiveImportWriter initializes a write stream for a path.  This is usually called
// for getting a writer during docker cp TO container.
func (a *ArchiveProxy) ArchiveImportWriter(ctx context.Context, store, deviceID string, filterSpec vicarchive.FilterSpec) (io.WriteCloser, error) {
	defer trace.End(trace.Begin(deviceID))

	if store == "" || deviceID == "" {
		return nil, fmt.Errorf("ArchiveImportWriter called with either empty store or deviceID")
	}

	var err error

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		// make sure we get out of io.Copy if context is canceled
		select {
		case <-ctx.Done():
			pipeWriter.Close()
		}
	}()

	go func() {
		// encodedFilter and destination are not required (from swagge spec) because
		// they are allowed to be empty.
		params := storage.NewImportArchiveParamsWithContext(ctx).
			WithStore(store).
			WithDeviceID(deviceID).
			WithArchive(pipeReader)

		// Encode the filter spec
		encodedFilter := ""
		if valueBytes, merr := json.Marshal(filterSpec); merr == nil {
			encodedFilter = base64.StdEncoding.EncodeToString(valueBytes)
			params = params.WithFilterSpec(&encodedFilter)
		}

		client := PortLayerClient()
		_, err = client.Storage.ImportArchive(params)
		if err != nil {
			switch err := err.(type) {
			case *storage.ImportArchiveInternalServerError:
				plErr := InternalServerError(fmt.Sprintf("Server error from archive writer for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeReader.CloseWithError(plErr)
			case *storage.ImportArchiveLocked:
				plErr := ResourceLockedError(fmt.Sprintf("Resource locked for device %s", deviceID))
				log.Errorf(plErr.Error())
				pipeReader.CloseWithError(plErr)
			default:
				//Check for EOF.  Since the connection, transport, and data handling are
				//encapsulated inside of Swagger, we can only detect EOF by checking the
				//error string
				if strings.Contains(err.Error(), swaggerSubstringEOF) {
					log.Errorf(err.Error())
					pipeReader.Close()
				} else {
					pipeReader.CloseWithError(err)
				}
			}
		} else {
			pipeReader.Close()
		}
	}()

	return pipeWriter, nil
}

func createGzipTarClient(connectTimeout, responseTimeout, responseHeaderTimeout time.Duration) (*client.PortLayer, *httpclient.Transport) {

	r := rc.New(PortLayerServer(), "/", []string{"http"})
	transport := &httpclient.Transport{
		ConnectTimeout:        connectTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
		RequestTimeout:        responseTimeout,
	}

	r.Transport = transport

	plClient := client.New(r, nil)
	bsc := runtime.ByteStreamConsumer()
	r.Consumers["application/octet-stream"] = bsc
	r.Producers["application/octet-stream"] = runtime.ByteStreamProducer()

	r.Consumers["application/x-tar"] = runtime.ConsumerFunc(func(rdr io.Reader, data interface{}) error {
		return bsc.Consume(rdr, data)
	})
	return plClient, transport
}
