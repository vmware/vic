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
Documentation  Test 23-05 - VCH Logs
Resource  ../../resources/Util.robot
Suite Setup  Start VIC Machine Server
Suite Teardown  Terminate All Processes  kill=True
Test Setup  Install VIC Appliance To Test Server
Test Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Start VIC Machine Server
    Start Process  ./bin/vic-machine-server --port 31337 --scheme http  shell=True  cwd=/go/src/github.com/vmware/vic
    Sleep  1s  for service to start

Curl No Datacenter
    [Arguments]  ${vch-id}  ${auth}
    ${rc}  ${output}=  Run And Return Rc And Output  curl -s -w "\%{http_code}\n" -X GET "http://127.0.0.1:31337/container/target/%{TEST_URL}/vch/${vch-id}/log?thumbprint=%{TEST_THUMBPRINT}" -H "authorization: Basic ${auth}"
    [Return]  ${rc}  ${output}

Curl Datacenter
    [Arguments]  ${vch-id}  ${auth}
    ${dcID}=  Get Datacenter ID
    ${rc}  ${output}=  Run And Return Rc And Output  curl -s -w "\%{http_code}\n" -X GET "http://127.0.0.1:31337/container/target/%{TEST_URL}/datacenter/${dcID}/vch/${vch-id}/log?thumbprint=%{TEST_THUMBPRINT}" -H "authorization: Basic ${auth}"
    [Return]  ${rc}  ${output}

Delete Log File From VCH Datastore
    ${filename}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls %{VCH-NAME} | grep vic-machine_
    Should Not Be Empty  ${filename}
    ${output}=  Run  govc datastore.rm "%{VCH-NAME}/${filename}"
    ${filename}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls %{VCH-NAME} | grep vic-machine_
    Should Be Empty  ${filename}


*** Test Cases ***
Get VCH Creation Log succeeds after installation completes
    ${id}=  Get VCH ID  %{VCH-NAME}
    ${auth}=  Evaluate  base64.b64encode("%{TEST_USERNAME}:%{TEST_PASSWORD}")  modules=base64
    ${rc}  ${output}=  Curl No Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Installer completed successfully
    ${status}=  Get Line  ${output}  -1
    Should Be Equal As Integers  200  ${status}
    ${rc}  ${outputDC}=  Curl Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  ${outputDC}

Get VCH Creation log errors with 404 after log file is deleted
    ${id}=  Get VCH ID  %{VCH-NAME}
    ${auth}=  Evaluate  base64.b64encode("%{TEST_USERNAME}:%{TEST_PASSWORD}")  modules=base64
    Delete Log File From VCH Datastore
    ${rc}  ${output}=  Curl No Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    ${status}=  Get Line  ${output}  -1
    Should Be Equal As Integers  404  ${status}
    ${rc}  ${outputDC}=  Curl Datacenter  ${id}  ${auth}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  ${outputDC}
