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
Suite Setup  Wait Until Keyword Succeeds  10x  10m  High Availability Setup
Suite Teardown  Nimbus Cleanup  ${list}
Test Teardown  Run Keyword If Test Failed  Gather Logs From Test Server

*** Variables ***
${esx_number}=  3
${datacenter}=  ha-datacenter
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

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rmi -f busybox
    Should Be Equal As Integers  ${rc}  0
    Check ImageStore
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  busybox

    Scrape Logs For The Password

High Availability Setup
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${vc}=  Evaluate  'VC-' + str(random.randint(1000,9999)) + str(time.clock())  modules=random,time
    ${pid}=  Deploy Nimbus vCenter Server Async  ${vc}
    Set Suite Variable  ${VC}  ${vc}

    Log To Console  \nStarting test...
    &{esxes}=  Deploy Multiple Nimbus ESXi Servers in Parallel  3  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    @{esx_names}=  Get Dictionary Keys  ${esxes}
    @{esx_ips}=  Get Dictionary Values  ${esxes}

    Set Suite Variable  @{list}  @{esx_names}[0]  @{esx_names}[1]  @{esx_names}[2]  %{NIMBUS_USER}-${vc}

    # Finish vCenter deploy
    ${output}=  Wait For Process  ${pid}
    Should Contain  ${output.stdout}  Overall Status: Succeeded

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  2 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc_ip}=  Get IP  ${vc}
    Close Connection

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${vc_ip}

    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ${datacenter}
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create cls
    Should Be Empty  ${out}

    Create A Distributed Switch  ${datacenter}

    Create Three Distributed Port Groups  ${datacenter}

    Log To Console  Add ESX hosts to the cluster
    :FOR  ${IDX}  IN RANGE  ${esx_number}
    \   ${out}=  Run  govc cluster.add -hostname=@{esx_ips}[${IDX}] -username=root -dc=${datacenter} -password=${NIMBUS_ESX_PASSWORD} -noverify=true
    \   Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    Wait Until Keyword Succeeds  5x  5min  Add Host To Distributed Switch  /${datacenter}/host/cls

    Log To Console  Enable HA and DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled -ha-enabled /${datacenter}/host/cls
    Should Be Empty  ${out}

    ${name}  ${ip}=  Deploy Nimbus NFS Datastore  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Append To List  ${list}  ${name}

    ${out}=  Run  govc datastore.create -mode readWrite -type nfs -name nfsDatastore -remote-host ${ip} -remote-path /store /ha-datacenter/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc_ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_RESOURCE  cls
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  nfsDatastore
    Set Environment Variable  TEST_TIMEOUT  30m

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
    \     Should Be Equal As Integers  ${rc}  0
    \     Append To List  ${running}  ${c}

    @{stopped}=  Create List
    :FOR  ${index}  IN RANGE  3
    \     ${rc}  ${c}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox ls
    \     Should Be Equal As Integers  ${rc}  0
    \     Append To List  ${stopped}  ${c}

    Sleep  2 minutes

    ${output}=  Run  govc vm.info %{VCH-NAME}
    @{output}=  Split To Lines  ${output}
    ${curHost}=  Fetch From Right  @{output}[-1]  ${SPACE}

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

    # Abruptly power off the host
    Open Connection  ${curHost}  prompt=:~]
    Login  root  e2eFunctionalTest
    ${out}=  Execute Command  poweroff -d 0 -f
    Close connection

    ${info}=  Run  govc vm.info \\*
    Log  ${info}

    # Really not sure what better to do here?  Otherwise, vic-machine-inspect returns the old IP address... maybe some sort of power monitoring? Can I pull uptime of the system?
    Sleep  4 minutes
    Run VIC Machine Inspect Command
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    ${info}=  Run  govc vm.info \\*
    Log  ${info}

    Verify Volume Inspect Info  After Host Power OFF  ${containerMountDataTestID}  ${checkList}

    # Remove Mount Data Test Container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${containerMountDataTestID}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${containerMountDataTestID}

    # check running containers are still running
    :FOR  ${c}  IN  @{running}
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format '{{.State.Status}}' ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  running
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c}
    \     Should Be Equal As Integers  ${rc}  0

    # check stopped containers are still stopped
    :FOR  ${c}  IN  @{stopped}
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format '{{.State.Status}}' ${c}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  exited
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c}
    \     Should Be Equal As Integers  ${rc}  0

Run Regression Tests
    Run Regression Test With More Log Information
