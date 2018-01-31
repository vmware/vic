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

BASE_DIR=$(readlink -f ../../)

# Grab vSphere's thumbprint for calling vic-machine
function get-thumbprint () {
    govc about.cert | grep "SHA-1 Thumbprint" | awk '{print $NF}'
}

# Run the command given on the VCH instead of locally
function on-vch() {
    sshpass -ppassword ssh -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no root@$VCH_IP -C $@ 2>/dev/null
}

# SCPs the component in $1 to the VCH, plops it in place, and brutally kills the previous running process
function replace-component() {
    sshpass -p 'password' scp -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no $BASEDIR/bin/$1 root@$VCH_IP:/tmp/$1
    on-vch mv /sbin/$1 /tmp/old-$1
    on-vch mv /tmp/$1 /sbin/$1
    on-vch chmod 755 /sbin/$1
    pid=$(on-vch ps -e --format='pid,args' | grep $1 | grep -v grep | awk '{print $1}')
    on-vch kill -9 $pid
    on-vch rm -f /tmp/old-$1
}

# Check GOVC vars
if [[ ! $(govc ls 2>/dev/null) ]]; then
        echo "GOVC environment variables are required to use this script. Set the necessary variables to allow govc to connect to your vSphere deployment:";
        echo "GOVC_USERNAME: username on vSphere target"
        echo "GOVC_PASSWORD: password on vSphere target"
        echo "GOVC_URL: IP or FQDN of your vSphere target"
        echo "GOVC_INSECURE: set to 1 to disable tls verify when using govc to talk to vSphere"
        exit 1
fi

# Check for our one required argument
if [[ -v $VIC_NAME ]]; then
    echo "Please set the following environment variable to specify the VCH which you would like to reconfigure:"
    echo "VIC_NAME: name of VCH; matches --name argument for vic-machine"
    exit 1
fi

echo "Enabling SSH access on your VCH"
$BASEDIR/bin/vic-machine-linux debug --target=$GOVC_URL --name=$VIC_NAME --user=$GOVC_USERNAME --password=$GOVC_PASSWORD --thumbprint=$(get-thumbprint)


VCH_IP="$(govc vm.ip $VIC_NAME)"


if [[ $1 == "" ]]; then # replace everything
    for x in port-layer-server docker-engine-server vicadmin vic-init; do
        replace-component $x
    done
else
    for x in $@; do # replace provided list
        replace-component $x
    done
fi

on-vch vic-init &
echo "IP address may change when appliance finishes re-initializing. Get the new IP with govc vm.ip $VIC_NAME"
