#!/bin/bash
# Copyright 2018 VMware, Inc. All Rights Reserved.
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

set -e

for x in $(echo GOVC_USERNAME GOVC_PASSWORD VCH_NAME GOVC_URL GOPATH GOVC_INSECURE); do
    if [[ ! -v $x ]]; then
        echo "Insufficient inputs. Please set $x environment variable and re-execute this script.";
        echo "GOVC_USERNAME: username on ESX/vCenter target"
        echo "GOVC_PASSWORD: password on ESX/vCenter target"
        echo "VCH_NAME: name of VCH; matches --name argument for vic-machine"
        echo "GOVC_URL: IP or FQDN of your vCenter/ESX target"
        echo "GOPATH: your GOPATH, obviously"
        echo "GOVC_INSECURE: set to 1 to disable tls verify when using govc to talk to ESX/vC"
        exit 1
    fi;
done

function get-thumbprint () {
    openssl s_client -connect $GOVC_URL:443 </dev/null 2>/dev/null \
        | openssl x509 -fingerprint -noout \
        | cut -d= -f2
}

$GOPATH/src/github.com/vmware/vic/bin/vic-machine-linux debug --target=$GOVC_URL --name=$VCH_NAME --user=$GOVC_USERNAME --password=$GOVC_PASSWORD --thumbprint=$(get-thumbprint)

on-vch() {
    sshpass -ppassword ssh -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no root@$VCH_IP -C $@ 2>/dev/null
}

VCH_IP="$(govc vm.ip $VCH_NAME)"

function replace-component() {
    sshpass -p 'password' scp -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no $GOPATH/src/github.com/vmware/vic/bin/$1 root@$VCH_IP:/tmp/$1
    on-vch mv /sbin/$1 /tmp/old-$1
    on-vch mv /tmp/$1 /sbin/$1
    on-vch chmod 755 /sbin/$1
    pid=$(on-vch ps -e --format='pid,args' | grep $1 | grep -v grep | awk '{print $1}')
    on-vch kill -9 $pid
    on-vch rm -f /tmp/old-$1

}

if [[ $1 == "" ]]; then
    for x in port-layer-server docker-engine-server vicadmin vic-init; do
        replace-component $x
    done
else
    replace-component $1
fi


on-vch vic-init &
echo "IP address may change when appliance finishes re-initializing. Get the new IP with govc vm.ip"
