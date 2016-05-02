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
ifeq ($(USER),vagrant)
	# assuming we are in a shared directory where host arch is different from the guest
	BIN_ARCH := -$(OS)
endif

BASE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

BIN ?= bin
IGNORE := $(shell mkdir -p $(BIN))

export GOPATH ?= $(shell echo $(CURDIR) | sed -e 's,/src/.*,,')
SWAGGER ?= $(GOPATH)/bin/swagger$(BIN_ARCH)
GOIMPORTS ?= $(GOPATH)/bin/goimports$(BIN_ARCH)
GOLINT ?= $(GOPATH)/bin/golint$(BIN_ARCH)
GVT ?= $(GOPATH)/bin/gvt$(BIN_ARCH)

.PHONY: all tools clean test check \
	goversion goimports gvt gopath \
	isos tethers apiservers copyright

.DEFAULT_GOAL := all

ifeq ($(ENABLE_RACE_DETECTOR),true)
	RACE := -race
else
	RACE :=
endif

# target aliases - environment variable definition
docker-engine-api := $(BIN)/docker-engine-server
portlayerapi := $(BIN)/port-layer-server
portlayerapi-client := apiservers/portlayer/client/port_layer_client.go
portlayerapi-server := apiservers/portlayer/restapi/server.go

imagec := $(BIN)/imagec
vicadmin := $(BIN)/vicadmin
rpctool := $(BIN)/rpctool
vic-machine := $(BIN)/vic-machine

tether-linux := $(BIN)/tether-linux
tether-windows := $(BIN)/tether-windows.exe

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

tether-linux: $(tether-linux)
tether-windows: $(tether-windows)

appliance: $(appliance)
appliance-staging: $(appliance-staging)
bootstrap: $(bootstrap)
bootstrap-staging: $(bootstrap-staging)
bootstrap-debug: $(bootstrap-debug)
bootstrap-staging-debug: $(bootstrap-staging-debug)
iso-base: $(iso-base)
vic-machine: $(vic-machine)

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
tethers: $(tether-linux) $(tether-windows)


# utility targets
goversion:
	@echo checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

$(GOIMPORTS): vendor/manifest
	@echo building $(GOIMPORTS)...
	@$(GO) build $(RACE) -o $(GOIMPORTS) ./vendor/golang.org/x/tools/cmd/goimports

$(GVT):
	@echo getting gvt
	@$(GO) get -u github.com/FiloSottile/gvt

$(GOLINT): vendor/manifest
	@echo building $(GOLINT)...
	@$(GO) build $(RACE) -o $(GOLINT) ./vendor/github.com/golang/lint/golint

$(SWAGGER): vendor/manifest
	@echo building $(SWAGGER)...
	@$(GO) build $(RACE) -o $(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

copyright:
	@echo "checking copyright in header..."
	scripts/header-check.sh

whitespace:
	@echo "checking whitespace..."
	scripts/whitespace-check.sh

# exit 1 if golint complains about anything other than comments
golintf = $(GOLINT) $(1) | sh -c "! grep -v 'should have comment'" | sh -c "! grep -v 'comment on exported'"

$(go-lint): $(GOLINT)
	@echo checking go lint...
	@$(call golintf,github.com/vmware/vic/cmd/...)
	@$(call golintf,github.com/vmware/vic/pkg/...)
	@$(call golintf,github.com/vmware/vic/install/...)
	@$(call golintf,github.com/vmware/vic/portlayer/...)
	@$(call golintf,github.com/vmware/vic/apiservers/portlayer/restapi/handlers/...)
	@$(call golintf,github.com/vmware/vic/apiservers/engine/server/...)
	@$(call golintf,github.com/vmware/vic/apiservers/engine/backends/...)
	@touch $@

# For use by external tools such as emacs or for example:
# GOPATH=$(make gopath) go get ...
gopath:
	@echo -n $(GOPATH)

$(go-imports): $(GOIMPORTS) $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "apiservers/portlayer") $(PORTLAYER_DEPS)
	@echo checking go imports...
	@! $(GOIMPORTS) -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | egrep -v '^$$'
	@touch $@

govet: $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "apiservers/portlayer") $(PORTLAYER_DEPS)
	@echo checking go vet...
	@$(GO) tool vet -all $$(find . -type f -name '*.go' -not -path "./vendor/*")
	@$(GO) tool vet -shadow $$(find . -type f -name '*.go' -not -path "./vendor/*")

vendor: $(GVT)
	@echo restoring vendor
	$(GOPATH)/bin/gvt restore

integration-tests: Dockerfile.integration-tests tests/imagec.bats tests/helpers/helpers.bash components
	@echo Running integration tests
	@docker build -t imagec_tests -f Dockerfile.integration-tests .
	docker run -e BIN=$(BIN) --rm imagec_tests

TEST_DIRS=github.com/vmware/vic/cmd/tether
TEST_DIRS+=github.com/vmware/vic/cmd/imagec
TEST_DIRS+=github.com/vmware/vic/cmd/vicadmin
TEST_DIRS+=github.com/vmware/vic/cmd/rpctool
TEST_DIRS+=github.com/vmware/vic/cmd/vic-machine
TEST_DIRS+=github.com/vmware/vic/portlayer
TEST_DIRS+=github.com/vmware/vic/pkg
TEST_DIRS+=github.com/vmware/vic/apiservers/portlayer
TEST_DIRS+=github.com/vmware/vic/install


