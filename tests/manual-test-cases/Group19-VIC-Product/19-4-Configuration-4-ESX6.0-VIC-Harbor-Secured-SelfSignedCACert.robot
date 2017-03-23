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
Documentation  Test 19-4 - Configuration 4 ESX 6.0 VIC Harbor Secured Self Signed
Resource  ../../resources/Util.robot
Suite Setup  19-4-Setup
Suite Teardown  19-4-Teardown

*** Keywords ***
19-4-Setup
    Set Environment Variable  ESX_VERSION  3620759
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Environment Variable  TEST_URL_ARRAY  ${esx1-ip}
    Set Environment Variable  TEST_USERNAME  root
    Set Environment Variable  TEST_PASSWORD  e2eFunctionalTest
    Set Environment Variable  BRIDGE_NETWORK  network
    Set Environment Variable  PUBLIC_NETWORK  'VM Network'
    Set Environment Variable  TEST_DATACENTER  ${SPACE}
    Set Environment Variable  TEST_TIMEOUT  30m
    
    Set Suite Variable  ${ESX1}  ${esx1}
    Set Suite Variable  ${ESX1-IP}  ${esx1-ip}
    Set Global Variable  @{list}  ${esx1}

    Install Harbor To Test Server  name=19-4-harbor
    Install Harbor Self Signed Cert
    Install VIC Appliance To Test Server  vol=default --registry-ca=/etc/docker/certs.d/%{HARBOR_IP}/ca.crt

19-4-Teardown
    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server
    ${out}=  Run Keyword And Continue On Failure  Run  govc vm.destroy 19-4-harbor
    Run Keyword And Continue On Failure  Nimbus Cleanup  ${list}

*** Test Cases ***
Pos001
    Log Into Harbor  user=admin
    Create A New User  user1  user1@email.com  user1  Password123
    Create A New User  user2  user2@email.com  user2  Password123 
    Create A New Project  vic-harbor  ${false}
    Close All Browsers
    ${out}=  Run  docker login -u admin -p %{TEST_PASSWORD} %{HARBOR_IP}
    ${out}=  Run  docker pull busybox
    ${out}=  Run  docker tag busybox %{HARBOR_IP}/vic-harbor/busybox
    ${out}=  Run  docker push %{HARBOR_IP}/vic-harbor/busybox
    
    ${out}=  Run  docker %{VCH-PARAMS} pull %{HARBOR_IP}/vic-harbor/busybox