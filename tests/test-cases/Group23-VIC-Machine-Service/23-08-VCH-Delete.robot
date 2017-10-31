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
Documentation     Test 23-03 - VCH Create
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Start VIC Machine Server
Suite Teardown    Terminate All Processes    kill=True
Test Setup        Install VIC Appliance To Test Server
Default Tags


*** Keywords ***
Get VCH ID ${name}
    Get Path Under Target    vch
    ${id}=    Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="${name}").id'
    [Return]    ${id}


Test Delete
    [Arguments]    ${path}

    Get Path Under Target       ${path}

    Verify Return Code
    Verify Status Ok

    Delete Path Under Target    ${path}

    Verify Return Code
    Verify Status Accepted

    Get Path Under Target       ${path}

    Verify Return Code
    Verify Status Not Found


*** Test Cases ***
Delete VCH
    ${id}=    Get VCH ID %{VCH-NAME}

    Test Delete    vch/${id}


Delete VCH Within Datacenter
    ${dc}=    Get Datacenter ID
    ${id}=    Get VCH ID %{VCH-NAME}

    Test Delete    datacenter/${dc}/vch/${id}


Delete Invalid VCH
    ${id}=    Get VCH ID %{VCH-NAME}

    Delete Path Under Target    vch/INVALID

    Verify Return Code
    Verify Status Not Found

    Get Path Under Target       vch/${id}

    Verify Return Code
    Verify Status Ok

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete VCH in Invalid Datacenter
    ${id}=    Get VCH ID %{VCH-NAME}

    Delete Path Under Target    datacenter/INVALID/vch/${id}

    Verify Return Code
    Verify Status Not Found

    Get Path Under Target       vch/${id}

    Verify Return Code
    Verify Status Ok

    [Teardown]    Cleanup VIC Appliance On Test Server
