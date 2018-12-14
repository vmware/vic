# Copyright 2018 VMware, Inc. All Rights Reserved.
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
Documentation  Test 5-29 - Opaque Network
Resource  ../../resources/Util.robot
Suite Setup  Vdnet NSXT Topology Setup
Suite Teardown  Vdnet NSXT Topology Cleanup  ${NIMBUS_POD}  ${testrunid}

*** Variables ***
${NIMBUS_LOCATION}  sc
${NIMBUS_LOCATION_FULL}  NIMBUS_LOCATION=${NIMBUS_LOCATION}
${VDNET_LAUNCHER_HOST}  10.160.201.180
${USE_LOCAL_TOOLCHAIN}  0
${VDNET_MC_SETUP}  0
${NIMBUS_BASE}  /mts/git
${NSX_VERSION}  ob-10536046
${VC_VERSION}  ob-7312210
${ESX_VERSION}  ob-7389243

${nsxt_username}  admin
${nsxt_password}  Admin\!23Admin

${NIMBUS_NSXT_USER}  svc.victestuser
${NIMBUS_NSXT_PASSWORD}  Q2D6EEi8k7hcWL2xj99

${nimbus_personal_user_cache}  %{NIMBUS_PERSONAL_USER}

*** Keywords ***
Vdnet NSXT Topology Setup
    [Timeout]    180 minutes
    ${TESTRUNID}=  Evaluate  'NSXT' + str(random.randint(1000,9999))  modules=random
    ${json}=  OperatingSystem.Get File  tests/resources/nimbus-testbeds/vic-vdnet-nsxt.json
    ${file}=  Replace Variables  ${json}
    Create File  /tmp/vic-vdnet-nsxt.json  ${file}

    # Unset NIMBUS_PERSONAL_USER, replace value to NSXT user so that all nimbus commands executed by NSXT user
    Remove Environment Variable  NIMBUS_PERSONAL_USER
    Set Environment Variable  NIMBUS_PERSONAL_USER  ${NIMBUS_NSXT_USER}

    Open Connection  ${VDNET_LAUNCHER_HOST}
    #Login account must be the value of variable in this test suit, otherwise other pulic accounts will have no right to perform on 10.160.201.180. 
    Wait Until Keyword Succeeds  2 min  30 sec  Login  ${NIMBUS_NSXT_USER}  ${NIMBUS_NSXT_PASSWORD}

    Put File  /tmp/vic-vdnet-nsxt.json  destination=/tmp/vic-vdnet-nsxt.json
    #Deploy the testbed
    ${result}=  Execute Command  cd /src/nsx-qe/automation && VDNET_MC_SETUP=${VDNET_MC_SETUP} USE_LOCAL_TOOLCHAIN=${USE_LOCAL_TOOLCHAIN} NIMBUS_BASE=${NIMBUS_BASE} vdnet3/test.py --config /tmp/vic-vdnet-nsxt.json /src/nsx-qe/automation/TDS/NSXTransformers/Openstack/HybridEdge/OSDeployTds.yaml::DeployOSMHTestbed --testbed save 2>&1
    Log  ${result}
    Should Not Contain  ${result}  failed
    #Run the test cases to customize NSXT and vCenter
    ${result}=  Execute Command  cd /src/nsx-qe/automation && VDNET_MC_SETUP=${VDNET_MC_SETUP} USE_LOCAL_TOOLCHAIN=${USE_LOCAL_TOOLCHAIN} NIMBUS_BASE=${NIMBUS_BASE} vdnet3/test.py --config /tmp/vic-vdnet-nsxt.json /src/nsx-qe/automation/TDS/NSXTransformers/Openstack/HybridEdge/OSDeployTds.yaml::DeployOSMHTestbed --testbed reuse 2>&1
    Log  ${result}
    Should Not Contain  ${result}  failed
    Close connection

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  ${NIMBUS_NSXT_USER}  ${NIMBUS_NSXT_PASSWORD}
    ${vc-ip}=  Get IP  vdnet-vc-${VC_VERSION}-1-${TESTRUNID}
    ${nsxt-mgr-ip}=  Get IP  vdnet-nsxmanager-${NSX_VERSION}-1-${TESTRUNID}
    ${pod}=  Fetch POD  vdnet-nsxmanager-${NSX_VERSION}-1-${TESTRUNID}
    Close Connection

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_URL  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  TEST_DATASTORE  vdnetSharedStorage
    Set Environment Variable  TEST_DATACENTER  /os-test-dc
    Set Environment Variable  TEST_RESOURCE  os-compute-cluster-1
    Set Environment Variable  TEST_TIMEOUT  30m

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Set Environment Variable  NSXT_MANAGER_URI  ${nsxt-mgr-ip}
    Set Suite Variable  ${NIMBUS_POD}  ${pod}
    Set Suite Variable  ${testrunid}  ${TESTRUNID}

    OperatingSystem.Remove File  /tmp/vic-vdnet-nsxt.json

Vdnet NSXT Topology Cleanup
    [Timeout]    180 minutes
    [Arguments]    ${pod_name}  ${testrunid}

    #Reset NIMBUS_PERSONAL_USER value to the initial value
    Remove Environment Variable  NIMBUS_PERSONAL_USER
    Set Environment Variable  NIMBUS_PERSONAL_USER  ${nimbus_personal_user_cache}

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  ${NIMBUS_NSXT_USER}  ${NIMBUS_NSXT_PASSWORD}
    Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl --nimbus=${pod_name} kill *${testrunid}*
    Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl --nimbus=${pod_name} kill *isolated-06-gw
    Close Connection

*** Test Cases ***
Basic VIC Tests with NSXT Topology
    Log To Console  \nStarting test...
    ${bridge_ls_name}=  Evaluate  'vic-bridge-ls'
    ${result}=  Create Nsxt Logical Switch  %{NSXT_MANAGER_URI}  ${nsxt_username}  ${nsxt_password}  ${bridge_ls_name}
    Should Not Be Equal    ${result}    null
    ${container_ls_name}=  Evaluate  'vic-container-ls'
    ${result}=  Create Nsxt Logical Switch  %{NSXT_MANAGER_URI}  ${nsxt_username}  ${nsxt_password}  ${container_ls_name}
    Should Not Be Equal    ${result}    null
    Set Environment Variable  BRIDGE_NETWORK  ${bridge_ls_name}
    Set Environment Variable  CONTAINER_NETWORK  ${container_ls_name}
    #Wait for 5 mins to make sure the logical switches is synchronized with vCenter
    Sleep  300

    Install VIC Appliance To Test Server  additional-args=--container-network-firewall=%{CONTAINER_NETWORK}:open --container-network-ip-range %{CONTAINER_NETWORK}:50.0.9.100-50.0.9.200 --container-network-gateway %{CONTAINER_NETWORK}:50.0.9.1/24
    Run Regression Tests

    #Validate containers' connection when using container network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it -d --net=public --name first ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${ip}=  Get Container IP  %{VCH-PARAMS}  first  public
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net=public ${debian} ping -c1 ${ip}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server
