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
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${esx1}  ${esx4-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${esx2}  ${esx5-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    Set Global Variable  @{list}  ${esx1}  ${esx2}

    Create a Simple VC Cluster  datacenter1  cls1

    Log To Console  Create cluster2 on the VC
    ${out}=  Run  govc cluster.create cls2
    Should Be Empty  ${out}
    ${out}=  Run  govc cluster.add -hostname=${esx4-ip} -username=root -dc=datacenter1 -cluster=cls2 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create cluster3 on the VC
    ${out}=  Run  govc cluster.create cls3
    Should Be Empty  ${out}
    ${out}=  Run  govc cluster.add -hostname=${esx5-ip} -username=root -dc=datacenter1 -cluster=cls3 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Regression Tests

    Cleanup VIC Appliance On Test Server
