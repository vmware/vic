GO ?= go
GOVERSION ?= go1.6
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ifeq ($(USER),vagrant)
	# assuming we are in a shared directory where host arch is different from the guest
	BIN_ARCH := -$(OS)
endif

export GOPATH ?= $(shell echo $(CURDIR) | sed -e 's,/src/.*,,')
SWAGGER ?= $(GOPATH)/bin/swagger$(BIN_ARCH)
GOVET ?= $(GOPATH)/bin/vet$(BIN_ARCH)
GOIMPORTS ?= $(GOPATH)/bin/goimports$(BIN_ARCH)
GOLINT ?= $(GOPATH)/bin/golint$(BIN_ARCH)
GVT ?= $(GOPATH)/bin/gvt$(BIN_ARCH)


.PHONY: all tools clean test check \
	goversion govet goimports gvt gopath \
	bootstrap apiservers  \


.DEFAULT_GOAL := all


# target aliases - environment variable definition
docker-engine-api := binary/docker-engine-server

portlayerapi := binary/port-layer-server
portlayerapi-client := apiservers/portlayer/client/port_layer_client.go
portlayerapi-server := apiservers/portlayer/restapi/server.go

imagec := binary/imagec
vicadmin := binary/vicadmin
rpctool := binary/rpctool

tether-linux := binary/tether-linux
tether-windows := binary/tether-windows.exe

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

swagger: $(SWAGGER)

# convenience targets
all: check bootstrap apiservers $(imagec) $(vicadmin)
tools: $(GOIMPORTS) $(GOVET) $(GVT) $(GOLINT) $(SWAGGER) goversion
check: goversion goimports govet golint
apiservers: $(portlayerapi) $(docker-engine-api)
bootstrap: $(tether-linux) $(tether-windows) $(rpctool)


goversion:
	@echo checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

$(GOIMPORTS):
	@echo building $(GOIMPORTS)...
	$(GO) build -o $(GOIMPORTS) ./vendor/golang.org/x/tools/cmd/goimports

$(GOVET):
	@echo building $(GOVET)...
	$(GO) build -o $(GOVET) ./vendor/golang.org/x/tools/cmd/vet

$(GVT):
	@echo getting gvt
	$(GO) get -u github.com/FiloSottile/gvt

$(GOLINT): vendor/manifest
	@echo building $(GOLINT)...
	$(GO) build -o $(GOLINT) ./vendor/github.com/golang/lint/golint

$(SWAGGER):
	@echo building $(SWAGGER)...
	@$(GO) build -o $(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

# exit 1 if golint complains about anything other than comments
golintf = $(GOLINT) $(1) | sh -c "! grep -v 'should have comment'"

golint: $(GOLINT)
	@echo checking go lint...
	@#$(call golintf,github.com/vmware/vic/bootstrap/...) # this is commented out due to number of warnings
	@$(call golintf,github.com/vmware/vic/imagec/...)
	@$(call golintf,github.com/vmware/vic/vicadmin/...)
	@$(call golintf,github.com/vmware/vic/pkg/...)
	@$(call golintf,github.com/vmware/vic/portlayer/...)
	@$(call golintf,github.com/vmware/vic/apiservers/portlayer/restapi/handlers/...)
	@$(call golintf,github.com/vmware/vic/apiservers/engine/server/...)
	@$(call golintf,github.com/vmware/vic/apiservers/engine/backends/...)

# For use by external tools such as emacs or for example:
# GOPATH=$(make gopath) go get ...
gopath:
	@echo -n $(GOPATH)

goimports: $(GOIMPORTS)
	@echo checking go imports...
	@! $(GOIMPORTS) -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | egrep -v '^$$'

govet: $(GOVET)
	@echo checking go vet...
	@$(GOVET) -all -shadow $$(find . -type f -name '*.go' -not -path "./vendor/*")

vendor: $(GVT)
	@echo restoring vendor
	$(GOPATH)/bin/gvt restore

integration-tests:
	docker build -t imagec_tests -f Dockerfile.integration-tests .
	docker run --rm imagec_tests

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
	@echo Building docker-engine-api server...
	@$(GO) build -o $@ ./apiservers/engine/server

# Common portlayer dependencies between client and server
PORTLAYER_DEPS ?= apiservers/portlayer/swagger.yml \
				  apiservers/portlayer/restapi/configure_port_layer.go \
				  apiservers/portlayer/restapi/options/*.go apiservers/portlayer/restapi/handlers/*.go

$(portlayerapi-client): $(PORTLAYER_DEPS)  $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API client...
	@$(SWAGGER) generate client -A PortLayer -t $(dir $<) -f $<


$(portlayerapi-server): $(PORTLAYER_DEPS) $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API server...
	@$(SWAGGER) generate server -A PortLayer -t $(dir $<) -f $<

$(portlayerapi): $(portlayerapi-server) $(shell find apiservers/engine/ -name '*.go')
	@echo building Portlayer API server...
	@$(GO) build -o $@ ./apiservers/portlayer/cmd/port-layer-server

clean:
	rm -rf ./binary

	@echo removing swagger generated files...
	rm -f ./apiservers/portlayer/restapi/doc.go
	rm -f ./apiservers/portlayer/restapi/embedded_spec.go
	rm -f ./apiservers/portlayer/restapi/server.go
	rm -rf ./apiservers/portlayer/client/
	rm -rf ./apiservers/portlayer/cmd/
	rm -rf ./apiservers/portlayer/models/
	rm -rf ./apiservers/portlayer/restapi/operations/
