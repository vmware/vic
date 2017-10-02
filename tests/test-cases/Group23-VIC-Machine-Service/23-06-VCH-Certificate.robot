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
Documentation  Test 23-06 - VCH Certificate
Resource  ../../resources/Util.robot
Suite Setup  Start VIC Machine Server
Suite Teardown  Terminate All Processes  kill=True
Test Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Install VIC Machine Without TLS
    [Arguments]  ${vic-machine}=bin/vic-machine-linux  ${appliance-iso}=bin/appliance.iso  ${bootstrap-iso}=bin/bootstrap.iso  ${certs}=${true}  ${vol}=default  ${cleanup}=${true}  ${debug}=1  ${additional-args}=${EMPTY}
    Set Test Environment Variables
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.esxcli network firewall set -e false
    # Attempt to cleanup old/canceled tests
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Networks On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling vSwitches On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Containers On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Resource Pools On Test Server
    Set Suite Variable  ${vicmachinetls}  --no-tls
    Log To Console  \nInstalling VCH to test server with tls disabled...
    ${output}=  Run VIC Machine Command  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}  ${debug}  ${additional-args}
    Log  ${output}
    Should Contain  ${output}  Installer completed successfully

    Get Docker Params  ${output}  ${certs}
    Log To Console  Installer completed successfully: %{VCH-NAME}...

    [Return]  ${output}

Start VIC Machine Server
    Start Process  ./bin/vic-machine-server --port 31337 --scheme http  shell=True  cwd=/go/src/github.com/vmware/vic
    Sleep  1s  for service to start

Curl No Datacenter
    [Arguments]  ${vch-id}  ${auth}
    ${rc}  ${output}=  Run And Return Rc And Output  curl -s -w "\%{http_code}\n" -X GET "http://127.0.0.1:31337/container/target/%{TEST_URL}/vch/${vch-id}/certificate?thumbprint=%{TEST_THUMBPRINT}" -H "authorization: Basic ${auth}"
    [Return]  ${rc}  ${output}

Curl Datacenter
    [Arguments]  ${vch-id}  ${auth}
    ${orig}=  Get Environment Variable  TEST_DATACENTER
    ${dc}=  Run Keyword If  '%{TEST_DATACENTER}' == '${SPACE}'  Get Datacenter Name
    Run Keyword If  '%{TEST_DATACENTER}' == '${SPACE}'  Set Environment Variable  TEST_DATACENTER  ${dc}
    ${dcID}=  Get Datacenter ID
    ${rc}  ${output}=  Run And Return Rc And Output  curl -s -w "\%{http_code}\n" -X GET "http://127.0.0.1:31337/container/target/%{TEST_URL}/datacenter/${dcID}/vch/${vch-id}/certificate?thumbprint=%{TEST_THUMBPRINT}" -H "authorization: Basic ${auth}"
    Set Environment Variable  TEST_DATACENTER  ${orig}
    [Return]  ${rc}  ${output}

*** Test Cases ***
Get VCH Certificate
    Install VIC Appliance To Test Server
    ${id}=  Get VCH ID  %{VCH-NAME}
    ${auth}=  Evaluate  base64.b64encode("%{TEST_USERNAME}:%{TEST_PASSWORD}")  modules=base64
    ${rc}  ${output}=  Curl No Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  BEGIN CERTIFICATE
    Should Contain  ${output}  END CERTIFICATE
    Should Not Contain  ${output}  "
    ${status}=  Get Line  ${output}  -1
    Should Be Equal As Integers  200  ${status}
    ${rc}  ${outputDC}=  Curl Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  ${outputDC}

Get VCH Certificate No TLS
    Install VIC Machine Without TLS
    ${id}=  Get VCH ID  %{VCH-NAME}
    ${auth}=  Evaluate  base64.b64encode("%{TEST_USERNAME}:%{TEST_PASSWORD}")  modules=base64
    ${rc}  ${output}=  Curl No Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  No certificate found for VCH
    Should Not Contain  ${output}  BEGIN CERTIFICATE
    Should Not Contain  ${output}  END CERTIFICATE
    ${status}=  Get Line  ${output}  -1
    Should Be Equal As Integers  404  ${status}
    ${rc}  ${outputDC}=  Curl Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  ${outputDC}

