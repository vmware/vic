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

mask2cdr ()
{
  set -- 0^^^128^192^224^240^248^252^254^ ${#1} ${1##*255.} 
  set -- $(( ($2 - ${#3})*2 )) ${1%%${3%%.*}*} 
  echo $(( $1 + (${#2}/4) )) 
}

network_address=$(ovfenv --key network.ip0)
network_conf_file=/etc/systemd/network/09-vic.network

if [[ x$network_address != "x" ]]; then
  # If IP is configured via OVF environment, we create a file for systemd-networkd to parse
  network_cidr=$(mask2cdr $(ovfenv --key network.netmask0))

  cat <<EOF > ${network_conf_file}
[Match]
Name=eth0

[Network]
Address=${network_address}/${network_cidr}
Gateway=$(ovfenv --key network.gateway)
DNS=$(ovfenv --key network.DNS)
Domains=$(ovfenv --key network.searchpath)
EOF

  chmod 644 ${network_conf_file}
else
  # If previous configuration exists, we remove it
  [ -e ${network_conf_file} ] && rm ${network_conf_file} || exit 0
fi