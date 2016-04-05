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
	"fmt"
	"net"
	"net/http"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/apiservers/portlayer/models"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/apiservers/portlayer/restapi/operations/scopes"
	"github.com/vmware/vic/portlayer/network"
)

// ScopesHandlersImpl is the receiver for all of the storage handler methods
type ScopesHandlersImpl struct {
	netCtx *network.Context
}

// Configure assigns functions to all the scopes api handlers
func (handler *ScopesHandlersImpl) Configure(api *operations.PortLayerAPI, netCtx *network.Context) {
	api.ScopesCreateHandler = scopes.CreateHandlerFunc(handler.ScopesCreate)
	api.ScopesListAllHandler = scopes.ListAllHandlerFunc(handler.ScopesListAll)
	api.ScopesListHandler = scopes.ListHandlerFunc(handler.ScopesList)

	handler.netCtx = netCtx
}

func parseScopeConfig(cfg *models.ScopeConfig) (subnet *net.IPNet, gateway net.IP, dns []net.IP, err error) {
	if cfg.Subnet != nil {
		if _, subnet, err = net.ParseCIDR(*cfg.Subnet); err != nil {
			return
		}
	}

	gateway = net.IPv4(0, 0, 0, 0)
	if cfg.Gateway != nil {
		if gateway = net.ParseIP(*cfg.Gateway); gateway == nil {
			err = fmt.Errorf("invalid gateway")
			return
		}
	}

	dns = make([]net.IP, len(cfg.DNS))
	for i, d := range cfg.DNS {
		dns[i] = net.ParseIP(d)
		if dns[i] == nil {
			err = fmt.Errorf("invalid dns entry")
			return
		}
	}

	return
}

func listScopes(ctx *network.Context, idName string) ([]*models.ScopeConfig, error) {
	_scopes, err := ctx.Scopes(&idName)
	if err != nil {
		return nil, err
	}

	cfgs := make([]*models.ScopeConfig, len(_scopes))
	for i, s := range _scopes {
		cfgs[i] = toScopeConfig(s)
	}

	return cfgs, nil
}

func (handler *ScopesHandlersImpl) ScopesCreate(params scopes.CreateParams) middleware.Responder {
	cfg := params.Config
	if cfg.ScopeType == "external" {
		return scopes.NewCreateDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "cannot create external networks"})
	}

	subnet, gateway, dns, err := parseScopeConfig(cfg)
	if err != nil {
		return scopes.NewCreateDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
	}

	s, err := handler.netCtx.NewScope(cfg.ScopeType, cfg.Name, subnet, gateway, dns, cfg.IPAM)
	if _, ok := err.(network.DuplicateResourceError); ok {
		return scopes.NewCreateConflict()
	}

	if err != nil {
		return scopes.NewCreateDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
	}

	return scopes.NewCreateCreated().WithPayload(toScopeConfig(s))
}

func (handler *ScopesHandlersImpl) ScopesListAll() middleware.Responder {
	cfgs, err := listScopes(handler.netCtx, "")
	if err != nil {
		return scopes.NewListDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: err.Error()})
	}

	return scopes.NewListAllOK().WithPayload(cfgs)
}

func (handler *ScopesHandlersImpl) ScopesList(params scopes.ListParams) middleware.Responder {
	cfgs, err := listScopes(handler.netCtx, params.IDName)
	if _, ok := err.(network.ResourceNotFoundError); ok {
		return scopes.NewListNotFound()
	}

	return scopes.NewListOK().WithPayload(cfgs)
}

func toScopeConfig(scope *network.Scope) *models.ScopeConfig {
	id := scope.ID()
	subnet := scope.Subnet().String()
	gateway := ""
	if !scope.Gateway().IsUnspecified() {
		gateway = scope.Gateway().String()
	}
	return &models.ScopeConfig{
		ID:        &id,
		Name:      scope.Name(),
		ScopeType: scope.Type(),
		IPAM:      scope.IPAM().Pools(),
		Subnet:    &subnet,
		Gateway:   &gateway,
	}
}
