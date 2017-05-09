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
echo "Downloading VIC Engine ${BUILD_VICENGINE_REVISION}"
curl -LO "https://storage.googleapis.com/vic-engine-releases/vic_${BUILD_VICENGINE_REVISION}.tar.gz"
tar xfz vic_${BUILD_VICENGINE_REVISION}.tar.gz -C $FILES_DIR vic/ui/vsphere-client-serenity/com.vmware.vic.ui-v${BUILD_VICENGINE_REVISION}.zip vic/ui/plugin-packages/com.vmware.vic-v${BUILD_VICENGINE_REVISION}.zip --strip-components=3
mv vic_${BUILD_VICENGINE_REVISION}.tar.gz $FILES_DIR
