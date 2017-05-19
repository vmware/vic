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
Suite Teardown  Nimbus Cleanup  ${list}

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

*** Test Cases ***
Test
    ${status}=  Get State Of Github Issue  4858
    Run Keyword If  '${status}' == 'closed'  Fail  Test 5-4-High-Availability.robot needs to be updated now that Issue #4858 has been resolved

    Log To Console  \nStarting test...
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${esx2}  ${esx2-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${esx3}  ${esx3-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    ${vc}  ${vc-ip}=  Deploy Nimbus vCenter Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${VC}  ${vc}

    Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${vc}

    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ha-datacenter
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create cls
    Should Be Empty  ${out}

    Log To Console  Add ESX host to the VC
    ${out}=  Run  govc cluster.add -hostname=${esx1-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx2-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx3-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create a distributed switch
    ${out}=  Run  govc dvs.create -dc=ha-datacenter test-ds
    Should Contain  ${out}  OK

    Log To Console  Create three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds bridge
    Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    Wait Until Keyword Succeeds  5x  5min  Add Host To Distributed Switch  /ha-datacenter/host/cls

    Log To Console  Enable HA on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled -ha-enabled /ha-datacenter/host/cls
    Should Be Empty  ${out}

    ${name}  ${ip}=  Deploy Nimbus NFS Datastore  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Append To List  ${list}  ${name}

    ${out}=  Run  govc datastore.create -mode readWrite -type nfs -name nfsDatastore -remote-host ${ip} -remote-path /store /ha-datacenter/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_DATASTORE  nfsDatastore
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Regression Tests

    # have a few containers running and stopped for when we
    # shut down the host and HA brings it up again
    # make sure we have busybox
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0

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

    ${output}=  Run  govc vm.info %{VCH-NAME}/%{VCH-NAME}
    @{output}=  Split To Lines  ${output}
    ${curHost}=  Fetch From Right  @{output}[-1]  ${SPACE}

    ${info}=  Get VM Info  \*
    Log  ${info}

    # Abruptly power off the host
    Open Connection  ${curHost}  prompt=:~]
    Login  root  e2eFunctionalTest
    ${out}=  Execute Command  poweroff -d 0 -f
    Close connection

    ${info}=  Get VM Info  \*
    Log  ${info}

    # Really not sure what better to do here?  Otherwise, vic-machine-inspect returns the old IP address... maybe some sort of power monitoring? Can I pull uptime of the system?
    Sleep  4 minutes
    Run VIC Machine Inspect Command
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    ${info}=  Get VM Info  \*
    Log  ${info}

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
