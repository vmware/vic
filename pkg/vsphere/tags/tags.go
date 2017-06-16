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
// +build linux

package tags

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

const (
	TagURL = "/com/vmware/cis/tagging/tag"
)

type TagCreateSpec struct {
	CreateSpec TagCreate `json:"create_spec"`
}

type TagCreate struct {
	CategoryId  string `json:"category_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Tag struct {
	Id          string   `json:"id"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	CategoryId  string   `json:"category_id"`
	UsedBy      []string `json:"used_by"`
}

func (c *RestClient) CreateTagIfNotExist(name string, description string, categoryId string) (*string, error) {
	tagCreate := TagCreate{categoryId, description, name}
	spec := TagCreateSpec{tagCreate}
	id, err := c.CreateTag(&spec)
	if err != nil {
		log.Debugf("Created tag %s failed for %s", errors.ErrorStack(err))
		// if already exists, query back
		if strings.Contains(err.Error(), "already_exists") {
			tagObjs, err := c.GetTagByNameForCategory(name, categoryId)
			if err != nil {
				log.Errorf("Failed to query tag %s for category %s", name, categoryId)
				return nil, errors.Trace(err)
			}
			if tagObjs != nil {
				return &tagObjs[0].Id, nil
			} else {
				// should not happen
				log.Debugf("Failed to create tag for it's exsited, but could not query back. Please check system")
				return nil, errors.Errorf("Failed to create tag for it's exsited, but could not query back. Please check system")
			}
		} else {
			log.Debugf("Failed to create inventory category for %s", errors.ErrorStack(err))
			return nil, errors.Trace(err)
		}
	}

	return id, nil
}

func (c *RestClient) DeleteTagIfNoObjectAttached(id string) error {
	objs, err := c.ListAttachedObjects(id)
	if err != nil {
		return errors.Trace(err)
	}
	if objs != nil && len(objs) > 0 {
		log.Debugf("tag %s related objects is not empty, do not delete it.", id)
		return nil
	}
	return c.DeleteTag(id)
}

func (c *RestClient) CreateTag(spec *TagCreateSpec) (*string, error) {
	log.Debugf("Create Tag %v", spec)
	stream, _, status, err := c.call("POST", TagURL, spec, nil)

	log.Debugf("Get status code: %d", status)
	if status != 200 || err != nil {
		log.Debugf("Create tag failed with status code: %d, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	type RespValue struct {
		Value string
	}

	var pId RespValue
	if err := json.NewDecoder(stream).Decode(&pId); err != nil {
		log.Debugf("Decode response body failed for: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}
	return &(pId.Value), nil
}

func (c *RestClient) GetTag(id string) (*Tag, error) {
	log.Debugf("Get tag %s", id)

	stream, _, status, err := c.call("GET", fmt.Sprintf("%s/id:%s", TagURL, id), nil, nil)

	if status != 200 || err != nil {
		log.Debugf("Get tag failed with status code: %s, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	type RespValue struct {
		Value Tag
	}

	var pTag RespValue
	if err := json.NewDecoder(stream).Decode(&pTag); err != nil {
		log.Debugf("Decode response body failed for: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}
	return &(pTag.Value), nil
}

func (c *RestClient) DeleteTag(id string) error {
	log.Debugf("Delete tag %s", id)

	_, _, status, err := c.call("DELETE", fmt.Sprintf("%s/id:%s", TagURL, id), nil, nil)

	if status != 200 || err != nil {
		log.Debugf("Delete tag failed with status code: %s, error message: %s", status, errors.ErrorStack(err))
		return errors.Errorf("Status code: %d, error: %s", status, err)
	}
	return nil
}

func (c *RestClient) ListTags() ([]string, error) {
	log.Debugf("List all tags")

	stream, _, status, err := c.call("GET", TagURL, nil, nil)

	if status != 200 || err != nil {
		log.Debugf("Get tags failed with status code: %s, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	return c.handleTagIdList(stream)
}

func (c *RestClient) ListTagsForCategory(id string) ([]string, error) {
	log.Debugf("List tags for category: %s", id)

	type PostCategory struct {
		CId string `json:"category_id"`
	}
	spec := PostCategory{id}
	stream, _, status, err := c.call("POST", fmt.Sprintf("%s/id:%s?~action=list-tags-for-category", TagURL, id), spec, nil)

	if status != 200 || err != nil {
		log.Debugf("List tags for category failed with status code: %s, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Status code: %d, error: %s", status, err)
	}

	return c.handleTagIdList(stream)
}

func (c *RestClient) handleTagIdList(stream io.ReadCloser) ([]string, error) {
	type Tags struct {
		Value []string
	}

	var pTags Tags
	if err := json.NewDecoder(stream).Decode(&pTags); err != nil {
		log.Debugf("Decode response body failed for: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}
	return pTags.Value, nil
}

// Get tag through tag name and category id
func (c *RestClient) GetTagByNameForCategory(name string, id string) ([]Tag, error) {
	log.Debugf("Get tag %s for category %s", name, id)
	tagIds, err := c.ListTagsForCategory(id)
	if err != nil {
		log.Debugf("Get tag failed for %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}

	var tags []Tag
	for _, tId := range tagIds {
		tag, err := c.GetTag(tId)
		if err != nil {
			log.Debugf("Get tag %s failed for %s", tId, errors.ErrorStack(err))
			return nil, errors.Trace(err)
		}
		if tag.Name == name {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

// Get attached tags through tag name pattern
func (c *RestClient) GetAttachedTagsByNamePattern(namePattern string, objId string, objType string) ([]Tag, error) {
	tagIds, err := c.ListAttachedTags(objId, objType)
	if err != nil {
		log.Debugf("Get attached tags failed for %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}

	var validName = regexp.MustCompile(namePattern)
	var tags []Tag
	for _, tId := range tagIds {
		tag, err := c.GetTag(tId)
		if err != nil {
			log.Debugf("Get tag %s failed for %s", tId, errors.ErrorStack(err))
		}
		if validName.MatchString(tag.Name) {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}
