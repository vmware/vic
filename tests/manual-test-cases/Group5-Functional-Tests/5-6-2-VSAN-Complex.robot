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
Documentation  Test 5-6-2 - VSAN-Complex
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Complex VSAN
    ${name}=  Evaluate  'vic-vsan-complex-' + str(random.randint(1000,9999))  modules=random
    Set Test Variable  ${user}  %{NIMBUS_USER}
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --vcvaBuild ${VC_VERSION} --esxPxeDir ${ESX_VERSION} --esxBuild ${ESX_VERSION} --testbedName vcqa-vsan-complex-pxeBoot-vcva --runName ${name}
    ${out}=  Split To Lines  ${out}
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  .vcva-${VC_VERSION}' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${vc-ip}  ${ip}
    \   Exit For Loop If  ${status}

    Set Global Variable  @{list}  ${user}-${name}.vcva-${VC_VERSION}  ${user}-${name}.esx.0  ${user}-${name}.esx.1  ${user}-${name}.esx.2  ${user}-${name}.esx.3  ${user}-${name}.esx.4  ${user}-${name}.esx.5  ${user}-${name}.esx.6  ${user}-${name}.esx.7  ${user}-${name}.nfs.0  ${user}-${name}.iscsi.0

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Wait Until Keyword Succeeds  10x  3 minutes  Create A Distributed Switch  vcqaDC

    Create Three Distributed Port Groups  vcqaDC

    Add Host To Distributed Switch  /vcqaDC/host/cluster-vsan-1

    Log To Console  Enable DRS and VSAN on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cluster-vsan-1
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    ${datastore}=  Run  govc ls -t Datastore host/cluster-vsan-1/* | grep -v local | cut -d '/' -f 6 | sort | uniq | grep vsan
    Set Environment Variable  TEST_DATASTORE  "${datastore}"
    Set Environment Variable  TEST_RESOURCE  cluster-vsan-1
    Set Environment Variable  TEST_TIMEOUT  30m

    ${out}=  Run  govc datastore.vsan.dom.ls -ds %{TEST_DATASTORE} -l -o
    Should Be Empty  ${out}

    Install VIC Appliance To Test Server

    Run Regression Tests

    Cleanup VIC Appliance On Test Server

    ${out}=  Run  govc datastore.vsan.dom.ls -ds %{TEST_DATASTORE} -l -o
    Should Be Empty  ${out}