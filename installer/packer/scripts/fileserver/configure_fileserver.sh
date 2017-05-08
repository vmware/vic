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

data_dir=/opt/vmware/fileserver
cert_dir=${data_dir}/cert
cert=${cert_dir}/server.crt
key=${cert_dir}/server.key

port=$(ovfenv -k fileserver.port)
fileserver_cert=$(ovfenv -k fileserver.ssl_cert)
fileserver_key=$(ovfenv -k fileserver.ssl_cert_key)

if [ -z "$port" ]; then
  port="9443"
fi

iptables -w -A INPUT -j ACCEPT -p tcp --dport $port

#Format cert file
function formatCert {
  content=$1
  file=$2
  echo $content | sed -r 's/(-{5}BEGIN [A-Z ]+-{5})/&\n/g; s/(-{5}END [A-Z ]+-{5})/\n&\n/g' | sed -r 's/.{64}/&\n/g; /^\s*$/d' > $file
}

if [[ x$fileserver_cert != "x" && x$fileserver_key != "x" ]]; then
  mkdir -p "$cert_dir"
  formatCert "$fileserver_cert" $cert
  formatCert "$fileserver_key" $key
fi
