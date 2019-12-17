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
Documentation  Test 5-17 - FC Datastore
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  FC Datastore Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Variables ***
${NIMBUS_LOCATION}  sc
${NIMBUS_LOCATION_FULL}  NIMBUS_LOCATION=${NIMBUS_LOCATION}

*** Keywords ***
FC Datastore Setup
    [Timeout]    110 minutes
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${name}=  Evaluate  'vic-fc-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  args=--plugin testng --customizeTestbed '/esx desiredPassword=e2eFunctionalTest' --noSupportBundles --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vcqa-sdrs-fc-fullInstall-vcva --runName ${name}
    Should Contain  ${out}  "deployment_result"=>"PASS"

    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${pod}=  Fetch Pod  ${name}
    Log  ${pod}
    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.vcva-${VC_VERSION} | grep %{NIMBUS_PERSONAL_USER}-${name}.vcva-${VC_VERSION}
    ${vc-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.esx.0 | grep %{NIMBUS_PERSONAL_USER}-${name}.esx.0
    ${esx0-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.esx.1 | grep %{NIMBUS_PERSONAL_USER}-${name}.esx.1
    ${esx1-ip}=  Fetch From Right  ${out}  ${SPACE}
    Close Connection

    Set Suite Variable  ${nimbus_pod}  ${pod}
    Set Suite Variable  ${testbedname}  ${name}

    Set Environment Variable  GOVC_URL  ${esx0-ip}
    Set Environment Variable  GOVC_USERNAME  root
    Set Environment Variable  GOVC_PASSWORD  e2eFunctionalTest
    Run  govc host.esxcli network firewall set -e false
    Set Environment Variable  GOVC_URL  ${esx1-ip}
    Run  govc host.esxcli network firewall set -e false

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Create A Distributed Switch  vcqaDC

    Create Three Distributed Port Groups  vcqaDC

    Add Host To Distributed Switch  /vcqaDC/host/cls

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cls
    Should Be Empty  ${out}

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
FC Datastore
    Install VIC Appliance To Test Server
    Run Regression Tests
