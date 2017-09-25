# Copyright 2017 VMware, Inc. All Rights Reserved.
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
Documentation  Test OVA-installer
Resource  ../../resources/Util.robot
Library  SSHLibrary
#Suite Setup  Setup Suite Environment
#Suite Teardown  Cleanup Suite Environment

*** Variables ***
${dhcp}  True
${ova_options}  --acceptAllEulas --name=vic-unified-ova-integration-test --X:injectOvfEnv --X:enableHiddenProperties -st=OVA --powerOn --X:waitForIp --noSSLVerify=true -ds=datastore1 -dm=thin  --prop:appliance.root_pwd="password"   --prop:appliance.permit_root_login=True   --prop:management_portal.port=8282   --prop:registry.port=443
${esx_number}  1

*** Keywords ***
Setup Suite Environment
    @{vms}=  Create A Simple VC Cluster  datacenter  cls  ${esx_number}  ${network}
    ${vm_num}=  Evaluate  ${esx_number}+1
    ${vm-names}=  Get Slice From List  ${vms}  0  ${vm_num}
    ${vm-ips}=  Get Slice From List  ${vms}  ${vm_num}
    Set Suite Variable  ${vc-ip}  @{vm-ips}[${esx_number}]
    Set Suite Variable  ${vm-names}  ${vm-names}
    Set Suite Variable  ${vm-ips}  ${vm-ips}
    GOVC Use VC Environment

Cleanup Suite Environment
        Nimbus Cleanup  ${vm-names}

GOVC Use VC Environment
        Set Environment Variable  GOVC_URL  ${vc-ip}
        Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
        Set Environment Variable  GOVC_PASSWORD  Admin\!23

GOVC Use ESXi Environment
        Set Environment Variable  GOVC_URL  @{vm-ips}[0]
        Set Environment Variable  GOVC_USERNAME  root
        Set Environment Variable  GOVC_PASSWORD

*** Test Cases ***
Deploy VIC-OVA To Test Server
  [Arguments]  ${dhcp}=True  ${protocol}=http  ${user}=%{TEST_USERNAME}  ${password}=%{TEST_PASSWORD}  ${host}=%{TEST_URL_ARRAY}  ${datastore}=%{TEST_DATASTORE}  ${cluster}=%{TEST_RESOURCE}  ${datacenter}=%{TEST_DATACENTER}
  Set Global Variable  ${ova_path}  bin/vic-v1.2.0-d0ea01c2.ova

  Log To Console  \nStarting to deploy unified-ova to test server...
  ${out}=  Run Keyword If  ${dhcp}  Run  ovftool --datastore=${datastore} ${ova_options} --net:"Network"="vm-network" ${ova_path} 'vi://Administrator@vSphere.local:Admin\!23@${host}/${datacenter}/host/${cluster}'
  Log To Console  ${out}

  Should Contain  ${out}  Received IP address:
  Should Not Contain  ${out}  None

  ${out}=  Split To Lines  ${out}
  :FOR  ${line}  IN  @{out}
  \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  Received IP address:
  \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
  \   Log To Console  ${ip}
  \   Exit For Loop If  ${status}

   Set Environment Variable  Appliance-IP  ${ip}
   Sleep  30s

   Log To Console  Waiting for Getting Started Page to Come Up...
   :FOR  ${i}  IN RANGE  20
    \  ${rc}=  Run And Return Rc  curl -sk https://%{Appliance-IP}:9443 -XPOST -F target=${host} -F user=Administrator@vsphere.local -F password=Admin!23
    \  Should Be Equal As Integers  ${rc}  0

   Log To Console  ssh into appliance...
    ${out}=  Run  sshpass -p password ssh -o StrictHostKeyChecking\=no root@%{Appliance-IP}

    Open Connection  %{Appliance-IP}
    Wait Until Keyword Succeeds  2 min  30 sec  Login  root  password

    Log To Console  status of harbor...
    ${out}=  Execute Command  systemctl status harbor
    Should Contain  ${out}  Active: active (running)

    Log To Console  status of admiral...
    ${out}=  Execute Command  systemctl status admiral
    Should Contain  ${out}  Active: active (running)

    Log To Console  status of fileserver...
    ${out}=  Execute Command  systemctl status fileserver
    Should Contain  ${out}  Active: active (running)

    Log To Console  status of engine_installer...
    ${out}=  Execute Command  systemctl status engine_installer
    Should Contain  ${out}  Active: active (running)

    Close connection

    [Teardown]  Cleanup VIC-OVA On Test Server
