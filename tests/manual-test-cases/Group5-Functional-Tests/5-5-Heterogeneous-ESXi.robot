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
Documentation  Test 5-1 - Distributed Switch
Resource  ../../resources/Util.robot
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${vc}=  Evaluate  'VC-' + str(random.randint(1000,9999))  modules=random
    ${pid-vc}=  Deploy Nimbus vCenter Server Async  ${vc}
    Set Global Variable  @{list}  '%{NIMBUS_USER}'-${vc}

    Run Keyword And Ignore Error  Cleanup Nimbus PXE folder  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'
    Append To List  ${list}  ${esx1}
    
    Run Keyword And Ignore Error  Cleanup Nimbus PXE folder  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'
    ${esx2}  ${esx2-ip}=  Deploy Nimbus ESXi Server  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'  3029944
    Append To List  ${list}  ${esx2}

    Run Keyword And Ignore Error  Cleanup Nimbus PXE folder  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'
    ${esx3}  ${esx3-ip}=  Deploy Nimbus ESXi Server  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'  4240417
    Append To List  ${list}  ${esx3}

    # Finish vCenter deploy
    ${output}=  Wait For Process  ${pid-vc}
    Should Contain  ${output.stdout}  Overall Status: Succeeded

    Open Connection  '%{NIMBUS_GW}'
    Wait Until Keyword Succeeds  2 min  30 sec  Login  '%{NIMBUS_USER}'  '%{NIMBUS_PASSWORD}'
    ${vc-ip}=  Get IP  ${vc}
    Close Connection

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${vc-ip}

    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ha-datacenter
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create cls
    Should Be Empty  ${out}

    Log To Console  Add ESX host to the VC
    Add Host To VCenter  ${esx1-ip}  root  ha-datacenter  e2eFunctionalTest
    Add Host To VCenter  ${esx2-ip}  root  ha-datacenter  e2eFunctionalTest
    ${vc-ver}=  Run  govc about | grep Version:
    ${vc-ver}=  Fetch From Right  ${vc-ver}  ${SPACE}
    Run Keyword If  '${vc-ver}' == '6.5.0'  Add Host To VCenter  ${esx3-ip}  root  ha-datacenter  e2eFunctionalTest

    Create A Distributed Switch  ha-datacenter

    Create Three Distributed Port Groups  ha-datacenter

    Add Host To Distributed Switch  /ha-datacenter/host/cls

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /ha-datacenter/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Regression Tests
