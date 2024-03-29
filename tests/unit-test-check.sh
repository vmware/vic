#!/bin/bash
# Copyright 2018 VMware, Inc. All Rights Reserved.
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
set -e -o pipefail +h && [ -n "$DEBUG" ] && set -x

prNumber=$(drone build info --format "{{ .Ref }}" vmware/vic "${DRONE_BUILD_NUMBER}" | cut -f 3 -d'/')
prBody=$(curl -H "Authorization: token $GITHUB_AUTOMATION_API_KEY" "https://api.github.com/repos/vmware/vic/pulls/${prNumber}" | jq -r ".body")

if ! (echo "${prBody}" | grep -q "\[skip unit\]"); then
  make test
  codecov --token "${CODECOV_TOKEN}" --file .cover/cover.out
fi
