GO ?= go
GOVERSION ?= go1.6
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ifeq ($(USER),vagrant)
	# assuming we are in a shared directory where host arch is different from the guest
	BIN_ARCH := -$(OS)
endif

export GOPATH ?= $(shell echo $(CURDIR) | sed -e 's,/src/.*,,')
SWAGGER ?= $(GOPATH)/bin/swagger$(BIN_ARCH)
VET ?= $(GOPATH)/bin/vet$(BIN_ARCH)
GOIMPORTS ?= $(GOPATH)/bin/goimports$(BIN_ARCH)
GOLINT ?= $(GOPATH)/bin/golint$(BIN_ARCH)

.DEFAULT_GOAL := all

goversion:
	@echo checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

all: check bootstrap apiservers

check: goversion goimports govet golint

bootstrap: binary/tether-linux binary/tether-windows binary/rpctool

apiservers: dockerapi portlayerapi

$(GOIMPORTS): vendor/manifest
	@echo building $(GOIMPORTS)...
	$(GO) build -o $(GOIMPORTS) ./vendor/golang.org/x/tools/cmd/goimports

goimports: $(GOIMPORTS)
	@echo checking go imports...
	@! $(GOIMPORTS) -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | egrep -v '^$$'

$(VET): vendor/manifest
	@echo building $(VET)...
	$(GO) build -o $(VET) ./vendor/golang.org/x/tools/cmd/vet

govet: $(VET)
	@echo checking go vet...
	@$(VET) -all -shadow -structtags=false -methods=false $$(find . -type f -name '*.go' -not -path "./vendor/*")

$(GOLINT): vendor/manifest
	@echo building $(GOLINT)...
	$(GO) build -o $(GOLINT) ./vendor/github.com/golang/lint/golint

# exit 1 if golint complains about anything other than comments
golintf = $(GOLINT) $(1) | sh -c "! grep -v 'should have comment'"

golint: $(GOLINT)
	@echo checking go lint...
	$(call golintf,github.com/vmware/vic/imagec/...)
	$(call golintf,github.com/vmware/vic/pkg/...)

# For use by external tools such as emacs or for example:
# GOPATH=$(make gopath) go get ...
gopath:
	@echo -n $(GOPATH)

gvt:
	@echo getting gvt
	$(GO) get -u github.com/FiloSottile/gvt

vendor:
	@echo restoring vendor
	$(GOPATH)/bin/gvt restore

test:
	# test everything but vendor
	$(GO) test -v $(TEST_OPTS) github.com/vmware/vic/bootstrap/...
	$(GO) test -v $(TEST_OPTS) github.com/vmware/vic/imagec
	$(GO) test -v $(TEST_OPTS) github.com/vmware/vic/portlayer/...
	$(GO) test -v $(TEST_OPTS) github.com/vmware/vic/pkg/...

binary/tether-linux: $(shell find ./bootstrap/tether -name '*.go')
	@echo building tether-linux
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -tags netgo -installsuffix netgo -o ./binary/tether-linux github.com/vmware/vic/bootstrap/tether/cmd/tether

binary/tether-windows: $(shell find ./bootstrap/tether -name '*.go')
	@echo building tether-windows
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -tags netgo -installsuffix netgo -o ./binary/tether-windows github.com/vmware/vic/bootstrap/tether/cmd/tether

binary/rpctool: $(find ./bootstrap/rpctool -name '*.go')
ifeq ($(OS),linux)
	@echo building rpctool
	@GOARCH=amd64 GOOS=linux $(GO) build -o ./binary/rpctool --ldflags '-extldflags "-static"' github.com/vmware/vic/bootstrap/rpctool
else
	@echo skipping rpctool, cannot cross compile cgo
endif

rpctool: binary/rpctool

imagec: portlayerapi-client
	@echo building imagec...
	@CGO_ENABLED=0 $(GO) build -o ./binary/imagec --ldflags '-extldflags "-static"' github.com/vmware/vic/imagec

$(SWAGGER): vendor/manifest
	@echo building $(SWAGGER)...
	@$(GO) build -o $(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

dockerapi-server: $(SWAGGER)
	@echo regenerating swagger models and operations for Docker API server...
	@$(SWAGGER) generate server -A docker -t ./apiservers/docker -f ./apiservers/docker/swagger.json

dockerapi: dockerapi-server
	@echo building Docker API server...
	@$(GO) build -o ./binary/docker-server ./apiservers/docker/cmd/docker-server

portlayerapi-client: $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API client...
	@$(SWAGGER) generate client -A PortLayer -t ./apiservers/portlayer -f ./apiservers/portlayer/swagger.yml

portlayerapi-server: $(SWAGGER)
	@echo regenerating swagger models and operations for Portlayer API server...
	@$(SWAGGER) generate server -A PortLayer -t ./apiservers/portlayer -f ./apiservers/portlayer/swagger.yml

portlayerapi: portlayerapi-server
	@echo building Portlayer API server...
	@$(GO) build -o ./binary/port-layer-server ./apiservers/portlayer/cmd/port-layer-server/

clean:
	rm -rf ./binary

	@echo removing swagger generated files...
	rm -f ./apiservers/docker/restapi/doc.go
	rm -f ./apiservers/docker/restapi/embedded_spec.go
	rm -f ./apiservers/docker/restapi/server.go
	rm -rf ./apiservers/docker/cmd
	rm -rf ./apiservers/docker/models
	rm -rf ./apiservers/docker/restapi/operations

	rm -f ./apiservers/portlayer/restapi/doc.go
	rm -f ./apiservers/portlayer/restapi/embedded_spec.go
	rm -f ./apiservers/portlayer/restapi/server.go
	rm -rf ./apiservers/portlayer/client/
	rm -rf ./apiservers/portlayer/cmd/
	rm -rf ./apiservers/portlayer/models/
	rm -rf ./apiservers/portlayer/restapi/operations/

.PHONY: test vendor imagec
