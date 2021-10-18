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
Documentation  Test 5-5 - Heterogeneous ESXi
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Heterogenous ESXi Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}
Force Tags  hetero

*** Keywords ***
Heterogenous ESXi Setup
    [Timeout]  60 minutes
    ${name}=  Evaluate  'vic-hetegeneous-' + str(random.randint(1000,9999))  modules=random
    Log To Console  Create a new simple vc cluster with spec vic-hetegeneous.rb...
    ${out}=  Deploy Nimbus Testbed  spec=vic-hetergeneous-esxi.rb  args=--noSupportBundles --plugin testng --vcvaBuild "${VC_VERSION}" --esxBuild "${ESX_VERSION}" --testbedName vic-hetegeneous --runName ${name}
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
    ${out}=  Run  govc dvs.create -product-version 5.5.0 -dc=dc1 test-ds
    Should Contain  ${out}  OK

    Create Three Distributed Port Groups  dc1

    Add Host To Distributed Switch  /dc1/host/cls1

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_RESOURCE  cls1
    Set Environment Variable  TEST_TIMEOUT  30m
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0 

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests
