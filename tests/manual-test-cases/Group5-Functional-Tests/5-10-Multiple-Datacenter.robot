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
Documentation  Test 5-10 - Multiple Datacenters
Resource  ../../resources/Util.robot
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${esx2}  ${esx2-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    ${esx3}  ${esx4}  ${esx5}  ${vc}  ${esx3-ip}  ${esx4-ip}  ${esx5-ip}  ${vc-ip}=  Create a Simple VC Cluster  datacenter1  cls1

    Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${esx4}  ${esx5}  ${vc}

    Log To Console  Create datacenter2 on the VC
    ${out}=  Run  govc datacenter.create datacenter2
    Should Be Empty  ${out}
    ${out}=  Run  govc host.add -hostname=${esx1-ip} -username=root -dc=datacenter2 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create datacenter3 on the VC
    ${out}=  Run  govc datacenter.create datacenter3
    Should Be Empty  ${out}
    ${out}=  Run  govc host.add -hostname=${esx2-ip} -username=root -dc=datacenter3 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Set Environment Variable  TEST_DATACENTER  /datacenter1
    Set Environment Variable  GOVC_DATACENTER  /datacenter1
    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Regression Tests

    Cleanup VIC Appliance On Test Server
