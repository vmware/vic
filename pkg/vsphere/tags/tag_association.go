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

package tags

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

const (
	TagAssociationURL = "/com/vmware/cis/tagging/tag-association"
)

type AssociatedObject struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
}

type TagAssociationSpec struct {
	ObjectID *AssociatedObject `json:"object_id,omitempty"`
	TagID    *string           `json:"tag_id,omitempty"`
}

func (c *RestClient) getAssociatedObject(objID *string, objType *string) *AssociatedObject {
	if objID == nil && objType == nil {
		return nil
	}
	object := AssociatedObject{
		ID:   objID,
		Type: objType,
	}
	return &object
}

func (c *RestClient) getAssociationSpec(tagID *string, objID *string, objType *string) *TagAssociationSpec {
	object := c.getAssociatedObject(objID, objType)
	spec := TagAssociationSpec{
		TagID:    tagID,
		ObjectID: object,
	}
	return &spec
}

func (c *RestClient) AttachTagToObject(ctx context.Context, tagID string, objID string, objType string) error {
	log.Debugf("Attach Tag %s to object id: %s, type: %s", tagID, objID, objType)

	spec := c.getAssociationSpec(&tagID, &objID, &objType)
	_, _, status, err := c.call(ctx, "POST", fmt.Sprintf("%s?~action=attach", TagAssociationURL), *spec, nil)

	log.Debugf("Get status code: %d", status)
	if status != 200 || err != nil {
		log.Debugf("Attach tag failed with status code: %d, error message: %s", status, errors.ErrorStack(err))
		return errors.Errorf("Get unexpected status code: %d", status)
	}
	return nil
}

func (c *RestClient) DetachTagFromObject(ctx context.Context, tagID string, objID string, objType string) error {
	log.Debugf("Detach Tag %s to object id: %s, type: %s", tagID, objID, objType)

	spec := c.getAssociationSpec(&tagID, &objID, &objType)
	_, _, status, err := c.call(ctx, "POST", fmt.Sprintf("%s?~action=detach", TagAssociationURL), *spec, nil)

	log.Debugf("Get status code: %d", status)
	if status != 200 || err != nil {
		log.Debugf("Detach tag failed with status code: %d, error message: %s", status, errors.ErrorStack(err))
		return errors.Errorf("Get unexpected status code: %d", status)
	}
	return nil
}

func (c *RestClient) ListAttachedTags(ctx context.Context, objID string, objType string) ([]string, error) {
	log.Debugf("List attached tags of object id: %s, type: %s", objID, objType)

	spec := c.getAssociationSpec(nil, &objID, &objType)
	stream, _, status, err := c.call(ctx, "POST", fmt.Sprintf("%s?~action=list-attached-tags", TagAssociationURL), *spec, nil)

	log.Debugf("Get status code: %d", status)
	if status != 200 || err != nil {
		log.Debugf("Detach tag failed with status code: %d, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Get unexpected status code: %d", status)
	}

	type RespValue struct {
		Value []string
	}

	var pTag RespValue
	if err := json.NewDecoder(stream).Decode(&pTag); err != nil {
		log.Debugf("Decode response body failed for: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}
	return pTag.Value, nil
}

func (c *RestClient) ListAttachedObjects(ctx context.Context, tagID string) ([]AssociatedObject, error) {
	log.Debugf("List attached objects of tag: %s", tagID)

	spec := c.getAssociationSpec(&tagID, nil, nil)
	log.Debugf("List attached objects for tag %v", *spec)
	//	stream, _, status, err := c.call("POST", fmt.Sprintf("%s?~action=list-attached-objects", TagAssociationURL), *spec, nil)
	stream, _, status, err := c.call(ctx, "POST", fmt.Sprintf("%s?~action=list-attached-objects", TagAssociationURL), *spec, nil)
	log.Debugf("Get status code: %d", status)
	if status != 200 || err != nil {
		log.Debugf("List object failed with status code: %d, error message: %s", status, errors.ErrorStack(err))
		return nil, errors.Errorf("Get unexpected status code: %d", status)
	}

	type RespValue struct {
		Value []AssociatedObject
	}

	var pTag RespValue
	if err := json.NewDecoder(stream).Decode(&pTag); err != nil {
		log.Debugf("Decode response body failed for: %s", errors.ErrorStack(err))
		return nil, errors.Trace(err)
	}
	return pTag.Value, nil
}
