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
${VDNET_LAUNCHER_HOST}  10.92.34.93
${VIC_SVC_USER}  svc.victestuser
${VIC_SVC_PWD}  Q2D6EEi8k7hcWL2xj99
${USE_LOCAL_TOOLCHAIN}  0
${VDNET_MC_SETUP}  0
${NIMBUS_BASE}  /mts/git
${NSX_VERSION}  ob-10536046
${VC_VERSION}  ob-7312210
${ESX_VERSION}  ob-7389243

${nsxt_username}  admin
${nsxt_password}  Admin\!23Admin

${NSXT_JAVIS_TOPOLOGY}  NSXTransformers.Openstack.HybridEdge.OSDeploy.DeployOSMHTestbed
${NSXT_JAVIS_CONFIG_FILE}  /src/nsx-qe/vdnet/automation/yaml/nsxtransformers/jarvis/topologies/NSX-Transformers/NSX-Transformers-OSHybridEdge.json

*** Keywords ***
Vdnet NSXT Topology Setup
    [Timeout]    180 minutes
    ${TESTRUNID}=  Evaluate  'NSXT' + str(random.randint(1000,9999))  modules=random

    ${json}=  OperatingSystem.Get File  tests/resources/nimbus-testbeds/vic-vdnet-nsxt.json
    ${file}=  Replace Variables  ${json}
    Create File  /tmp/vic-vdnet-nsxt.json  ${file}

    Open Connection  ${VDNET_LAUNCHER_HOST}
    Wait Until Keyword Succeeds  2 min  30 sec  Login  ${VIC_SVC_USER}  ${VIC_SVC_PWD}

    Put File  /tmp/vic-vdnet-nsxt.json  destination=/tmp/vic-vdnet-nsxt.json
    ${result}=  Execute Command  cd /src/nsx-qe/automation && VDNET_MC_SETUP=${VDNET_MC_SETUP} USE_LOCAL_TOOLCHAIN=${USE_LOCAL_TOOLCHAIN} NIMBUS_BASE=${NIMBUS_BASE} vdnet3/test.py --config /tmp/vic-vdnet-nsxt.json /src/nsx-qe/automation/TDS/NSXTransformers/Openstack/HybridEdge/OSDeployTds.yaml::DeployOSMHTestbed --testbed save
    Log  ${result}
    ${status}=  Run Keyword And Return Status  Should Not Contain  ${result}  failed
    Should Be True  ${status}
    Close connection

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  ${VIC_SVC_USER}  ${VIC_SVC_PWD}
    ${vc-ip}=  Get IP  vdnet-vc-${VC_VERSION}-1-${TESTRUNID}
    ${nsxt-mgr-ip}=  Get IP  vdnet-nsxmanager-${NSX_VERSION}-1-${TESTRUNID}
    ${pod}=  Fetch POD  vdnet-nsxmanager-${NSX_VERSION}-1-${TESTRUNID}
    Close Connection

    Log To Console  Set environment variables up for GOVC
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
    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  ${VIC_SVC_USER}  ${VIC_SVC_PWD}
    Execute Command  ${NIMBUS_LOCATION} nimbus-ctl --nimbus=${pod_name} kill *${testrunid}*
    Close Connection

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${ls_name}=  Evaluate  'vic-bridge-ls'

    Create Nsxt Logical Switch  %{NSXT_MANAGER_URI} ${nsxt_username} ${nsxt_password} ${ls_name}
    Install VIC Appliance To Test Server  additional-args=--bridge-network ${ls_name}
    Run Regression Tests

    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server

