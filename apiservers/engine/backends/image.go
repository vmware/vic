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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	derr "github.com/docker/docker/errors"
	v1 "github.com/docker/docker/image"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/registry"
	"github.com/vmware/vic/apiservers/portlayer/client/storage"
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

	result := []*types.Image{}

	host, err := os.Hostname()
	if err != nil {
		return result,
			derr.NewErrorWithStatusCode(fmt.Errorf("image.Images got unexpected error getting hostname"),
				http.StatusInternalServerError)
	}

	params := storage.NewListImagesParams().WithStoreName(host)

	client := PortLayerClient()
	if client == nil {
		return result,
			derr.NewErrorWithStatusCode(fmt.Errorf("image.Images failed to create a portlayer client"),
				http.StatusInternalServerError)
	}

	images, err := client.Storage.ListImages(params)
	if err != nil {
		log.Warning(err)
		return result, nil
	}

	// build a map from image id to image v1Compatibility metadata
	v1Map := make(map[string]*v1.V1Image)
	for _, image := range images.Payload {

		if image.ID == "scratch" {
			continue
		}

		v1Image := &v1.V1Image{}
		decoder := json.NewDecoder(strings.NewReader(image.Metadata["v1Compatibility"]))
		if err := decoder.Decode(v1Image); err != nil {
			log.Fatal(err)
		}
		v1Map[image.ID] = v1Image
	}

	resolveImageSizes(v1Map)

	for _, v1Image := range v1Map {
		newImage := &types.Image{
			ID:       v1Image.ID,
			ParentID: v1Image.Parent,
			RepoTags: []string{"<fixme>:<fixme>"}, // TODO(jzt): replace with actual tags
			//RepoDigests: []string{"<fixme>:<fixme>"},
			Created:     v1Image.Created.Unix(),
			VirtualSize: v1Image.Size,
			Labels:      v1Image.Config.Labels,
		}
		result = append(result, newImage)
	}

	sort.Sort(sort.Reverse(byCreated(result)))

	return result, nil
}

func resolveImageSizes(v1Map map[string]*v1.V1Image) {
	resolved := map[string]int64{}
	for id := range v1Map {
		resolveImageSize(v1Map, resolved, id)
	}
}

func resolveImageSize(v1Map map[string]*v1.V1Image, resolved map[string]int64, id string) {

	// this image's size has already been resolved by a previous traversal
	if _, ok := resolved[id]; ok {
		return
	}

	parentID := v1Map[id].Parent
	if parentID != "" {
		// If we have already resolved the parent image's size, we don't need to traverse up the chain
		if parentSize, ok := resolved[parentID]; ok {
			v1Map[id].Size += parentSize
		} else {
			resolveImageSize(v1Map, resolved, parentID)
			v1Map[id].Size += resolved[parentID]
		}
	}
	resolved[id] = v1Map[id].Size
}

func (i *Image) LookupImage(name string) (*types.ImageInspect, error) {
	return nil, fmt.Errorf("%s does not implement image.LookupImage", i.ProductName)
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

func (i *Image) PullImage(ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	log.Printf("PullImage: ref = %+v, metaheaders = %+v\n", ref, metaHeaders)

	binImageC := "imagec"

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

	fetcherPath := getImageFetcherPath(binImageC)

	log.Printf("PullImage: cmd = %s %+v\n", fetcherPath, cmdArgs)

	cmd := exec.Command(fetcherPath, cmdArgs...)
	cmd.Stdout = outStream
	cmd.Stderr = outStream

	// Execute
	err := cmd.Start()

	if err != nil {
		log.Printf("Error starting %s - %s\n", fetcherPath, err)
		return fmt.Errorf("Error starting %s - %s\n", fetcherPath, err)
	}

	err = cmd.Wait()

	if err != nil {
		log.Println("imagec exit code:", err)
		return err
	}

	return nil
}

func (i *Image) PushImage(ref reference.Named, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return fmt.Errorf("%s does not implement image.PushImage", i.ProductName)
}

func (i *Image) SearchRegistryForImages(term string, authConfig *types.AuthConfig, metaHeaders map[string][]string) (*registry.SearchResults, error) {
	return nil, fmt.Errorf("%s does not implement image.SearchRegistryForImages", i.ProductName)
}

func getImageFetcherPath(fetcherName string) string {
	fullpath := "./" + fetcherName

	dir, ferr := filepath.Abs(filepath.Dir(os.Args[0]))

	if ferr == nil {
		fullpath = dir + "/" + fetcherName
	}

	return fullpath
}
