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
Documentation  Test 6-07 - Verify vic-machine create network function
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server
Test Timeout  20 minutes

*** Keywords ***
Cleanup Container Firewalls Test Networks
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove bridge
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove open-net
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove closed-net
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove published-net
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove outbound-net
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove peers-net-1
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove peers-net-2

    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  bridge
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  open-net
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  closed-net
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  published-net
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  outbound-net
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  peers-net-1
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  peers-net-2

Main Test
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove bridge
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove vm-network
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  bridge
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  vm-network

    Log To Console  Create a public portgroup.
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vswitch vSwitchLAN vm-network
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Add VC Distributed Portgroup  test-ds  vm-network

    Log To Console  Create a bridge portgroup.
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vswitch vSwitchLAN bridge
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Add VC Distributed Portgroup  test-ds  bridge

    ${output}=  Run  bin/vic-machine-linux create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --force=true --bridge-network=bridge --public-network=vm-network --compute-resource=%{TEST_RESOURCE} --container-network vm-network --container-network-firewall vm-network:published --no-tlsverify --insecure-registry harbor.ci.drone.local

    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    # this container will listen on :8000 and we're passing the -p option to the VCH so it should be exposed
    Log To Console  Creating public container.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=vm-network -p 8000 --name p1 ${busybox} nc -l -p 8000
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Getting IP for public container
    ${ip}=  Run  docker %{VCH-PARAMS} inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress }}{{end}}' p1

    Log To Console  Connecting to container on external network from container bridged network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net bridge ${busybox} nc ${ip} 8000
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:

    # nc is listening, but since we didn't pass the -p flag to docker, the port should not be exposed.
    Log To Console  Creating public container with no ports exposed.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=vm-network --name p2 ${busybox} nc -l -p 8000
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Getting IP for public container
    ${ip}=  Run  docker %{VCH-PARAMS} inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress }}{{end}}' p2

    # we expect this to fail since the port wasn't exposed
    Log To Console  Connecting to container on external network from container bridged network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net bridge ${busybox} nc ${ip} 8000
    Should Not Be Equal As Integers  ${rc}  0

    Log To Console  Port connection test from bridge to public networks succeeded.

    Cleanup VIC Appliance On Test Server


Cleanup Container Firewalls Test
    Cleanup VIC Appliance On Test Server
    Cleanup Container Firewalls Test Networks

*** Test Cases ***
Connectivity Bridge to Public1
    Main Test


Connectivity Bridge to Public2
    Main Test


Connectivity Bridge to Public3
    Main Test


Connectivity Bridge to Public4
    Main Test


Connectivity Bridge to Public5
    Main Test


Connectivity Bridge to Public6
    Main Test


Connectivity Bridge to Public7
    Main Test
