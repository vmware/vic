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
Test Teardown     Cleanup
Default Tags


*** Keywords ***
Cleanup
    Run Keyword And Continue On Failure    Remove Group     %{VCH-NAME}

    Cleanup VIC Appliance On Test Server

Remove Group
    [Arguments]    ${name}

    ${rc}  ${out}=    Run And Return Rc And Output     govc cluster.group.remove -name ${name} --json 2>&1
    Should Be Equal As Integers    ${rc}    0


Verify Group Not Found
    [Arguments]    ${name}

    ${out}=    Run     govc cluster.group.ls -name ${name} --json 2>&1
    Should Be Equal As Strings    ${out}    govc: group "${name}" not found

Verify Group Contains VMs
    [Arguments]    ${name}    ${count}

    ${out}=    Run    govc cluster.group.ls -name ${name} --json | jq 'length'
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
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} run --name ${RUN_CONTAINER_NAME} ${busybox} /bin/top

    Set Test Variable    ${RUN_CONTAINER_NAME}

Delete Containers
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm --name ${POWERED_OFF_CONTAINER_NAME} ${busybox} /bin/top
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm --name ${POWERED_ON_CONTAINER_NAME} ${busybox} /bin/top
    ${rc}  ${out}=    Run And Return Rc And Output    docker %{VCH-PARAMS} rm --name ${RUN_CONTAINER_NAME} ${busybox} /bin/top


*** Test Cases ***
Creating a VCH creates a VM group and container VMs get added to it
    Set Test Environment Variables

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Verify Group Contains VMs    %{VCH-NAME}    1

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4


Deleting a VCH deletes its VM group
    Set Test Environment Variables

    Verify Group Not Found         %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Verify Group Contains VMs    %{VCH-NAME}    1

    Cleanup VIC Appliance On Test Server

    Verify Group Not Found         %{VCH-NAME}


Deleting a container cleans up its VM group
    Set Test Environment Variables

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4

    Delete Containers

    Verify Group Contains VMs    %{VCH-NAME}    1
