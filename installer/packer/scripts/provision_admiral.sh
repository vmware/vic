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

mkdir /etc/vmware/admiral

BUILD_ADMIRAL_REVISION="${BUILD_ADMIRAL_REVISION:-dev}"

# start docker
echo "starting Docker .."
systemctl daemon-reload
systemctl start docker
echo "Docker started"

# pull admiral image
ADMIRAL_IMAGE="vmware/admiral:vic_${BUILD_ADMIRAL_REVISION}"
echo "Pulling Admiral Docker image.."
echo "Downloading Admiral ${ADMIRAL_IMAGE}"
docker pull ${ADMIRAL_IMAGE}
docker tag ${ADMIRAL_IMAGE} vmware/admiral:ova
echo "Pulled Admiral image"

# stop docker
echo "stopping Docker .."
systemctl stop docker
echo "Docker stopped"
