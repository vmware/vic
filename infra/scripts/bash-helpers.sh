#!/bin/bash
# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

BASE_DIR=$(dirname $(readlink -f "$BASH_SOURCE"))
OS=$(uname | tr '[:upper:]' '[:lower:]')

unset-vic () {
    unset TARGET_URL 
    MAPPED_NETWORKS NETWORKS IMAGE_STORE DATASTORE COMPUTE VOLUME_STORES IPADDR GOVC_INSECURE TLS THUMBPRINT OPS_CREDS VIC_NAME PRESERVE_VOLUMES
}

vic-path () {
    echo "${GOPATH}/src/github.com/vmware/vic"
}

vic-create () {
    vicProfileTranscode

    base=$(pwd)
    (
        cd "$(vic-path)"/bin || return
        "$(vic-path)"/bin/vic-machine-"$OS" create --target="$TARGET_URL" "${OPS_CREDS[@]}" --image-store="$IMAGE_STORE" --compute-resource="$COMPUTE" "${TLS[@]}" ${TLS_OPTS} --name="${VIC_NAME:-${USER}test}" "${MAPPED_NETWORKS[@]}" "${VOLUME_STORES[@]}" "${NETWORKS[@]}" ${IPADDR} ${TIMEOUT} --thumbprint="$THUMBPRINT" "$@"
    )

    vic-select

    cd "$base" || exit
}

vic-delete () {
    vicProfileTranscode

    force="true"
    if [ -n "${PRESERVE_VOLUMES}" ]; then
        force="false"
    fi

    "$(vic-path)"/bin/vic-machine-"$OS" delete --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --thumbprint="$THUMBPRINT" --force=${force} "$@"
}

vic-select () {
    vicProfileTranscode

    base=$(pwd)

    unset DOCKER_HOST DOCKER_CERT_PATH DOCKER_TLS_VERIFY
    unalias docker 2>/dev/null

    envfile=$(vic-path)/bin/${VIC_NAME:-${USER}test}/${VIC_NAME:-${USER}test}.env
    if [ -f "$envfile" ]; then
        set -a
        source "$envfile"
        set +a
    fi

    # Something of a hack, but works for --no-tls so long as that's enabled via TLS_OPTS
    if [ -z "${DOCKER_TLS_VERIFY+x}" ] && [[ "${DOCKER_HOST+x}" != "*:2375" ]]; then
        alias docker='docker --tls'
    fi
}

vic-inspect () {
    vicProfileTranscode

    "$(vic-path)"/bin/vic-machine-"$OS" inspect --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --thumbprint="$THUMBPRINT" "$@"
}

vic-upgrade () {
    vicProfileTranscode

    "$(vic-path)"/bin/vic-machine-"$OS" upgrade --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --thumbprint="$THUMBPRINT" "$@"
}

vic-ls () {
    vicProfileTranscode

    "$(vic-path)"/bin/vic-machine-"$OS" ls --target="$TARGET_URL" --thumbprint="$THUMBPRINT" "$@"
}

vic-ssh () {
    vicProfileTranscode

    unset keyarg
    if [ -e "$HOME"/.ssh/authorized_keys ]; then
        keyarg="--authorized-key=$HOME/.ssh/authorized_keys"
    fi

    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh $keyarg --rootpw=password --thumbprint="$THUMBPRINT" "$@")
    host=$(echo "$out" | grep DOCKER_HOST | awk -F"DOCKER_HOST=" '{print $2}' | cut -d ":" -f1 | cut -d "=" -f2)

    echo "SSH to ${host}"
    sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@"${host}"
}

vic-admin () {
    vicProfileTranscode

    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh "$keyarg" --rootpw=password --thumbprint="$THUMBPRINT" "$@")
    host=$(echo "$out" | grep DOCKER_HOST | sed -n 's/.*DOCKER_HOST=\([^:\s*\).*/\1/p')

    open http://"${host}":2378
}

addr-from-dockerhost () {
    echo "$DOCKER_HOST" | sed -e 's/:[0-9]*$//'
}

vic-tail-portlayer() {
    vicProfileTranscode

    unset keyarg
    if [ -e "$HOME"/.ssh/authorized_keys ]; then
        keyarg="--authorized-key=$HOME/.ssh/authorized_keys"
    fi

    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh "$keyarg" --rootpw=password --thumbprint="$THUMBPRINT" "$@")
    host=$(echo "$out" | grep DOCKER_HOST | awk -F"DOCKER_HOST=" '{print $2}' | cut -d ":" -f1 | cut -d "=" -f2)

    sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@"${host}" tail -f /var/log/vic/port-layer.log
}

