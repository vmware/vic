# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  This resource contains all keywords related to creating, deleting, maintaining Docker in VIC

*** Keywords ***
Install DinV Into VCH
    [Arguments]  ${port}=3000
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  VCH-PARAMS
    Run Keyword If  '${status}' == 'FAIL'  Fail  VCH needs to be installed before you can install DinV    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull vmware/dinv
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${cid}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p ${port}:2375 vmware/dinv
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  5x  3s  Ensure DinV Daemon Starts  ${cid}
    Set Environment Variable  DINV-IP  %{VCH-IP}
    Set Environment Variable  DINV-URL  %{VCH-IP}:${port}
    Set Environment Variable  DINV-PARAMS  -H %{VCH-IP}:${port}
    [Return]  ${cid}

Ensure DinV Daemon Starts
    [Arguments]  ${cid}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs ${cid}
    Should Contain  ${output}  Daemon has completed initialization
    Should Be Equal As Integers  ${rc}  0
