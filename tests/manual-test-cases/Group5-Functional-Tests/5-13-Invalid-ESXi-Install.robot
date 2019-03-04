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
Documentation  Test 5-13 - Invalid ESXi Install
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Invalid ESXi Install Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Keywords ***
Invalid ESXi Install Setup
    [Timeout]  60 minutes
    ${name}=  Evaluate  'vic-iscsi-cluster-' + str(random.randint(1000,9999))  modules=random
    Log To Console  Create a new simple vc cluster with spec vic-cluster-2esxi-iscsi.rb...
    ${out}=  Deploy Nimbus Testbed  spec=vic-cluster-2esxi-iscsi.rb  args=--noSupportBundles --plugin testng --vcvaBuild "${VC_VERSION}" --esxBuild "${ESX_VERSION}" --testbedName vic-iscsi-cluster --runName ${name}
    Log  ${out}
    Should Contain  ${out}  "deployment_result"=>"PASS"
    Log To Console  Finished creating cluster ${name}
    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    Log  ${vc-ip}
    ${esx1-ip}=  Get IP  ${name}.esx.1
    Log  ${esx1-ip}
    ${pod}=  Fetch Pod  ${name}
    Log  ${pod}
    Close Connection

    Set Suite Variable  ${nimbus_pod}  ${pod}
    Set Suite Variable  ${testbedname}  ${name}
    Set Suite Variable  ${esx1-ip}  ${esx1-ip}
    Set Suite Variable  ${vc-ip}  ${vc-ip}     

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

*** Test Cases ***
Test
    ${out}=  Run  bin/vic-machine-linux create --target ${esx1-ip} --user root --password e2eFunctionalTest --no-tls --name VCH-invalid-test --force --timeout 30m
    Should Contain  ${out}  Target is managed by vCenter server \\"${vc-ip}\\", please change --target to vCenter server address or select a standalone ESXi
