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
#

# runs the given command in a container with iso build dependencies
set -e && [ -n "$DEBUG" ] && set -x

function main {
  PKGMGR=$(cat isos/base/repos/$1/repo-spec.json | jq -r '.packagemanager')
  shift

  rmArg=""
  if [ -z "$DEBUG" ]; then
    rmArg="--rm"
  fi

  docker run \
  -it \
  ${rmArg} \
  -v $GOPATH/bin:/go/bin:ro \
  -v $GOPATH/src/github.com/vmware/vic:/go/src/github.com/vmware/vic:ro \
  -v $GOPATH/src/github.com/vmware/vic/bin:/go/src/github.com/vmware/vic/bin \
  -e DEBUG=${DEBUG} \
  -e BUILD_NUMBER=${BUILD_NUMBER} \
  gcr.io/eminent-nation-87317/vic-build-image:${PKGMGR:-tdnf} "$*"
}

REPO="photon-2.0"
# Find the dependency manager. The d stands for distro.
while getopts ':d:' flag; do
  case "${flag}" in
    d) REPO="${OPTARG}" ;;
  esac
done
shift $((OPTIND-1))

# Check if jq is available - we need this on either path
which jq >/dev/null 2>&1
[ $? -ne 0 ] && "Echo please install 'jq' to continue..." && exit 1

echo "building $REPO"
# Check if docker installed
if ! docker info >/dev/null 2>&1; then
   /bin/bash -c "$*" # prevent docker in docker
else
   main "${REPO}" "$@"
fi