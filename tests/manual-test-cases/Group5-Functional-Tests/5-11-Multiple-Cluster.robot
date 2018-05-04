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
Suite Setup  Wait Until Keyword Succeeds  10x  10m  Multiple Cluster Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  '*5-11-multiple-cluster*'

*** Keywords ***
# Insert elements from dict2 into dict1, overwriting conflicts in dict1 & returning new dict
Combine Dictionaries
    [Arguments]  ${dict1}  ${dict2}
    ${dict2keys}=  Get Dictionary Keys  ${dict2}
    :FOR  ${key}  IN  @{dict2keys}
    \    ${elem}=  Get From Dictionary  ${dict2}  ${key}
    \    Set To Dictionary  ${dict1}  ${key}  ${elem}
    [Return]  ${dict1}

Multiple Cluster Setup
    [Timeout]    110 minutes
    Run Keyword And Ignore Error  Nimbus Cleanup  '*5-11-multiple-cluster*'  ${false}
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --plugin testng --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-multiple-cluster --testbedSpecRubyFile /dbc/pa-dbc1111/mhagen/nimbus-testbeds/testbeds/vic-multiple-cluster.rb --runName 5-11-multiple-cluster
    Log  ${out}

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  5-11-multiple-cluster.vc.0
    Close Connection
    
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
    Set Environment Variable  TEST_DATASTORE  local-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  15m

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    Custom Testbed Keepalive  /dbc/pa-dbc1111/mhagen

    Install VIC Appliance To Test Server  certs=${false}  vol=default
    Run Regression Tests
    Cleanup VIC Appliance On Test Server
