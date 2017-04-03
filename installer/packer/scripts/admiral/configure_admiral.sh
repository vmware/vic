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

deploy=$(ovfenv -k management_portal.deploy)
port=$(ovfenv -k management_portal.port)

if [ ${deploy,,} != "true" ]; then
  echo "Not configuring Admiral and disabling startup"
  systemctl stop admiral
  exit 0
fi

data_dir=/data/admiral
conf_dir=/etc/vmware/admiral
javadir=/var/opt/OpenJDK-1.8.0.112-bin/bin

cert_dir=${data_dir}/cert
flag=${conf_dir}/cert_gen_type
admiral_start_script=${conf_dir}/start_admiral.sh
cfg=${data_dir}/admiral.cfg

ca_download_dir=${data_dir}/ca_download
mkdir -p {${cert_dir},${ca_download_dir}}

cert=${cert_dir}/server.crt
key=${cert_dir}/server.key
jks=${cert_dir}/trustedcertificates.jks
csr=${cert_dir}/server.csr
ca_cert=${cert_dir}/ca.crt
ca_key=${cert_dir}/ca.key
ext=${cert_dir}/extfile.cnf

rm -rf $ca_download_dir/*

#Configure attr in start_admiral.sh
function configureAdmiralStart {
  cfg_key=$1
  cfg_value=$2

  basedir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

  if [ -n "$cfg_key" ]; then
    cfg_value=$(echo "$cfg_value" | sed -r -e 's%[\/&%]%\\&%g')
    sed -i -r "s%#?$cfg_key\s*=\s*.*%$cfg_key=$cfg_value%" $admiral_start_script
  fi
}

#Format cert file
function formatCert {
  content=$1
  file=$2
  echo $content | sed -r "s/-----BEGIN CERTIFICATE-----/-----BEGIN CERTIFICATE-----\n/g" | sed -r "s/-----END CERTIFICATE-----/\n-----END CERTIFICATE-----\n/g"
}

#Format key file
function formatKey {
  content=$1
  file=$2
  echo $content | sed -r "s/-----BEGIN.*PRIVATE KEY-----/-----BEGIN RSA PRIVATE KEY-----\n/g" | sed -r "s/-----END.*PRIVATE KEY-----/\n-----END RSA PRIVATE KEY-----\n/g"
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

  echo "creating java keystore with self-signed CA"
  $javadir/keytool -import -noprompt -v -trustcacerts -alias selfsignedca -file $ca_cert -keystore $jks -keypass changeit -storepass changeit
}

function secure {
  ssl_cert=$(ovfenv -k management_portal.ssl_cert)
  ssl_cert_key=$(ovfenv -k management_portal.ssl_cert_key)
  if [ -n "$ssl_cert" ] && [ -n "$ssl_cert_key" ]; then
    echo "ssl_cert and ssl_cert_key are both set, using customized certificate"
    formatCert "$ssl_cert" $cert
    formatKey "$ssl_cert_key" $key
    echo "customized" > $flag
    echo "creating java keystore with provided cert for xenon"
    $javadir/keytool -import -noprompt -v -trustcacerts -alias selfsignedca -file $cert -keystore $jks -keypass changeit -storepass changeit
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
  if [ -n "$hostname" ]; then
    if [ "$hostname" = "localhost.localdomain" ]; then
      hostname=""
      return
    fi
    echo "Get hostname from command 'hostnamectl status --static': $hostname"
    return
  fi
}

hostname=""
ip_address=$(ip addr show dev eth0 | sed -nr 's/.*inet ([^ ]+)\/.*/\1/p')

#Modify hostname
detectHostname
if [[ x$hostname != "x" ]]; then
  echo "Hostname: ${hostname}"
  configureAdmiralStart "hostname" ${hostname}
else
  echo "Hostname is null, set it to IP"
  hostname=${ip_address}
fi

# Init certs
secure

configureAdmiralStart ADMIRAL_DATA_LOCATION $data_dir
configureAdmiralStart ADMIRAL_CERT_LOCATION $cert
configureAdmiralStart ADMIRAL_KEY_LOCATION $key
configureAdmiralStart ADMIRAL_JKS_LOCATION $jks
configureAdmiralStart ADMIRAL_EXPOSED_PORT $port

iptables -A INPUT -j ACCEPT -p tcp --dport $port

touch $data_dir/custom.conf

harbor_deploy=$(ovfenv -k registry.deploy)

if [ ${harbor_deploy,,} == "true" ]; then
  harbor_port=$(ovfenv -k registry.port)
  # If harbor is deployed, configure the integration URL
  echo "harbor.tab.url=https://${hostname}:${harbor_port}" > $data_dir/custom.conf
fi
