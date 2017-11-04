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
Documentation     Test 23-03 - VCH Create
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Start VIC Machine Server
Suite Teardown    Terminate All Processes    kill=True
Test Setup        Install And Prepare VIC Appliance
Default Tags


*** Keywords ***
Pull Busybox
    Run Docker Command    pull ${busybox}

    Verify Return Code
    Output Should Not Contain    Error

Install And Prepare VIC Appliance
    Install VIC Appliance To Test Server
    Pull Busybox


Install And Prepare VIC Appliance With Volume Stores
    Set Test Environment Variables
    Install VIC Appliance To Test Server    additional-args=--volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL-foo:foo
    Pull Busybox


Get VCH ID ${name}
    Get Path Under Target    vch
    ${id}=    Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="${name}").id'
    [Return]    ${id}


Run Docker Command
    [Arguments]    ${command}

    ${RC}  ${OUTPUT}=    Run And Return Rc And Output    docker %{VCH-PARAMS} ${command}
    Set Test Variable    ${RC}
    Set Test Variable    ${OUTPUT}


Populate VCH with Powered Off Container
    ${POWERED_OFF_CONTAINER_NAME}=  Generate Random String  15
    Run Docker Command    create --name ${POWERED_OFF_CONTAINER_NAME} ${busybox} /bin/top

    Verify Return Code
    Output Should Not Contain    Error

    Set Test Variable    ${POWERED_OFF_CONTAINER_NAME}

Populate VCH with Powered On Container
    ${POWERED_ON_CONTAINER_NAME}=  Generate Random String  15
    Run Docker Command    create --name ${POWERED_ON_CONTAINER_NAME} ${busybox} /bin/top

    Verify Return Code
    Output Should Not Contain    Error

    Run Docker Command    start ${OUTPUT}

    Verify Return Code
    Output Should Not Contain    Error:

    Set Test Variable    ${POWERED_ON_CONTAINER_NAME}


Verify Container Exists
    [Arguments]    ${name}

    ${vm}=    Run    govc vm.info -json=true ${name}-* | jq '.VirtualMachines | length'

    Should Be Equal As Integers       ${vm}    1

Verify Container Not Exists
    [Arguments]    ${name}

    ${vm}=    Run    govc vm.info -json=true ${name}-* | jq '.VirtualMachines | length'
    ${ds}=    Run    govc datastore.ls VIC/${name}-*

    Should Be Equal As Integers       ${vm}    0
    Should Contain                    ${ds}    was not found


Verify VCH Exists
    [Arguments]    ${path}

    Get Path Under Target             ${path}

    Verify Return Code
    Verify Status Ok

