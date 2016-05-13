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

package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/portlayer/attach"
)

// AttachHandlersImpl is the receiver for all container attach methods.
type AttachHandlersImpl struct {
	s *attach.Server
}

func (a *AttachHandlersImpl) Configure(api *operations.PortLayerAPI, _ *HandlerContext) {
	// XXX this needs to live on the mgmt netwerk
	a.s = attach.NewAttachServer("", 0)

	if err := a.s.Start(); err != nil {
		log.Fatalf("Attach server unable to start: %s", err)
		return
	}
}
