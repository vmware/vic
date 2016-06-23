# Copyright 2016 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO ?= go
GOVERSION ?= go1.6.2
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ifeq (vagrant, $(filter vagrant,$(USER) $(SUDO_USER)))
	# assuming we are in a shared directory where host arch is different from the guest
	BIN_ARCH := -$(OS)
endif

BASE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BASE_PKG := github.com/vmware/vic/

BIN ?= bin
IGNORE := $(shell mkdir -p $(BIN))

export GOPATH ?= $(shell echo $(CURDIR) | sed -e 's,/src/.*,,')
SWAGGER ?= $(GOPATH)/bin/swagger$(BIN_ARCH)
GOIMPORTS ?= $(GOPATH)/bin/goimports$(BIN_ARCH)
GOLINT ?= $(GOPATH)/bin/golint$(BIN_ARCH)
GVT ?= $(GOPATH)/bin/gvt$(BIN_ARCH)
GOVC ?= $(GOPATH)/bin/govc$(BIN_ARCH)

.PHONY: all tools clean test check \
	goversion goimports gopath govet \
	isos tethers apiservers copyright

.DEFAULT_GOAL := all

ifeq ($(ENABLE_RACE_DETECTOR),true)
	RACE := -race
else
	RACE :=
endif

# utility function to dynamically generate go dependencies
define godeps
	$(wildcard $1) $(shell $(BASE_DIR)/infra/scripts/go-deps.sh $(dir $1) $(MAKEFLAGS))
endef

# target aliases - environment variable definition
docker-engine-api := $(BIN)/docker-engine-server
portlayerapi := $(BIN)/port-layer-server
portlayerapi-client := lib/apiservers/portlayer/client/port_layer_client.go
portlayerapi-server := lib/apiservers/portlayer/restapi/server.go

imagec := $(BIN)/imagec
vicadmin := $(BIN)/vicadmin
rpctool := $(BIN)/rpctool
vic-machine-linux := $(BIN)/vic-machine-linux
vic-machine-windows := $(BIN)/vic-machine-windows.exe
vic-machine-darwin := $(BIN)/vic-machine-darwin
vch-init := $(BIN)/vch-init
# NOT BUILT WITH make all TARGET
# vic-dns variants to create standalone DNS service.
vic-dns-linux := $(BIN)/vic-dns-linux
vic-dns-windows := $(BIN)/vic-dns-windows.exe
vic-dns-darwin := $(BIN)/vic-dns-darwin

tether-linux := $(BIN)/tether-linux
tether-windows := $(BIN)/tether-windows.exe
tether-darwin := $(BIN)/tether-darwin

appliance := $(BIN)/appliance.iso
appliance-staging := $(BIN)/appliance-staging.tgz
bootstrap := $(BIN)/bootstrap.iso
bootstrap-staging := $(BIN)/bootstrap-staging.tgz
bootstrap-staging-debug := $(BIN)/bootstrap-staging-debug.tgz
bootstrap-debug := $(BIN)/bootstrap-debug.iso
iso-base := $(BIN)/iso-base.tgz

go-lint := $(BIN)/.golint
go-imports := $(BIN)/.goimports

# target aliases - target mapping
docker-engine-api: $(docker-engine-api)
portlayerapi: $(portlayerapi)
portlayerapi-client: $(portlayerapi-client)
portlayerapi-server: $(portlayerapi-server)

imagec: $(imagec)
vicadmin: $(vicadmin)
rpctool: $(rpctool)
vch-init: $(vch-init)

tether-linux: $(tether-linux)
tether-windows: $(tether-windows)
tether-darwin: $(tether-darwin)

appliance: $(appliance)
appliance-staging: $(appliance-staging)
bootstrap: $(bootstrap)
bootstrap-staging: $(bootstrap-staging)
bootstrap-debug: $(bootstrap-debug)
bootstrap-staging-debug: $(bootstrap-staging-debug)
iso-base: $(iso-base)
vic-machine: $(vic-machine-linux) $(vic-machine-windows) $(vic-machine-darwin)
# NOT BUILT WITH make all TARGET
# vic-dns variants to create standalone DNS service.
vic-dns: $(vic-dns-linux) $(vic-dns-windows) $(vic-dns-darwin)

swagger: $(SWAGGER)

golint: $(go-lint)
goimports: $(go-imports)


# convenience targets
all: components tethers isos vic-machine
tools: $(GOIMPORTS) $(GVT) $(GOLINT) $(SWAGGER) goversion
check: goversion goimports govet golint copyright whitespace
apiservers: $(portlayerapi) $(docker-engine-api)
components: check apiservers $(imagec) $(vicadmin) $(rpctool)
isos: $(appliance) $(bootstrap)
tethers: $(tether-linux) $(tether-windows) $(tether-darwin)

# utility targets
goversion:
	@echo checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

$(GOIMPORTS): vendor/manifest
	@echo building $(GOIMPORTS)...
	@$(GO) build $(RACE) -o $(GOIMPORTS) ./vendor/golang.org/x/tools/cmd/goimports

