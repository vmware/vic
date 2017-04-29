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
Documentation  Test 5-3 - Enhanced Linked Mode
Resource  ../../resources/Util.robot
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    ${name}=  Evaluate  'els-' + str(random.randint(1000,9999))  modules=random
    Set Test Variable  ${user}  %{NIMBUS_USER}
    ${output}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --plugin test-vpx --testbedName test-vpx-m2n2-vcva-3esx-pxeBoot-8gbmem --vcvaBuild ${VC_VERSION} --esxPxeDir ${ESX_VERSION} --runName ${name}

    ${output}=  Split To Lines  ${output}
    :FOR  ${line}  IN  @{output}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  ${name}.vc.0' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${vc1-ip}  ${ip}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  ${name}.vc.1' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${vc2-ip}  ${ip}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  ${name}.esx.0' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${esx1-ip}  ${ip}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  ${name}.esx.1' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${esx2-ip}  ${ip}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  ${name}.esx.2' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${esx3-ip}  ${ip}

    ${esx1}  ${esx4-ip}  ${esx2}  ${esx5-ip}  ${esx3}  ${esx6-ip}=  Deploy Multiple Nimbus ESXi Servers in Parallel  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${user}-${name}.vc.0  ${user}-${name}.vc.1  ${user}-${name}.vc.2  ${user}-${name}.vc.3  ${user}-${name}.nfs.0  ${user}-${name}.esx.0  ${user}-${name}.esx.1  ${user}-${name}.esx.2

    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${esx1-ip}
    Wait Until Keyword Succeeds  10x  3 minutes  Change ESXi Server Password  e2eFunctionalTest
    Set Environment Variable  GOVC_URL  root:@${esx2-ip}
    Wait Until Keyword Succeeds  10x  3 minutes  Change ESXi Server Password  e2eFunctionalTest
    Set Environment Variable  GOVC_URL  root:@${esx3-ip}
    Wait Until Keyword Succeeds  10x  3 minutes  Change ESXi Server Password  e2eFunctionalTest

    Set Environment Variable  GOVC_URL  ${vc1-ip}
    Set Environment Variable  GOVC_USERNAME  administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23

    # First VC cluster
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

    Create A Distributed Switch  ha-datacenter

    Create Three Distributed Port Groups  ha-datacenter

    Add Host To Distributed Switch  /ha-datacenter/host/cls

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /ha-datacenter/host/cls
    Should Be Empty  ${out}

    # Second VC cluster
    Set Environment Variable  GOVC_URL  ${vc2-ip}
    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ha-datacenter
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create cls
    Should Be Empty  ${out}

    Log To Console  Add ESX host to the VC
    ${out}=  Run  govc cluster.add -hostname=${esx4-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx5-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx6-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Create A Distributed Switch  ha-datacenter

    Create Three Distributed Port Groups  ha-datacenter

    Add Host To Distributed Switch  /ha-datacenter/host/cls

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /ha-datacenter/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  GOVC_URL  ${vc1-ip}
    Set Environment Variable  TEST_URL_ARRAY  ${vc1-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server

    Run Regression Tests
