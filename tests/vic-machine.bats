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

installer="/sbin/vic-machine"
docker="/usr/bin/docker"
log="/sbin"

setup () {
    assert [ -x "${installer}" ]
    assert [ -x "$docker" ]

    export GOVC_URL=$VIC_ESX_TEST_URL
    export GOVC_INSECURE=true
    export GOVC_PERSIST_SESSION=true

    datacenter=`govc ls -k / | awk -F/ '{print $2}'`
    echo $datacenter
    host=`govc ls -k /${datacenter}/host/`
    host=`echo ${host} | awk '{print $1}'`
    echo $host
    POOL_PATH=`govc ls $host/Resources`
    echo $POOL_PATH
    datastore=`govc ls /$datacenter/datastore | awk -F/ '{print $4}'`
    IMAGE_STORE_NAME=`echo $datastore | awk '{print $1}'`
    echo $IMAGE_STORE_NAME
    USER=$(awk -F: '{print $1}' <<<"$GOVC_URL")
    echo $USER
    passurl=$(awk -F: '{print $2}' <<<"$GOVC_URL")
    echo $passurl
    PASSWORD=$(awk -F@ '{print $1}' <<<"$passurl")
    echo $PASSWORD
    TEST_URL=$(awk -F@ '{print $2}' <<<"$passurl")
    echo $TEST_URL
}

clearVCH() {
    echo "Clear VCH deployments if exists"
    
    govc vm.destroy "test-vch"
    govc pool.destroy "${POOL_PATH}/test-vch"
    govc datastore.rm -ds ${IMAGE_STORE_NAME} "test-vch"
    govc datastore.rm -ds ${IMAGE_STORE_NAME} "VIC"
}

@test "vic-machine usage" {
    run "${installer}"
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "Target argument must be specified" ]]
}

@test "vic-machine user is missing" {
    run "${installer}" -target ${TEST_URL}
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-user User to login target must be specified" ]]
}

@test "vic-machine pool is missing" {
    run "${installer}" -target ${TEST_URL} -user root
    [ $status -eq 1 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-pool Compute resource path must be specified" ]]
}

@test "vic-machine iStore is missing" {
    run "${installer}" -target ${TEST_URL} -user root -pool $POOL_PATH
    [ $status -ne 0 ]
    echo ${lines[0]}
    [[ ${lines[0]} =~ "-iStore Image datastore name must be specified" ]]
}

@test "vic-machine generate certificates" {
    echo "${installer}" -target $TEST_URL -name test-vch -user ${USER} -passwd ${PASSWORD} -pool $POOL_PATH -iStore $IMAGE_STORE_NAME -gen
    run "${installer}" -target $TEST_URL -name test-vch -user ${USER} -passwd ${PASSWORD} -pool $POOL_PATH -iStore $IMAGE_STORE_NAME -gen
	# -bdgNet="hasan-test"
    len=${#lines[@]}
    echo ${lines[len-1]}
    [[ ${lines[len-1]} =~ "Installer completed successfully..." ]]
    docker_cmd=`echo ${lines[len-2]} | awk '{$1=$2=""; $(NF--)=""; print$0}'`
    echo "docker_cmd="$docker_cmd
# execute sample docker command
    run "$docker" $docker_cmd info
    echo ${lines[0]}
}

@test "vic-machine run with same name" {
    echo "${installer}" -target ${TEST_URL} -name test-vch -user ${USER} -passwd ${PASSWORD} -pool $POOL_PATH -iStore $IMAGE_STORE_NAME -force -key ./test-vch-key.pem -cert ./test-vch-cert.pem
    run "${installer}" -target ${TEST_URL} -name test-vch -user ${USER} -passwd ${PASSWORD} -pool $POOL_PATH -iStore $IMAGE_STORE_NAME -force -key ./test-vch-key.pem -cert ./test-vch-cert.pem
	# -bdgNet="hasan-test"
    len=${#lines[@]}
    echo ${lines[len-1]}
    [[ ${lines[len-1]} =~ "Installer completed successfully..." ]]
    docker_cmd=`echo ${lines[len-2]} | awk '{$1=$2=""; $(NF--)=""; print$0}'`
# execute docker create/start command
    run "$docker" $docker_cmd info
    [ $status -eq 0 ]
    sleep 5
    run "$docker" $docker_cmd pull busybox
    len=${#lines[@]}
    echo ${lines[len-2]}
    echo ${lines[len-1]}
    run "$docker" $docker_cmd create busybox
    len=${#lines[@]}
    echo ${lines[len-1]}
    run "$docker" $docker_cmd start ${lines[0]}
    len=${#lines[@]}
    echo ${lines[len-1]}
}

@test "clear env" {
	clearVCH
}