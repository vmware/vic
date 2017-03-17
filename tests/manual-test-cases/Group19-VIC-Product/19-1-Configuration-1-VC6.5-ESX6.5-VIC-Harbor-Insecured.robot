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
Documentation  Test 19-1 - Configuration 1 VC6.5 ESX6.5 VIC Harbor Insecured
Resource  ../../resources/Util.robot
Suite Setup  19-1-Setup
Suite Teardown  19-1-Teardown

*** Keywords ***
19-1-Setup
    ${esx1}  ${esx2}  ${esx3}  ${vc}  ${vc-ip}=  Create a Simple VC Cluster
    Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${vc}
    Install Harbor To Test Server  name=19-1-harbor
    Install VIC Appliance To Test Server  vol=default --insecure-registry %{HARBOR_IP}

19-1-Teardown
    Cleanup VIC Appliance On Test Server
    ${out}=  Run  govc vm.destroy 19-1-harbor
    Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    Log Into Harbor  user=admin
    Create A New Project  vic-harbor
    Close All Browsers
    Restart Docker With Insecure Registry Option
    ${out}=  Run  docker login -u admin -p %{TEST_PASSWORD} %{HARBOR_IP}
    ${out}=  Run  docker pull busybox
    ${out}=  Run  docker tag busybox %{HARBOR_IP}/vic-harbor/busybox
    ${out}=  Run  docker push %{HARBOR_IP}/vic-harbor/busybox