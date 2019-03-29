# Copyright 2016-2019 VMware, Inc. All Rights Reserved.
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
Documentation  Test 11-06- Upgrade To latest Check storage quota feature
Resource  ../../resources/Util.robot
Suite Setup  Disable Ops User And Install VIC To Test Server
Suite Teardown  Re-Enable Ops User And Clean Up VIC Appliance
Default Tags

*** Variables ***
${run-as-ops-user}=  ${EMPTY}
${old_version}=  v1.4.3

*** Keywords ***
Disable Ops User And Install VIC To Test Server
    ${run-as-ops-user}=  Get Environment Variable  RUN_AS_OPS_USER  0
    Set Environment Variable  RUN_AS_OPS_USER  0
    Install VIC with version to Test Server  ${old_version}  additional-args=--container-name-convention VCH_1_{name} --cpu-reservation 1 --cpu-shares normal --memory-reservation 1 --memory-shares normal --endpoint-cpu 1 --endpoint-memory 2048 --base-image-size 8GB --bridge-network-range 172.16.0.0/12 --container-network-firewall vm-network:published --certificate-key-size 2048

Re-Enable Ops User And Clean Up VIC Appliance
    Set Environment Variable  RUN_AS_OPS_USER  ${run-as-ops-user}
    Log  vchName1:${vchName1} vchName2:${vchName2} vchName3:${vchName3}

    Reload VCH Related Environment Variables  ${vchName1}  ${vchIP1}  ${vchParams1}  ${vchAdmin1}
    Clean up VIC Appliance And Local Binary
    Reload VCH Related Environment Variables  ${vchName2}  ${vchIP2}  ${vchParams2}  ${vchAdmin2}
    Cleanup VIC Appliance On Test Server
    Reload VCH Related Environment Variables  ${vchName3}  ${vchIP3}  ${vchParams3}  ${vchAdmin3}
    Cleanup VIC Appliance On Test Server

Create Containers For VCH_1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm1 -id busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm2 -id nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name c-vm3 -i nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt Capacity=5G volume1
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Custom Install VCH
    [Arguments]  ${version}  ${prefix}  ${imageStore}=%{TEST_DATASTORE}  ${baseImageSize}=8G  ${cleanup}=${False}  ${bridgeNetworkIPRange}=172.17.0.0/12
    Custom VIC With Version To Test Server  ${version}  imageStore=${imageStore}  cleanup=${cleanup}  additional-args=--container-name-convention ${prefix}_{name} --cpu-reservation 1 --cpu-shares normal --memory-reservation 1 --memory-shares normal --endpoint-cpu 1 --endpoint-memory 2048 --base-image-size ${baseImageSize} --volume-store %{TEST_DATASTORE}/shared:shared --bridge-network-range ${bridgeNetworkIPRange} --container-network-firewall vm-network:published --certificate-key-size 2048

Custom VIC With Version To Test Server
    [Arguments]  ${version}  ${imageStore}=%{TEST_DATASTORE}  ${cleanup}=${true}  ${additional-args}=
    Log To Console  \nDownloading vic ${version} from gcp...
    ${time}=  Evaluate  time.clock()  modules=time
    ${rc}  ${output}=  Run And Return Rc And Output  wget https://storage.googleapis.com/vic-engine-releases/vic_${version}.tar.gz -O vic-${time}.tar.gz
    Create Directory  vic-${time}
    ${rc}  ${output}=  Run And Return Rc And Output  tar zxvf vic-${time}.tar.gz -C vic-${time}
    Custom Install VIC Appliance To Test Server  vic-machine=./vic-${time}/vic/vic-machine-linux  appliance-iso=./vic-${time}/vic/appliance.iso  bootstrap-iso=./vic-${time}/vic/bootstrap.iso  imageStore=${imageStore}  certs=${false}  cleanup=${cleanup}  additional-args=${additional-args}  vol=default
    Set Environment Variable  VIC-ADMIN  %{VCH-IP}:2378
    Set Environment Variable  INITIAL-VERSION  ${version}
    Run  rm -rf vic-${time}.tar.gz vic-${time}

