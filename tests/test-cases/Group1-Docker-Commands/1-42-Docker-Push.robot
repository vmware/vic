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
Documentation   Test 1-42 - Docker Push
Resource        ../../resources/Util.robot
Suite Setup     Push Setup
Suite Teardown  Push Cleanup

*** Keywords ***
Push Setup
    Install VIC Appliance To Test Server
    Set Environment Variable  REGISTRY-IP  %{VCH-IP}:5000
    Set Environment Variable  REGISTRY-VCH-NAME  %{VCH-NAME}
    Set Environment Variable  REGISTRY-VIC-ADMIN  %{VIC-ADMIN}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p 5000:5000 --name registry registry
    Should Be Equal As Integers  ${rc}  0
    Install VIC Appliance To Test Server  cleanup=${false}  additional-args=--insecure-registry %{REGISTRY-IP}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox_1_26_0}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${ubuntu}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${alpine}
    Should Be Equal As Integers  ${rc}  0

Push Cleanup
    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server
    Set Environment Variable  VCH-NAME  %{REGISTRY-VCH-NAME}
    Set Environment Variable  VIC-ADMIN  %{REGISTRY-VIC-ADMIN}
    Cleanup VIC Appliance On Test Server

*** Test Cases ***
Push busybox image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag ${busybox} %{REGISTRY-IP}/busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push %{REGISTRY-IP}/busybox
    Should Be Equal As Integers  ${rc}  0

Push busybox image with tag
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag ${busybox_1_26_0} %{REGISTRY-IP}/busybox:1.26.0
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push %{REGISTRY-IP}/busybox:1.26.0
    Should Be Equal As Integers  ${rc}  0

Push ubuntu image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag ${ubuntu} %{REGISTRY-IP}/ubuntu
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push %{REGISTRY-IP}/ubuntu
    Should Be Equal As Integers  ${rc}  0

Push image with content trust
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag ${alpine} %{REGISTRY-IP}/alpine
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push --disable-content-trust %{REGISTRY-IP}/alpine
    Should Be Equal As Integers  ${rc}  0

Push newer image with small change to it
    Log  Need commit first
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --name push1 ${busybox} tail -f /dev/null
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec push1 touch test.file
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop -t1 push1
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} commmit push1 busybox-test
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag busybox-test %{REGISTRY-IP}/busybox-test
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push %{REGISTRY-IP}/busybox-test
    #Should Be Equal As Integers  ${rc}  0
    # Verify small push?

Try to push a fake image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push %{REGISTRY-IP}/fakeimage
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  An image does not exist locally with the tag: %{REGISTRY-IP}/fakeimage

Try to push to a fake repo
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} tag ${busybox} fakerepo:5000/busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} push fakerepo:5000/busybox
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Get https://fakerepo:5000/v1/_ping: dial tcp: lookup fakerepo on 10.118.81.1:53: no such host
