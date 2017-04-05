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

data_dir=/data/fileserver
cert_dir=${data_dir}/cert
cert=${cert_dir}/server.crt
key=${cert_dir}/server.key
fileserver_cert=$(ovfenv -k fileserver.ssl_cert)
fileserver_key=$(ovfenv -k fileserver.ssl_cert_key)

FILESERVER_EXPOSED_PORT=$(ovfenv -k fileserver.port)

if [[ x$fileserver_cert != "x" && x$fileserver_key != "x" ]]; then
  /usr/local/bin/ova-webserver --addr ":${FILESERVER_EXPOSED_PORT}" --cert "${cert}" --key "${key}"
else
  /usr/local/bin/ova-webserver --addr ":${FILESERVER_EXPOSED_PORT}"
fi