vic-tail-docker() {
    vicProfileTranscode    

    unset keyarg
    if [ -e "$HOME"/.ssh/authorized_keys ]; then
        keyarg="--authorized-key=$HOME/.ssh/authorized_keys"
    fi

    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh "$keyarg" --rootpw=password --thumbprint="$THUMBPRINT" "$@")
    host=$(echo "$out" | grep DOCKER_HOST | awk -F"DOCKER_HOST=" '{print $2}' | cut -d ":" -f1 | cut -d "=" -f2)

    sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@"${host}" tail -f /var/log/vic/docker-personality.log
}

# turns the configuration variables into the ones expected by the vic-x functions
vicProfileTranscode () {
    export COMPUTE=cluster/pool
    export DATASTORE=datastore1
    export TIMEOUT="--timeout=10m"
    export IPADDR="--client-network-ip=vch-hostname.domain.com --client-network-gateway=x.x.x.x/22 --dns-server=y.y.y.y --dns-server=z.z.z.z"
    export VIC_NAME="MyVCH"

    TLS=("--tls-cname=vch-hostname.domain.com" "--organization=MyCompany")
    OPS_CREDS=("--ops-user=<user>" "--ops-password=<password>")
    NETWORKS=("--bridge-network=private-dpg-vlan" "--public-network=extern-dpg")
    MAPPED_NETWORKS=("--container-network=VM Network:external" "--container-network=SomeOtherNet:elsewhere")
    VOLUME_STORES=("--volume-store=$DATASTORE:default")

    ## additional boiler plate to complete vic-machine setup and ease use of govc with the profiles
    #
    # configure thumbprint for vic-machine and known hosts for govc
    known_host=$(govc about.cert -u host -k -thumbprint)
    export THUMBPRINT="$(echo ${known_host##* })"
    echo "${known_host}" >> "${GOVC_TLS_KNOWN_HOSTS}"

    #   set vSphere target for vic-machine and govc
    export TARGET_URL="https://${user}:${password}@${vsphere}/${datacenter}"
    export GOVC_URL="https://${vsphere}"
    export GOVC_USERNAME="${user}"
    export GOVC_PASSWORD="${password}"

    #   set other govc variables
    export GOVC_DATACENTER="${datacenter}"
    export GOVC_DATASTORE="${datastore}"

    #   set image store to a predictable path on the datastore based on profile name
    IMAGE_STORE="${datastore}/${vch_name}"

    COMPUTE="${compute}"
    DATASTORE="${datastore}"
    TIMEOUT="${timeout}"
    VIC_NAME="${vch_name}"
    TLS="${tls[@]}"
    VOLUME_STORES="${volumestores}"
    MAPPED_NETWORKS="${containernet}"
    PRESERVE_VOLUMES="${preserveVolumestores}"

    if [ ! -z "${opsuser} ]; then
        OPS_CREDS="--ops-user=${opsuser} --ops-password="${opspass}"
    fi
    
    #   configure IPADDR from the raw argument values
    for ns in "$dns[@]"; do
        IPADDR=" ${IPADDR} --dns-server=${ns}"
    done

    if [ -z "${clientIP}" ]; then
        IPADDR=" --client-network-ip=${clientIP}"
    fi
    if [ -z "${clientGW}" ]; then
        IPADDR=" --client-network-gateway=${clientGW}"
    fi
    if [ -z "${managementIP}" ]; then
        IPADDR=" --management-network-ip=${managementIP}"
    fi
    if [ -z "${managementGW}" ]; then
        IPADDR=" --management-network-gateway=${managementGW}"
    fi
    if [ -z "${bridgeRange}" ]; then
        IPADDR=" --bridge-network-range=${bridgeRange}"
    fi

    NETWORKS=()
    if [ -z "${bridgeNet}" ]; then
        NETWORKS+=("--bridge-network=${bridgeNet}")
    fi
    if [ -z "${publicNet}" ]; then
        NETWORKS+=("--public-network=${publicNet}")
    fi
    if [ -z "${clientNet}" ]; then
        NETWORKS+=("--client-network=${clientNet}")
    fi
    if [ -z "${managementNet}" ]; then
        NETWORKS+=("--management-network=${managementNet}")
    fi

    # export all of the vic-machine control variables
    export TARGET_URL COMPUTE DATASTORE TIMEOUT IPADDR VIC_NAME NETWORKS MAPPED_NETWORKS VOLUME_STORES OPS_CREDS TLS PRESERVE_VOLUMES
}

# import the custom sites
if [ ! -r ~/.vic ]; then
    echo "The bash-helpers depend on a \"$HOME/.vic\" file that contains profiles for different VCH configurations."
    echo "There is a sample profile file $BASE_DIR/sample-helper.profiles"
else
    . $HOME/.vic
fi
