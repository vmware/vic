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
Documentation  Test 1-45 - Docker Container Network
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  additional-args=--container-network-firewall=%{PUBLIC_NETWORK}:open
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Curl tomcat endpoint
    [Arguments]  ${endpoint}
    ${rc}  ${output}=  Run And Return Rc And Output  curl ${endpoint}
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${output}

*** Test Cases ***
Tomcat with port mapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name tomcat1 -d -p 8082:8080 tomcat:alpine
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl tomcat endpoint  %{VCH-IP}:8082
    Should Contain  ${output}  Apache Tomcat

Tomcat in container-network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name tomcat2 -d --net=public tomcat:alpine
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${ip}=  Get Container IP  %{VCH-PARAMS}  tomcat2  public
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl tomcat endpoint  ${ip}:8080
    Should Contain  ${output}  Apache Tomcat

Tomcat with port mapping in container-network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name tomcat3 -d -p 8083:8080 --net=public tomcat:alpine
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${ip}=  Get Container IP  %{VCH-PARAMS}  tomcat3  public
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl tomcat endpoint  ${ip}:8083
    Should Contain  ${output}  Apache Tomcat
