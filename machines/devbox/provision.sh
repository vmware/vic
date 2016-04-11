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
pro="/home/"${BASH_ARGV[0]}"/.profile"
echo "export GOPATH="${BASH_ARGV[1]} >> $pro

# add GOPATH/bin to the PATH
echo "export PATH=$PATH:"${BASH_ARGV[1]}"/bin" >> $pro

# vmwaretools automatic kernel update
echo "answer AUTO_KMODS_ENABLED yes" | tee -a /etc/vmware-tools/locations

packages=(curl lsof strace git shellcheck tree mc silversearcher-ag jq)
for package in "${packages[@]}" ; do
    apt-get -y install "$package"
done

if [ ! -d "/usr/local/go" ] ; then
    (cd /usr/local &&
        (curl --silent -L https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz | tar -zxf -) &&
        ln -s /usr/local/go/bin/* /usr/local/bin/)
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

systemctl daemon-reload
systemctl restart docker
