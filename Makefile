GO ?= go
GOVERSION ?= go1.5
OS := $(shell uname)
GO15VENDOREXPERIMENT=1

export GO15VENDOREXPERIMENT

goversion:
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )


all: check gvt vendor test bootstrap

check: goversion goimports govet

bootstrap: tether.linux tether.windows rpctool

goimports:
	@echo checking go imports...
	@! goimports -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | egrep -v '^$$'

govet:
	@echo checking go vet...
	@go tool vet -structtags=false -methods=false $$(find . -type f -name '*.go' -not -path "./vendor/*")

gvt:
	@echo getting gvt
	go get -u github.com/FiloSottile/gvt

vendor:
	@echo restoring vendor
	$(GOPATH)/bin/gvt restore

test:
	# test everything but vendor
	go test -v $(TEST_OPTS) github.com/vmware/vic/bootstrap/...

tether.linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -x -tags netgo -installsuffix netgo -o ./binary/tether-linux github.com/vmware/vic/bootstrap/tether/cmd/tether

tether.windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a -x -tags netgo -installsuffix netgo -o ./binary/tether-windows github.com/vmware/vic/bootstrap/tether/cmd/tether

rpctool.linux:
	GOARCH=amd64 GOOS=linux go build -a -x -o ./binary/rpctool --ldflags '-extldflags "-static"' github.com/vmware/vic/bootstrap/rpctool

rpctool: rpctool.linux

clean:
	rm -rf ./binary ./vendor

.PHONY: test vendor
