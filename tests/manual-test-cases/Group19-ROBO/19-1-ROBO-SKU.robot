# Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
Documentation  Test 19-1 - ROBO SKU
Resource  ../../resources/Util.robot
Suite Setup  Wait Until Keyword Succeeds  10x  10m  ROBO SKU Setup
#Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
ROBO SKU Setup
    [Timeout]    110 minutes
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${name}=  Evaluate  'vic-vsan-' + str(random.randint(1000,9999))  modules=random
    Set Suite Variable  ${user}  %{NIMBUS_USER}
    Log To Console  Deploying testbed
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --plugin testng --vcfvtBuildPath /dbc/pa-dbc1111/mhagen/ --noSupportBundles --vcvaBuild ${VC_VERSION} --esxPxeDir ${ESX_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-vsan-simple-pxeBoot-vcva --runName ${name}
    Should Contain  ${out}  "deployment_result"=>"PASS"
    Log To Console  Retrieving IP for ${user}-${name}.vcva-${VC_VERSION}
    ${vc-ip}=  Get IP  ${name}.vcva-${VC_VERSION}
    Log To Console  ${user}-${name}.vcva-${VC_VERSION} IP: ${vc-ip}

    Set Suite Variable  @{list}  ${user}-${name}.vcva-${VC_VERSION}  ${user}-${name}.esx.0  ${user}-${name}.esx.1  ${user}-${name}.esx.2  ${user}-${name}.esx.3  ${user}-${name}.nfs.0  ${user}-${name}.iscsi.0

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23
    Set Environment Variable  GOVC_INSECURE  1

    Add Host To Distributed Switch  /vcqaDC/host/cls

    Log To Console  Enable HA
    ${out}=  Run  govc cluster.change -ha-enabled cls
    Should Be Empty  ${out}

    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  vsanDatastore
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  15m

    # TODO - add ROBO license(s) to the test suite
#    Add Vsphere License  %{ROBO_LICENSE}
#
#    :FOR  ${IDX}  IN RANGE  0  4
#    \   Log To Console  Getting IP for ${user}-${name}.esx.${IDX}
#    \   ${esx-ip}=  Get IP  ${user}-${name}.esx.${IDX}
#    \   Log To Console  Applying ROBO license to host ${esx-ip}
#    \   Assign Vsphere License  %{ROBO_LICENSE}  ${esx-ip}

*** Test Cases ***
Test
    Log To Console  VIC does not support ROBO SKU yet, waiting on a valid license with serial support for this to work
    #Log To Console  \nStarting test...
    #Install VIC Appliance To Test Server
    #Run Regression Tests
