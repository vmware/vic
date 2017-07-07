#!/bin/bash -e
# Copyright 2017 VMware, Inc. All Rights Reserved.
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
#

show_error () {
    echo "Provide the path to the H5 Client SDK. e.g. $0 -s /path/to/sdk" >&2
    exit 1
}

while getopts :s:b: flag ; do
    case $flag in
        s)
            VSPHERE_H5C_SDK_HOME=$OPTARG
            ;;
        b)
            BUILD_MODE=$OPTARG
            ;;
        \?)
            show_error
            ;;
        *)
            show_error
            ;;
    esac
done

if [[ -z "$VSPHERE_H5C_SDK_HOME" ]] ; then
    show_error
fi

yarn install

if [[ "$BUILD_MODE" = "prod" ]] ; then
    echo "Building in production mode"
    rm -rf ../main/webapp/resources/build-dev 2>/dev/null
    npm run build:prod
else
    echo "Building in development mode"
    rm -rf ../main/webapp/resources/dist 2>/dev/null
    npm run build:dev
fi
