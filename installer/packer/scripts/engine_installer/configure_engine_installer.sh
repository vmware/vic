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

data_dir="/opt/vmware/engine_installer"
mkdir -p ${data_dir}

port=$(ovfenv -k engine_installer.port)
iptables -w -A INPUT -j ACCEPT -p tcp --dport $port

FILESERVER_DIR="/opt/vmware/fileserver/files"
FILE_COUNT=$(find ${FILESERVER_DIR} -name "vic*.tar.gz" | wc -l)
if [ ! -d "${FILESERVER_DIR}" ] || [ ${FILE_COUNT} -ne 1 ] ; then
	echo "Fileserver files not present. Unable to get VIC Engine tarball"
	systemctl stop engine_installer
	exit 1
fi

# Extract vic-machine
VIC_DIR="${FILESERVER_DIR}/vic"
if [ -d "${VIC_DIR}" ]; then
  rm -rf ${VIC_DIR}
fi

find ${FILESERVER_DIR} -name "vic*.tar.gz" | xargs -I {} tar xvf {} --directory ${data_dir}
