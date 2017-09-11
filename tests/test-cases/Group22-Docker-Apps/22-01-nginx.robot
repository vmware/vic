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
Documentation  Test 22-01 - nginx
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Get IP Address of Container
    [Arguments]  ${container}
    ${ip}=  Run  docker %{VCH-PARAMS} inspect ${container} | jq -r ".[].NetworkSettings.Networks.bridge.IPAddress"
    [Return]  ${ip}

*** Test Cases ***
Simple background nginx
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name nginx1 -d ${nginx}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${ip}=  Get IP Address of Container  nginx1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run ${busybox} sh -c "wget ${ip} && cat index.html"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Welcome to nginx!

Nginx with port mapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name nginx2 -d -p 8080:80 ${nginx}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl %{VCH-IP}:8080
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Welcome to nginx!

Nginx with a read only content file
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=vol1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v vol1:/mydata ${busybox} echo '<p>HelloWorld</p>' > /mydata/test.html
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name nginx3 -v vol1:/usr/share/nginx/html:ro -d -p 8080:80 ${nginx}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl %{VCH-IP}:8080/test.html
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  HelloWorld

Nginx with a read only config file
    Log To Console  Not implemented yet...