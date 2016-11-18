// Copyright 2016 VMware, Inc. All Rights Reserved.
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
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/registry"

	"github.com/vmware/vic/lib/apiservers/engine/backends/cache"
	"github.com/vmware/vic/lib/apiservers/portlayer/client/storage"
	"github.com/vmware/vic/lib/imagec"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/uid"
	"github.com/vmware/vic/pkg/vsphere/sys"
)

// byCreated is a temporary type used to sort a list of images by creation
// time.
type byCreated []*types.Image

func (r byCreated) Len() int           { return len(r) }
func (r byCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }

type Image struct {
}

func (i *Image) Commit(name string, config *types.ContainerCommitConfig) (imageID string, err error) {
	return "", fmt.Errorf("%s does not implement image.Commit", ProductName())
}

func (i *Image) Exists(containerName string) bool {
	return false
}

// TODO fix the errors so the client doesnt print the generic POST or DELETE message
func (i *Image) ImageDelete(imageRef string, force, prune bool) ([]types.ImageDelete, error) {
	defer trace.End(trace.Begin(imageRef))

	var deleted []types.ImageDelete
	var userRefIsID bool
	var imageRemoved bool

	// Use the image cache to go from the reference to the ID we use in the image store
	img, err := cache.ImageCache().Get(imageRef)
	if err != nil {
		return nil, err
	}

	// Get the tags from the repo cache for this image
	// TODO: remove this -- we have it in the image above
	tags := cache.RepositoryCache().Tags(img.ImageID)

	// did the user pass an id or partial id
	userRefIsID = cache.ImageCache().IsImageID(imageRef)
	// do we have any reference conflicts
	if len(tags) > 1 && userRefIsID && !force {
		t := uid.Parse(img.ImageID).Truncate()
		return nil,
			fmt.Errorf("conflict: unable to delete %s (must be forced) - image is referenced in one or more repositories", t)
	}

	// if we have an ID or only 1 tag lets delete the vmdk(s) via the PL
	if userRefIsID || len(tags) == 1 {
		log.Infof("Deleting image via PL %s (%s)", img.ImageID, img.ID)

		// needed for image store
		host, err := sys.UUID()
		if err != nil {
			return nil, err
		}

		params := storage.NewDeleteImageParamsWithContext(ctx).WithStoreName(host).WithID(img.ID)
		// TODO: This will fail if any containerVMs are referencing the vmdk - vanilla docker
		// allows the removal of an image (via force flag) even if a container is referencing it
		// should vic?
		_, err = PortLayerClient().Storage.DeleteImage(params)
		if err != nil {
			switch err := err.(type) {
			case *storage.DeleteImageLocked:
				return nil, fmt.Errorf("Failed to remove image %q: %s", imageRef, err.Payload.Message)
			default:
				return nil, err
			}
		}

		// we've deleted the image so remove from cache
		cache.ImageCache().RemoveImageByConfig(img)
		imagec.LayerCache().Remove(img.ID)
		imageRemoved = true

	} else {

		// only untag the ref supplied
		n, err := reference.ParseNamed(imageRef)
		if err != nil {
			return nil, fmt.Errorf("unable to parse reference(%s): %s", imageRef, err.Error())
		}
		tag := reference.WithDefaultTag(n)
		tags = []string{tag.String()}
	}
	// loop thru and remove from repoCache
	for i := range tags {
		// remove from cache, but don't save -- we'll do that afer all
		// updates
		refNamed, _ := cache.RepositoryCache().Remove(tags[i], false)
		dd := types.ImageDelete{Untagged: refNamed}
		deleted = append(deleted, dd)
	}

	// save repo now -- this will limit the number of PL
	// calls to one per rmi call
	err = cache.RepositoryCache().Save()
	if err != nil {
		return nil, fmt.Errorf("Untag error: %s", err.Error())
	}

	if imageRemoved {
		imageDeleted := types.ImageDelete{Deleted: img.ImageID}
		deleted = append(deleted, imageDeleted)
	}

	return deleted, err
}

func (i *Image) ImageHistory(imageName string) ([]*types.ImageHistory, error) {
	return nil, fmt.Errorf("%s does not implement image.History", ProductName())
}

func (i *Image) Images(filterArgs string, filter string, all bool) ([]*types.Image, error) {
	defer trace.End(trace.Begin("Images"))

	images := cache.ImageCache().GetImages()

	result := make([]*types.Image, 0, len(images))

	for _, image := range images {
		result = append(result, convertV1ImageToDockerImage(image))
	}

	// sort on creation time
	sort.Sort(sort.Reverse(byCreated(result)))

	return result, nil
}

