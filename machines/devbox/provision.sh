#!/bin/bash -e
# Copyright 2016 VMware, Inc. All Rights Reserved.
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


apt-get update && apt-get -y dist-upgrade

# set GOPATH based on shared folder of vagrant
pro="/home/${BASH_ARGV[0]}/.profile"
echo "export GOPATH=${BASH_ARGV[1]}" >> "$pro"

# add GOPATH/bin to the PATH
echo "export PATH=$PATH:${BASH_ARGV[1]}/bin" >> "$pro"

# vmwaretools automatic kernel update
echo "answer AUTO_KMODS_ENABLED yes" | tee -a /etc/vmware-tools/locations

packages=(curl lsof strace git shellcheck tree mc silversearcher-ag jq htpdate)
for package in "${packages[@]}" ; do
    apt-get -y install "$package"
done

# install / upgrade go
go_file="https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz"
go_version=$(basename $go_file | cut -d. -f1-3)
go_current=$(go version | awk '{print $(3)}')

if [[ ! -d "/usr/local/go" || "$go_current" != "$go_version" ]] ; then
    (cd /usr/local &&
        (curl --silent -L $go_file | tar -zxf -) &&
        ln -fs /usr/local/go/bin/* /usr/local/bin/)
fi

cat << EOF > /etc/systemd/system/docker.service
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network.target docker.socket
Requires=docker.socket

[Service]
Type=notify
ExecStart=/usr/bin/docker daemon -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock -D
MountFlags=slave
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=infinity

[Install]
WantedBy=multi-user.target
EOF

timedatectl set-timezone US/Pacific-New

systemctl daemon-reload
systemctl restart docker
