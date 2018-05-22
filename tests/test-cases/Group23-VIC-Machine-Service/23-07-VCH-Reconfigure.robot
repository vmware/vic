# Copyright 2018 VMware, Inc. All Rights Reserved.
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
Documentation     Test 23-07 - VCH Reconfigure
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Setup
Suite Teardown    Teardown
Default Tags

*** Keywords ***
Setup
    Start VIC Machine Server
    Install VIC Appliance To Test Server  debug=0

    ${id}=  Get VCH ID  %{VCH-NAME}
    ${dc-id}=  Get Datacenter ID

    Set Suite Variable  ${VCH-ID}  ${id}
    Set Suite Variable  ${DC-ID}  ${dc-id}

Teardown
    Cleanup VIC Appliance On Test Server
    Cleanup Test Server
    Terminate All Processes    kill=True

Inspect VCH
    Get Path Under Target  vch/${VCH-ID}
    Verify Return Code
    Verify Status Ok

Put VCH
    [Arguments]    ${data}
    Put Path Under Target    vch/${VCH-ID}    ${data}

Put VCH Within Datacenter
    [Arguments]    ${data}
    Put Path Under Target  datacenter/${DC-ID}/vch/${VCH-ID}   ${data}

Patch VCH
    [Arguments]    ${data}
    Patch Path Under Target  vch/${VCH-ID}    ${data}

Patch VCH Within Datacenter
    [Arguments]    ${data}
    Patch Path Under Target  datacenter/${DC-ID}/vch/${VCH-ID}   ${data}


*** Test Cases ***
Put VCH Debug Level
    Inspect VCH
    Property Should Be Equal  .debug  null

    Put VCH  '{"name":"%{VCH-NAME}", "debug": 3}'
    Verify Return Code
    Verify Status Accepted

    Inspect VCH
    Property Should Be Equal  .debug   3

Put Fails For Nonconsistent Immutable Fields
    Put VCH  '{"name": "IMMUTABLE"}'
    Verify Return Code
    Verify Status Conflict

    Inspect VCH
    Property Should Be Equal  .name   %{VCH-NAME}

Put VCH Debug Level Within Datacenter
    Inspect VCH
    Property Should Be Equal  .debug  3

    Put VCH Within Datacenter  '{"name": "%{VCH-NAME}", "debug": 0}'
    Verify Return Code
    Verify Status Accepted

    Inspect VCH
    Property Should Be Equal  .debug  null

Patch VCH Debug Level
    Inspect VCH
    Property Should Be Equal  .debug  null

    Patch VCH  '{"debug": 3}'
    Verify Return Code
    Verify Status Accepted

    Inspect VCH
    Property Should Be Equal  .debug  3

Patch Fails If Immutable Fields Present
    Patch VCH  '{"name": "%{VCH-NAME}", "debug": 0}'
    Verify Return Code
    Verify Status Bad Request

Patch VCH Debug Level Within Datacenter
    Inspect VCH
    Property Should Be Equal  .debug  3

    Patch VCH Within Datacenter  '{"debug": 0}'
    Verify Return Code
    Verify Status Accepted

    Inspect VCH
    Property Should Be Equal  .debug  null
