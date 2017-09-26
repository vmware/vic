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
Documentation  Test 12-01 - Delete
Resource  ../../resources/Util.robot
Suite Setup  Install VIC 0.6.0 to Test Server
Test Teardown  Run Keyword If Test Failed  Clean up VIC Appliance And Local Binary

*** Keywords ***
Clean up VIC Appliance And Local Binary
    Cleanup VIC Appliance On Test Server
    Run  rm -rf vic.tar.gz vic

Install VIC 0.6.0 to Test Server
    Log To Console  \nDownloading vic 4890 from gcp...
    ${rc}  ${output}=  Run And Return Rc And Output  wget https://storage.googleapis.com/vic-engine-builds/vic_4890.tar.gz -O vic.tar.gz
    ${rc}  ${output}=  Run And Return Rc And Output  tar zxvf vic.tar.gz
    Set Test Environment Variables

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  ./vic/vic-machine-linux create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=./vic/appliance.iso --bootstrap-iso=./vic/bootstrap.iso --password=%{TEST_PASSWORD} --bridge-network=%{BRIDGE_NETWORK} --external-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Installer completed successfully
    Get 0.6.0 VIC Docker Params  ${output}  false
    Log To Console  Installer completed successfully: %{VCH-NAME}

Get 0.6.0 VIC Docker Params
    # Get VCH docker params e.g. "-H 192.168.218.181:2376 --tls"
    [Arguments]  ${output}  ${certs}
    @{output}=  Split To Lines  ${output}
    :FOR  ${item}  IN  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  DOCKER_HOST=
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}

    # Ensure we start from a clean slate with docker env vars
    Remove Environment Variable  DOCKER_HOST  DOCKER_TLS_VERIFY  DOCKER_CERT_PATH

    # Split the log log into pieces, discarding the initial log decoration, and assign to env vars
    ${logdeco}  ${vars}=  Split String  ${line}  ${SPACE}  1
    ${vars}=  Split String  ${vars}
    :FOR  ${var}  IN  @{vars}
    \   ${varname}  ${varval}=  Split String  ${var}  =
    \   Set Environment Variable  ${varname}  ${varval}

    ${dockerHost}=  Get Environment Variable  DOCKER_HOST

    @{hostParts}=  Split String  ${dockerHost}  :
    ${ip}=  Strip String  @{hostParts}[0]
    ${port}=  Strip String  @{hostParts}[1]
    Set Environment Variable  VCH-IP  ${ip}
    Set Environment Variable  VCH-PORT  ${port}

    ${proto}=  Set Variable If  ${port} == 2376  "https"  "http"
    Set Suite Variable  ${proto}

    Run Keyword If  ${port} == 2376  Set Environment Variable  VCH-PARAMS  -H ${dockerHost} --tls
    Run Keyword If  ${port} == 2375  Set Environment Variable  VCH-PARAMS  -H ${dockerHost}


    :FOR  ${index}  ${item}  IN ENUMERATE  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  http
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  Published ports can be reached at
    \   ${idx} =  Evaluate  ${index} + 1
    \   Run Keyword If  '${status}' == 'PASS'  Set Environment Variable  EXT-IP  @{output}[${idx}]

*** Test Cases ***
Delete VCH with new vic-machine
    Log To Console  \nRunning docker pull busybox...
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.11 %{VCH-PARAMS} pull busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${name}=  Generate Random String  15
    ${rc}  ${container-id}=  Run And Return Rc And Output  docker1.11 %{VCH-PARAMS} create --name ${name} busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container-id}  Error
    Set Suite Variable  ${containerName}  ${name}

    # Get VCH uuid and container VM uuid, to check if resources are removed correctly
    Run Keyword And Ignore Error  Gather Logs From Test Server
    ${uuid}=  Run  govc vm.info -json\=true %{VCH-NAME} | jq -r '.VirtualMachines[0].Config.Uuid'
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME}
    Should Contain  ${ret}  is different than installer version

    # Delete with force
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME} --force
    Should Contain  ${ret}  Completed successfully
    Should Not Contain  ${ret}  delete failed

    # Check VM is removed
    ${ret}=  Run  govc vm.info -json=true ${containerName}-*
    Should Contain  ${ret}  {"VirtualMachines":null}
    ${ret}=  Run  govc vm.info -json=true %{VCH-NAME}
    Should Contain  ${ret}  {"VirtualMachines":null}

    # Check resource pool is removed
    ${ret}=  Run  govc pool.info -json=true host/*/Resources/%{VCH-NAME}
	Should Contain  ${ret}  {"ResourcePools":null}
    Run  rm -rf vic.tar.gz vic

    Run Keyword And Ignore Error  Cleanup VCH Bridge Network  %{VCH-NAME}
