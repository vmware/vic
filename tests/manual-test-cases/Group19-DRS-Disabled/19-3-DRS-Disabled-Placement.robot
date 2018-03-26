# Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
Documentation  Test 19-3 - DRS-Disabled-Placement
Resource  ../../resources/Util.robot
Suite Setup  Wait Until Keyword Succeeds  10x  10m  ROBO SKU Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
Deploy Stress Container And Return Host IP
    [Arguments]  ${name}
    Log To Console  Deploying stress container ${name}...
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --name ${name} progrium/stress --vm 5 --vm-bytes 256M --vm-hang 0
    Should Be Equal As Integers  ${rc}  0

    ${vm_name}=  Get VM display name  ${id}
    ${host_ip}=  Get VM Host Name  ${vm_name}
    [Return]  ${host_ip}

Create Busybox Container VM And Return ID
    [Arguments]
    Log To Console  Deploying busybox container VM..
    Pull Image  busybox
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${id}

Relocate VM To Host
    [Arguments]  ${host}  ${vm_name}
    Log To Console  Relocating VM ${vm_name} to host ${host}...
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.migrate -host ${host} ${vm_name}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  OK

*** Test Cases ***

# TODO(jzt): we need to test against a single ESX host

Simple Placement
    Install VIC Appliance To Test Server
    ${vch_host}=  Get VM Host Name  %{VCH-NAME}

    Pull Image  progrium/stress

    ${stressed_hosts}=  Create List  ${vch_host}
    :FOR  ${i}  IN RANGE  2
    \   ${ip}=  Deploy Stress Container And Return Host IP  stresser${i}
    \   Should Not Contain  ${stressed_hosts}  ${ip}
    \   Append To List  ${stressed_hosts}  ${ip}

    # 1 VCH host + 2 stressed hosts out of 4 total, leaving one
    # clean host to which a new container should relocate
    ${len}=  Get Length  ${stressed_hosts}
    Should Be Equal As Integers  ${len}  3

    ${id}=  Create Busybox Container VM And Return ID
    ${vm_name}=  Get VM display name  ${id}

    # Move it onto the same host as the VC
    Relocate VM To Host  ${vch_host}  ${vm_name}

    ${host_ip}=  Get VM Host Name  ${vm_name}
    Should Be Equal  ${host_ip}  ${vch_host}

    # power on the busybox container - it should relocate to a non-stressed host
    ${rc}  ${cid}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${id}
    Should Be Equal As Integers  ${rc}  0
    ${host_ip}=  Get VM Host Name  ${vm_name}

    Should Not Contain  ${stressed_hosts}  ${host_ip}
    Cleanup VIC Appliance On Test Server
