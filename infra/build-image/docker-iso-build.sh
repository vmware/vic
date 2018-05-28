#!/bin/bash -e
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
  which jq >/dev/null 2>&1
  [ ! $? -eq 0 ] && "Echo please install 'jq' to continue..." && exit 1

  REPO=$(cat isos/base/repos/$1/repo-spec.json | jq -r '.packagemanager')
  shift

  docker run \
  -it \
  --rm \
  -v $GOPATH/bin:/go/bin:ro \
  -v $GOPATH/src/github.com/vmware/vic:/go/src/github.com/vmware/vic:ro \
  -v $GOPATH/src/github.com/vmware/vic/bin:/go/src/github.com/vmware/vic/bin \
  -w /go/src/github.com/vmware/vic \
  -e TERM=linux \
  -e DEBUG=${DEBUG} \
  gcr.io/eminent-nation-87317/vic-build-image:${REPO:-tdnf} /bin/bash -c "$*"
}

REPO=""
# Find the dependency manager. The d stands for distro.
while getopts ':d:' flag; do
  case "${flag}" in
    d) REPO="${OPTARG}" ;;
  esac
done
shift $((OPTIND-1))

echo "building $REPO"
if [[ -f "/proc/1/cgroup" && -n "$(grep docker /proc/1/cgroup)" ]]; then
  /bin/bash -c "$*" # prevent docker in docker
else
  main "${REPO:-photon-2.0}" "$@"
fi