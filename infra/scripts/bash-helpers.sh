#!/bin/bash
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
# limitations under the License.!/bin/bash

tls () {
    unset TLS_OPTS
}

no-tls () {
    export TLS_OPTS="--no-tls"
}

unset-vic () {
    unset MAPPED_NETWORKS NETWORKS IMAGE_STORE DATASTORE COMPUTE VOLUME_STORES IPADDR GOVC_INSECURE TLS
    unset DOCKER_CERT_PATH DOCKER_TLS_VERIFY
    unalias docker 2>/dev/null
}

vic-path () {
    echo "${GOPATH}/src/github.com/vmware/vic"
}

vic-create () {
    pushd $(vic-path)/bin/

    $(vic-path)/bin/vic-machine-linux create --target="$GOVC_URL" --image-store="$IMAGE_STORE" --compute-resource="$COMPUTE" ${TLS} ${TLS_OPTS} --name=${VIC_NAME:-${USER}test} ${MAPPED_NETWORKS} ${VOLUME_STORES} ${NETWORKS} ${IPADDR} ${TIMEOUT} --thumbprint=$THUMBPRINT $*

    envfile=${VIC_NAME:-${USER}test}/${VIC_NAME:-${USER}test}.env
    if [ -f "$envfile" ]; then
        set -a
        source $envfile
        set +a
    fi

    if [ -z ${DOCKER_TLS_VERIFY+x} ]; then
        alias docker='docker --tls'
    fi

    popd
}

vic-delete () {
    $(vic-path)/bin/vic-machine-linux delete --target="$GOVC_URL" --compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} --thumbprint=$THUMBPRINT --force $*
}

vic-inspect () {
    $(vic-path)/bin/vic-machine-linux inspect --target="$GOVC_URL" --compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} --thumbprint=$THUMBPRINT $*
}

vic-ls () {
    $(vic-path)/bin/vic-machine-linux ls --target="$GOVC_URL" --thumbprint=$THUMBPRINT $*
}

vic-ssh () {
    unset keyarg
    if [ -e $HOME/.ssh/authorized_keys ]; then
        keyarg="--authorized-key=$HOME/.ssh/authorized_keys"
    fi

    out=$($(vic-path)/bin/vic-machine-linux debug --target="$GOVC_URL" --compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} --enable-ssh $keyarg --rootpw=password --thumbprint=$THUMBPRINT $*)
    host=$(echo $out | grep DOCKER_HOST | sed -n 's/.*DOCKER_HOST=\([^i:]*\).*/\1/p')

    echo "SSH to ${host}"
    sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@${host}
}

addr-from-dockerhost () {
    echo $DOCKER_HOST | sed -e 's/:[0-9]*$//'
}

# import the custom sites
# example entry, actived by typing "example"
#example () {
#    target='https://user:password@host.domain.com/datacenter'
#    unset-vic
#
#    export GOVC_URL=$target
#
#    eval "export THUMBPRINT=$(govc-linux about.cert -k -json | jq -r .ThumbprintSHA1)"
#    export COMPUTE=cluster/pool
#    export DATASTORE=datastore1
#    export IMAGE_STORE=$DATASTORE/image/path
#    export NETWORKS="--bridge-network=private-dpg-vlan --external-network=extern-dpg"
#    export TIMEOUT="--timeout=10m"
#    export IPADDR="--client-network-ip=vch-hostname.domain.com --client-network-gateway=x.x.x.x/22 --dns-server=y.y.y.y --dns-server=z.z.z.z"
#    export TLS="--tls-cname=vch-hostname.domain.com --organisation=MyCompany"
#    export VIC_NAME="MyVCH"
#}

. ~/.vic
