#!/bin/bash -e
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

# Run robot integration tests locally, no .yml files required.
# Set GITHUB_TOKEN once and switch environments just by changing GOVC_URL

if [ -z "$GITHUB_TOKEN" ] || [ -z "$GOVC_URL" ]; then
    echo "usage: GITHUB_TOKEN=... GOVC_URL=... $0 test.robot..."
    exit 1
fi

cd "$(git rev-parse --show-toplevel)"

tests=${*#${PWD}/}

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.6
    pull: true
    environment:
      GITHUB_AUTOMATION_API_KEY: $GITHUB_TOKEN
      TEST_URL_ARRAY:   $(govc env -x GOVC_URL_HOST)
      TEST_USERNAME:    $(govc env GOVC_USERNAME)
      TEST_PASSWORD:    $(govc env GOVC_PASSWORD)
      TEST_DATASTORE:   ${GOVC_DATASTORE:-$(basename "$(govc ls datastore)")}
      TEST_RESOURCE:    ${GOVC_RESOURCE_POOL:-$(govc ls host/*/Resources)}
      BRIDGE_NETWORK:   $BRIDGE_NETWORK
      EXTERNAL_NETWORK: $EXTERNAL_NETWORK
      BIN: bin
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"
      TEST_TIMEOUT: 60s
      GOVC_INSECURE: true
    commands:
      - pybot ${tests:-tests/test-cases}
CONFIG
)
