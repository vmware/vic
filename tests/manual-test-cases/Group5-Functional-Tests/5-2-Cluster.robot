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
Documentation  Test 5-2 - Cluster
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Cluster Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Keywords ***
Cluster Setup
    [Timeout]  60 minutes
    ${name}=  Evaluate  'vic-iscsi-cluster-' + str(random.randint(1000,9999))  modules=random
    Log To Console  Create a new simple vc cluster with spec vic-cluster-2esxi-iscsi.rb...
    ${out}=  Deploy Nimbus Testbed  spec=vic-cluster-2esxi-iscsi.rb  args=--noSupportBundles --plugin testng --vcvaBuild "${VC_VERSION}" --esxBuild "${ESX_VERSION}" --testbedName vic-iscsi-cluster --runName ${name}
    Log  ${out}
    Log To Console  Finished creating cluster ${name}
    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    Log  ${vc-ip}
    ${pod}=  Fetch Pod  ${name}
    Log  ${pod}
    Close Connection
    
    Set Suite Variable  ${nimbus_pod}  ${pod}
    Set Suite Variable  ${testbedname}  ${name}

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${vc-ip}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

Verify LS Output For Busybox
       [Arguments]  ${output}
       Should Contain  ${output}  bin
       Should Contain  ${output}  dev
       Should Contain  ${output}  etc
       Should Contain  ${output}  home
       Should Contain  ${output}  lib
       Should Contain  ${output}  lost+found
       Should Contain  ${output}  proc
       Should Contain  ${output}  root
       Should Contain  ${output}  run
       Should Contain  ${output}  sys
       Should Contain  ${output}  tmp
       Should Contain  ${output}  usr
       Should Contain  ${output}  var

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    Install VIC Appliance To Test Server
    Run Regression Tests

Concurrent Simple Exec
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${suffix}=  Evaluate  '%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
    Set Test Variable  ${ExecSimpleContainer}  Exec-simple-${suffix}
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name ${ExecSimpleContainer} ${busybox} /bin/top
    Should Be Equal As Integers  ${rc}  0

    :FOR  ${idx}  IN RANGE  1  51
    \   ${docker_cmd}=  Set Variable  docker %{VCH-PARAMS} exec -e idx=${idx} ${id} sh -c 'echo index is:${idx};/bin/ls'
    \   Start Process  ${docker_cmd}  alias=exec-simple-%{VCH-NAME}-${idx}  shell=true

    ${error_count}=  Set Variable  ${0}

    :FOR  ${idx}  IN RANGE  1  51
    \   ${result}=  Wait For Process  exec-simple-%{VCH-NAME}-${idx}  timeout=300s
    \   Run Keyword If  ${result.rc} == 0  Log To Console  rc=0 is expected.  ELSE  Log To Console  rc should be 0.
    \   ${status}=  Run Keyword And Return Status  Verify LS Output For Busybox  ${result.stdout}
    \   Log  ${result.stdout}
    \   Run Keyword If  "${status}" == "${True}"  Log To Console  ${result.stdout} expected result has been found.  ELSE  Log To Console  ${status} is not expected data.
    \   ${error_count}=  Run Keyword If  "${status}" == "${False}"  Evaluate  int(${error_count}) + 1  ELSE  Evaluate  int(${error_count})

    Log To Console  failed count:${error_count}
    # stop the container now that we have a successful series of concurrent execs
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} stop ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Integers  ${error_count}  0
