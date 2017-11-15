# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 22-1 - ESXi Restart
Resource  ../../resources/Util.robot
Suite Setup  ESXi Restart Setup
#Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
ESXi Restart Setup
    Log To Console  \nStarting test...
    ${name}=  Evaluate  'ESX-' + str(random.randint(1000,9999)) + str(time.clock())  modules=random,time
    Set Suite Variable  ${esx}  %{NIMBUS_USER}-${name}
    Log To Console  \nDeploying Nimbus ESXi server: ${name}
    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  2 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    ${out}=  Execute Command  nimbus-esxdeploy ${name} --disk=48000000 --ssd=24000000 --memory=8192 --lease=1 --nics 2 ${ESX_VERSION}
    Log  ${out}
    # Make sure the deploy actually worked
    ${status}=  Run Keyword And Return Status  Should Contain  ${out}  To manage this VM use
    Set Global Variable  @{list}  ${esx}
    Close Connection

    # Now grab the IP address and return the name and ip for later use
    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  IP is
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    @{gotIP}=  Split String  ${line}  ${SPACE}
    ${ip}=  Remove String  @{gotIP}[5]  ,
    Set Suite Variable  ${esx-ip}  ${ip}

    Set Environment Variable  GOVC_URL  root:@${esx-ip}
    Set Environment Variable  TEST_URL_ARRAY  ${esx-ip}
    Set Environment Variable  TEST_URL  ${esx-ip}
    Set Environment Variable  TEST_USERNAME  root
    Set Environment Variable  TEST_PASSWORD  ''
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_TIMEOUT  30m
    Set Environment Variable  HOST_TYPE  ESXi
    Remove Environment Variable  TEST_DATACENTER
    Remove Environment Variable  TEST_RESOURCE
    Remove Environment Variable  BRIDGE_NETWORK
    Remove Environment Variable  PUBLIC_NETWORK

*** Test Cases ***
Test
    ${name}=  Evaluate  'VCH-' + str(random.randint(1000,9999))  modules=random
    Set Environment Variable  VCH-NAME  ${name}
    ${out}=  Run  bin/vic-machine-linux create -t ${esx-ip} -u root -p '' --force --image-store=datastore1 --name %{VCH-NAME} --no-tls
    #${out}=  Run  bin/vic-machine-linux create --debug 1 --name=VCH-0-6576 --target=${esx-ip} --user=root --image-store=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=${NIMBUS_ESX_PASSWORD} --force=true --bridge-network=VCH-0-6576-bridge --public-network='VM Network' --timeout 30m --insecure-registry harbor.ci.drone.local --volume-store=datastore1/VCH-0-6576-VOL:default --container-network='VM Network':public --no-tlsverify
    Get Docker Params  ${out}  ${false}
    #Install VIC Appliance To Test Server
    Run Regression Tests

    Reset Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  ${esx}
    Wait Until vSphere Is Powered On

    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on %{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    Wait Until VM Powers On  %{VCH-NAME}

    Log To Console  Getting VCH IP ...
    ${new-vch-ip}=  Get VM IP  %{VCH-NAME}
    Log To Console  New VCH IP is ${new-vch-ip}
    Replace String  %{VCH-PARAMS}  %{VCH-IP}  ${new-vch-ip}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    Run Regression Tests
