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

# TODO(frapposelli): Parametrize download url
curl -L "http://192.168.90.1:8000/installer.tgz" | tar xz -C /var/tmp

# TODO(frapposelli): Parametrize file name
tar xvfz /var/tmp/harbor-offline-installer-dev.tgz -C /var/tmp

# Start docker service
systemctl start docker.service 
sleep 2
# Load Containers in local registry cache
# TODO(frapposelli): parametrize file name
docker load -i /var/tmp/harbor/harbor.dev.tgz

# Copy configuration data from tarball
mkdir /etc/vmware/harbor
cp -pr /var/tmp/harbor/{prepare,common} /data/harbor

# Stop docker service
systemctl stop docker.service
