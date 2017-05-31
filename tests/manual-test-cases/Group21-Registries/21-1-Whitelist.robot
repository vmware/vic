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
Documentation  Test 21-01 - Whitelist
Resource  ../../resources/Util.robot
Resource  ../../resources/Harbor-Util.robot
Suite Setup  Setup Harbor
Suite Teardown  Harbor Test Cleanup
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Keywords ***
Setup Harbor
    Set Test Environment Variables

    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi as Harbor is only supported on VC
    
    # Install a Harbor server with HTTPS a Harbor server with HTTP
    Install Harbor To Test Server  protocol=https  name=harbor-https
    Set Environment Variable  HTTPS_HARBOR_IP  %{HARBOR-IP}

    Install Harbor To Test Server  protocol=http  name=harbor-http
    Set Environment Variable  HTTP_HARBOR_IP  %{HARBOR-IP}

    Get HTTPS Harbor Certificate

Harbor Test Cleanup
    ${out}=  Run  govc vm.destroy harbor-http
    ${out}=  Run  govc vm.destroy harbor-https
    Log To Console  Cleaning up Harbor servers

Get HTTPS Harbor Certificate
    [Arguments]  ${HARBOR_IP}=%{HTTPS_HARBOR_IP}
    # Get the certificates from the HTTPS server
    ${out}=  Run  wget --tries=10 --connect-timeout=10 --auth-no-challenge --no-check-certificate --user admin --password %{TEST_PASSWORD} https://${HARBOR_IP}/api/systeminfo/getcert
    Log  ${out}
    Move File  getcert  ./ca.crt


*** Test Cases ***
Basic Whitelisting
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi as Harbor is only supported on VC
    
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Install VCH with registry CA for whitelisted registry
    ${output}=  Install VIC Appliance To Test Server  vol=default --whitelist-registry=%{HTTPS_HARBOR_IP} --registry-ca=./ca.crt
    Should Contain  ${output}  Secure registry %{HTTPS_HARBOR_IP} confirmed
    Should Contain  ${output}  Whitelist registries =

    # Check docker info for whitelist info
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Registry Whitelist Mode: enabled
    Should Contain  ${output}  Whitelisted Registries:
    Should Contain  ${output}  Registry: registry-1.docker.io

    # Try to login and pull from the HTTPS whitelisted registry (should succeed)
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login -u admin -p %{TEST_PASSWORD} %{HTTPS_HARBOR_IP}
    Should Contain  ${output}  Succeeded
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull %{HTTPS_HARBOR_IP}/library/photon:1.0
    Should Be Equal As Integers  ${rc}  0

    # Try to login and pull from the HTTPS whitelisted registry with :443 tacked on at the end (should succeed)
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login -u admin -p %{TEST_PASSWORD} %{HTTPS_HARBOR_IP}:443
    Should Contain  ${output}  Succeeded
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull %{HTTPS_HARBOR_IP}:443/library/photon:1.0
    Should Be Equal As Integers  ${rc}  0

    # Try to login and pull from docker hub (should fail)
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Access denied to unauthorized registry
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull victest/busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Access denied to unauthorized registry

    Cleanup VIC Appliance On Test Server
