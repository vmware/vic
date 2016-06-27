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

package vicbackends

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/context"

	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/registry"
	"github.com/vmware/vic/lib/metadata"
	"github.com/vmware/vic/pkg/trace"
)

// byCreated is a temporary type used to sort a list of images by creation
// time.
type byCreated []*types.Image

func (r byCreated) Len() int           { return len(r) }
func (r byCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }

type Image struct {
	ProductName string
}

func (i *Image) Commit(name string, config *types.ContainerCommitConfig) (imageID string, err error) {
	return "", fmt.Errorf("%s does not implement image.Commit", i.ProductName)
}

func (i *Image) Exists(containerName string) bool {
	return false
}

func (i *Image) ImageDelete(imageRef string, force, prune bool) ([]types.ImageDelete, error) {
	return []types.ImageDelete{}, fmt.Errorf("%s does not implement image.Delete", i.ProductName)
}

func (i *Image) ImageHistory(imageName string) ([]*types.ImageHistory, error) {
	return nil, fmt.Errorf("%s does not implement image.History", i.ProductName)
}

func (i *Image) Images(filterArgs string, filter string, all bool) ([]*types.Image, error) {
	defer trace.End(trace.Begin("Images"))

	imageCache := ImageCache()

	images, err := imageCache.GetImages()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving image list: %s", err)
	}

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

	imageConfig, err := getImageConfigFromCache(name)
	if err != nil {
		return nil, err
	}

	return imageConfigToDockerImageInspect(imageConfig, i.ProductName), nil
}

func (i *Image) TagImage(newTag reference.Named, imageName string) error {
	return fmt.Errorf("%s does not implement image.Tag", i.ProductName)
}

func (i *Image) LoadImage(inTar io.ReadCloser, outStream io.Writer, quiet bool) error {
	return fmt.Errorf("%s does not implement image.LoadImage", i.ProductName)
}

func (i *Image) ImportImage(src string, newRef reference.Named, msg string, inConfig io.ReadCloser, outStream io.Writer, config *container.Config) error {
	return fmt.Errorf("%s does not implement image.ImportImage", i.ProductName)
}

func (i *Image) ExportImage(names []string, outStream io.Writer) error {
	return fmt.Errorf("%s does not implement image.ExportImage", i.ProductName)
}

func (i *Image) PullImage(ctx context.Context, ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	defer trace.End(trace.Begin("PullImage"))

	log.Debugf("PullImage: ref = %+v, metaheaders = %+v\n", ref, metaHeaders)

	var cmdArgs []string

	cmdArgs = append(cmdArgs, "-reference", ref.String())

	if authConfig != nil {
		if len(authConfig.Username) > 0 {
			cmdArgs = append(cmdArgs, "-username", authConfig.Username)
		}
		if len(authConfig.Password) > 0 {
			cmdArgs = append(cmdArgs, "-password", authConfig.Password)
		}
	}

	portLayerServer := PortLayerServer()

	if portLayerServer != "" {
		cmdArgs = append(cmdArgs, "-host", portLayerServer)
	}

	// intruct imagec to use os.TempDir
	cmdArgs = append(cmdArgs, "-destination", os.TempDir())

	log.Debugf("PullImage: cmd = %s %+v\n", Imagec, cmdArgs)

	cmd := exec.Command(Imagec, cmdArgs...)
	cmd.Stdout = outStream
	cmd.Stderr = outStream

	// Execute
	err := cmd.Start()

	if err != nil {
		log.Debugf("Error starting %s - %s\n", Imagec, err)
		return fmt.Errorf("Error starting %s - %s\n", Imagec, err)
	}

	err = cmd.Wait()

	if err != nil {
		log.Debugf("imagec exit code: %s", err)
		return err
	}

	client := PortLayerClient()
	ImageCache().Update(client)
	return nil
}

func (i *Image) PushImage(ctx context.Context, ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return fmt.Errorf("%s does not implement image.PushImage", i.ProductName)
}

func (i *Image) SearchRegistryForImages(ctx context.Context, term string, authConfig *types.AuthConfig, metaHeaders map[string][]string) (*registry.SearchResults, error) {
	return nil, fmt.Errorf("%s does not implement image.SearchRegistryForImages", i.ProductName)
}

// Utility functions

func convertV1ImageToDockerImage(image *metadata.ImageConfig) *types.Image {
	var labels map[string]string
	if image.Config != nil {
		labels = image.Config.Labels
	}

	// TODO(jzt): change ImageConfig to contain a map from image name to all of its tags
	repoTag := fmt.Sprintf("%s:%s", image.Name, image.Tag)
	repoDigest := fmt.Sprintf("%s:%s", image.Name, image.Digest)

	return &types.Image{
		ID:          image.ImageID,
		ParentID:    image.Parent,
		RepoTags:    []string{repoTag},
		RepoDigests: []string{repoDigest},
		Created:     image.Created.Unix(),
		Size:        image.Size,
		VirtualSize: image.Size,
		Labels:      labels,
	}
}

// Retrieve the image metadata from the image cache.
func getImageConfigFromCache(image string) (*metadata.ImageConfig, error) {
	// Use docker reference code to validate the id's format
	digest, named, err := reference.ParseIDOrReference(image)
	if err != nil {
		return nil, err
	}

	// Try to get the image config from our image cache
	imageCache := ImageCache()

	if digest != "" {
		config, err := imageCache.GetImageByDigest(digest)
		if err != nil {
			log.Errorf("Inspect lookup failed for image %s: %s.  Returning no such image.", image, err)
			return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", image))
		}
		if config != nil {
			return config, nil
		}
	} else {
		config, err := imageCache.GetImageByNamed(named)
		if err != nil {
			log.Errorf("Inspect lookup failed for image %s: %s.  Returning no such image.", image, err)
			return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", image))
		}
		if config != nil {
			return config, nil
		}
	}

	return nil, derr.NewRequestNotFoundError(fmt.Errorf("No such image: %s", image))
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
		RepoTags:        make([]string, 0),
		RepoDigests:     make([]string, 0),
		Parent:          imageConfig.Parent,
		Comment:         imageConfig.Comment,
		Created:         imageConfig.Created.String(),
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

	//FIXME: Image tags storage is not yet fully implemented in the portlayer.
	//The following code dealing with tags needs to be revisited.
	taggedName := imageConfig.Name + ":" + imageConfig.Tag
	inspectData.RepoTags = append(inspectData.RepoTags, taggedName)
	inspectData.GraphDriver.Name = productName + " " + PortlayerName

	//imageid is currently stored within VIC without "sha256:" so we add it to
	//match Docker
	inspectData.ID = "sha256:" + imageConfig.ImageID

	return inspectData
}
