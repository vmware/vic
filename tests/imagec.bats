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

imagec="$GOPATH/src/github.com/vmware/vic/$BIN/imagec"
portlayer="$GOPATH/src/github.com/vmware/vic/$BIN/port-layer-server"
IMAGES_DIR="images"
DEFAULT_IMAGE="https/registry-1.docker.io/v2/library/photon/latest"
ALT_IMAGE="https/registry-1.docker.io/v2/tatsushid/tinycore"

setup () {
    assert [ -x "$imagec" ]
    assert [ -x "$portlayer" ]

    # create temp directories & start the port layer server
    mkdir -p /tmp/imagec_test_dir
    mkdir -p /tmp/portlayer
    cd /tmp/imagec_test_dir
#    start_port_layer
#    assert [ ! -z "$port_layer_pid" ]
}

teardown() {
    # stop the port layer between tests
#    kill_port_layer
    cd >/dev/null
    # nuke everything between tests
    rm -rf /tmp/imagec_test_dir
    rm -rf /tmp/portlayer
}

@test "validate helper functions" {
    skip "Skipping until the portlayer can be brought up without a valid SDK target"
    assert [ ! -z "$port_layer_pid" ]
    assert [ kill_port_layer ]
}

@test "imagec run without arguments downloads photon into default destination from Docker Hub" {
    run "$imagec" -standalone
    assert_success
    assert [ -d "$IMAGES_DIR/$DEFAULT_IMAGE" ] # check that the correct dir is created
    assert [ -e imagec.log ] # check the existence of the default logfile
    assert [ -n imagec.log ] # logfile shouldn't be empty either
    assert verify_checksums "$IMAGES_DIR/$DEFAULT_IMAGE" # make sure we got the right image
}

@test "imagec -help should show usage information" {
    run $imagec --help
    assert_failure # default 'you called -help' error code
    [[ ${lines[0]} =~ Usage.*imagec.* ]]
}

@test "imagec -debug should enable debugging and -stdout outputs logs to stdout" {
    # should get at least one line with 'level=output' present
    run [ $("$imagec" -standalone -stdout -debug | grep "level=debug" | wc -l) -ge 1 ]
    assert_success
    assert verify_checksums "$IMAGES_DIR/$DEFAULT_IMAGE"
}

@test "imagec -destination allows us to change where the image is saved" {
    run "$imagec" -standalone -destination foo
    assert_success
    assert verify_checksums "foo/$DEFAULT_IMAGE"
}

@test "imagec -digest should specify tag name or image digest to download" {
    run "$imagec" -standalone -digest 7.0-x86_64 -image tatsushid/tinycore
    assert_success
    assert verify_checksums "$IMAGES_DIR/$ALT_IMAGE"/7.0-x86_64
}

#@test "imagec -host should allow us to specify which host runs the portlayer API" {
#    assert kill_port_layer # it was started on 8080 by setup()
#    assert start_port_layer 1337 # restart it on 1337 cause we're elite
#    run "$imagec" -host localhost:1337
#    assert_success
#    assert verify_checksums "$IMAGES_DIR/$DEFAULT_IMAGE"
#}


@test "imagec -standalone should allow us to run imagec without portlayer API" {
#    assert kill_port_layer # it was started on 8080 by setup()
    run "$imagec" -standalone
    assert_success
    assert verify_checksums "$IMAGES_DIR/$DEFAULT_IMAGE"
}

@test "imagec -image should allow specifying a specific image to download" {
    run "$imagec" -standalone -image tatsushid/tinycore
    assert_success
    assert verify_checksums "$IMAGES_DIR/$ALT_IMAGE/latest"
}

@test "imagec -logfile should change the path of the installer log file (default \"imagec.log\")" {
    run "$imagec" -standalone -logfile foo.log
    assert_success
    assert [ $(wc -l foo.log | awk '{print $1}') -ge 1 ] # logfile shouldn't be empty
}

@test "imagec -registry should allow us to change the registry" {
    skip "test not implemented"
}

@test "imagec -timeout duration should change the HTTP timeout (default 1m0s)" {
    skip "test not implemented"
}

@test "imagec -username and -password should allow us to specify a username and password for a private registry" {
    skip "test not implemented"
}
