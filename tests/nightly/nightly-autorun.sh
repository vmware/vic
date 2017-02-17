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

# Get the latest code from vmware/vic repo
cd ~/go/src/github.com/vmware/vic
git checkout master
git remote update
git rebase upstream/master

# Kick off the nightly
now=$(date +"%m_%d_%Y")
sudo ./tests/nightly/nightly-kickoff.sh > nightly_$now.txt
