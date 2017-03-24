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
#
# Create a VM and boot stateless ESXi via cdrom/iso

set -o pipefail

usage() {
    cat <<'EOF'
Usage: $0 <host> <port-layer.log | docker.log>

GOVC_* environment variables also apply, see https://github.com/vmware/govmomi/tree/master/govc#usage
If GOVC_USERNAME is set, it is used to create an account on the ESX vm.  Default is to use the existing root account.
If GOVC_PASSWORD is set, the account password will be set to this value.  Default is to use the given ESX_URL password.
EOF
}

if [ $# -ne 2 ] ; then
    usage
    exit 1
fi

username=$GOVC_USERNAME
password=$GOVC_PASSWORD
unset GOVC_USERNAME GOVC_PASSWORD

if [ -z "$password" ] ; then
    # extract password from $GOVC_URL
    password=$(govc env | grep GOVC_PASSWORD= | cut -d= -f 2-)
fi

if [ -z "$username" ] ; then
    username=$(govc env | grep GOVC_USERNAME= | cut -d= -f 2-)
fi

jar=$(mktemp -t cookie-XXXX)

curl -k -c ${jar} --form username=${username} --form password=${password} https://$1:2378/authentication
curl -k -b ${jar} https://$1:2378/logs/tail/$2

rm $jar
