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

# Run integration tests

# exit on failure and configure debug, include util functions
set -e && [ -n "$DEBUG" ] && set -x
DIR=$(dirname $(readlink -f "$0"))
. $DIR/../isos/base/utils.sh

function usage() {
  echo "Usage: $0 [-c <test case file>]" 1>&2
  exit 1
}

function install_bats() {
  git clone https://github.com/sstephenson/bats ${DIR}/bats && \
        git clone https://github.com/ztombol/bats-assert ${DIR}/helpers/bats-assert && \
        git clone https://github.com/ztombol/bats-support ${DIR}/helpers/bats-support
  exe="$(which bats)" || true
  if [ -z "$exe" -o ! -x "$exe" ]; then
    ${DIR}/bats/install.sh /usr/local
  fi
}

function install_docker() {
  exe="$(which docker)" || true
  if [ -z "$exe" -o ! -x "$exe" ]; then
     curl -sSL https://get.docker.com/ | sh
  fi
}

clean_bats_code() {
  if [ ! -z ${DIR}/bats ]; then
    rm -rf ${DIR}/bats 
  fi
  if [ ! -z  ${DIR}/helpers/bats-assert ]; then
    rm -rf  ${DIR}/helpers/bats-assert
  fi
  if [ ! -z  ${DIR}/helpers/bats-support ]; then
    rm -rf  ${DIR}/helpers/bats-support
  fi
}

while getopts "c:" flag
do
    case $flag in

        c)
            # Optional. Specify test case file. Default to all files in local directory
            test="$OPTARG"
            ;;

        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" ]; then
    usage
fi

# prep the test system
ensure_apt_packges jq git
install_docker
go get github.com/vmware/govmomi/govc

clean_bats_code
install_bats

pushd ${DIR}
  if [ -z "$test" ]; then
    bats -t . 
  else
    bats -t "$test"
  fi
popd

clean_bats_code
