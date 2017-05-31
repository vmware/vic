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

data_dir=/opt/vmware/engine_installer
cert_dir=${data_dir}/cert
cert=${cert_dir}/server.crt
key=${cert_dir}/server.key

engine_installer_cert=$(ovfenv -k engine_installer.ssl_cert)
engine_installer_key=$(ovfenv -k engine_installer.ssl_cert_key)
port=$(ovfenv -k engine_installer.port)
target=$(ovfenv -k engine_installer.vcenter_addr)
username=$(ovfenv -k engine_installer.admin_user)
password=$(ovfenv -k engine_installer.admin_password)

if [[ x$engine_installer_cert != "x" && x$engine_installer_key != "x" ]]; then
  /usr/local/bin/ova-engine-installer --addr ":${port}" --cert "${cert}" --key "${key}" --target "${target}" --username "${username}" --password "${password}"
else
  /usr/local/bin/ova-engine-installer --addr ":${port}" --target "${target}" --username "${username}" --password "${password}"

fi

