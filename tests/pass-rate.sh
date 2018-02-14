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

curl -s https://api.github.com/repos/vmware/vic/commits?access_token=$GITHUB_AUTOMATION_API_KEY | jq -r ".[].sha" | xargs -n1 -I{} curl -s https://api.github.com/repos/vmware/vic/statuses/{}?access_token=$GITHUB_AUTOMATION_API_KEY | jq -r ".[-0].state" > status.out

failures=$(cat status.out | grep -c failure)
echo "Number of failed merges to master in the last 30 merges: $failures"
successes=$(cat status.out | grep -c success)
echo "Number of successful merges to master in the last 30 merges: $successes"

let total=$successes+$failures
passrate=$(bc -l <<< "scale=2;100 * ($successes / $total)")

echo "Current CI passrate: $passrate"
curl --max-time 10 --retry 3 -s -d "payload={'channel': '#vic-bots', 'text': 'Current CI passrate: $passrate%'}" "$SLACK_URL"
