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
# limitations under the License.

set -e

VICPATH="$(git rev-parse --show-toplevel)"
TEST_SRC="${VICPATH}/vendor/github.com/docker/docker/integration-cli"
TEST_DEST="${VICPATH}/tests/docker-tests/integration-tests"
RUN_PATH="$(pwd)"

if [ "$#" -gt "0" ]; then
    if [ $1 = "clean" ];
    then
        rm -rf ${TEST_DEST}
        exit 0
    fi
fi

if [ -z ${DOCKER_HOST+x} ];
then
    echo "DOCKER_HOST not set.  Please set this and run again."
    echo "  examples:"
    echo "   from <VIC>/docker-tests: $> DOCKER_HOST=tcp://<host>:<port> ./run-tests.sh"
    echo "   from <VIC>: $> DOCKER_HOST=tcp://<host>:<port> make docker-integration-tests"
else
    echo
    echo Setting up pre-requisites
    docker pull busybox:latest

    if [ ! -d "${TEST_DEST}" ];
    then
        mkdir ${TEST_DEST}
    fi

    echo Copying Docker integration tests utils from VIC\'s vendored sources
    cp ${TEST_SRC}/check_test.go ${TEST_DEST}
    cp ${TEST_SRC}/daemon.go ${TEST_DEST}
    cp ${TEST_SRC}/docker_cli_cp_utils.go ${TEST_DEST}
    cp ${TEST_SRC}/docker_test_vars.go ${TEST_DEST}
    cp ${TEST_SRC}/docker_utils.go ${TEST_DEST}
    cp ${TEST_SRC}/npipe.go ${TEST_DEST}
    cp ${TEST_SRC}/registry.go ${TEST_DEST}
    cp ${TEST_SRC}/requirements.go ${TEST_DEST}
    cp ${TEST_SRC}/test_vars_exec.go ${TEST_DEST}
    cp ${TEST_SRC}/test_vars_noexec.go ${TEST_DEST}
    cp ${TEST_SRC}/test_vars_noseccomp.go ${TEST_DEST}
    cp ${TEST_SRC}/test_vars_unix.go ${TEST_DEST}
    cp ${TEST_SRC}/test_vars_windows.go ${TEST_DEST}
    cp ${TEST_SRC}/trust_server.go ${TEST_DEST}
    cp ${TEST_SRC}/utils.go ${TEST_DEST}

    echo Copying appropriate Docker API integration tests from VIC\'s vendored sources
    cp ${TEST_SRC}/docker_api_containers_test.go ${TEST_DEST}
    cp ${TEST_SRC}/docker_api_create_test.go ${TEST_DEST}

    echo Copying appropriate Docker CLI integration tests utils from VIC\'s vendored sources
    cp ${TEST_SRC}/docker_cli_create_test.go ${TEST_DEST}
    cp ${TEST_SRC}/docker_cli_start_test.go ${TEST_DEST}

    echo Running the Docker integration tests
    export DOCKER_REMOTE_DAEMON=$(echo $DOCKER_HOST | grep -oE '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+')

    echo "  DOCKER_HOST = ${DOCKER_HOST}"
    echo "  DOCKER_REMOTE_DAEMON = ${DOCKER_REMOTE_DAEMON}"
    
    # Limit the tests we run until we have complete implementation.
    # Docker integration tests uses gocheck.  -check.f to run selective tests
    # instead of standard -run <regex>.
    cd ${TEST_DEST}
    go test -check.v -test.timeout=360m -check.f "Create|Start"
    cd ${RUN_PATH}
fi
