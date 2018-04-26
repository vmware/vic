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
Documentation    Test 6-16 - Verify vic-machine configure
Resource         ../../resources/Util.robot
Suite Setup      Install VIC Appliance To Test Server
Suite Teardown   Cleanup VIC Appliance On Test Server
Test Timeout     20 minutes

*** Keywords ***
Wait For DNS Update
    [Arguments]    ${shouldContain}
    ${output}=    Run And Verify RC    govc vm.info -e %{VCH-NAME} | grep dns

    Run Keyword If        ${shouldContain}    Should Contain        ${output}    network/dns
    Run Keyword If        ${shouldContain}    Should Not Contain    ${output}    assigned.dns

    Run Keyword Unless    ${shouldContain}    Should Contain        ${output}    assigned.dns
    Run Keyword Unless    ${shouldContain}    Should Not Contain    ${output}    network/dns


Guestinfo Should Contain
    [Arguments]    ${vm}    ${key}    ${expected}

    ${out}=    Check VM Guestinfo    ${vm}     ${key}
               Should Contain        ${out}    ${expected}

Configure VCH
    [Arguments]    ${additional-args}

    ${out}=    Run And Verify Rc     bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout=%{TEST_TIMEOUT} ${additional-args}
               Should Contain        ${out}    Completed successfully


*** Test Cases ***
Configure VCH debug state
    ${out}=    Run And Verify Rc      bin/vic-machine-linux configure --help
               Should Contain                       ${out}    --debug

               Run And Verify Rc      docker %{VCH-PARAMS} pull ${busybox}
    ${id1}=    Run And Verify Rc      docker %{VCH-PARAMS} run -itd ${busybox}
    ${vm1}=    Get VM display name    ${id1}

               Inspect Config Should Contain        --debug=1
               Guestinfo Should Contain             %{VCH-NAME}    guestinfo.vice./init/diagnostics/debug    1
               Guestinfo Should Contain             ${vm1}         guestinfo.vice./diagnostics/debug         1

               Configure VCH          --debug=0

               Inspect Config Should Not Contain    --debug=1
               Guestinfo Should Contain             %{VCH-NAME}    guestinfo.vice./init/diagnostics/debug    0
               Guestinfo Should Contain             ${vm1}         guestinfo.vice./diagnostics/debug         1

    ${id2}=    Run And Verify Rc      docker %{VCH-PARAMS} run -itd ${busybox}
    ${vm2}=    Get VM display name    ${id2}

               Guestinfo Should Contain             ${vm2}         guestinfo.vice./diagnostics/debug         0

               Configure VCH          --debug=1

               Inspect Config Should Contain        --debug=1
               Guestinfo Should Contain             %{VCH-NAME}    guestinfo.vice./init/diagnostics/debug    1
               Guestinfo Should Contain             ${vm2}         guestinfo.vice./diagnostics/debug         0

Configure VCH Container Networks
    ${vlan}=  Get Public Network VLAN ID
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove cn-network
    ${vswitch}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.vswitch.info -json | jq -r ".Vswitch[0].Name"
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vlan=${vlan} -vswitch ${vswitch} cn-network
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  cn-network
    ${dvs}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc find -type DistributedVirtualSwitch | head -n1
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc dvs.portgroup.add -vlan=${vlan} -dvs ${dvs} cn-network

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  vmnet

    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --container-network=cn-network:vmnet

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -dit --net=vmnet ${busybox}
    Should Be Equal As Integers  ${rc}  0

    # Test that configure fails if an existing container-network is not specified
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove management
    ${vswitch}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.vswitch.info -json | jq -r ".Vswitch[0].Name"
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vswitch ${vswitch} management
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  management
    ${dvs}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc find -type DistributedVirtualSwitch | head -n1
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc dvs.portgroup.add -dvs ${dvs} management

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network management:mgmt --container-network-ip-range=management:10.10.10.0/24 --container-network-gateway=management:10.10.10.1/24
    Should Contain  ${output}  all existing container networks must also be specified
    Should Not Contain  ${output}  Completed successfully

    # Add another container network while specifying the existing one
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet --container-network management:mgmt --container-network-ip-range=management:10.10.10.0/24 --container-network-gateway=management:10.10.10.1/24 --container-network-firewall=management:open
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  vmnet
    Should Contain  ${output}  mgmt

    ${stripped}=  Remove String  %{PUBLIC_NETWORK}  '
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --container-network=${stripped}:public
    Should Contain  ${output}  --container-network=cn-network:vmnet
    Should Contain  ${output}  --container-network=management:mgmt
    Should Contain  ${output}  --container-network-ip-range=management:10.10.10.0/24
    Should Contain  ${output}  --container-network-gateway=management:10.10.10.1/24
    Should Contain  ${output}  --container-network-firewall=management:open

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -dit --net=mgmt ${busybox}
    Should Be Equal As Integers  ${rc}  0

    # Test that changes to existing networks are not supported
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet --container-network management:mgmt --container-network-ip-range=management:10.10.10.0/24 --container-network-gateway=management:10.10.10.2/24
    Should Contain  ${output}  changes to existing container networks are not supported
    Should Not Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet --container-network management:mgmt --container-network-ip-range=management:10.10.10.0/16 --container-network-gateway=management:10.10.10.1/24
    Should Contain  ${output}  changes to existing container networks are not supported
    Should Not Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet --container-network management:mgmt
    Should Contain  ${output}  changes to existing container networks are not supported
    Should Not Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --container-network=%{PUBLIC_NETWORK}:public --container-network cn-network:vmnet --container-network management:mgmt --container-network-firewall=management:closed
    Should Contain  ${output}  changes to existing container networks are not supported
    Should Not Contain  ${output}  Completed successfully


    # Clean up portgroups
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove cn-network
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  cn-network
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove management
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Remove VC Distributed Portgroup  management

