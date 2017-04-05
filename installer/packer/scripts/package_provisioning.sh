#!/bin/sh
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

# Install sudo and ESX-optimized kernel
tdnf install -y sudo linux-esx rsync lvm2 docker gawk parted tar openjre

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/1.11.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Remove unused packages
tdnf remove -y cloud-init

# Create directory to host VMware-specific scripts
mkdir /etc/vmware
mkdir "/usr/lib/systemd/system/getty@tty2.service.d"

# Create Data Dir
mkdir /data

# Create filesystem on /dev/sdb to be mounted as /data
parted -a optimal --script /dev/sdb 'mklabel gpt mkpart primary ext4 0% 100%'
mkfs.ext4 /dev/sdb1
tune2fs -L vic-data-v1 /dev/sdb1

# Seed directories in /data
mount /dev/sdb1 /data -t ext4
mkdir -p /data/{admiral,harbor,fileserver}
mkdir -p /data/fileserver/files
