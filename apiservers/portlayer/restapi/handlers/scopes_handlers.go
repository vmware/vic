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
	"net/http"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/scopes"
)

// ScopesHandlersImpl is the receiver for all of the storage handler methods
type ScopesHandlersImpl struct {
}

// Configure assigns functions to all the storage api handlers
func (handler *ScopesHandlersImpl) Configure(api *operations.PortLayerAPI) {
	api.ScopesCreateHandler = scopes.CreateHandlerFunc(handler.ScopesCreate)
	api.ScopesListAllHandler = scopes.ListAllHandlerFunc(handler.ScopesListAll)
	api.ScopesListHandler = scopes.ListHandlerFunc(handler.ScopesList)
}

func (handler *ScopesHandlersImpl) ScopesCreate(params scopes.CreateParams) middleware.Responder {
	return scopes.NewCreateDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "Not implemented"})
}

func (handler *ScopesHandlersImpl) ScopesListAll() middleware.Responder {
	return scopes.NewListAllDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "Not implemented"})
}

func (handler *ScopesHandlersImpl) ScopesList(params scopes.ListParams) middleware.Responder {
	return scopes.NewListDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "Not implemented"})
}
