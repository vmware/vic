#!/usr/bin/bash
# Copyright 2017 VMware, Inc. All Rights Reserved.
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
set -euf -o pipefail

ADMIRAL_EXPOSED_PORT=""
ADMIRAL_DATA_LOCATION=""
ADMIRAL_KEY_LOCATION=""
ADMIRAL_CERT_LOCATION=""
ADMIRAL_JKS_LOCATION=""

/usr/bin/docker run -p ${ADMIRAL_EXPOSED_PORT}:8282 \
  --name vic-admiral \
  -e ADMIRAL_PORT=8282 \
  -e JAVA_OPTS="-Ddcp.net.ssl.trustStore=/tmp/trusted_certificates.jks -Ddcp.net.ssl.trustStorePassword=changeit" \
  -e XENON_OPTS="--port=-1 --securePort=8282 --certificateFile=/tmp/server.crt --keyFile=/tmp/server.key" \
  -v "$ADMIRAL_CERT_LOCATION:/tmp/server.crt" \
  -v "$ADMIRAL_KEY_LOCATION:/tmp/server.key" \
  -v "$ADMIRAL_JKS_LOCATION:/tmp/trusted_certificates.jks" \
  -v "$ADMIRAL_DATA_LOCATION/custom.conf:/admiral/config/configuration.properties" \
  -v "$ADMIRAL_DATA_LOCATION:/var/admiral" \
  --log-driver=json-file \
  --log-opt max-size=1g \
  --log-opt max-file=10 \
  vmware/admiral:ova
  