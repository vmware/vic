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
Documentation  Test 5-28 - Failed Install
Resource  ../../resources/Util.robot
Suite Setup  Wait Until Keyword Succeeds  10x  10m  Failed Install Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
Failed Install Setup
    [Timeout]    110 minutes
    Log To Console  Starting failed install test...
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${name}=  Evaluate  'vic-5-28-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --plugin testng --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-simple-cluster --testbedSpecRubyFile /dbc/pa-dbc1111/mhagen/nimbus-testbeds/testbeds/vic-simple-cluster.rb --runName ${name}
    Set Suite Variable  @{list}  %{NIMBUS_USER}-${name}.esx.0  %{NIMBUS_USER}-${name}.esx.1  %{NIMBUS_USER}-${name}.esx.2  %{NIMBUS_USER}-${name}.nfs.0  %{NIMBUS_USER}-${name}.vc.0
    Log To Console  Finished Creating Cluster ${name}

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    Close Connection

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    ${rc}  ${output}=  Run And Return Rc And Output  govc dvs.portgroup.add -vlan=2000 -dvs test-ds invalid-dpg
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  invalid-dpg
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  nfs0-1
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  5m

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    Custom Testbed Keepalive  /dbc/pa-dbc1111/mhagen

    Run Keyword And Ignore Error  Install VIC Appliance To Test Server
    ${out}=  Run  govc ls vm
    Log  ${out}
    Should Not Contain  ${out}  %{VCH-NAME}
    ${out}=  Run  govc ls host/*/Resources/*
    Log  ${out}
    Should Not Contain  ${out}  %{VCH-NAME}
    ${out}=  Run  govc datastore.ls
    Log  ${out}
    Should Not Contain  ${out}  %{VCH-NAME}
