#!/bin/bash
# Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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

set -eo pipefail

SSH="ssh -o StrictHostKeyChecking=no"
SCP="scp -q -o StrictHostKeyChecking=no"

REMOTE_DLV_ATTACH=/usr/local/bin/dlv-attach-headless.sh
REMOTE_DLV_DETACH=/usr/local/bin/dlv-detach-headless.sh

function usage() {
    echo "Usage: $0 -h vch-address [-a/-d] [attach/detach] target" >&2
    echo "Valid targets are: "
    echo "    vic-init"
    echo "    vic-admin"
    echo "    docker-engine"
    echo "    port-layer"
    echo "    vic-machine"
    exit 1
}

while getopts "h:adl" flag
do
    case $flag in

        h)
            # Optional
            export DLV_TARGET_HOST="$OPTARG"
            ;;

        a)
            export COMMAND="attach"
            ;;

        d)
            export COMMAND="detach"
            ;;

        l)
            export COMMAND="log"
            ;;

        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ -n "${DOCKER_HOST}" -a -z "${DLV_TARGET_HOST}" ]; then
    export DLV_TARGET_HOST=$(echo $DOCKER_HOST | cut -d ':' -f 1)
fi

if [ -z "${DLV_TARGET_HOST}" ]; then
    usage
fi

if [[ -z "${COMMAND}" &&  $# != 2 ]]; then
    usage
elif [[ -n "${COMMAND}" && $# != 1 ]]; then
    usage
fi

if [ -z "${COMMAND}" ]; then
    COMMAND=$1
    TARGET=$2
else
    TARGET=$1
fi

case ${TARGET} in

    vic-init)
        PORT=2345
        ;;

    vic-admin)
        # Change target to vicadmin
        TARGET=vicadmin
        PORT=2346
        ;;

    docker-engine)
        PORT=2347
        ;;

    port-layer)
        PORT=2348
        ;;

    vic-machine)
        PORT=2349
        ;;

    *)
        usage
        ;;
esac

if [ -z "${DLV_TARGET_HOST}" ]; then
    usage
fi

if [ ${COMMAND} == "attach" ]; then
    ${SSH} root@${DLV_TARGET_HOST} "nohup /usr/local/bin/dlv-attach-headless.sh $TARGET $PORT > /var/tmp/${TARGET}.log 2>&1 &"
elif [ ${COMMAND} == "detach" ]; then
    ${SSH} root@${DLV_TARGET_HOST} "/usr/local/bin/dlv-detach-headless.sh $PORT"
elif [ ${COMMAND} == "log" ]; then
    ${SSH} root@${DLV_TARGET_HOST} "cat /var/tmp/${TARGET}.log"
else
    usage
fi

echo $DLV_TARGET_HOST:$PORT
