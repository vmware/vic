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
Documentation  Test 13-2 - vMotion Container
Resource  ../../resources/Util.robot
Suite Setup  Create a VSAN Cluster
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}
Test Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Test
    Set Test Variable  ${user}  %{NIMBUS_USER}
    Set Global Variable  @{list}  ${user}-vic-vmotion.vcva-${VC_VERSION}  ${user}-vic-vmotion.esx.0  ${user}-vic-vmotion.esx.1  ${user}-vic-vmotion.esx.2  ${user}-vic-vmotion.esx.3  ${user}-vic-vmotion.nfs.0  ${user}-vic-vmotion.iscsi.0
    Install VIC Appliance To Test Server
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container3}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container3}
    Should Be Equal As Integers  ${rc}  0
    
    ${vmName1}=  Get VM display name  ${container1}
    ${vmName2}=  Get VM display name  ${container2}
    ${vmName3}=  Get VM display name  ${container3}
    
    vMotion A VM  %{VCH-NAME}/${vmName1}
    vMotion A VM  %{VCH-NAME}/${vmName2}
    vMotion A VM  %{VCH-NAME}/${vmName3}
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container1}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs ${container3}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container3}
    Should Be Equal As Integers  ${rc}  0