$(GVT): vendor/manifest
	@echo building $(GVT)...
	@$(GO) build $(RACE) -o $(GVT) ./vendor/github.com/FiloSottile/gvt

$(GOLINT): vendor/manifest
	@echo building $(GOLINT)...
	@$(GO) build $(RACE) -o $(GOLINT) ./vendor/github.com/golang/lint/golint

$(SWAGGER): vendor/manifest
	@echo building $(SWAGGER)...
	@$(GO) build $(RACE) -o $(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

$(GOVC): vendor/manifest
	@echo building $(GOVC)...
	@$(GO) build $(RACE) -o $(GOVC) ./vendor/github.com/vmware/govmomi/govc

copyright:
	@echo "checking copyright in header..."
	@infra/scripts/header-check.sh

whitespace:
	@echo "checking whitespace..."
	@infra/scripts/whitespace-check.sh

# exit 1 if golint complains about anything other than comments
golintf = $(GOLINT) $(1) | sh -c "! grep -v 'should have comment'" | sh -c "! grep -v 'comment on exported'"

$(go-lint): $(GOLINT)
	@echo checking go lint...
	@$(call golintf,github.com/vmware/vic/cmd/...)
	@$(call golintf,github.com/vmware/vic/pkg/...)
	@$(call golintf,github.com/vmware/vic/lib/install/...)
	@$(call golintf,github.com/vmware/vic/lib/portlayer/...)
	@$(call golintf,github.com/vmware/vic/lib/apiservers/portlayer/restapi/handlers/...)
	@$(call golintf,github.com/vmware/vic/lib/apiservers/engine/backends/...)
	@touch $@

# For use by external tools such as emacs or for example:
# GOPATH=$(make gopath) go get ...
gopath:
	@echo -n $(GOPATH)

$(go-imports): $(GOIMPORTS) $(find . -type f -name '*.go' -not -path "./vendor/*") $(PORTLAYER_DEPS)
	@echo checking go imports...
	@! $(GOIMPORTS) -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | egrep -v '^$$'
	@touch $@

govet:
	@echo checking go vet...
	@$(GO) tool vet -all $$(find . -mindepth 1 -maxdepth 1 -type d -not -name vendor)

vendor: $(GVT)
	@echo restoring vendor
	$(GVT) restore

TEST_DIRS=github.com/vmware/vic/cmd/tether
TEST_DIRS+=github.com/vmware/vic/cmd/imagec
TEST_DIRS+=github.com/vmware/vic/cmd/vicadmin
TEST_DIRS+=github.com/vmware/vic/cmd/rpctool
TEST_DIRS+=github.com/vmware/vic/cmd/vic-machine
TEST_DIRS+=github.com/vmware/vic/lib/apiservers/portlayer
TEST_DIRS+=github.com/vmware/vic/lib/install
TEST_DIRS+=github.com/vmware/vic/lib/portlayer
TEST_DIRS+=github.com/vmware/vic/pkg


test:
	@echo Running unit tests
	# test everything but vendor
ifdef DRONE
	@echo Generating coverage data
	@$(TIME) infra/scripts/coverage.sh $(TEST_DIRS)
else
	@echo Generating local html coverage report
	@$(TIME) infra/scripts/coverage.sh --html $(TEST_DIRS)
endif

$(vch-init): $(call godeps,cmd/vch-init/*.go)
	@echo building vch-init
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./$(dir $<)

$(tether-linux): $(call godeps,cmd/tether/*.go)
	@echo building tether-linux
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(TIME) $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./$(dir $<)

$(tether-windows): $(call godeps,cmd/tether/*.go)
	@echo building tether-windows
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(TIME) $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./$(dir $<)

# CGO is disabled for darwin otherwise build fails with "gcc: error: unrecognized command line option '-mmacosx-version-min=10.6'"
$(tether-darwin): $(call godeps,cmd/tether/*.go)
	@echo building tether-darwin
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(TIME) $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./$(dir $<)

$(rpctool): $(call godeps,cmd/rpctool/*.go)
ifeq ($(OS),linux)
	@echo building rpctool
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)
else
	@echo skipping rpctool, cannot cross compile cgo
endif

$(vicadmin): $(call godeps,cmd/vicadmin/*.go)
	@echo building vicadmin
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)

$(imagec): $(call godeps,cmd/imagec/*.go) $(portlayerapi-client)
	@echo building imagec...
	@$(TIME) $(GO) build $(RACE) -o ./$@ ./$(dir $<)


$(docker-engine-api): $(call godeps,cmd/docker/*.go) $(portlayerapi-client)
ifeq ($(OS),linux)
	@echo Building docker-engine-api server...
	@$(TIME) $(GO) build $(RACE) -o $@ ./cmd/docker
else
	@echo skipping docker-engine-api server, cannot build on non-linux
endif

# Common portlayer dependencies between client and server
PORTLAYER_DEPS ?= lib/apiservers/portlayer/swagger.yml \
				  lib/apiservers/portlayer/restapi/configure_port_layer.go \
				  lib/apiservers/portlayer/restapi/options/*.go \
				  lib/apiservers/portlayer/restapi/handlers/*.go

$(portlayerapi-client): $(PORTLAYER_DEPS)  $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API client...
	@$(SWAGGER) generate client -A PortLayer --template-dir lib/apiservers/templates  -t $(realpath $(dir $<)) -f $<


$(portlayerapi-server): $(PORTLAYER_DEPS) $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API server...
	@$(SWAGGER) generate server --exclude-main -A PortLayer --template-dir lib/apiservers/templates -t $(realpath $(dir $<)) -f $<

$(portlayerapi): $(call godeps,cmd/port-layer-server/*.go) $(portlayerapi-server) $(portlayerapi-client)
	@echo building Portlayer API server...
	@$(TIME) $(GO) build $(RACE) -o $@ ./cmd/port-layer-server

$(iso-base): isos/base.sh isos/base/*.repo isos/base/isolinux/** isos/base/xorriso-options.cfg
	@echo building iso-base docker image
	@$(TIME) $< -c $(BIN)/yum-cache.tgz -p $@

# appliance staging - allows for caching of package install
$(appliance-staging): isos/appliance-staging.sh $(iso-base)
	@echo staging for VCH appliance
	@$(TIME) $< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@

# main appliance target - depends on all top level component targets
$(appliance): isos/appliance.sh isos/appliance/* $(rpctool) $(vicadmin) $(imagec) $(vch-init) $(portlayerapi) $(docker-engine-api) $(appliance-staging)
	@echo building VCH appliance ISO
	@$(TIME) $< -p $(appliance-staging) -b $(BIN)

# main bootstrap target
$(bootstrap): isos/bootstrap.sh $(tether-linux) $(rpctool) $(bootstrap-staging) isos/bootstrap/*
	@echo "Making bootstrap iso"
	@$(TIME) $< -p $(bootstrap-staging) -b $(BIN)

$(bootstrap-debug): isos/bootstrap.sh $(tether-linux) $(rpctool) $(bootstrap-staging-debug) isos/bootstrap/*
	@echo "Making bootstrap-debug iso"
	@$(TIME) $< -p $(bootstrap-staging-debug) -b $(BIN) -d true

$(bootstrap-staging): isos/bootstrap-staging.sh $(iso-base)
	@echo staging for bootstrap
	@$(TIME) $< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@

$(bootstrap-staging-debug): isos/bootstrap-staging.sh $(iso-base)
	@echo staging debug for bootstrap
	@$(TIME) $< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@ -d true


$(vic-machine-linux): $(call godeps,cmd/vic-machine/*.go)
	@echo building vic-machine linux...
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -ldflags "-X github.com/vmware/vic/cmd/vic-machine.BuildID=$(BUILD_NUMBER)" -o ./$@ ./$(dir $<)

$(vic-machine-windows): $(call godeps,cmd/vic-machine/*.go)
	@echo building vic-machine windows...
	@GOARCH=amd64 GOOS=windows $(TIME) $(GO) build $(RACE) -ldflags "-X github.com/vmware/vic/cmd/vic-machine.BuildID=$(BUILD_NUMBER)" -o ./$@ ./$(dir $<)

$(vic-machine-darwin): $(call godeps,cmd/vic-machine/*.go)
	@echo building vic-machine darwin...
	@GOARCH=amd64 GOOS=darwin $(TIME) $(GO) build $(RACE) -ldflags "-X github.com/vmware/vic/cmd/vic-machine.BuildID=$(BUILD_NUMBER)" -o ./$@ ./$(dir $<)

$(vic-dns-linux): $(call godeps,cmd/vic-dns/*.go)
	@echo building vic-dns linux...
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -o ./$@ ./$(dir $<)

$(vic-dns-windows): $(call godeps,cmd/vic-dns/*.go)
	@echo building vic-dns windows...
	@GOARCH=amd64 GOOS=windows $(TIME) $(GO) build $(RACE) -o ./$@ ./$(dir $<)

$(vic-dns-darwin): $(call godeps,cmd/vic-dns/*.go)
	@echo building vic-dns darwin...
	@GOARCH=amd64 GOOS=darwin $(TIME) $(GO) build $(RACE) -o ./$@ ./$(dir $<)

clean:
	rm -rf $(BIN)

	@echo removing swagger generated files...
	rm -f ./lib/apiservers/portlayer/restapi/doc.go
	rm -f ./lib/apiservers/portlayer/restapi/embedded_spec.go
	rm -f ./lib/apiservers/portlayer/restapi/server.go
	rm -rf ./lib/apiservers/portlayer/client/
	rm -rf ./lib/apiservers/portlayer/cmd/
	rm -rf ./lib/apiservers/portlayer/models/
	rm -rf ./lib/apiservers/portlayer/restapi/operations/

	rm -f lib/apiservers/docker/restapi/doc.go
	rm -f lib/apiservers/docker/restapi/embedded_spec.go
	rm -f lib/apiservers/docker/restapi/server.go
	rm -fr lib/apiservers/docker/cmd
	rm -fr lib/apiservers/docker/models
	rm -fr lib/apiservers/docker/restapi/operations
