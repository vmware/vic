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

package restapi

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/swag"

	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/handlers"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/operations"
	"github.com/vmware/vic/lib/apiservers/portlayer/restapi/options"
	"github.com/vmware/vic/lib/portlayer"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/version"
	"github.com/vmware/vic/pkg/vsphere/session"

	"golang.org/x/net/context"
)

// This file is safe to edit. Once it exists it will not be overwritten

func init() {
	trace.Logger.Level = log.DebugLevel
}

type handler interface {
	Configure(api *operations.PortLayerAPI, handlerCtx *handlers.HandlerContext)
}

var portlayerhandlers = []handler{
	&handlers.StorageHandlersImpl{},
	&handlers.MiscHandlersImpl{},
	&handlers.ScopesHandlersImpl{},
	&handlers.ContainersHandlersImpl{},
	&handlers.InteractionHandlersImpl{},
	&handlers.LoggingHandlersImpl{},
	&handlers.KvHandlersImpl{},
}

func configureFlags(api *operations.PortLayerAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			LongDescription:  "Port Layer Options",
			Options:          options.PortLayerOptions,
			ShortDescription: "Port Layer Options",
		},
	}
}

func configureAPI(api *operations.PortLayerAPI) http.Handler {
	if options.PortLayerOptions.Debug {
		log.SetLevel(log.DebugLevel)
	}

	ctx := context.Background()

	sessionconfig := &session.Config{
		Service:        options.PortLayerOptions.SDK,
		Insecure:       options.PortLayerOptions.Insecure,
		Keepalive:      options.PortLayerOptions.Keepalive,
		DatacenterPath: options.PortLayerOptions.DatacenterPath,
		ClusterPath:    options.PortLayerOptions.ClusterPath,
		PoolPath:       options.PortLayerOptions.PoolPath,
		DatastorePath:  options.PortLayerOptions.DatastorePath,
		UserAgent:      version.UserAgent("vic-engine"),
	}

	sess, err := session.NewSession(sessionconfig).Create(ctx)
	if err != nil {
		log.Fatalf("configure_port_layer ERROR: %s", err)
	}

	// Configure the func invoked if the PL panics or is restarted by vic-init
	api.ServerShutdown = func() {
		log.Infof("Shutting down port-layer-server")

		// Logout the session
		if err := sess.Logout(ctx); err != nil {
			log.Warnf("unable to log out of session: %s", err)
		}
	}

	// initialize the port layer
	if err = portlayer.Init(ctx, sess); err != nil {
		log.Fatalf("could not initialize port layer: %s", err)
	}

	// configure the api here
	api.ServeError = errors.ServeError

	api.BinConsumer = httpkit.ByteStreamConsumer()

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.TxtProducer = httpkit.TextProducer()

	handlerCtx := &handlers.HandlerContext{
		Session: sess,
	}
	for _, handler := range portlayerhandlers {
		handler.Configure(api, handlerCtx)
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
