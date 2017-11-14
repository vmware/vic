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
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
ESXi Restart Setup
    Log To Console  \nStarting test...
    ${esx}  ${esx-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Global Variable  @{list}  ${esx}
    Set Suite Variable  ${esx}  ${esx}
    Set Suite Variable  ${esx-ip}  ${esx-ip}
    Set Environment Variable  GOVC_URL  root:${NIMBUS_ESX_PASSWORD}@${esx-ip}
    Set Environment Variable  TEST_URL_ARRAY  ${esx-ip}
    Set Environment Variable  TEST_URL  ${esx-ip}
    Set Environment Variable  TEST_USERNAME  root
    Set Environment Variable  TEST_PASSWORD  ${NIMBUS_ESX_PASSWORD}
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_TIMEOUT  30m
    Set Environment Variable  HOST_TYPE  ESXi
    Remove Environment Variable  TEST_DATACENTER
    Remove Environment Variable  TEST_RESOURCE
    Remove Environment Variable  BRIDGE_NETWORK
    Remove Environment Variable  PUBLIC_NETWORK

*** Test Cases ***
Test
    ${out}=  Run  bin/vic-machine-linux create -t ${esx-ip} -u root -p ${NIMBUS_ESX_PASSWORD} --force --image-store=datastore1 --no-tls
    #${out}=  Run  bin/vic-machine-linux create --debug 1 --name=VCH-0-6576 --target=${esx-ip} --user=root --image-store=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=${NIMBUS_ESX_PASSWORD} --force=true --bridge-network=VCH-0-6576-bridge --public-network='VM Network' --timeout 30m --insecure-registry harbor.ci.drone.local --volume-store=datastore1/VCH-0-6576-VOL:default --container-network='VM Network':public --no-tlsverify
    Get Docker Params  ${out}  ${false}
    #Install VIC Appliance To Test Server
    #Run Regression Tests

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