Custom Install VIC Appliance To Test Server
    [Arguments]  ${vic-machine}=bin/vic-machine-linux  ${appliance-iso}=bin/appliance.iso  ${bootstrap-iso}=bin/bootstrap.iso  ${imageStore}=%{TEST_DATASTORE}  ${certs}=${true}  ${vol}=default  ${cleanup}=${true}  ${debug}=1  ${additional-args}=${EMPTY}
    Set Test Environment Variables
    ${opsuser-args}=  Get Ops User Args
    ${output}=  Custom Install VIC With Current Environment Variables  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${imageStore}  ${certs}  ${vol}  ${cleanup}  ${debug}  ${opsuser-args}  ${additional-args}
    Log  ${output}
    [Return]  ${output}

Custom Install VIC With Current Environment Variables
    [Arguments]  ${vic-machine}=bin/vic-machine-linux  ${appliance-iso}=bin/appliance.iso  ${bootstrap-iso}=bin/bootstrap.iso  ${imageStore}=%{TEST_DATASTORE}  ${certs}=${true}  ${vol}=default  ${cleanup}=${true}  ${debug}=1  ${opsuser-args}=${EMPTY}  ${additional-args}=${EMPTY}
    # disable firewall
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.esxcli network firewall set -e false
    # Attempt to cleanup old/canceled tests
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Networks On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling vSwitches On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Containers On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Folders On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling Resource Pools On Test Server
    Run Keyword If  ${cleanup}  Run Keyword And Ignore Error  Cleanup Dangling VM Groups On Test Server

    # Install the VCH now
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Custom Run VIC Machine Command  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${imageStore}  ${vol}  ${debug}  ${opsuser-args}  ${additional-args}
    Log  ${output}

    ${noIP}=  Set Variable If  '"waiting for IP"' in '''${output}'''  True  False
    Run Keyword If  ${noIP}  Log             Possible DHCP failure
    Run Keyword If  ${noIP}  Log To Console  Possible DHCP failure

    Should Contain  ${output}  Installer completed successfully

    Get Docker Params  ${output}  ${certs}
    Log To Console  Installer completed successfully: %{VCH-NAME}...

    [Return]  ${output}

Custom Run VIC Machine Command
    [Tags]  secret
    [Arguments]  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${imageStore}  ${vol}  ${debug}  ${opsuser-args}  ${additional-args}
    ${output}=  Run Keyword If  ${certs}  Run  ${vic-machine} create --debug ${debug} --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=${imageStore} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --insecure-registry wdc-harbor-ci.eng.vmware.com --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL:${vol} --container-network=%{PUBLIC_NETWORK}:public ${vicmachinetls} ${opsuser-args} ${additional-args}
    Run Keyword If  ${certs}  Should Contain  ${output}  Installer completed successfully
    Return From Keyword If  ${certs}  ${output}

    ${output}=  Run Keyword Unless  ${certs}  Run  ${vic-machine} create --debug ${debug} --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=${imageStore} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --insecure-registry wdc-harbor-ci.eng.vmware.com --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL:${vol} --container-network=%{PUBLIC_NETWORK}:public --no-tlsverify ${opsuser-args} ${additional-args}
    Run Keyword Unless  ${certs}  Should Contain  ${output}  Installer completed successfully
    [Return]  ${output}
    
Get VCH Related Environment Variables
    ${vchName}=  Get Environment Variable  VCH-NAME
    ${vchIP}=  Get Environment Variable  VCH-IP
    ${vchParams}=  Get Environment Variable  VCH-PARAMS
    ${vchAdmin}=  Get Environment Variable  VIC-ADMIN
    [Return]  ${vchName}  ${vchIP}  ${vchParams}  ${vchAdmin}

Reload VCH Related Environment Variables
    [Arguments]  ${vchName}  ${vchIP}  ${vchParams}  ${vchAdmin}
    Set Environment Variable  VCH-NAME  ${vchName}
    Set Environment Variable  VCH-IP  ${vchIP}
    Set Environment Variable  VCH-PARAMS  ${vchParams}
    Set Environment Variable  VIC-ADMIN  ${vchAdmin}
 
