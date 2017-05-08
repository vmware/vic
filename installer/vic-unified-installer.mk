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

PACKER ?= packer
OVFTOOL ?= ovftool
SHA256SUM ?= sha256sum
SED ?= sed
RM ?= rm
CP ?= cp

VER :=$(shell git describe --abbrev=0 --tags)

PHOTON_ISO := https://storage.googleapis.com/vic-product-ova-build-deps/photon-1.0-62c543d.iso
PHOTON_ISO_SHA1SUM := c4c6cb94c261b162e7dac60fdffa96ddb5836d66

.PHONY: ova-release ova-debug vagrant-local

ovfenv := $(BIN)/ovfenv
vic-ova-ui := $(BIN)/vic-ova-ui
ova-webserver := $(BIN)/ova-webserver
ova-engine-installer := $(BIN)/ova-engine-installer
ovfenv: $(ovfenv)
vic-ova-ui: $(vic-ova-ui)
ova-webserver: $(ova-webserver)
ova-engine-installer: $(ova-engine-installer)

$(ovfenv): $$(call godeps,installer/ovatools/ovfenv/*.go)
	@echo building ovfenv linux...
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -ldflags "$(LDFLAGS)" -o ./$@ ./$(dir $<)

$(vic-ova-ui): $$(call godeps,installer/ovatools/vic-ova-ui/*.go)
	@echo building vic-ova-ui
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -ldflags "$(LDFLAGS)" -o ./$@ ./$(dir $<)

$(ova-webserver): $$(call godeps,installer/fileserver/*.go)
	@echo building ova-webserver
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -ldflags "$(LDFLAGS)" -o ./$@ ./$(dir $<)

$(ova-engine-installer): $$(call godeps,installer/engine_installer/*.go)
	@echo building ova-engine-installer
	@GOARCH=amd64 GOOS=linux $(TIME) $(GO) build $(RACE) -ldflags "$(LDFLAGS)" -o ./$@ ./$(dir $<)

ova-release: $(ovfenv) $(vic-ova-ui) $(ova-webserver) $(ova-engine-installer)
	@echo building vic-unified-installer OVA using packer...
	@cd $(BASE_DIR)installer/packer && $(PACKER) build \
			-only=ova-release \
			-var 'iso_sha1sum=$(PHOTON_ISO_SHA1SUM)'\
			-var 'iso_file=$(PHOTON_ISO)'\
			-var 'esx_host=$(PACKER_ESX_HOST)'\
			-var 'remote_username=$(PACKER_USERNAME)'\
			-var 'remote_password=$(PACKER_PASSWORD)'\
			-var 'build_vicengine_revision=$(BUILD_VICENGINE_REVISION)' \
			-var 'build_vicengine_dev_revision=$(BUILD_VICENGINE_DEV_REVISION)' \
			-var 'build_kov_cli_revision=$(BUILD_KOV_CLI_REVISION)' \
			-var 'build_kovd_revision=$(BUILD_KOVD_REVISION)' \
			-var 'build_admiral_revision=$(BUILD_ADMIRAL_REVISION)' \
			-var 'build_harbor_revision=$(BUILD_HARBOR_REVISION)' \
			packer-vic.json
	@echo adding proper vic OVF file...
	@cd $(BASE_DIR)installer/packer/vic/vic && $(RM) vic.ovf && $(CP) ../../vic-unified.ovf vic.ovf
	@echo rebuilding OVF manifest...
	@cd $(BASE_DIR)installer/packer/vic/vic && $(RM) vic.mf && $(SHA256SUM) --tag * | $(SED) s/SHA256\ \(/SHA256\(/ > vic.mf
	@echo packaging OVA...
	@$(OVFTOOL) -st=ovf -tt=ova $(BASE_DIR)installer/packer/vic/vic/vic.ovf $(BASE_DIR)$(BIN)/vic-$(VER)-$(REV).ova
	@echo cleaning packer directory...
	@cd $(BASE_DIR)installer/packer && $(RM) -rf vic

ova-debug: $(ovfenv) $(vic-ova-ui) $(ova-webserver) $(ova-engine-installer)
	@echo building vic-unified-installer OVA using packer...
	cd $(BASE_DIR)installer/packer && PACKER_LOG=1 $(PACKER) build \
			-only=ova-release \
			-var 'iso_sha1sum=$(PHOTON_ISO_SHA1SUM)'\
			-var 'iso_file=$(PHOTON_ISO)'\
			-var 'esx_host=$(PACKER_ESX_HOST)'\
			-var 'remote_username=$(PACKER_USERNAME)'\
			-var 'remote_password=$(PACKER_PASSWORD)'\
			-var 'build_vicengine_revision=$(BUILD_VICENGINE_REVISION)' \
			-var 'build_vicengine_dev_revision=$(BUILD_VICENGINE_DEV_REVISION)' \
			-var 'build_kov_cli_revision=$(BUILD_KOV_CLI_REVISION)' \
			-var 'build_kovd_revision=$(BUILD_KOVD_REVISION)' \
			-var 'build_admiral_revision=$(BUILD_ADMIRAL_REVISION)' \
			-var 'build_harbor_revision=$(BUILD_HARBOR_REVISION)' \
			--on-error=abort packer-vic.json
	@echo adding proper vic OVF file...
	cd $(BASE_DIR)installer/packer/vic/vic && $(RM) vic.ovf && $(CP) ../../vic-unified.ovf vic.ovf
	@echo rebuilding OVF manifest...
	cd $(BASE_DIR)installer/packer/vic/vic && $(RM) vic.mf && $(SHA256SUM) --tag * | $(SED) s/SHA256\ \(/SHA256\(/ > vic.mf
	@echo packaging OVA...
	$(OVFTOOL) -st=ovf -tt=ova $(BASE_DIR)installer/packer/vic/vic/vic.ovf $(BASE_DIR)$(BIN)/vic-$(VER)-$(REV).ova
	@echo cleaning packer directory...
	cd $(BASE_DIR)installer/packer && $(RM) -rf vic

vagrant-local: $(ovfenv) $(vic-ova-ui) $(ova-webserver) $(ova-engine-installer)
	@echo building vic-unified-installer Vagrant box using packer...
	@cd $(BASE_DIR)installer/packer && $(PACKER) build \
			-only=vagrant-local \
			-var 'iso_sha1sum=$(PHOTON_ISO_SHA1SUM)'\
			-var 'iso_file=$(PHOTON_ISO)'\
			-var 'build_vicengine_revision=$(BUILD_VICENGINE_REVISION)' \
			-var 'build_vicengine_dev_revision=$(BUILD_VICENGINE_DEV_REVISION)' \
			-var 'build_kov_cli_revision=$(BUILD_KOV_CLI_REVISION)' \
			-var 'build_kovd_revision=$(BUILD_KOVD_REVISION)' \
			-var 'build_admiral_revision=$(BUILD_ADMIRAL_REVISION)' \
			-var 'build_harbor_revision=$(BUILD_HARBOR_REVISION)' \
			--on-error=abort packer-vic.json
