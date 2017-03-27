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
set -x -euf -o pipefail

deploy=$(ovfenv -k harbor.deploy)

if [ ${deploy,,} != "true" ]; then
  echo "Not configuring Harbor and disabling startup"
  systemctl stop harbor
  exit 0
fi

data_dir=/data/harbor
conf_dir=/etc/vmware/harbor

harbor_compose_file=${conf_dir}/docker-compose.yml

cert_dir=${data_dir}/cert
flag=${conf_dir}/cert_gen_type
cfg=${data_dir}/harbor.cfg

ca_download_dir=${data_dir}/ca_download
mkdir -p {${cert_dir},${ca_download_dir}}

cert=${cert_dir}/server.crt
key=${cert_dir}/server.key
csr=${cert_dir}/server.csr
ca_cert=${cert_dir}/ca.crt
ca_key=${cert_dir}/ca.key
ext=${cert_dir}/extfile.cnf

rm -rf $ca_download_dir/*

#Configure attr in harbor.cfg
function configureHarborCfg {
  cfg_key=$1
  cfg_value=$2

  basedir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

  if [ -n "$cfg_key" ]; then
    cfg_value=$(echo "$cfg_value" | sed -r -e 's%[\/&%]%\\&%g')
    sed -i -r "s%#?$cfg_key\s*=\s*.*%$cfg_key = $cfg_value%" $cfg
  fi
}

function format {
  file=$1
  head=$(sed -rn 's/(-+[A-Za-z ]*-+)([^-]*)(-+[A-Za-z ]*-+)/\1/p' $file)
  body=$(sed -rn 's/(-+[A-Za-z ]*-+)([^-]*)(-+[A-Za-z ]*-+)/\2/p' $file)
  tail=$(sed -rn 's/(-+[A-Za-z ]*-+)([^-]*)(-+[A-Za-z ]*-+)/\3/p' $file)
  echo $head > $file
  echo $body | sed  's/\s\+/\n/g' >> $file
  echo $tail >> $file
}

function genCert {
  if [ ! -e $ca_cert ] || [ ! -e $ca_key ]
  then
    openssl req -newkey rsa:4096 -nodes -sha256 -keyout $ca_key \
      -x509 -days 365 -out $ca_cert -subj \
      "/C=US/ST=California/L=Palo Alto/O=VMware, Inc./OU=Containers on vSphere/CN=Self-signed by VMware, Inc."
  fi
  openssl req -newkey rsa:4096 -nodes -sha256 -keyout $key \
    -out $csr -subj \
    "/C=US/ST=California/L=Palo Alto/O=VMware/OU=Containers on vSphere/CN=$hostname"

  echo "Add subjectAltName = IP: $ip_address to certificate"
  echo subjectAltName = IP:$ip_address > $ext
  openssl x509 -req -days 365 -in $csr -CA $ca_cert -CAkey $ca_key -CAcreateserial -extfile $ext -out $cert

  echo "self-signed" > $flag
  echo "Copy CA certificate to $ca_download_dir"
  cp $ca_cert $ca_download_dir/
}

function secure {
  ssl_cert=$(ovfenv -k harbor.ssl_cert)
  ssl_cert_key=$(ovfenv -k harbor.ssl_cert_key)
  if [ -n "$ssl_cert" ] && [ -n "$ssl_cert_key" ]; then
    echo "ssl_cert and ssl_cert_key are both set, using customized certificate"
    echo $ssl_cert > $cert
    format $cert
    echo $ssl_cert_key > $key
    format $key
    echo "customized" > $flag
    return
  fi

  if [ ! -e $ca_cert ] || [ ! -e $cert ] || [ ! -e $key ]; then
    echo "CA, Certificate or key file does not exist, will generate a self-signed certificate"
    genCert
    return
  fi

  if [ ! -e $flag ]; then
    echo "The file which records the way generating certificate does not exist, will generate a new self-signed certificate"
    genCert
    return
  fi

  if [ ! $(cat $flag) = "self-signed" ]; then
    echo "The way generating certificate changed, will generate a new self-signed certificate"
    genCert
    return
  fi

  cn=$(openssl x509 -noout -subject -in $cert | sed -n '/^subject/s/^.*CN=//p') || true
  if [ "$hostname" !=  "$cn" ]; then
    echo "Common name changed: $cn -> $hostname , will generate a new self-signed certificate"
    genCert
    return
  fi

  ip_in_cert=$(openssl x509 -noout -text -in $cert | sed -n '/IP Address:/s/.*IP Address://p') || true
  if [ "$ip_address" !=  "$ip_in_cert" ]; then
    echo "IP changed: $ip_in_cert -> $ip_address , will generate a new self-signed certificate"
    genCert
    return
  fi

  echo "Use the existing CA, certificate and key file"
  echo "Copy CA certificate to $ca_download_dir"
  cp $ca_cert $ca_download_dir/
}

function detectHostname {
  hostname=$(hostnamectl status --static) || true
  if [ -n $hostname ]; then
    if [ "$hostname" = "localhost.localdomain" ]; then
      hostname=""
      return
    fi
    echo "Get hostname from command 'hostnamectl status --static': $hostname"
    return
  fi
}

function setPortInYAML {
FILE="$1" HPORT="$2" NPORT="$3" python - <<END
import yaml, os
harbor_port = os.environ['HPORT']
notary_port = os.environ['NPORT']
file = os.environ['FILE']
f = open(file, "r+")
dataMap = yaml.safe_load(f)
newports = ['{0}:443'.format(harbor_port),'{0}:4443'.format(notary_port)]
dataMap["services"]["proxy"]["ports"] = newports
f.seek(0)
yaml.dump(dataMap, f, default_flow_style=False)
f.truncate()
f.close()
END
}

attrs=( 
  appliance.email_server 
  appliance.email_server_port 
  appliance.email_username 
  appliance.email_password 
  appliance.email_from 
  appliance.email_ssl
  harbor.admin_password
  harbor.db_password
  harbor.gc_enabled
)

hostname=""
ip_address=$(ip addr show dev eth0 | sed -nr 's/.*inet ([^ ]+)\/.*/\1/p')

#Modify hostname
detectHostname
if [[ x$hostname != "x" ]]; then
  echo "Hostname: ${hostname}"
else
  echo "Hostname is null, set it to IP"
  hostname=${ip_address}
fi
configureHarborCfg "hostname" ${hostname}:$(ovfenv -k harbor.port)

configureHarborCfg ui_url_protocol https
secure

configureHarborCfg ssl_cert $cert
configureHarborCfg ssl_cert_key $key
configureHarborCfg secretkey_path $data_dir

for attr in "${attrs[@]}"
do
  echo "Read attribute using ovfenv: [ $attr ]"
  value=$(ovfenv -k $attr)
  
  configureHarborCfg $(echo ${attr} | cut -d. -f2) "$value"
done

setPortInYAML $harbor_compose_file $(ovfenv -k harbor.port) $(ovfenv -k harbor.notary_port)

admiral_deploy=$(ovfenv -k admiral.deploy)

if [ ${admiral_deploy,,} == "true" ]; then
  # If admiral is deployed, configure the integration URL
  configureHarborCfg admiral_url https://${hostname}:$(ovfenv -k admiral.port)
fi
