#!/usr/bin/env bats
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

load helpers/helpers

setup () {
    installer="$VIC_DIR/bin/vic-machine"
    docker="$(which docker)"

    assert [ -n "$VIC_ESX_TEST_URL" ]
    assert [ -x "$installer" ]
    assert [ -x "$docker" ]

    export GOVC_URL=$VIC_ESX_TEST_URL GOVC_INSECURE=true
    export GOVC_USERNAME GOVC_PASSWORD GOVC_RESOURCE_POOL GOVC_DATASTORE

    if [ -z "$GOVC_RESOURCE_POOL" ] ; then
        GOVC_RESOURCE_POOL=$(govc ls /*/host/*/Resources | head -1)
    fi

    if [ -z "$GOVC_DATASTORE" ] ; then
        GOVC_DATASTORE=$(govc ls /*/datastore | head -1 | xargs basename)
    fi

    # split GOVC_URL into GOVC_{URL,USERNAME,PASSWORD}
    eval "$(govc env | sed -e 's/\$/\\$/g')"

    vch_name="VCH-${RANDOM}"
}

teardown() {
    if [ -n "$(govc vm.info $vch_name)" ] ; then
        govc vm.destroy $vch_name
        govc pool.destroy "$GOVC_RESOURCE_POOL/$vch_name"
        govc datastore.rm -f $vch_name
        govc datastore.rm -f VIC
        if govc host.vswitch.info | grep Name: | grep -q $vch_name ; then
            govc host.vswitch.remove $vch_name
        fi
    fi

    rm -f ./$vch_name-*.pem
}

@test "vic-machine usage" {
    run "$installer"
    assert_failure
    assert_line -e "-target argument must be specified"
}

@test "vic-machine user is missing" {
    run "$installer" -target "$GOVC_URL"
    assert_failure
    assert_line -e "-user User to login target must be specified"
}

@test "vic-machine compute-resource is missing" {
    run "$installer" -target "$GOVC_URL" -user root
    assert_failure
    assert_line -e "-compute-resource Compute resource path must be specified"
}

@test "vic-machine image-store is missing" {
    run "$installer" -target "$GOVC_URL" -user root -compute-resource "$GOVC_RESOURCE_POOL"
    assert_failure
    assert_line -e "-image-store Image datastore name must be specified"
}

@test "vic-machine deployment and docker commands" {
    params=(-name $vch_name -target $GOVC_URL -user $GOVC_USERNAME -passwd $GOVC_PASSWORD
            -compute-resource $GOVC_RESOURCE_POOL -image-store $GOVC_DATASTORE)

    if [ "$(govc about -json | jq -r .About.ApiType)" = "VirtualCenter" ] ; then
        assert [ -n "$BRIDGE_NETWORK" ]
        params+=(-bridge-network $BRIDGE_NETWORK)
    fi

    local output status
    run "$installer" "${params[@]}"
    assert_success
    assert_line -e "Installer completed successfully"
    docker_cmd=($(grep 'docker -H' <<<"$output" | cut -d' ' -f3-7))

    # execute docker info command
    for _ in $(seq 1 5) ; do
      run "$docker" "${docker_cmd[@]}" info

      if [ "$status" -eq 0 ] ; then
          assert_line "Name: VIC"
          break
      fi

      # retry as listener may not be up yet
      sleep 1
    done

    assert_success

    # execute docker pull command
    run "$docker" "${docker_cmd[@]}" pull busybox
    assert_success
    assert_line -e "Status: .* for library/busybox:latest"

    # execute docker create/start command
    run "$docker" "${docker_cmd[@]}" create busybox
    assert_success
    name="$output"

    run "$docker" "${docker_cmd[@]}" start "$name"
    assert_success

    run govc vm.destroy "$name"
    assert_success

    run govc datastore.rm -f "$name"
    assert_success

    # test reinstall with same name
    params+=(-force -key ./${vch_name}-key.pem -cert ./${vch_name}-cert.pem)

    run "$installer" "${params[@]}"
    assert_success
    assert_line -e "Installer completed successfully"
}