Get Container Image Info
    [Arguments]  ${info}
    ${containers-line}=  Get Lines Containing String  ${info}  Containers:
    ${running-line}=  Get Lines Containing String  ${info}  Running:
    ${paused-line}=  Get Lines Containing String  ${info}  Paused:
    ${stopped-line}=  Get Lines Containing String  ${info}  Stopped:
    ${images-line}=  Get Lines Containing String  ${info}  Images:
    @{containers-line}=  Split String  ${containers-line}
    Length Should Be  ${containers-line}  2
    ${containersVal}=  Convert To Number  @{containers-line}[1]   

    @{running-line}=  Split String  ${running-line}
    Length Should Be  ${running-line}  2
    ${runningVal}=  Convert To Number  @{running-line}[1]

    @{paused-line}=  Split String  ${paused-line}
    Length Should Be  ${paused-line}  2
    ${pausedVal}=  Convert To Number  @{paused-line}[1]
    
    @{stopped-line}=  Split String  ${stopped-line}
    Length Should Be  ${stopped-line}  2
    ${stoppedVal}=  Convert To Number  @{stopped-line}[1]

    @{images-line}=  Split String  ${images-line}
    Length Should Be  ${images-line}  2
    ${imagesVal}=  Convert To Number  @{images-line}[1]

    [Return]  ${containersVal}  ${runningVal}  ${pausedVal}  ${stoppedVal}  ${imagesVal}

Get Containers Storage Usage
    [Arguments]  ${info}
    ${containerUsageline}=  Get Lines Containing String  ${info}  VCH containers storage usage:
    @{usageline}=  Split String  ${ContainerUsageline}
    Length Should Be  ${usageline}  6
    ${usageval}=  Convert To Number  @{usageline}[4]

    [Return]  ${usageval}

Get Storage Quota Limit
    [Arguments]  ${info}
    ${limitline}=  Get Lines Containing String  ${info}  VCH storage limit:
    @{limitline}=  Split String  ${limitline}
    Length Should Be  ${limitline}  5
    ${limitval}=  Convert To Number  @{limitline}[3]

    [Return]  ${limitval}

Create Containers And Volume For VCH_2
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm1 -id busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm2 -id nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name c-vm3 -i nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt Capacity=5G volume2
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt Capacity=5G --opt VolumeStore=shared volume1
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Create Containers For VCH_3
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm1 -id busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name c-vm2 -id nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt Capacity=5G --opt VolumeStore=shared volume3
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Run Docker Info Cmd
    [Arguments]  ${vchParams}
    ${rc}  ${output}=  Run And Return Rc And Output  time docker ${vchParams} info
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    
    [Return]  ${output}

Configure VCH With Storage Quota
    [Arguments]  ${vchName}  ${value}
    ${output}=  Run  bin/vic-machine-linux configure --name=${vchName} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --storage-quota ${value}
    Log  ${output}
    Should Contain  ${output}  Completed successfully

Configure Rollback
    [Arguments]  ${vchName}
    ${output}=  Run  bin/vic-machine-linux configure --name=${vchName} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --rollback
    Log  ${output}
    Should Contain  ${output}  Completed successfully

Get Docker Info Exec Time
    [Arguments]  ${info}
    ${execTimeline}=  Get Lines Containing String  ${info}  elapsed
    @{realTime}=  Split String  ${execTimeline}
    @{execTime}=  Split String  @{realTime}[2]  e
    Length Should Be  ${execTime}  3
    ${execTime}=  Fetch From Right  @{execTime}[0]  :
    ${execTime}=  Convert To Number  ${execTime}

    [Return]  ${execTime}

Check After Upgrade
    [Arguments]  ${info}  ${containerNum}  ${runningNum}  ${pausedNum}  ${stoppedNum}  ${imagesNum}  ${containersStorageUsageLow}  ${containersStorageUsageHigh}  ${execTime}=2
    ${containersVal}  ${runningVal}  ${pausedVal}  ${stoppedVal}  ${imagesVal}=  Get Container Image Info  ${info}
    Log  containersVal:${containersVal} runningVal:${runningVal} pausedVal:${pausedVal} stoppedVal:${stoppedVal} imagesVal:${imagesVal}
    Should Be Equal As Integers  ${containersVal}  ${containerNum}
    Should Be Equal As Integers  ${runningVal}  ${runningNum}
    Should Be Equal As Integers  ${pausedVal}  ${pausedVal}
    Should Be Equal As Integers  ${stoppedVal}  ${stoppedNum}
    Should Be Equal As Integers  ${imagesVal}  ${imagesNum}
    
    ${usageval}=  Get Containers Storage Usage  ${info}
    Log  usageval:${usageval}
    Should Be True  ${containersStorageUsageLow} < ${usageval} < ${containersStorageUsageHigh}

    ${time}=  Get Docker Info Exec Time  ${info}
    Log  time:${time}s
    Should Be True  ${time} < ${execTime}

