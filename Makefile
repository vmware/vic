GO ?= go
GOVERSION ?= go1.6
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
GOVET ?= $(GOPATH)/bin/vet$(BIN_ARCH)
GOIMPORTS ?= $(GOPATH)/bin/goimports$(BIN_ARCH)
GOLINT ?= $(GOPATH)/bin/golint$(BIN_ARCH)
GVT ?= $(GOPATH)/bin/gvt$(BIN_ARCH)

REPO ?= https://dl.bintray.com/vmware/vic-repo/kernel/
KERNEL ?= linux-esx-4.2.0-10.x86_64.rpm

.PHONY: all tools clean test check \
	goversion govet goimports gvt gopath \
	isos tethers apiservers copyright

.DEFAULT_GOAL := all


# target aliases - environment variable definition
docker-engine-api := $(BIN)/docker-engine-server
portlayerapi := $(BIN)/port-layer-server
portlayerapi-client := apiservers/portlayer/client/port_layer_client.go
portlayerapi-server := apiservers/portlayer/restapi/server.go

imagec := $(BIN)/imagec
vicadmin := $(BIN)/vicadmin
rpctool := $(BIN)/rpctool

tether-linux := $(BIN)/tether-linux
tether-windows := $(BIN)/tether-windows.exe

appliance := $(BIN)/appliance.iso
appliance-staging := $(BIN)/appliance-staging.tgz
bootstrap := $(BIN)/bootstrap.iso
iso-base := $(BIN)/iso-base.tgz
kernel := $(BIN)/$(KERNEL)

install := $(BIN)/install.sh

go-lint := $(BIN)/.golint
go-vet := $(BIN)/.govet
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
iso-base: $(iso-base)
kernel: $(kernel)
install: $(install)

swagger: $(SWAGGER)

golint: $(go-lint)
govet: $(go-vet)
goimports: $(go-imports)


# convenience targets
all: components isos install
tools: $(GOIMPORTS) $(GOVET) $(GVT) $(GOLINT) $(SWAGGER) goversion
check: goversion goimports govet golint copyright
apiservers: $(portlayerapi) $(docker-engine-api)
bootstrap: $(tether-linux) $(tether-windows) $(rpctool)
components: check apiservers $(imagec) $(vicadmin) $(rpctool)
isos: $(appliance) $(bootstrap)
tethers: $(tether-linux) $(tether-windows)


# utility targets
goversion:
	@echo checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

$(GOIMPORTS): vendor/manifest
	@echo building $(GOIMPORTS)...
	$(GO) build -o $(GOIMPORTS) ./vendor/golang.org/x/tools/cmd/goimports

$(GOVET): vendor/manifest
	@echo building $(GOVET)...
	$(GO) build -o $(GOVET) ./vendor/golang.org/x/tools/cmd/vet

$(GVT):
	@echo getting gvt
	$(GO) get -u github.com/FiloSottile/gvt

$(GOLINT): vendor/manifest
	@echo building $(GOLINT)...
	$(GO) build -o $(GOLINT) ./vendor/github.com/golang/lint/golint

$(SWAGGER): vendor/manifest
	@echo building $(SWAGGER)...
	@$(GO) build -o $(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

copyright:
	@echo "checking copyright in header..."
	scripts/header-check.sh

# exit 1 if golint complains about anything other than comments
golintf = $(GOLINT) $(1) | sh -c "! grep -v 'should have comment'"

$(go-lint): $(GOLINT)
	@echo checking go lint...
	@ #$(call golintf,github.com/vmware/vic/bootstrap/...) # this is commented out due to number of warnings
	@$(call golintf,github.com/vmware/vic/imagec/...)
	@$(call golintf,github.com/vmware/vic/vicadmin/...)
	@$(call golintf,github.com/vmware/vic/pkg/...)
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

$(go-vet): $(GOVET) $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "apiservers/portlayer") $(PORTLAYER_DEPS)
	@echo checking go vet...
	@$(GOVET) -all -shadow $$(find . -type f -name '*.go' -not -path "./vendor/*")
	@touch $@

vendor: $(GVT)
	@echo restoring vendor
	$(GOPATH)/bin/gvt restore

integration-tests: Dockerfile.integration-tests tests/imagec.bats tests/helpers/helpers.bash components
	@echo Running integration tests
	@docker build -t imagec_tests -f Dockerfile.integration-tests .
	docker run -e BIN=$(BIN) --rm imagec_tests

