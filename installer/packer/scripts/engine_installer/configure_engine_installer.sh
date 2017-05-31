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

umask 077

run_engine_installer=$(ovfenv -k engine_installer.wizard_enabled)

if [ ${deploy,,} != "true" ]; then
  echo "Not configuring Engine Installer and disabling service startup"
  systemctl stop engine_installer
  exit 0
fi

data_dir=/opt/vmware/engine_installer
cert_dir=${data_dir}/cert
cert=${cert_dir}/server.crt
key=${cert_dir}/server.key

engine_installer_cert=$(ovfenv -k engine_installer.ssl_cert)
engine_installer_key=$(ovfenv -k engine_installer.ssl_cert_key)
port=$(ovfenv -k engine_installer.port)

iptables -w -A INPUT -j ACCEPT -p tcp --dport $port

# Format cert file
function formatCert {
  content=$1
  file=$2
  echo ${content} | sed -r 's/(-{5}BEGIN [A-Z ]+-{5})/&\n/g; s/(-{5}END [A-Z ]+-{5})/\n&\n/g' | sed -r 's/.{64}/&\n/g; /^\s*$/d' > ${file}
}

if [[ x${engine_installer_cert} != "x" && x${engine_installer_key} != "x" ]]; then
  mkdir -p "${cert_dir}"
  formatCert "${engine_installer_cert}" ${cert}
  formatCert "${engine_installer_key}" ${key}
fi

FILESERVER_DIR="/opt/vmware/fileserver/files"
FILE_COUNT=$(find ${FILESERVER_DIR} -name "vic*.tar.gz" | wc -l)
if [ ! -d "${FILESERVER_DIR}" ] || [ ${FILE_COUNT} -ne 1] ; then
	echo "Fileserver files not present. Unable to get VIC Engine tarball"
	systemctl stop engine_installer
	exit 1
fi

# Extract vic-machine
find ${FILESERVER_DIR} -name "vic*.tar.gz" | xargs -I {} tar xvf {} --directory ${data_dir}
