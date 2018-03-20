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
# limitations under the License.

[ -n "$DEBUG" ] && set -x

BASE_DIR=$(dirname $(readlink -f "$BASH_SOURCE"))
REPO_DIR=$(pwd "${BASE_DIR}/../../")
OS=$(uname | tr '[:upper:]' '[:lower:]')

# global settings that reduce friction
if [ -z "${GOVC_TLS_KNOWN_HOSTS}" ]; then
    export GOVC_TLS_KNOWN_HOSTS=~/.govmomi/known_hosts
fi

init-profile () {
    unset-vic
    vch_name=${FUNCNAME[1]}
}

unset-vic () {
    unset TARGET_URL MAPPED_NETWORKS NETWORKS IMAGE_STORE DATASTORE COMPUTE VOLUME_STORES IPADDR TLS THUMBPRINT OPS_CREDS VIC_NAME PRESERVE_VOLUMES
    unset GOVC_URL GOVC_INSECURE GOVC_DATACENTER GOVC_USERNAME GOVC_PASSWORD GOVC_DATASTORE GOVC_CERTIFICATE 

    unset vsphere datacenter user password opsuser opspass opsgrant timeout compute datastore dns publicNet publicIP publicGW bridgeNet bridgeRange
    unset clientNet clientIP clientGW managementNet managementIP managementGW tls volumestores preserveVolumestores containernet
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

    # change to the bin directory as that's where our certs would have been generated by vic-create
    (
        cd "$(vic-path)"/bin || return
        "$(vic-path)"/bin/vic-machine-"$OS" inspect --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --thumbprint="$THUMBPRINT" "$@"
    )
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

    # cannot execute in subshell as we want access to the environment variables that result
    pushd "$(vic-path)"/bin >/dev/null
    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh $keyarg --rootpw=password --thumbprint="$THUMBPRINT")
    host=$(echo "$out" | grep DOCKER_HOST | awk -F"DOCKER_HOST=" '{print $2}' | cut -d ":" -f1 | cut -d "=" -f2)
    popd >/dev/null

    echo "SSH to ${host}"
    sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@"${host}" "$@"
}

vic-admin () {
    vicProfileTranscode

    out=$("$(vic-path)"/bin/vic-machine-"$OS" debug --target="$TARGET_URL" --compute-resource="$COMPUTE" --name="${VIC_NAME:-${USER}test}" --enable-ssh $keyarg --rootpw=password --thumbprint="$THUMBPRINT" "$@")
    host=$(echo "$out" | grep DOCKER_HOST | sed -n 's/.*DOCKER_HOST=\([^:\s*\).*/\1/p')

    open http://"${host}":2378
}

addr-from-dockerhost () {
    echo "$DOCKER_HOST" | sed -e 's/:[0-9]*$//'
}

vic-tail-portlayer() {
    vic-ssh tail -f /var/log/vic/port-layer.log
}

vic-tail-docker() {
    vic-ssh tail -f /var/log/vic/docker-personality.log
}

