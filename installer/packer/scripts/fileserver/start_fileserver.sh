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

FILESERVER_EXPOSED_PORT=$(ovfenv -k fileserver.port)
FILESERVER_CERT_LOCATION=$(ovfenv -k fileserver.ssl_cert)
FILESERVER_KEY_LOCATION=$(ovfenv -k fileserver.ssl_cert_key)

/usr/local/bin/ova-webserver --addr ":${FILESERVER_EXPOSED_PORT}" --cert "${FILESERVER_CERT_LOCATION}" --key "${FILESERVER_KEY_LOCATION}"
