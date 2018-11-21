# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 23-13 - VCH-tdnf
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Gpgcheck Being Enabled
    Enable VCH SSH
    ${rc}  ${output}=  Run And Return Rc and Output  sshpass -p %{TEST_PASSWORD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@%{VCH-IP} grep gpgcheck=1 /etc/yum.repos.d/photon.repo
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Tdnf Install Which
    ${rc}  ${output}=  Run And Return Rc and Output  sshpass -p %{TEST_PASSWORD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@%{VCH-IP} tdnf install -y which
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