*** Test Cases ***
Test Storage Quota
    ${vchName}  ${vchIP}  ${vchParams}  ${vchAdmin}=  Get VCH Related Environment Variables
    Set Suite Variable  ${vchName1}  ${vchName}
    Set Suite Variable  ${vchIP1}  ${vchIP}
    Set Suite Variable  ${vchParams1}  ${vchParams}
    Set Suite Variable  ${vchAdmin1}  ${vchAdmin}
    Create Containers For VCH_1

    Custom Install VCH  ${old_version}  VCH_2  imageStore=%{TEST_DATASTORE}/shared_imageFolder
    ${vchName}  ${vchIP}  ${vchParams}  ${vchAdmin}=  Get VCH Related Environment Variables
    Set Suite Variable  ${vchName2}  ${vchName}
    Set Suite Variable  ${vchIP2}  ${vchIP}
    Set Suite Variable  ${vchParams2}  ${vchParams}
    Set Suite Variable  ${vchAdmin2}  ${vchAdmin}
    Create Containers And Volume For VCH_2

    Custom Install VCH  ${old_version}  VCH_3  imageStore=%{TEST_DATASTORE}/shared_imageFolder  baseImageSize=10G  bridgeNetworkIPRange=127.18.0.0/12  
    ${vchName}  ${vchIP}  ${vchParams}  ${vchAdmin}=  Get VCH Related Environment Variables
    Set Suite Variable  ${vchName3}  ${vchName}
    Set Suite Variable  ${vchIP3}  ${vchIP}
    Set Suite Variable  ${vchParams3}  ${vchParams}
    Set Suite Variable  ${vchAdmin3}  ${vchAdmin}
    Create Containers For VCH_3
    
    Reload VCH Related Environment Variables  ${vchName1}  ${vchIP1}  ${vchParams1}  ${vchAdmin1}
    Upgrade
    Check Upgraded Version

    Reload VCH Related Environment Variables  ${vchName2}  ${vchIP2}  ${vchParams2}  ${vchAdmin2}
    Upgrade
    Check Upgraded Version
    
    Reload VCH Related Environment Variables  ${vchName3}  ${vchIP3}  ${vchParams3}  ${vchAdmin3}
    Upgrade
    Check Upgraded Version

    ${output}=  Run Docker Info Cmd  ${vchParams1}
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    Check After Upgrade  ${output}  3  2  0  1  2  27  30
    Configure VCH With Storage Quota  ${vchName1}  25
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    Check After Upgrade  ${output}  3  2  0  1  2  27  30
    ${storageLimitValue}=  Get Storage Quota Limit  ${output}
    Should Be Equal As Integers  ${storageLimitValue}  25
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${vchParams1} create --name=fail_container -it busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Storage quota exceeds
    
    Configure VCH With Storage Quota  ${vchName1}  50
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    ${storageLimitValue}=  Get Storage Quota Limit  ${output}
    Should Be Equal As Integers  ${storageLimitValue}  50
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${vchParams1} create --name=success_container -it busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    Check After Upgrade  ${output}  4  2  0  2  2  36  40
    
    ${output}=  Run Docker Info Cmd  ${vchParams2}
    Check After Upgrade  ${output}  3  2  0  1  2  27  30
    Configure VCH With Storage Quota  ${vchName2}  60
    ${output}=  Run Docker Info Cmd  ${vchParams2}
    ${storageLimitValue}=  Get Storage Quota Limit  ${output}
    Should Be Equal As Integers  ${storageLimitValue}  60
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${vchParams2} create --name=test_container -it busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Run Docker Info Cmd  ${vchParams2}
    Check After Upgrade  ${output}  4  2  0  2  2  36  40

    ${output}=  Run Docker Info Cmd  ${vchParams3}
    Check After Upgrade  ${output}  2  2  0  0  2  22  24
    Should Not Contain  ${output}  VCH storage limit:

    Reload VCH Related Environment Variables  ${vchName1}  ${vchIP1}  ${vchParams1}  ${vchAdmin1}
    Configure Rollback  ${vchName1}
    Rollback
    Check Original Version
    Upgrade with ID
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    ${output}=  Run Docker Info Cmd  ${vchParams1}
    Check After Upgrade  ${output}  4  2  0  2  2  36  40