# Transcodes the configuration variables from the new profile format into the ones expected by the vic-x functions
# This exists so that profiles using the old variable structure will still function as expected
vicProfileTranscode() {
    # temporary check to see if we're using an old profile format
    if [ -z "${vsphere}" ]; then
        export TARGET_URL="${GOVC_URL}"
        return
    fi

    # configure thumbprint for vic-machine and known hosts for govc
    known_host=$(govc about.cert -u ${vsphere} -k -thumbprint)
    THUMBPRINT="$(echo ${known_host##* })"
    echo "${known_host}" >> "${GOVC_TLS_KNOWN_HOSTS}"

    # set vSphere target for vic-machine and govc
    TARGET_URL="https://${user}:${password}@${vsphere}/${datacenter}"
    GOVC_URL="https://${vsphere}"
    GOVC_USERNAME="${user}"
    GOVC_PASSWORD="${password}"

    # set other govc variables
    GOVC_DATACENTER="${datacenter}"
    GOVC_DATASTORE="${datastore}"
    export GOVC_URL GOVC_USERNAME GOVC_PASSWORD GOVC_DATACENTER GOVC_DATASTORE

    # set image store to a predictable path on the datastore based on profile name
    IMAGE_STORE="${datastore}/${vch_name}"

    COMPUTE="${compute}"
    DATASTORE="${datastore}"
    TIMEOUT="${timeout}"
    VIC_NAME="${vch_name}"
    TLS=("${tls[@]}")
    VOLUME_STORES=("${volumestores[@]}")
    MAPPED_NETWORKS=("${containernet[@]}")
    PRESERVE_VOLUMES="${preserveVolumestores}"

    if [ -n "${opsuser}" ]; then
        OPS_CREDS=("--ops-user=${opsuser}" "--ops-password=${opspass}")

        if [ -n "${opsgrant}" ]; then
            OPS_CREDS+=("--ops-grant-perms")
        fi
    fi

    #   configure IPADDR from the raw argument values
    IPADDR=""
    for ns in "${dns[@]}"; do
        IPADDR+=" --dns-server=${ns}"
    done

    if [ -n "${publicIP}" ]; then
        IPADDR+=" --public-network-ip=${publicIP}"
    fi
    if [ -n "${publicGW}" ]; then
        IPADDR+=" --public-network-gateway=${publicGW}"
    fi
    if [ -n "${clientIP}" ]; then
        IPADDR+=" --client-network-ip=${clientIP}"
    fi
    if [ -n "${clientGW}" ]; then
        IPADDR+=" --client-network-gateway=${clientGW}"
    fi
    if [ -n "${managementIP}" ]; then
        IPADDR+=" --management-network-ip=${managementIP}"
    fi
    if [ -n "${managementGW}" ]; then
        IPADDR+=" --management-network-gateway=${managementGW}"
    fi
    if [ -n "${bridgeRange}" ]; then
        IPADDR+=" --bridge-network-range=${bridgeRange}"
    fi

    if [ -z "${publicIP}" -a -z "${clientIP}" ]; then
        noverify="--no-tlsverify"
        for tlsopt in "${tls[@]}"; do
            if [ "${tlsopt}" == "--tls-cname=*" -o "${tlsopt}" == "--no-tls" ]; then
                unset noverify
            fi
        done
        TLS+=($noverify)
    fi

    NETWORKS=()
    if [ -n "${bridgeNet}" ]; then
        NETWORKS+=("--bridge-network=${bridgeNet}")
    fi
    if [ -n "${publicNet}" ]; then
        NETWORKS+=("--public-network=${publicNet}")
    fi
    if [ -n "${clientNet}" ]; then
        NETWORKS+=("--client-network=${clientNet}")
    fi
    if [ -n "${managementNet}" ]; then
        NETWORKS+=("--management-network=${managementNet}")
    fi

    # export all of the vic-machine control variables
    export TARGET_URL THUMBPRINT COMPUTE DATASTORE TIMEOUT IPADDR VIC_NAME NETWORKS MAPPED_NETWORKS VOLUME_STORES OPS_CREDS TLS PRESERVE_VOLUMES
}

# reset profile if resourced
unset-vic

# import the custom sites
if [ ! -r "$HOME/.vic" -a ! -r "$REPO_DIR/.vic.profiles" ]; then
    echo "The bash-helpers depend on files that contains profiles for different VCH configurations:"
    echo " \"$HOME/.vic\" and/or"
    echo " \"$REPO_DIR/.vic.profiles\""
    echo
    echo "The profiles from your home directory will override the others if name collisions exist"
    echo "There is a sample profile file $BASE_DIR/sample-helper.profiles"
else
    echo "Loading profiles - a new profile selection must be made"
    [ -r "$REPO_DIR/.vic.profiles" ] && . "$REPO_DIR/.vic.profiles"
    [ -r "$HOME/.vic" ] && . "$HOME/.vic"
fi