Verify VCH Not Exists
    [Arguments]    ${path}

    Get Path Under Target             ${path}

    Verify Return Code
    Verify Status Not Found

    ${rp}=    Run    govc pool.info -json=true host/*/Resources/%{VCH-NAME} | jq '.ResourcePools | length'
    ${vm}=    Run    govc vm.info -json=true %{VCH-NAME} | jq '.VirtualMachines | length'
    ${ds}=    Run    govc datastore.ls %{VCH-NAME}

    Should Be Equal As Integers       ${rp}    0
    Should Be Equal As Integers       ${vm}    0
    Should Contain                    ${ds}    was not found


Verify Volume Exists
    [Arguments]    ${volume}    ${name}

    ${ds}=    Run    govc datastore.ls ${volume}/volumes/${name}

    Should Not Contain                ${ds}    was not found

Verify Volume Not Exists
    [Arguments]    ${volume}    ${name}

    ${ds}=    Run    govc datastore.ls ${volume}/volumes/${name}

    Should Contain                    ${ds}    was not found

Verify Volume Store Exists
    [Arguments]    ${name}

    ${ds}=    Run    govc datastore.ls ${name}

    Should Not Contain                ${ds}    was not found

Verify Volume Store Not Exists
    [Arguments]    ${name}

    ${ds}=    Run    govc datastore.ls ${name}

    Should Contain                    ${ds}    was not found


*** Test Cases ***
Delete VCH
    ${id}=    Get VCH ID %{VCH-NAME}

    Verify VCH Exists                 vch/${id}

    Delete Path Under Target          vch/${id}

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}


Delete VCH within datacenter
    ${dc}=    Get Datacenter ID
    ${id}=    Get VCH ID %{VCH-NAME}

    Verify VCH Exists                 datacenter/${dc}/vch/${id}

    Delete Path Under Target          datacenter/${dc}/vch/${id}

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             datacenter/${dc}/vch/${id}


Delete invalid VCH
    ${id}=    Get VCH ID %{VCH-NAME}

    Delete Path Under Target           vch/INVALID

    Verify Return Code
    Verify Status Not Found

    Verify VCH Exists                  vch/${id}

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete VCH in invalid datacenter
    ${id}=    Get VCH ID %{VCH-NAME}

    Delete Path Under Target           datacenter/INVALID/vch/${id}

    Verify Return Code
    Verify Status Not Found

    Verify VCH Exists                  vch/${id}

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete with invalid bodies
    ${id}=    Get VCH ID %{VCH-NAME}

    Delete Path Under Target           vch/${id}    '{"invalid"}'

    Verify Return Code
    Verify Status Bad Request

    Delete Path Under Target           vch/${id}    '{"containers":"invalid"}'

    Verify Return Code
    Verify Status Unprocessable Entity
    Output Should Contain              containers

    Delete Path Under Target           vch/${id}    '{"volume_stores":"invalid"}'

    Verify Return Code
    Verify Status Unprocessable Entity
    Output Should Contain              volume_stores

    Delete Path Under Target           vch/${id}    '{"containers":"invalid", "volume_stores":"all"}'

    Verify Return Code
    Verify Status Unprocessable Entity
    Output Should Contain              containers

    Delete Path Under Target           vch/${id}    '{"containers":"all", "volume_stores":"invalid"}'

    Verify Return Code
    Verify Status Unprocessable Entity
    Output Should Contain              volume_stores

    Verify VCH Exists                  vch/${id}

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete VCH with powered off container
    ${id}=    Get VCH ID %{VCH-NAME}

    Populate VCH with Powered Off Container

    Verify Container Exists           ${POWERED_OFF_CONTAINER_NAME}
    Verify VCH Exists                 vch/${id}

    Delete Path Under Target          vch/${id}

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${POWERED_OFF_CONTAINER_NAME}


Delete VCH without deleting powered on container
    ${id}=    Get VCH ID %{VCH-NAME}

    Populate VCH with Powered On Container
    Populate VCH with Powered Off Container

    Verify Container Exists           ${POWERED_ON_CONTAINER_NAME}
    Verify Container Exists           ${POWERED_OFF_CONTAINER_NAME}
    Verify VCH Exists                 vch/${id}

    Delete Path Under Target          vch/${id}

    Verify Return Code
    Verify Status Internal Server Error

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${POWERED_ON_CONTAINER_NAME}
    Verify Container Not Exists       ${POWERED_OFF_CONTAINER_NAME}

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete VCH explicitly without deleting powered on container
    ${id}=    Get VCH ID %{VCH-NAME}

    Populate VCH with Powered On Container
    Populate VCH with Powered Off Container

    Verify Container Exists           ${POWERED_ON_CONTAINER_NAME}
    Verify Container Exists           ${POWERED_OFF_CONTAINER_NAME}
    Verify VCH Exists                 vch/${id}

    Delete Path Under Target          vch/${id}    '{"containers":"off"}'

    Verify Return Code
    Verify Status Internal Server Error

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${POWERED_ON_CONTAINER_NAME}
    Verify Container Not Exists       ${POWERED_OFF_CONTAINER_NAME}

    [Teardown]    Cleanup VIC Appliance On Test Server


Delete VCH and delete powered on container
    ${id}=    Get VCH ID %{VCH-NAME}

    Populate VCH with Powered On Container
    Populate VCH with Powered Off Container

    Verify Container Exists           ${POWERED_ON_CONTAINER_NAME}
    Verify Container Exists           ${POWERED_OFF_CONTAINER_NAME}
    Verify VCH Exists                 vch/${id}

    Delete Path Under Target          vch/${id}    '{"containers":"all"}'

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${POWERED_ON_CONTAINER_NAME}
    Verify Container Not Exists       ${POWERED_OFF_CONTAINER_NAME}


Delete VCH and powered off containers and volumes
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered off container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code
    #

    # Setup powered off container VM with named volume on named volume store
    Run Docker Command    volume create --name ${NAME}-volume-foo --opt VolumeStore=foo
    Verify Return Code

    Run Docker Command    create --name ${NAME}-cfoo -v ${NAME}-volume-foo:/volume ${busybox} /bin/top
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Container Exists           ${NAME}-cfoo
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Store Exists        %{VCH-NAME}-VOL-foo
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Exists              %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo

    Delete Path Under Target          vch/${id}    '{"containers":"off","volume_stores":"all"}'

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${NAME}-container
    Verify Container Not Exists       ${NAME}-cfoo
    Verify Volume Store Not Exists    %{VCH-NAME}-VOL
    Verify Volume Store Not Exists    %{VCH-NAME}-VOL-foo
    Verify Volume Not Exists          %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Not Exists          %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo


Delete VCH and powered on containers and volumes
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered on container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Run Docker Command    start ${OUTPUT}
    Verify Return Code
    #

    # Setup powered on container VM with named volume on named volume store
    Run Docker Command    volume create --name ${NAME}-volume-foo --opt VolumeStore=foo
    Verify Return Code

    Run Docker Command    create --name ${NAME}-cfoo -v ${NAME}-volume-foo:/volume ${busybox} /bin/top
    Verify Return Code

    Run Docker Command    start ${OUTPUT}
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Container Exists           ${NAME}-cfoo
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Store Exists        %{VCH-NAME}-VOL-foo
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Exists              %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo

    Delete Path Under Target          vch/${id}    '{"containers":"all","volume_stores":"all"}'

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${NAME}-container
    Verify Volume Store Not Exists    %{VCH-NAME}-VOL
    Verify Volume Store Not Exists    %{VCH-NAME}-VOL-foo
    Verify Volume Not Exists          %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Not Exists          %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo


Delete VCH and powered off container and preserve volumes
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered off container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code
    #

    # Setup powered off container VM with named volume on named volume store
    Run Docker Command    volume create --name ${NAME}-volume-foo --opt VolumeStore=foo
    Verify Return Code

    Run Docker Command    create --name ${NAME}-cfoo -v ${NAME}-volume-foo:/volume ${busybox} /bin/top
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Container Exists           ${NAME}-cfoo
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Store Exists        %{VCH-NAME}-VOL-foo
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Exists              %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo

    Delete Path Under Target          vch/${id}    '{"containers":"off","volume_stores":"none"}'

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${NAME}-container
    Verify Container Not Exists       ${NAME}-cfoo
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Store Exists        %{VCH-NAME}-VOL-foo
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume
    Verify Volume Exists              %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo

    # Re-use preserved volumes
    Install And Prepare VIC Appliance

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume

    Run Docker Command    create --name ${NAME}-cfoo -v ${NAME}-volume-foo:/volume ${busybox} /bin/top
    Verify Return Code

    Verify Container Exists           ${NAME}-cfoo
    Verify Volume Store Exists        %{VCH-NAME}-VOL-foo
    Verify Volume Exists              %{VCH-NAME}-VOL-foo    ${NAME}-volume-foo


Delete VCH and powered on container but preserve volume
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered on container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Run Docker Command    start ${OUTPUT}
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume

    Delete Path Under Target          vch/${id}    '{"containers":"all","volume_stores":"none"}'

    Verify Return Code
    Verify Status Accepted

    Verify VCH Not Exists             vch/${id}
    Verify Container Not Exists       ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume

    # Re-use preserved volume
    Install And Prepare VIC Appliance

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume


Delete VCH and preserve powered on container and volumes
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered on container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Run Docker Command    start ${OUTPUT}
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume

    Delete Path Under Target          vch/${id}    '{"containers":"off","volume_stores":"none"}'

    Verify Return Code
    Verify Status Internal Server Error

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume


Delete VCH and preserve powered on container and fail to delete volumes
    [Setup]    Install And Prepare VIC Appliance With Volume Stores
    ${id}=    Get VCH ID %{VCH-NAME}

    ${NAME}=  Generate Random String  15

    # Setup powered on container VM with named volume on default volume store
    Run Docker Command    volume create --name ${NAME}-volume
    Verify Return Code

    Run Docker Command    create --name ${NAME}-container -v ${NAME}-volume:/volume ${busybox} /bin/top
    Verify Return Code

    Run Docker Command    start ${OUTPUT}
    Verify Return Code
    #

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume

    Delete Path Under Target          vch/${id}    '{"containers":"off","volume_stores":"all"}'

    Verify Return Code
    Verify Status Internal Server Error

    Verify VCH Exists                 vch/${id}
    Verify Container Exists           ${NAME}-container
    Verify Volume Store Exists        %{VCH-NAME}-VOL
    Verify Volume Exists              %{VCH-NAME}-VOL        ${NAME}-volume
