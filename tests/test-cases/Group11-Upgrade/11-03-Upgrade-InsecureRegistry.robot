# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 11-03 - Upgrade-InsecureRegistry
Suite Setup  Install VIC with version to Test Server  7315  --insecure-registry 10.10.10.10:1234
Suite Teardown  Clean up VIC Appliance And Local Binary
Resource  ../../resources/Util.robot

*** Keywords ***
Get host and path from guestinfo
    ${host}=  Run  govc vm.info -e -json %{VCH-NAME} | jq -c '.VirtualMachines[0].Config.ExtraConfig[] | select(.Key | . and contains("guestinfo.vice./registry/insecure_registries|0/Host")) .Value'
    ${path}=  Run  govc vm.info -e -json %{VCH-NAME} | jq -c '.VirtualMachines[0].Config.ExtraConfig[] | select(.Key | . and contains("guestinfo.vice./registry/insecure_registries|0/Path")) .Value'
    [Return]  ${host}  ${path}

*** Test Cases ***
Upgrade VCH with InsecureRegistry
    ${oldHost}  ${oldPath}=  Get host and path from guestinfo
    Log  ${oldHost} ${oldPath}
    Should Be Equal As Strings  ${oldHost}  "<nil>"
    Should Be Equal As Strings  ${oldPath}  "10.10.10.10:1234"
    Upgrade
    Check Upgraded Version
    ${newHost}  ${newPath}=  Get host and path from guestinfo
    Log  ${newHost} ${newPath}
    Should Be Equal As Strings  ${newPath}  "<nil>"
    Should Be Equal As Strings  ${newHost}  "10.10.10.10:1234"
