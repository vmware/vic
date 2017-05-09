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

BUILD_HARBOR_REVISION="${BUILD_HARBOR_REVISION:-dev}"

# Download Build
echo "Downloading Harbor ${BUILD_HARBOR_REVISION}"
curl -L "https://storage.googleapis.com/harbor-dev-builds/harbor-offline-installer-${BUILD_HARBOR_REVISION}.tgz" | tar xz -C /var/tmp

# Start docker service
systemctl start docker.service
sleep 2
# Load Containers in local registry cache
harbor_containers_bundle=$(find /var/tmp -size +20M -type f -regextype sed -regex ".*/harbor\..*\.t.*z$")
docker load -i "$harbor_containers_bundle"

# Copy configuration data from tarball
mkdir /etc/vmware/harbor
cp -p /var/tmp/harbor/harbor.cfg /data/harbor
cp -pr /var/tmp/harbor/{prepare,common,docker-compose.yml,docker-compose.notary.yml} /etc/vmware/harbor

# Stop docker service
systemctl stop docker.service

function overrideDataDirectory {
FILE="$1" DIR="$2"  python - <<END
import yaml, os
dir = os.environ['DIR']
file = os.environ['FILE']
f = open(file, "r+")
dataMap = yaml.safe_load(f)
for _, s in enumerate(dataMap["services"]):
  if "restart" in dataMap["services"][s]:
      if "always" in dataMap["services"][s]["restart"]:
        dataMap["services"][s]["restart"] = "on-failure"
  if "volumes" in dataMap["services"][s]:
    for kvol, vol in enumerate(dataMap["services"][s]["volumes"]):
      if vol.startswith( '/data' ):
        dataMap["services"][s]["volumes"][kvol] = vol.replace("/data", dir, 1)
f.seek(0)
yaml.dump(dataMap, f, default_flow_style=False)
f.truncate()
f.close()
END
}

# Replace default DataDirectories in the harbor-provided compose files
overrideDataDirectory /etc/vmware/harbor/docker-compose.yml /data/harbor
overrideDataDirectory /etc/vmware/harbor/docker-compose.notary.yml /data/harbor

chmod 600 /data/harbor/harbor.cfg
chmod -R 600 /etc/vmware/harbor/common