test:
	@echo Running unit tests
	# test everything but vendor
ifdef DRONE
	@echo Generating coverage data
	scripts/coverage.sh $(TEST_DIRS)
else
	@echo Generating local html coverage report
	scripts/coverage.sh --html $(TEST_DIRS)
endif

docker-integration-tests:
	@echo Running Docker integration tests
	tests/docker-tests/run-tests.sh

$(tether-linux): $(shell find cmd/tether -name '*.go') metadata/*.go
	@echo building tether-linux
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./cmd/tether

$(tether-windows): $(shell find cmd/tether -name '*.go') metadata/*.go
	@echo building tether-windows
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GO) build $(RACE) -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o ./$@ ./cmd/tether


$(rpctool): cmd/rpctool/*.go
ifeq ($(OS),linux)
	@echo building rpctool
	@GOARCH=amd64 GOOS=linux $(GO) build $(RACE) -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)
else
	@echo skipping rpctool, cannot cross compile cgo
endif

$(vicadmin): cmd/vicadmin/*.go pkg/vsphere/session/*.go
	@echo building vicadmin
	@GOARCH=amd64 GOOS=linux $(GO) build $(RACE) -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)

$(imagec): cmd/imagec/*.go $(portlayerapi-client)
	@echo building imagec...
	@$(GO) build $(RACE) -o ./$@ ./$(dir $<)


$(docker-engine-api): $(portlayerapi-client) apiservers/engine/server/*.go apiservers/engine/backends/*.go
ifeq ($(OS),linux)
	@echo Building docker-engine-api server...
	@$(GO) build $(RACE) -o $@ ./apiservers/engine/server
else
	@echo skipping docker-engine-api server, cannot build on non-linux
endif

# Common portlayer dependencies between client and server
PORTLAYER_DEPS ?= apiservers/portlayer/swagger.yml \
				  apiservers/portlayer/restapi/configure_port_layer.go \
				  apiservers/portlayer/restapi/options/*.go apiservers/portlayer/restapi/handlers/*.go

$(portlayerapi-client): $(PORTLAYER_DEPS)  $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API client...
	@$(SWAGGER) generate client -A PortLayer --template-dir templates  -t $(realpath $(dir $<)) -f $<


$(portlayerapi-server): $(PORTLAYER_DEPS) $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API server...
	@$(SWAGGER) generate server -A PortLayer -t $(realpath $(dir $<)) -f $<

$(portlayerapi): $(portlayerapi-server) $(shell find pkg/ apiservers/engine/ -name '*.go') metadata/*.go
	@echo building Portlayer API server...
	@$(GO) build $(RACE) -o $@ ./apiservers/portlayer/cmd/port-layer-server

$(iso-base): isos/base.sh isos/base/*.repo isos/base/isolinux/** isos/base/xorriso-options.cfg
	@echo building iso-base docker image
	@$< -c $(BIN)/yum-cache.tgz -p $@

# appliance staging - allows for caching of package install
$(appliance-staging): isos/appliance-staging.sh $(iso-base)
	@echo staging for VCH appliance
	@$< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@

# main appliance target - depends on all top level component targets
$(appliance): isos/appliance.sh isos/appliance/* $(rpctool) $(vicadmin) $(imagec) $(portlayerapi) $(docker-engine-api) $(appliance-staging)
	@echo building VCH appliance ISO
	@$< -p $(appliance-staging) -b $(BIN)

# main bootstrap target
$(bootstrap): isos/bootstrap.sh $(tether-linux) $(rpctool) $(bootstrap-staging) isos/bootstrap/*
	@echo "Making bootstrap iso"
	@$< -p $(bootstrap-staging) -b $(BIN)

$(bootstrap-debug): isos/bootstrap.sh $(tether-linux) $(rpctool) $(bootstrap-staging-debug) isos/bootstrap/*
	@echo "Making bootstrap-debug iso"
	@$< -p $(bootstrap-staging-debug) -b $(BIN) -d true

$(bootstrap-staging): isos/bootstrap-staging.sh $(iso-base)
	@echo staging for bootstrap
	@$< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@

$(bootstrap-staging-debug): isos/bootstrap-staging.sh $(iso-base)
	@echo staging debug for bootstrap
	@$< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@ -d true

$(vic-machine): cmd/vic-machine/*.go install/**
	@echo building vic-machine...
	@$(GO) build $(RACE) -o ./$@ ./$(dir $<)

clean:
	rm -rf $(BIN)

	@echo removing swagger generated files...
	rm -f ./apiservers/portlayer/restapi/doc.go
	rm -f ./apiservers/portlayer/restapi/embedded_spec.go
	rm -f ./apiservers/portlayer/restapi/server.go
	rm -rf ./apiservers/portlayer/client/
	rm -rf ./apiservers/portlayer/cmd/
	rm -rf ./apiservers/portlayer/models/
	rm -rf ./apiservers/portlayer/restapi/operations/
	rm -fr ./tests/helpers/bats-assert/
	rm -fr ./tests/helpers/bats-support/

	@tests/docker-tests/run-tests.sh clean
	
