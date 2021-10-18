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
Documentation  Test 5-8 - DRS
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  DRS Setup
Suite Teardown  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Keywords ***
DRS Setup
    [Timeout]  60 minutes 
    ${name}=  Evaluate  'vic-iscsi-cluster-' + str(random.randint(1000,9999))  modules=random
    Log To Console  Create a new simple vc cluster with spec vic-cluster-2esxi-iscsi.rb...
    ${out}=  Deploy Nimbus Testbed  spec=vic-cluster-2esxi-iscsi.rb  args=--noSupportBundles --plugin testng --vcvaBuild "${VC_VERSION}" --esxBuild "${ESX_VERSION}" --testbedName vic-iscsi-cluster --runName ${name}
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
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  /dc1/host/cls
    Set Environment Variable  TEST_TIMEOUT  30m

*** Test Cases ***
Test
    Log To Console  Disable DRS on the cluster    
    ${rc}  ${out}=  Run And Return Rc And Output  govc cluster.change -drs-enabled=false %{TEST_RESOURCE}
    Should Be Empty  ${out}
    Should Be Equal As Integers  ${rc}  0

    Log To Console  \nStarting test...
    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled %{TEST_RESOURCE}
    Should Be Empty  ${out}

    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests
