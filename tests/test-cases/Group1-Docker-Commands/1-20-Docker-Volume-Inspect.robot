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
Documentation  Test 1-20 - Docker Volume Inspect
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Verify Docker Volume Inspect
    [Arguments]  ${volume}
    Log To Console  \nInspecting Docker Volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume inspect ${volume}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Name
    Should Be Equal As Strings  ${id}  ${volume}

*** Test Cases ***
Simple docker volume inspect
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name test
    Should Be Equal As Integers  ${rc}  0
    Verify Docker Volume Inspect  test

Docker volume inspect invalid object
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume inspect fakeVolume
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such volume: fakeVolume

Restart VCH and Docker Volume Inspect Test
    Verify Docker Volume Inspect  test

    #Reboot VM and Verify VCH Info
    Log To Console  Rebooting VCH\n - %{VCH-NAME}
    Reboot VM  %{VCH-NAME}
    Log To Console  Getting VCH IP ...
    ${new_vch_ip}=  Get VM IP  %{VCH-NAME}
    Log To Console  New VCH IP is ${new_vch_ip}
    Replace String  %{VCH-PARAMS}  %{VCH-IP}  ${new_vch_ip}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    Verify Docker Volume Inspect  test