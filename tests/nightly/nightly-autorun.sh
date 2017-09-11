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

# Get the latest code from vic-internal repo for nightly_test_secrets.yml file
cd ~/internal-repo/vic-internal
git clean -fd
git fetch
git pull

# Get the latest code from vmware/vic repo
cd ~/go/src/github.com/vmware/vic
git clean -fd
git fetch https://github.com/vmware/vic master
git pull

# Kick off the nightly
echo "Removing VIC directory if present"
echo "Cleanup logs from previous run"

sudo rm -rf *.zip *.log
sudo rm -rf bin 60 65

sudo -E ./tests/nightly/nightly-kickoff.sh > ./nightly_console.txt 2>&1
