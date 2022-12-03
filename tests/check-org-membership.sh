#!/bin/bash
# Copyright 2020-2021 VMware, Inc. All Rights Reserved.
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

set -e

echo ${DRONE_COMMIT_AUTHOR}
echo $SKIP_CHECK_MEMBERSHIP

if [ "$SKIP_CHECK_MEMBERSHIP" == "true" ]; then
  echo 'check-org-membership step skipped'
  exit 0
fi

# assuming that the logic here is both a successful curl AND no output.
# testing with a bad auth token shows a 0 return code for curl, but json blob with error
# good auth token is just empty response
result=$(curl --silent -H "Authorization: token $GITHUB_AUTOMATION_API_KEY" "https://api.github.com/orgs/vmware/members/${DRONE_COMMIT_AUTHOR}")
if [ "$?" -eq 0 -o -n "$result" ]; then
  echo "checked origin membership successfully"
else
  echo "failed to check origin membership"
  exit 1
fi
