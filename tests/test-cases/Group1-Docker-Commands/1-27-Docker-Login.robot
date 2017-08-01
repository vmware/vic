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
Documentation  Test 1-27 - Docker Login
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Cleanup Harbor and VIC Appliance
    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server
    Run Keyword And Continue On Failure  Run  govc vm.destroy 19-4-harbor
    

*** Test Cases ***
Docker login and pull from docker.io
	[Teardown]  Cleanup Harbor and VIC Appliance	
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull victest/busybox
    Should Be Equal As Integers  ${rc}  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull victest/public-hello-world
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=incorrectPassword
    Should Contain  ${output}  incorrect username or password
    Should Be Equal As Integers  ${rc}  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=%{REGISTRY_PASSWORD}
    Should Contain  ${output}  Login Succeeded
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull victest/busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logout
    Should Be Equal As Integers  ${rc}  0

    ${ip}=  Install Harbor To Test Server  name=19-4-harbor  protocol=https

    # Install Harbor Self Signed Cert
    ${out}=  Run  wget --tries=10 --connect-timeout=10 --auth-no-challenge --no-check-certificate --user %{TEST_USERNAME} --password ${TEST_PASSWORD} https://%{ip}/api/systeminfo/getcert
    Log  ${out}
    ${out}=  Run  mkdir -p /etc/docker/certs.d/%{ip}
    Move File  getcert  /etc/docker/certs.d/%{ip}/ca.crt
    ${out}=  Run  systemctl daemon-reload
    ${out}=  Run  systemctl restart docker

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${ip}
    Should Contain  ${output}  Login Succeeded
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=%{TEST_USERNAME} --password=%{TEST_PASSWORD}1 ${ip}
    Should Contain  ${output}  incorrect username or password
    Should Be Equal As Integers  ${rc}  1

