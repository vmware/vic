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
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/registry"
)

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
	return nil, fmt.Errorf("%s does not implement image.Images", i.ProductName)
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