Configure VCH https-proxy
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --http-proxy http://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  http://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --http-proxy=http://proxy.vmware.com:3128
    Should Not Contain  ${output}  --https-proxy

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --https-proxy https://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --https-proxy=https://proxy.vmware.com:3128
    Should Not Contain  ${output}  --http-proxy

Configure VCH ops user credentials and thumbprint
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --ops-user=%{TEST_USERNAME} --ops-password=%{TEST_PASSWORD}
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --ops-user=%{TEST_USERNAME}
    Should Contain  ${output}  --thumbprint=%{TEST_THUMBPRINT}

Configure VCH https-proxy through vch id
    ${vch-id}=  Get VCH ID  %{VCH-NAME}
    ${output}=  Run  bin/vic-machine-linux configure --id=${vch-id} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --https-proxy ""
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128

Configure VCH DNS server
    ${status}=  Get State Of Github Issue  7775
    Run Keyword If  '${status}' == 'closed'  Fail  Test 6-16-Config.robot needs to be updated now that Issue #7775 has been resolved
    Log  Issue \#7775 is blocking implementation  WARN
#    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
#    Should Not Contain  ${output}  --dns-server
#    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --dns-server 10.118.81.1 --dns-server 10.118.81.2
#    Should Contain  ${output}  Completed successfully
#    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
#    Should Contain  ${output}  --dns-server=10.118.81.1
#    Should Contain  ${output}  --dns-server=10.118.81.2
#    Wait Until Keyword Succeeds  10x  6s  Wait For DNS Update  ${true}
#    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --dns-server ""
#    Should Contain  ${output}  Completed successfully
#    Should Not Contain  ${output}  --dns-server
#    Wait Until Keyword Succeeds  10x  6s  Wait For DNS Update  ${false}

Configure VCH resources
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --cpu 5129 --cpu-reservation 10 --cpu-shares 8000 --memory 4096 --memory-reservation 10 --memory-shares 163840
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --cpu=5129
    Should Contain  ${output}  --cpu-reservation=10
    Should Contain  ${output}  --cpu-shares=8000
    Should Contain  ${output}  --memory=4096
    Should Contain  ${output}  --memory-reservation=10
    Should Contain  ${output}  --memory-shares=163840

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --cpu 1 --cpu-shares 1000 --memory 1 --memory-shares 1000
    Should Not Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --cpu=5129
    Should Contain  ${output}  --cpu-reservation=10
    Should Contain  ${output}  --cpu-shares=8000
    Should Contain  ${output}  --memory=4096
    Should Contain  ${output}  --memory-reservation=10
    Should Contain  ${output}  --memory-shares=163840

Configure VCH volume stores
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL:default --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-conf:configure
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-VOL:default
    Should Contain  ${output}  --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-conf:configure

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    ${volstores}=  Get Lines Containing String  ${output}  VolumeStores:
    Should Contain  ${volstores}  default
    Should Contain  ${volstores}  configure
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create defaultVol
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create confVol --opt VolumeStore=configure
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  defaultVol
    Should Contain  ${output}  confVol

    # Test that configure fails if an existing volume store is not specified
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-conf:configure
    Should Contain  ${output}  all existing volume stores must also be specified
    Should Not Contain  ${output}  Completed successfully

    # Test that changes to existing volume stores are not supported
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL:default --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-badpath:configure
    Should Contain  ${output}  changes to existing volume stores are not supported
    Should Not Contain  ${output}  Completed successfully

    # Add a new volume store while specifying the URL scheme
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-VOL:default --volume-store=%{TEST_DATASTORE}/%{VCH-NAME}-conf:configure --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-scheme:scheme
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-VOL:default
    Should Contain  ${output}  --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-conf:configure
    Should Contain  ${output}  --volume-store=ds://%{TEST_DATASTORE}/%{VCH-NAME}-scheme:scheme

Configure Present in vic-machine
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux
    Should Contain  ${output}  configure
