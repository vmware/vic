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
Documentation  Test 5-4 - High Availability
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  High Availability Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}
Test Teardown  Run Keyword If Test Failed  Gather vSphere Logs

*** Variables ***
${namedVolume}=  named-volume
${mntDataTestContainer}=  mount-data-test
${mntTest}=  /mnt/test
${mntNamed}=  /mnt/named

*** Keywords ***
Run Regression Test With More Log Information
    Check ImageStore
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore
    # Pull an image that has been pulled already
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  busybox
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Container Stops  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Exited

    ${vmName}=  Get VM Display Name  ${container}
    Wait Until Keyword Succeeds  5x  10s  Check For The Proper Log Files  ${vmName}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  /bin/top
    Check ImageStore

    # Check for regression for #1265
    ${rc}  ${container1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox
    Should Be Equal As Integers  ${rc}  0
    ${shortname}=  Get Substring  ${container2}  1  12
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    ${lines}=  Get Lines Containing String  ${output}  ${shortname}
    Should Not Contain  ${lines}  /bin/top
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} rm ${container1}
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} rm ${container2}
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rmi busybox
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  busybox

    Scrape Logs For The Password

High Availability Setup
    [Timeout]  60 minutes
    ${name}=  Evaluate  'vic-iscsi-cluster-' + str(random.randint(1000,9999))  modules=random
    Log To Console  Create a new simple vc cluster with spec vic-cluster-ha.rb...
    ${out}=  Deploy Nimbus Testbed  spec=vic-cluster-ha.rb  args=--noSupportBundles --plugin testng --vcvaBuild "${VC_VERSION}" --esxBuild "${ESX_VERSION}" --testbedName vic-cluster-ha --runName ${name}
    Log  ${out}
    Log To Console  Finished creating cluster ${name}
    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    Log  ${vc-ip}
    ${pod}=  Fetch Pod  ${name}
    Log  ${pod}
    Close Connection
    
    Set Suite Variable  ${nimbus_pod}  ${pod}
    Set Suite Variable  ${testbedname}  ${name}

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${vc-ip}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_RESOURCE  cls
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_TIMEOUT  30m

Check All VM Migration Succeed
    [Arguments]  ${poweroff_host_ip}
    :FOR  ${index}  IN RANGE  30
    \     ${info}=  Run  govc vm.info \\*
    \     ${status}=  Run Keyword And Return Status  Should Not Contain  ${info}  ${poweroff_host_ip}
    \     Exit For Loop If  ${status}
    \     Sleep  1m
    Log  ${info}
    Return From Keyword If  ${status}
    Fail  the vms migration failed.

*** Test Cases ***
Test
    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests

    # have a few containers running and stopped for when we
    # shut down the host and HA brings it up again
    # make sure we have busybox
    Pull image  busybox

    @{running}=  Create List
    :FOR  ${index}  IN RANGE  3
    \     ${rc}  ${c}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd busybox
    \     Log  ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Append To List  ${running}  ${c}

    @{stopped}=  Create List
    :FOR  ${index}  IN RANGE  3
    \     ${rc}  ${c}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox ls
    \     Log  ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Append To List  ${stopped}  ${c}

    # Speculating that this is here to add some stability if VMs vmotion immediately after power on
    # If so it should be replaced with a check for active tasks running against the VM or a DRS rule to avoid it completely
    Sleep  2 minutes
    ${curHost}=  Get VM Host Name  %{VCH-NAME}

    ${info}=  Run  govc vm.info \\*
    Log  ${info}

    Log To Console  \nCreate a named volume and mount it to a container (Mount Inspect Test 1 of 2 - before VCH restart)\n
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=${namedVolume}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${container}  ${namedVolume}

    ${rc}  ${containerMountDataTestID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=${mntDataTestContainer} -v ${mntTest} -v ${namedVolume}:${mntNamed} busybox
    Should Be Equal As Integers  ${rc}  0

    # Create check list for Volume Inspect
    @{checkList}=  Create List  ${mntTest}  ${mntNamed}  ${namedVolume}
    Verify Volume Inspect Info  Before Host Power OFF  ${containerMountDataTestID}  ${checkList}

    Power Off Host  ${curHost}

    ${info}=  Run  govc vm.info \\*
    Log  ${info}

    # It can take a while for the host to power down and for HA to kick in
    Wait Until Keyword Succeeds  24x  10s  VM Host Has Changed  ${curHost}  %{VCH-NAME}

    # Wait for the VCH to come back up fully - if it's not completely reinitialized it will still report the old IP address 
    Wait For VCH Initialization  12x  20 seconds

    ${info}=  Run  govc vm.info \\*
    Log  ${info}

    Verify Volume Inspect Info  After Host Power OFF  ${containerMountDataTestID}  ${checkList}

    #Check if container and VCH are on the same host
    ${shortContainerID}=  Get Substring  ${containerMountDataTestID}  0  12
    ${testContainerName}=  Set Variable  ${mntDataTestContainer}-${shortContainerID}
    ${testContainerHost}=  Get VM Host Name  ${testContainerName}
    Log  ${testContainerHost} ${curHost}
    Run Keyword If  "${testContainerHost}" == "${curHost}"  Wait Until Keyword Succeeds  30x  10s  VM Host Has Changed  ${curHost}  ${testContainerName}

    # Remove Mount Data Test Container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${containerMountDataTestID}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${containerMountDataTestID}
    
    Check All VM Migration Succeed  ${curHost}
    # check running containers are still running
    :FOR  ${c}  IN  @{running}
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format '{{.State.Status}}' ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  running
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c}
    \     Log To Console  ${output}
    \     Should Be Equal As Integers  ${rc}  0

    # check stopped containers are still stopped
    :FOR  ${c}  IN  @{stopped}
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format '{{.State.Status}}' ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  exited
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c}
    \     Log To Console  ${output}
    \     Should Be Equal As Integers  ${rc}  0

Run Regression Tests
    Run Regression Test With More Log Information
