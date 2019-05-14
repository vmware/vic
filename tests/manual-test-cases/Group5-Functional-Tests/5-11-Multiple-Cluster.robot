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
Documentation  Test 5-11 - Multiple Clusters
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Multiple Cluster Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}
Test Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Multiple Cluster Setup
    [Timeout]  60 minutes
    ${name}=  Evaluate  'vic-multi-cluster-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nStarting testbed deploy...
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  spec=vic-multi-cls.rb  args=--noSupportBundles --plugin testng --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-multi-cls --runName ${name}
    Log  ${out}
    Should Contain  ${out}  "deployment_result"=>"PASS"

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    Log  ${vc-ip}
    ${esx0_ip}=  Get IP  ${name}.esx.0
    Log  ${esx0_ip}
    ${pod}=  Fetch Pod  ${name}
    Log  ${pod}
    Close Connection
    
    Set Suite Variable  ${nimbus_pod}  ${pod}
    Set Suite Variable  ${testbedname}  ${name}
    
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
    Set Environment Variable  TEST_TIMEOUT  15m

    # Get one of the hosts in the cluster we want and make sure we use the correct  datastore

    ${rc}  ${test_resource}=  Run And Return Rc And Output  govc host.info ${esx0_ip} | grep Path | awk -F: '{print $2}'
    Log  ${test_resource}
    ${test_resource}=  Remove String  ${test_resource}  /${esx0_ip}
    ${test_resource}=  Strip String  ${test_resource}
    Log  ${test_resource}
    Set Environment Variable  TEST_RESOURCE  ${test_resource}
    
    ${rc}  ${datastore}=  Run And Return Rc And Output  govc host.info -host.ip=${esx0_ip} -json | jq -r '.HostSystems[].Config.FileSystemVolume.MountInfo[].Volume | select(.Type == "VMFS" and .Local == false) | .Name'
    Log  ${datastore}
    Set Environment Variable  TEST_DATASTORE  ${datastore}

*** Test Cases ***
Test
    Log To Console  \nStarting test...

    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests
