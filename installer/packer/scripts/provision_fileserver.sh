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

FILES_DIR="/opt/vmware/fileserver/files"

mkdir -p /etc/vmware/fileserver  # Fileserver config scripts
mkdir -p $FILES_DIR              # Files to serve

[[ x$BUILD_VICENGINE_REVISION == "x" ]] && ( echo "VIC Engine build not set, failing"; exit 1 )

# Download Build
cd /var/tmp
VIC_ENGINE_FILE=""
VIC_ENGINE_URL=""

set +u
if [ -z "${BUILD_VICENGINE_DEV_REVISION}" ]; then
  VIC_ENGINE_FILE="vic_${BUILD_VICENGINE_REVISION}.tar.gz"
  VIC_ENGINE_URL="https://storage.googleapis.com/vic-engine-releases/${VIC_ENGINE_FILE}"
else
  VIC_ENGINE_FILE="vic_${BUILD_VICENGINE_DEV_REVISION}.tar.gz"
  VIC_ENGINE_URL="https://storage.googleapis.com/vic-engine-builds/${VIC_ENGINE_FILE}"
fi
set -u

echo "Downloading VIC Engine ${VIC_ENGINE_FILE}: ${VIC_ENGINE_URL}"
curl -LO ${VIC_ENGINE_URL}
tar xfz ${VIC_ENGINE_FILE} -C $FILES_DIR vic/ui/vsphere-client-serenity/com.vmware.vic.ui-v${BUILD_VICENGINE_REVISION}.zip vic/ui/plugin-packages/com.vmware.vic-v${BUILD_VICENGINE_REVISION}.zip --strip-components=3
mv ${VIC_ENGINE_FILE} $FILES_DIR
