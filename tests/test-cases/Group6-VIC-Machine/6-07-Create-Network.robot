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

*** Test Cases ***
Public network - default
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    ${info}=  Get VM Info  %{VCH-NAME}
    Should Contain  ${info}  VM Network

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Public network - invalid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Guarantee port group doesn't already exist
    Run  govc host.portgroup.remove 'AAAAAAAAAA'

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --public-network=AAAAAAAAAA ${vicmachinetls}

    Should Contain  ${output}  --public-network: network 'AAAAAAAAAA' not found
    Should Contain  ${output}  vic-machine-linux create failed

Public network - invalid vCenter
    Pass execution  Test not implemented

Public network - DHCP
    Pass execution  Test not implemented

Public network - valid
    Pass execution  Test not implemented

Management network - none
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    ${status}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "management" is sharing NIC with "public"
    ${status2}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "public" is sharing NIC with "management"
    ${status3}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "public" is sharing NIC with "client"
    ${status4}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "management" is sharing NIC with "client"
    Should Be True  ${status} | ${status2} | ${status3} | ${status4}
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Management network - invalid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Guarantee port group doesn't already exist
    Run  govc host.portgroup.remove 'AAAAAAAAAA'

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --management-network=AAAAAAAAAA ${vicmachinetls}

    Should Contain  ${output}  --management-network: network 'AAAAAAAAAA' not found
    Should Contain  ${output}  vic-machine-linux create failed

Management network - invalid vCenter
    Pass execution  Test not implemented

Management network - unreachable
    Pass execution  Test not implemented

Management network - valid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --management-network=%{PUBLIC_NETWORK} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Connectivity Bridge to Public
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${out}=  Run  govc host.portgroup.remove bridge
    ${out}=  Run  govc host.portgroup.remove vm-network

    Log To Console  Create a public portgroup.
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN vm-network

    Log To Console  Create a bridge portgroup.
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN bridge

    ${output}=  Run  bin/vic-machine-linux create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --force=true --bridge-network=bridge --public-network=vm-network --compute-resource=%{TEST_RESOURCE} --container-network vm-network --no-tlsverify

    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    # this container will listen on :8000 and we're passing the -p option to the VCH so it should be exposed
    Log To Console  Creating public container.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=vm-network -p 8000 --name p1 busybox nc -l -p 8000
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Getting IP for public container
    ${ip}=  Run  docker %{VCH-PARAMS} inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress }}{{end}}' p1

    Log To Console  Connecting to container on external network from container bridged network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net bridge busybox nc ${ip} 8000
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:

    # nc is listening, but since we didn't pass the -p flag to docker, the port should not be exposed.
    Log To Console  Creating public container with no ports exposed.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=vm-network --name p2 busybox nc -l -p 8000
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Getting IP for public container
    ${ip}=  Run  docker %{VCH-PARAMS} inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress }}{{end}}' p2

    # we expect this to fail since the port wasn't exposed
    Log To Console  Connecting to container on external network from container bridged network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net bridge busybox nc ${ip} 8000
    Should Not Be Equal As Integers  ${rc}  0

    Log To Console  Port connection test from bridge to public networks succeeded.

    Cleanup VIC Appliance On Test Server

Connectivity Bridge to Management
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${out}=  Run  govc host.portgroup.remove bridge
    ${out}=  Run  govc host.portgroup.remove management

    Log To Console  Create a bridge portgroup
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN bridge

    Log To Console  Create a management portgroup.
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN management

    ${output}=  Run  bin/vic-machine-linux create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --force=true --bridge-network=bridge --compute-resource=%{TEST_RESOURCE} --container-network management --container-network vm-network --container-network-ip-range=management:10.10.10.0/24 --container-network-gateway=management:10.10.10.1/24 --no-tlsverify

    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Log To Console  Creating management container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=management --name m1 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Starting management container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start m1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:

    Log To Console  Creating bridge container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net=bridge --name b1 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    Log To Console  Starting bridge container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start b1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:

    Log To Console  Getting IP for management container
    ${ip}=  Run  docker %{VCH-PARAMS} inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress }}{{end}}' m1

    Log To Console  Pinging from bridge to management container.
    ${id}=  Run  docker %{VCH-PARAMS} run -d busybox ping -c 30 ${ip}

    Log To Console  Attach to running container.
    ${out}=  Run  docker %{VCH-PARAMS} attach ${id}

    Should Contain  ${out}  100% packet loss
    Log To Console  Ping test succeeded.

    Cleanup VIC Appliance On Test Server

