# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 5-7 - NSX
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    Log To Console  TODO
    #${out}=  Deploy Nimbus Testbed  --noSupportBundles --vcvaBuild 3634791 --esxBuild 3620759 --testbedName test-vpx-4esx-virtual-fullInstall-vcva-8gbmem-nsx1m1c --runName VIC-NSX-Test --build nsx-transformers:beta:ob-3586094:master