TEST_DIRS=github.com/vmware/vic/bootstrap/...
TEST_DIRS+=github.com/vmware/vic/imagec
TEST_DIRS+=github.com/vmware/vic/vicadmin
TEST_DIRS+=github.com/vmware/vic/portlayer/...
TEST_DIRS+=github.com/vmware/vic/pkg/...
TEST_DIRS+=github.com/vmware/vic/apiservers/portlayer/...

test:
	# test everything but vendor
ifdef DRONE
	@echo generate coverage report
	./coverage $(TEST_DIRS)
else
	$(foreach var,$(TEST_DIRS), $(GO) test -v $(TEST_OPTS) $(var);)
endif

$(tether-linux): $(shell find bootstrap/tether -name '*.go')
	@echo building tether-linux
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -tags netgo -installsuffix netgo -o ./$@ ./bootstrap/tether/cmd/tether

$(tether-windows): $(shell find bootstrap/tether -name '*.go')
	@echo building tether-windows
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -tags netgo -installsuffix netgo -o ./$@ ./bootstrap/tether/cmd/tether


$(rpctool): pkg/vsphere/rpctool/*.go
ifeq ($(OS),linux)
	@echo building rpctool
	@GOARCH=amd64 GOOS=linux $(GO) build -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)
else
	@echo skipping rpctool, cannot cross compile cgo
endif

$(vicadmin): vicadmin/*.go pkg/vsphere/session/*.go
	@echo building vicadmin
	@GOARCH=amd64 GOOS=linux $(GO) build -o ./$@ --ldflags '-extldflags "-static"' ./$(dir $<)

$(imagec): imagec/*.go $(portlayerapi-client)
	@echo building imagec...
	@CGO_ENABLED=0 $(GO) build -o ./$@ --ldflags '-extldflags "-static"'  ./$(dir $<)


$(docker-engine-api): $(portlayerapi-client) apiservers/engine/server/*.go apiservers/engine/backends/*.go
ifeq ($(OS),linux)
	@echo Building docker-engine-api server...
	@$(GO) build -o $@ ./apiservers/engine/server
else
	@echo skipping docker-engine-api server, cannot build on non-linux
endif



# Common portlayer dependencies between client and server
PORTLAYER_DEPS ?= apiservers/portlayer/swagger.yml \
				  apiservers/portlayer/restapi/configure_port_layer.go \
				  apiservers/portlayer/restapi/options/*.go apiservers/portlayer/restapi/handlers/*.go

$(portlayerapi-client): $(PORTLAYER_DEPS)  $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API client...
	@$(SWAGGER) generate client -A PortLayer -t $(realpath $(dir $<)) -f $<


$(portlayerapi-server): $(PORTLAYER_DEPS) $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API server...
	@$(SWAGGER) generate server -A PortLayer -t $(realpath $(dir $<)) -f $<

$(portlayerapi): $(portlayerapi-server) $(shell find pkg/ apiservers/engine/ -name '*.go')
	@echo building Portlayer API server...
	@$(GO) build -o $@ ./apiservers/portlayer/cmd/port-layer-server

# cache the kernel locally
$(kernel):
	@ # ensure that the reference file exists, but date it at epoc inception
	@-[ ! -e $(kernel) ] && touch --date="@0" $@
	@curl -L --insecure $(REPO)/$(KERNEL) -o $(kernel) -z $(kernel)

$(iso-base): isos/base.sh isos/base/*.repo isos/base/isolinux/** isos/base/xorriso-options.cfg $(kernel)
	@echo building iso-base docker image
	@$< -c $(BIN)/yum-cache.tgz -p $@ -k $(BIN)/$(KERNEL)

# appliance staging - allows for caching of package install
$(appliance-staging): isos/appliance-staging.sh $(iso-base)
	@echo staging for VCH appliance
	@$< -c $(BIN)/yum-cache.tgz -p $(iso-base) -o $@

# main appliance target - depends on all top level component targets
$(appliance): isos/appliance.sh isos/appliance/* $(rpctool) $(vicadmin) $(imagec) $(portlayerapi) $(docker-engine-api) $(appliance-staging)
	@echo building VCH appliance ISO
	@$< -p $(appliance-staging) -b $(BIN)

# main bootstrap target
$(bootstrap): $(tether-linux) $(rpctool) $(iso-base)
	@echo "Placeholder target for bootstrap"

$(install): install/install.sh
	@echo Building installer
	@cp $< $@

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