Bridge network - vCenter none
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} ${vicmachinetls}
    Should Contain  ${output}  ERROR
    Should Contain  ${output}  An existing distributed port group must be specified for bridge network on vCenter

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - ESX none
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Pass Execution  Test skipped on VC

    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Bridge network - create bridge network if it doesn't exist
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Pass Execution  Test not applicable on vCenter
    # ESX should automatically create the bridge switch & port group AAAAAAAAAA, but vCenter would fail with unknown network error

    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Guarantee port group doesn't already exist
    Run  govc host.portgroup.remove 'AAAAAAAAAA'
    Run  govc host.vswitch.remove 'AAAAAAAAAA'

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=AAAAAAAAAA ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

    Run  govc host.portgroup.remove 'AAAAAAAAAA'
    Run  govc host.vswitch.remove 'AAAAAAAAAA'

Bridge network - invalid vCenter
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Pass execution  Test not implemented

Bridge network - non-DPG
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Pass execution  Test not implemented

Bridge network - valid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Bridge network - reused port group
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --management-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --client-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - invalid IP settings
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --bridge-network-range 1.1.1.1 ${vicmachinetls}
    Should Contain  ${output}  Error parsing bridge network ip range

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - invalid bridge network range
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --bridge-network-range 1.1.1.1/17 ${vicmachinetls}
    Should Contain  ${output}  --bridge-network-range must be /16 or larger network

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - valid with IP range
    Pass execution  Test not implemented

Container network - space in network name invalid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=bridge --container-network 'VM Network With Spaces' ${vicmachinetls}
    Should Contain  ${output}  A network alias must be supplied when network name "VM Network With Spaces" contains spaces.
    Should Contain  ${output}  vic-machine-linux create failed

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=bridge --container-network 'VM Network With Spaces': ${vicmachinetls}
    Should Contain  ${output}  A network alias must be supplied when network name "VM Network With Spaces:" contains spaces.
    Should Contain  ${output}  vic-machine-linux create failed

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=bridge --container-network 'vm-network':'vm network' ${vicmachinetls}
    Should Contain  ${output}  The network alias supplied in "vm-network:vm network" cannot contain spaces.
    Should Contain  ${output}  vic-machine-linux create failed

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Container network - space in network name valid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    Log To Console  Create a portgroup with a space in it's name
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN 'VM Network With Spaces'

    Log To Console  Create a bridge portgroup.
    ${out}=  Run  govc host.portgroup.add -vswitch vSwitchLAN bridge

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=bridge --container-network 'VM Network With Spaces':vmnet --insecure-registry harbor.ci.drone.local ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests

    ${output}=  Run  docker %{VCH-PARAMS} network ls
    Should Contain  ${output}  vmnet

    # Clean up port groups
    ${out}=  Run  govc host.portgroup.remove 'VM Network With Spaces'
    ${out}=  Run  govc host.portgroup.remove 'bridge'

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}
    Cleanup VIC Appliance On Test Server

Container network invalid 1
    Pass execution  Test not implemented

Container network invalid 2
    Pass execution  Test not implemented

Container network 1
    Pass execution  Test not implemented

Container network 2
    Pass execution  Test not implemented

Network mapping invalid
    Pass execution  Test not implemented

Network mapping gateway invalid
    Pass execution  Test not implemented

Network mapping IP invalid
    Pass execution  Test not implemented

DNS format invalid
    Pass execution  Test not implemented

Network mapping
    Pass execution  Test not implemented

VCH static IP - Static public
    Pass execution  Test not implemented

VCH static IP - Static client
    Pass execution  Test not implemented

VCH static IP - Static management
    Pass execution  Test not implemented

VCH static IP - different port groups 1
    Pass execution  Test not implemented

VCH static IP - different port groups 2
    Pass execution  Test not implemented

VCH static IP - same port group
    Pass execution  Test not implemented

VCH static IP - same subnet for multiple port groups
    Pass execution  Test not implemented