// Docker Inspect.  LookupImage looks up an image by name and returns it as an
// ImageInspect structure.
func (i *Image) LookupImage(name string) (*types.ImageInspect, error) {
	defer trace.End(trace.Begin("LookupImage (docker inspect)"))

	imageConfig, err := cache.ImageCache().Get(name)
	if err != nil {
		return nil, err
	}

	return imageConfigToDockerImageInspect(imageConfig, ProductName()), nil
}

func (i *Image) TagImage(newTag reference.Named, imageName string) error {
	return fmt.Errorf("%s does not implement image.Tag", ProductName())
}

func (i *Image) LoadImage(inTar io.ReadCloser, outStream io.Writer, quiet bool) error {
	return fmt.Errorf("%s does not implement image.LoadImage", ProductName())
}

func (i *Image) ImportImage(src string, newRef reference.Named, msg string, inConfig io.ReadCloser, outStream io.Writer, config *container.Config) error {
	return fmt.Errorf("%s does not implement image.ImportImage", ProductName())
}

func (i *Image) ExportImage(names []string, outStream io.Writer) error {
	return fmt.Errorf("%s does not implement image.ExportImage", ProductName())
}

func (i *Image) PullImage(ctx context.Context, ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	defer trace.End(trace.Begin(ref.String()))

	log.Debugf("PullImage: ref = %+v, metaheaders = %+v\n", ref, metaHeaders)

	options := imagec.Options{
		Destination: os.TempDir(),
		Reference:   ref.String(),
		Timeout:     imagec.DefaultHTTPTimeout,
		Outstream:   outStream,
		RegistryCAs: RegistryCertPool,
	}

	if authConfig != nil {
		if len(authConfig.Username) > 0 {
			options.Username = authConfig.Username
		}
		if len(authConfig.Password) > 0 {
			options.Password = authConfig.Password
		}
	}

	portLayerServer := PortLayerServer()

	if portLayerServer != "" {
		options.Host = portLayerServer
	}

	insecureRegistries := InsecureRegistries()
	for _, registry := range insecureRegistries {
		if registry == ref.Hostname() {
			options.InsecureAllowHTTP = true
			break
		}
	}

	log.Infof("PullImage: reference: %s, %s, portlayer: %#v",
		options.Reference,
		options.Host,
		portLayerServer)

	ic := imagec.NewImageC(options, streamformatter.NewJSONStreamFormatter())
	err := ic.PullImage()
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) PushImage(ctx context.Context, ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return fmt.Errorf("%s does not implement image.PushImage", ProductName())
}

func (i *Image) SearchRegistryForImages(ctx context.Context, term string, authConfig *types.AuthConfig, metaHeaders map[string][]string) (*registry.SearchResults, error) {
	return nil, fmt.Errorf("%s does not implement image.SearchRegistryForImages", ProductName())
}

// Utility functions

func convertV1ImageToDockerImage(image *metadata.ImageConfig) *types.Image {
	var labels map[string]string
	if image.Config != nil {
		labels = image.Config.Labels
	}

	return &types.Image{
		ID:          image.ImageID,
		ParentID:    image.Parent,
		RepoTags:    image.Tags,
		RepoDigests: image.Digests,
		Created:     image.Created.Unix(),
		Size:        image.Size,
		VirtualSize: image.Size,
		Labels:      labels,
	}
}

// Converts the data structure retrieved from the portlayer.  This src datastructure
// represents the unmarshalled data saved in the storage port layer.  The return
// data is what the Docker CLI understand and returns to user.
func imageConfigToDockerImageInspect(imageConfig *metadata.ImageConfig, productName string) *types.ImageInspect {
	if imageConfig == nil {
		return nil
	}

	rootfs := types.RootFS{
		Type:      "layers",
		Layers:    make([]string, 0, len(imageConfig.History)),
		BaseLayer: "",
	}

	for k := range imageConfig.DiffIDs {
		rootfs.Layers = append(rootfs.Layers, k)
	}

	inspectData := &types.ImageInspect{
		RepoTags:        imageConfig.Tags,
		RepoDigests:     imageConfig.Digests,
		Parent:          imageConfig.Parent,
		Comment:         imageConfig.Comment,
		Created:         imageConfig.Created.Format(time.RFC3339Nano),
		Container:       imageConfig.Container,
		ContainerConfig: &imageConfig.ContainerConfig,
		DockerVersion:   imageConfig.DockerVersion,
		Author:          imageConfig.Author,
		Config:          imageConfig.Config,
		Architecture:    imageConfig.Architecture,
		Os:              imageConfig.OS,
		Size:            imageConfig.Size,
		VirtualSize:     imageConfig.Size,
		RootFS:          rootfs,
	}

	inspectData.GraphDriver.Name = productName + " " + PortlayerName

	//imageid is currently stored within VIC without "sha256:" so we add it to
	//match Docker
	inspectData.ID = "sha256:" + imageConfig.ImageID

	return inspectData
}
