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
Documentation  Test 5-28 - VICAdmin Isolated
Resource  ../../resources/Util.robot
Suite Setup  Setup VCH With No WAN
Suite Teardown  Teardown VCH With No WAN
Default Tags

*** Keywords ***
Static IP Address Create
    [Timeout]    110 minutes
    Log To Console  Starting Static IP Address setup...
    Set Suite Variable  ${NIMBUS_LOCATION}  NIMBUS_LOCATION=wdc
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${name}=  Evaluate  'vic-5-28-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --plugin testng --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-simple-cluster --testbedSpecRubyFile /dbc/pa-dbc1111/mhagen/nimbus-testbeds/testbeds/vic-simple-cluster.rb --runName ${name}

    Open Connection  %{NIMBUS_GW}
    Wait Until Keyword Succeeds  10 min  30 sec  Login  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${vc-ip}=  Get IP  ${name}.vc.0
    ${pod}=  Fetch POD  ${name}.vc.0
    Set Suite Variable  ${NIMBUS_POD}  ${pod}
    Close Connection

    Set Suite Variable  @{list}  %{NIMBUS_USER}-${name}.esx.0  %{NIMBUS_USER}-${name}.esx.1  %{NIMBUS_USER}-${name}.esx.2  %{NIMBUS_USER}-${name}.nfs.0  %{NIMBUS_USER}-${name}.vc.0
    Log To Console  Finished Creating Cluster ${name}

    ${out}=  Get Static IP Address
    Set Suite Variable  ${static}  ${out}
    Append To List  ${list}  %{STATIC_WORKER_NAME}
    
    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  nfs0-1
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  15m

Setup VCH With No WAN
    Wait Until Keyword Succeeds  10x  10m  Static IP Address Create

    Log To Console  Create a vch with a public network on a no-wan portgroup.

    ${vlan}=  Evaluate  str(random.randint(1, 195))  modules=random

    ${dvs}=  Run  govc find -type DistributedVirtualSwitch | head -n1
    ${rc}  ${output}=  Run And Return Rc And Output  govc dvs.portgroup.add -vlan=${vlan} -dvs ${dvs} dpg-no-wan
    
    ${output}=  Run  bin/vic-machine-linux create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --no-tlsverify --bridge-network=%{BRIDGE_NETWORK} --management-network=%{PUBLIC_NETWORK} --client-network=%{PUBLIC_NETWORK} --client-network-ip &{static}[ip]/&{static}[netmask] --client-network-gateway &{static}[gateway] --public-network dpg-no-wan --public-network-ip 192.168.100.2/24 --public-network-gateway 192.168.100.1 --dns-server 10.170.16.48 --insecure-registry wdc-harbor-ci.eng.vmware.com

    Get Docker Params  ${output}

    Set Environment Variable  VIC-ADMIN  %{VCH-IP}:2378

Teardown VCH With No WAN
    Cleanup VIC Appliance On Test Server
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

Login And Save Cookies
    [Tags]  secret
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/authentication -XPOST -F username=%{TEST_USERNAME} -F password=%{TEST_PASSWORD} -D /tmp/cookies-%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

*** Test Cases ***
Display HTML
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN} -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  <title>VIC: %{VCH-NAME}</title>

WAN Status
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN} -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  <div class="sixty">Registry and Internet Connectivity<span class="error-message">

Get Portlayer Log
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/port-layer.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/init.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  reaping child processes

Get Docker Personality Log
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/docker-personality.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  docker personality

Get Container Logs
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container}=  Run And Return Rc and Output  docker %{VCH-PARAMS} create ${busybox} /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} start ${container}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${vmName}=  Get VM Display Name  ${container}
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk %{VIC-ADMIN}/container-logs.tar.gz -b /tmp/cookies-%{VCH-NAME} | (cd /tmp; tar xvzf - ${vmName}/tether.debug ${vmName}/vmware.log)
    Log  ${output}
    ${rc}  ${output}=  Run And Return Rc and Output  ls -l /tmp/${vmName}/vmware.log
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc and Output  ls -l /tmp/${vmName}/tether.debug
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc and Output  grep 'prepping for switch to container filesystem' /tmp/${vmName}/tether.debug
    Should Be Equal As Integers  ${rc}  0
    Run  rm -f /tmp/${vmName}/tether.debug /tmp/${vmName}/vmware.log

Get VICAdmin Log
    Login And Save Cookies
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/vicadmin.log -b /tmp/cookies-%{VCH-NAME}
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server