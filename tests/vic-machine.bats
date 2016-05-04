# Copyright 2016 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/usr/bin/env bats

load helpers/helpers
load 'helpers/bats-support/load'
load 'helpers/bats-assert/load'

installer="${PWD}/../bin/vic-machine"
docker="$(which docker)"
log="/sbin"

setup () {
    assert [ -x "${installer}" ]
    assert [ -x "$docker" ]

    export GOVC_URL=$VIC_ESX_TEST_URL
    export GOVC_INSECURE=true
    export GOVC_PERSIST_SESSION=true

    datacenter=`govc ls / | awk -F/ '{print $2}'`
    host=`govc ls /${datacenter}/host/`
    host=`echo ${host} | awk '{print $1}'`
    POOL_PATH=`govc ls $host/Resources`
    datastore=`govc ls /$datacenter/datastore | awk -F/ '{print $4}'`
    IMAGE_STORE_NAME=`echo $datastore | awk '{print $1}'`
    USER=$(awk -F: '{print $1}' <<<"$GOVC_URL")
    passurl=$(awk -F: '{print $2}' <<<"$GOVC_URL")
    PASSWORD=$(awk -F@ '{print $1}' <<<"$passurl")
    TEST_URL=$(awk -F@ '{print $2}' <<<"$passurl")
    vch_name="VCH-${RANDOM}"
    about=`govc about`
    if [[ "$about" == *"ESX"* ]]; then
      ISVC=0
    else
      ISVC=1
    fi
}

teardown() {
    echo "Clear VCH deployments if exists"

    govc vm.destroy "${vch_name}" || echo "no need to delete VCH"
    govc pool.destroy "${POOL_PATH}/${vch_name}" || echo "no need to delete resource pool"
    govc datastore.rm -ds ${IMAGE_STORE_NAME} "${vch_name}" || echo "no need to delete datastore files"
    govc datastore.rm -ds ${IMAGE_STORE_NAME} "VIC" || echo "no need to delete datastore file VIC"
    govc host.vswitch.remove "${vch_name}" || echo "no need to remove vswitch"
    rm *.pem || echo "no pem file found"
    rm *.log || echo "no log file found"
}

@test "vic-machine usage" {
    run "${installer}"
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-target argument must be specified" ]]
}

@test "vic-machine user is missing" {
    run "${installer}" -target ${TEST_URL}
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-user User to login target must be specified" ]]
}

@test "vic-machine compute-resource is missing" {
    run "${installer}" -target ${TEST_URL} -user root
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-compute-resource Compute resource path must be specified" ]]
}

@test "vic-machine image-store is missing" {
    run "${installer}" -target ${TEST_URL} -user root -compute-resource $POOL_PATH
    [ $status -ne 0 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-image-store Image datastore name must be specified" ]]
}

@test "vic-machine deployment and docker commands" {
    params="-target $TEST_URL -name ${vch_name} -user ${USER} -compute-resource $POOL_PATH -image-store $IMAGE_STORE_NAME"
    if [ $ISVC == 1 ]; then
	  assert [ -n $BRIDGE_NETWORK ]
	  params="$params -bridge-network $BRIDGE_NETWORK"
	fi
    echo "${installer}" $params
    params="$params -passwd ${PASSWORD}"
    run "${installer}" $params
    len=${#lines[@]}
    echo ${lines[len-2]}
    echo ${lines[len-1]}
    [[ ${lines[len-1]} =~ "Installer completed successfully..." ]]
    docker_cmd=`echo ${lines[len-2]} | awk '{$1=$2=""; $(NF--)=""; print$0}'`
    echo "docker_cmd="$docker_cmd
    # execute sample docker command
    run "$docker" $docker_cmd info
    echo ${lines[0]}

    #reinstall with same name
    params="-target $TEST_URL -name ${vch_name} -user ${USER} -compute-resource $POOL_PATH -image-store $IMAGE_STORE_NAME -force -key ./${vch_name}-key.pem -cert ./${vch_name}-cert.pem"
    if [ $ISVC == 1 ]; then
	  params="$params -bridge-network $BRIDGE_NETWORK"
	fi
    echo "${installer}" $params
    params="$params -passwd ${PASSWORD}"
    run "${installer}" $params
    len=${#lines[@]}
    echo ${lines[len-1]}
    [[ ${lines[len-1]} =~ "Installer completed successfully..." ]]
    docker_cmd=`echo ${lines[len-2]} | awk '{$1=$2=""; $(NF--)=""; print$0}'`
    # execute docker pull command
    run "$docker" $docker_cmd info
    [ $status -eq 0 ]
    sleep 5
    run "$docker" $docker_cmd pull busybox
    len=${#lines[@]}
    echo ${lines[len-2]}
    echo ${lines[len-1]}
    # execute docker create/start command
    run "$docker" $docker_cmd create busybox
    len=${#lines[@]}
    echo ${lines[len-1]}
    name=${lines[0]}
    run "$docker" $docker_cmd start ${name}
    len=${#lines[@]}
    echo ${lines[len-1]}
    govc vm.destroy ${name}
}
