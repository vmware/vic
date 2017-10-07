# Copyright 2017 VMware, Inc. All Rights Reserved.
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
Documentation     Test 23-02 - VCH List
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Setup
Suite Teardown    Teardown
Default Tags


*** Keywords ***
Setup
    Start VIC Machine Server
    Install VIC Appliance To Test Server


Teardown
    Cleanup VIC Appliance On Test Server
    Terminate All Processes    kill=True


Get VCH List
    Get Path Under Target    vch


Get VCH List Within Datacenter
    ${dcID}=    Get Datacenter ID
    Get Path Under Target    datacenter/${dcID}/vch


Verify VCH List
    ${expectedId}=    Get VCH ID    %{VCH-NAME}
    ${actualId}=  Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="%{VCH-NAME}").id'
    Should Be Equal    ${expectedId}    ${actualId}

    ${admin_portal}=  Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="%{VCH-NAME}").admin_portal'
    Should Not Be Empty    ${admin_portal}

    ${docker_host}=  Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="%{VCH-NAME}").docker_host'
    Should Not Be Empty    ${docker_host}

    ${upgrade_status}=  Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="%{VCH-NAME}").upgrade_status'
    Should Not Be Empty    ${upgrade_status}

    ${version}=  Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="%{VCH-NAME}").version'
    Should Not Be Empty    ${version}


*** Test Cases ***
Get VCH List
    Get VCH List

    Verify Return Code
    Verify Status Ok
    Verify VCH List


Get VCH List Within Datacenter
    Get VCH List Within Datacenter

    Verify Return Code
    Verify Status Ok
    Verify VCH List

# TODO: Add test for compute resource (once relevant code is updated to use ID instead of name)
# TODO: Add test for compute resource within datacenter (once relevant code is updated to use ID instead of name)

Get VCH List Within Invalid Datacenter
    Get Path Under Target    datacenter/INVALID/vch

    Verify Return Code
    Verify Status    404


Get VCH List Within Invalid Compute Resource
    Get Path Under Target    vch    compute-resource=INVALID

    Verify Return Code
    Verify Status    400


Get VCH List Within Invalid Datacenter and Compute Resource
    Get Path Under Target    datacenter/INVALID/vch    compute-resource=INVALID

    Verify Return Code
    Verify Status    404
