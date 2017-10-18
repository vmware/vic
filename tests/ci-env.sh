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

if [ ! -f /ci/test_env ]; then
    echo "'/ci/test_env' file not found! Make sure its available on host and mounted.";
    exit 1;
fi

# update file permissions
chmod 777 /ci/test_env

# source env variables from volume mounted test_env file
. /ci/test_env

# exit if TEST_URL is not set
if [ -z "${TEST_URL}" ]; then
    echo "ERROR: TEST_URL is not set";
    exit 1;
else
    echo "TEST_URL is set to '$TEST_URL'";
fi

# check if DRONE_MACHINE is set
# TODO: use DRONE_MACHINE from 0.8 drone version
if [ -z "${DRONE_MACHINE}" ]; then
    echo "WARN: DRONE_MACHINE is not set";
else
    echo "DRONE_MACHINE is set to '$DRONE_MACHINE'";
fi

# exit if ACTIVE_BUILD_NUMBER is not set
if [ -z "${ACTIVE_BUILD_NUMBER}" ]; then
    echo "ERROR: ACTIVE_BUILD_NUMBER is not set";
    exit 1;
else
    echo "ACTIVE_BUILD_NUMBER is set to '$ACTIVE_BUILD_NUMBER'";
fi
