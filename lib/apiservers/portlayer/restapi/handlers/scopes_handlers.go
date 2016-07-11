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

	log "github.com/Sirupsen/logrus"

	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/vmware/vic/lib/apiservers/portlayer/models"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations/scopes"
	"github.com/vmware/vic/lib/portlayer/exec"
	"github.com/vmware/vic/lib/portlayer/network"

	"github.com/vmware/vic/pkg/trace"
)

// ScopesHandlersImpl is the receiver for all of the storage handler methods
type ScopesHandlersImpl struct {
	netCtx     *network.Context
	handlerCtx *HandlerContext
}

// Configure assigns functions to all the scopes api handlers
func (handler *ScopesHandlersImpl) Configure(api *operations.PortLayerAPI, handlerCtx *HandlerContext) {
	api.ScopesCreateScopeHandler = scopes.CreateScopeHandlerFunc(handler.ScopesCreate)
	api.ScopesDeleteScopeHandler = scopes.DeleteScopeHandlerFunc(handler.ScopesDelete)
	api.ScopesListAllHandler = scopes.ListAllHandlerFunc(handler.ScopesListAll)
	api.ScopesListHandler = scopes.ListHandlerFunc(handler.ScopesList)
	api.ScopesAddContainerHandler = scopes.AddContainerHandlerFunc(handler.ScopesAddContainer)
	api.ScopesRemoveContainerHandler = scopes.RemoveContainerHandlerFunc(handler.ScopesRemoveContainer)
	api.ScopesBindContainerHandler = scopes.BindContainerHandlerFunc(handler.ScopesBindContainer)
	api.ScopesUnbindContainerHandler = scopes.UnbindContainerHandlerFunc(handler.ScopesUnbindContainer)

	err := network.Init()
	if err != nil {
		log.Fatalf("failed to create network context: %s", err)
	}

	handler.netCtx = network.DefaultContext
	handler.handlerCtx = handlerCtx
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

func errorPayload(err error) *models.Error {
	return &models.Error{Message: err.Error()}
}

func (handler *ScopesHandlersImpl) ScopesCreate(params scopes.CreateScopeParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesCreate"))

	cfg := params.Config
	if cfg.ScopeType == "external" {
		return scopes.NewCreateScopeDefault(http.StatusServiceUnavailable).WithPayload(&models.Error{Message: "cannot create external networks"})
	}

	subnet, gateway, dns, err := parseScopeConfig(cfg)
	if err != nil {
		return scopes.NewCreateScopeDefault(http.StatusServiceUnavailable).WithPayload(errorPayload(err))
	}

	s, err := handler.netCtx.NewScope(cfg.ScopeType, cfg.Name, subnet, gateway, dns, cfg.IPAM)
	if _, ok := err.(network.DuplicateResourceError); ok {
		return scopes.NewCreateScopeConflict()
	}

	if err != nil {
		return scopes.NewCreateScopeDefault(http.StatusServiceUnavailable).WithPayload(errorPayload(err))
	}

	return scopes.NewCreateScopeCreated().WithPayload(toScopeConfig(s))
}

func (handler *ScopesHandlersImpl) ScopesDelete(params scopes.DeleteScopeParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesDelete"))

	if err := handler.netCtx.DeleteScope(params.IDName); err != nil {
		switch err := err.(type) {
		case network.ResourceNotFoundError:
			return scopes.NewDeleteScopeNotFound().WithPayload(errorPayload(err))

		default:
			return scopes.NewDeleteScopeInternalServerError().WithPayload(errorPayload(err))
		}
	}

	return scopes.NewDeleteScopeOK()
}

func (handler *ScopesHandlersImpl) ScopesListAll() middleware.Responder {
	defer trace.End(trace.Begin("ScopesListAll"))

	cfgs, err := listScopes(handler.netCtx, "")
	if err != nil {
		return scopes.NewListDefault(http.StatusServiceUnavailable).WithPayload(errorPayload(err))
	}

	return scopes.NewListAllOK().WithPayload(cfgs)
}

func (handler *ScopesHandlersImpl) ScopesList(params scopes.ListParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesList"))

	cfgs, err := listScopes(handler.netCtx, params.IDName)
	if _, ok := err.(network.ResourceNotFoundError); ok {
		return scopes.NewListNotFound().WithPayload(errorPayload(err))
	}

	return scopes.NewListOK().WithPayload(cfgs)
}

func (handler *ScopesHandlersImpl) ScopesAddContainer(params scopes.AddContainerParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesAddContainer"))

	h := exec.GetHandle(params.Config.Handle)
	if h == nil {
		return scopes.NewAddContainerNotFound().WithPayload(&models.Error{Message: "container not found"})
	}

	err := func() error {
		var ip *net.IP
		if params.Config.NetworkConfig.Address != nil && *params.Config.NetworkConfig.Address != "" {
			i := net.ParseIP(*params.Config.NetworkConfig.Address)
			if i == nil {
				return fmt.Errorf("invalid ip address %q", *params.Config.NetworkConfig.Address)
			}

			ip = &i
		}

		if params.Config.NetworkConfig.Aliases != nil {
			log.Debugf("Links/Aliases: %#v", params.Config.NetworkConfig.Aliases)
		}

		options := &network.AddContainerOptions{
			Scope:   params.Config.NetworkConfig.NetworkName,
			IP:      ip,
			Aliases: params.Config.NetworkConfig.Aliases,
		}
		return handler.netCtx.AddContainer(h, options)
	}()

	if err != nil {
		if _, ok := err.(*network.ResourceNotFoundError); ok {
			return scopes.NewAddContainerNotFound().WithPayload(errorPayload(err))
		}

		return scopes.NewAddContainerInternalServerError().WithPayload(errorPayload(err))
	}

	return scopes.NewAddContainerOK().WithPayload(h.String())
}

func (handler *ScopesHandlersImpl) ScopesRemoveContainer(params scopes.RemoveContainerParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesRemoveContainer"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return scopes.NewRemoveContainerNotFound().WithPayload(&models.Error{Message: "container not found"})
	}

	if err := handler.netCtx.RemoveContainer(h, params.Scope); err != nil {
		if _, ok := err.(*network.ResourceNotFoundError); ok {
			return scopes.NewRemoveContainerNotFound().WithPayload(errorPayload(err))
		}

		return scopes.NewRemoveContainerInternalServerError().WithPayload(errorPayload(err))
	}

	return scopes.NewRemoveContainerOK().WithPayload(h.String())
}

func (handler *ScopesHandlersImpl) ScopesBindContainer(params scopes.BindContainerParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesBindContainer"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return scopes.NewBindContainerNotFound().WithPayload(&models.Error{Message: "container not found"})
	}

	if _, err := handler.netCtx.BindContainer(h); err != nil {
		return scopes.NewBindContainerInternalServerError().WithPayload(errorPayload(err))
	}

	return scopes.NewBindContainerOK().WithPayload(h.String())
}

func (handler *ScopesHandlersImpl) ScopesUnbindContainer(params scopes.UnbindContainerParams) middleware.Responder {
	defer trace.End(trace.Begin("ScopesUnbindContainer"))

	h := exec.GetHandle(params.Handle)
	if h == nil {
		return scopes.NewUnbindContainerNotFound()
	}

	if err := handler.netCtx.UnbindContainer(h); err != nil {
		return scopes.NewUnbindContainerInternalServerError().WithPayload(errorPayload(err))
	}

	return scopes.NewUnbindContainerOK().WithPayload(h.String())
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
