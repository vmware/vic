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
Documentation     Suite 24-01 - Basic
Resource          ../../resources/Util.robot
Test Setup        Set Test Environment Variables
Test Teardown     Cleanup
Default Tags


*** Keywords ***
Cleanup
    Run Keyword And Continue On Failure    Remove Group     %{VCH-NAME}

    Cleanup VIC Appliance On Test Server


Create Group
    [Arguments]    ${name}

    ${rc}  ${out}=    Run And Return Rc And Output     govc cluster.group.create -name "${name}" -vm --json 2>&1
    Should Be Equal As Integers    ${rc}    0

Remove Group
    [Arguments]    ${name}

    ${rc}  ${out}=    Run And Return Rc And Output     govc cluster.group.remove -name "${name}" --json 2>&1
    Should Be Equal As Integers    ${rc}    0


Verify Group Not Found
    [Arguments]    ${name}

    ${out}=    Run     govc cluster.group.ls -name "${name}" --json 2>&1
    Should Be Equal As Strings    ${out}    govc: group "${name}" not found

Verify Group Empty
    [Arguments]    ${name}

    ${out}=    Run     govc cluster.group.ls -name "${name}" --json 2>&1
    Should Be Equal As Strings    ${out}    null

Verify Group Contains VMs
    [Arguments]    ${name}    ${count}

    ${out}=    Run    govc cluster.group.ls -name "${name}" --json | jq 'length'
    Should Be Equal As Integers    ${out}    ${count}


Create Three Containers
    ${POWERED_OFF_CONTAINER_NAME}=    Generate Random String  15
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} create --name ${POWERED_OFF_CONTAINER_NAME} ${busybox} /bin/top

    Set Test Variable    ${POWERED_OFF_CONTAINER_NAME}

    ${POWERED_ON_CONTAINER_NAME}=    Generate Random String  15
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} create --name ${POWERED_ON_CONTAINER_NAME} ${busybox} /bin/top
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} start ${out}

    Set Test Variable    ${POWERED_ON_CONTAINER_NAME}

    ${RUN_CONTAINER_NAME}=    Generate Random String  15
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} run -d --name ${RUN_CONTAINER_NAME} ${busybox} /bin/top

    Set Test Variable    ${RUN_CONTAINER_NAME}

Delete Containers
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm ${POWERED_OFF_CONTAINER_NAME}
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm -f ${POWERED_ON_CONTAINER_NAME}
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm -f ${RUN_CONTAINER_NAME}


*** Test Cases ***
Creating a VCH creates a VM group and container VMs get added to it
    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Log To Console  %{VCH-NAME} created

    Verify Group Contains VMs    %{VCH-NAME}    1

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4


Deleting a VCH deletes its VM group
    [Teardown]    Run Keyword If Test Failed    Cleanup VIC Appliance On Test Server

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Verify Group Contains VMs    %{VCH-NAME}    1

    Run VIC Machine Delete Command

    Verify Group Not Found       %{VCH-NAME}


Deleting a container cleans up its VM group
    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4

    Delete Containers

    Verify Group Contains VMs    %{VCH-NAME}    1


Create a VCH without a VM group
    Verify Group Not Found       %{VCH-NAME}

    Create Group                 %{VCH-NAME}

    Verify Group Empty           %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    cleanup=${false}

    Verify Group Empty           %{VCH-NAME}

    Create Three Containers

    Verify Group Empty           %{VCH-NAME}


Attempt to create a VCH when a VM group with the same name already exists
    [Teardown]    Remove Group   %{VCH-NAME}

    Verify Group Not Found       %{VCH-NAME}

    Create Group                 %{VCH-NAME}

    Verify Group Empty           %{VCH-NAME}

    Run Keyword and Expect Error    *    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group    cleanup=${false}

    Verify Group Empty           %{VCH-NAME}


Deleting a VCH gracefully handles missing VM group
    [Teardown]    Run Keyword If Test Failed    Cleanup VIC Appliance On Test Server

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Verify Group Contains VMs    %{VCH-NAME}    1

    Remove Group                 %{VCH-NAME}

    Verify Group Not Found       %{VCH-NAME}

    Run VIC Machine Delete Command
