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
Documentation  Test 1-11 - Docker RM
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Check That VM Is Removed
    [Arguments]  ${container}
    ${id}=  Get container shortID  ${container}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc ls %{VCH-NAME}/vm
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Not Contain  ${output}  ${id}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc ls vm
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Not Contain  ${output}  ${id}

Check That Datastore Is Cleaned
    [Arguments]  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${container}

*** Test Cases ***
Basic docker remove container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox dmesg
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${container}
    Wait Until Keyword Succeeds  10x  6s  Check That Datastore Is Cleaned  ${container}

Remove a stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Container Stops  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${container}
    Wait Until Keyword Succeeds  10x  6s  Check That Datastore Is Cleaned  ${container}

Remove a running container
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: You cannot remove a running container. Stop the container before attempting removal or use -f

Force remove a running container
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${container}
    Wait Until Keyword Succeeds  10x  6s  Check That Datastore Is Cleaned  ${container}

Remove a fake container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer

Remove a container deleted out of band
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name testRMOOB -p 80:8080 busybox
    Should Be Equal As Integers  ${rc}  0
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "testRMOOB*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy %{VCH-NAME}/"testRMOOB*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm testRMOOB
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: testRMOOB
    # now recreate the same container to ensure it's completely deleted
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name testRMOOB -p 80:8080 busybox
    Should Be Equal As Integers  ${rc}  0

Remove a container created with unknown executable
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox xxxx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${container}
