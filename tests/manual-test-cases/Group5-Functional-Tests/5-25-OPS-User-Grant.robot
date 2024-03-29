# Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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
Documentation  Test 5-25 - OPS-User-Grant
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Ops User Create
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Keywords ***
Ops User Create
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
    Set Suite Variable  ${datacenter}  dc1
    Set Suite Variable  ${cluster}  cls
    
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
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    ${vc}=  Set Variable  vcname

    Set Suite Variable  ${ops_user_base_name}  vch-user
    Set Suite Variable  ${ops_user_domain}  vsphere.local
    ${ops_user_name}=  Catenate  SEPARATOR=@  ${ops_user_base_name}  ${ops_user_domain}
    Log To Console  Base User Name: ${ops_user_base_name}
    Log To Console  Full User Name: ${ops_user_name}

    Set Suite Variable  ${ops_user_name}
    Set Suite Variable  ${ops_user_password}  Admin!23
    Set Suite Variable  ${vc_admin_password}  Admin!23

    Log To Console  Setting up ops-user: ${ops_user_name}
    ${rc}  ${output}=  Run And Return Rc And Output  sshpass -p vmware ssh -o StrictHostKeyChecking=no root@${vc-ip} /usr/lib/vmware-vmafd/bin/dir-cli user create --account ${ops_user_base_name} --user-password ${ops_user_password} --first-name ${ops_user_base_name} --last-name ${ops_user_domain} --password ${vc_admin_password}
    Log  User Create ${ops_user_name}, rc: ${rc}, output: ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  sshpass -p vmware ssh -o StrictHostKeyChecking=no root@${vc-ip} /usr/lib/vmware-vmafd/bin/dir-cli user find-by-name --account ${ops_user_base_name} --password ${vc_admin_password}
    Log  User Find ${ops_user_base_name}, rc: ${rc}, output: ${output}
    Should Be Equal As Integers  ${rc}  0

    ${out}=  Run  govc role.usage
    Log  Output, govc role.usage: ${out}

Run privilege-dependent docker operations
    [Timeout]  15 minutes
    # Run containers with volumes and container networks to test scenarios requiring containerVMs
    # to have the highest privileges.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${c1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net public ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${c2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v fooVol:/dir ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${c3}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --net public -v barVol:/dir ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create existingvol
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${c4}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v existingvol:/dir ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c1} ${c2} ${c3} ${c4}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm existingvol
    Should Be Equal As Integers  ${rc}  0

    # Verify that containers cannot be destroyed out-of-band, i.e., the Destroy_VM task is
    # successfully disabled with an ops-user.
    ${c5}=  Evaluate  'cvm-' + str(random.randint(1000,9999))  modules=random
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name ${c5} ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.destroy "${c5}*"
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  The method is disabled by 'VIC'
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f ${c5}
    Should Be Equal As Integers  ${rc}  0

    # Verify that the required privileges for docker cp operations with a running container are present.
    Create File  ${CURDIR}/on-host.txt   hello world
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -dit --name online-cont ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online-cont sh -c 'echo "goodbye world" > /on-cont.txt'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # Copy from host to container.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/on-host.txt online-cont:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online-cont ls /
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  on-host.txt

    # Copy from container to host.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online-cont:/on-cont.txt ${CURDIR}/on-cont.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  OperatingSystem.Get File  ${CURDIR}/on-cont.txt
    Should Contain  ${content}  goodbye world

    # Clean up both files.
    Run Keyword and Ignore Error  Remove File  ${CURDIR}/on-cont.txt
    Run Keyword and Ignore Error  Remove File  ${CURDIR}/on-host.txt

Reconfigure VCH With Ops User
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux configure --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME} --ops-user=${ops_user_name} --ops-password=${ops_user_password} --ops-grant-perms --thumbprint=%{TEST_THUMBPRINT} --debug=1
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully

Attempt To Disable DRS
    Log To Console  Running govc to set drs-enabled, it should fail
    ${rc}  ${output}=  Run And Return Rc And Output  GOVC_USERNAME=${ops_user_name} GOVC_PASSWORD=${ops_user_password} govc cluster.change -drs-enabled /${datacenter}/host/${cluster}
    Log  Govc output: ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Permission to perform this operation was denied

Attempt To Create Resource Pool
    Log To Console  Running govc to create a resource pool named "5-25-OPS-User-Grant-%{DRONE_BUILD_NUMBER}", it should fail
    ${rc}  ${output}=  Run And Return Rc And Output  GOVC_USERNAME=${ops_user_name} GOVC_PASSWORD=${ops_user_password} govc pool.create */Resources/5-25-OPS-User-Grant-%{DRONE_BUILD_NUMBER}
    Log  Govc output: ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Permission to perform this operation was denied

*** Test Cases ***
vic-machine create grants ops-user perms
    Install VIC Appliance To Test Server  additional-args=--ops-user ${ops_user_name} --ops-password ${ops_user_password} --ops-grant-perms

    # Run a govc test to check that access is denied on some resources
    Attempt To Disable DRS

    Run Regression Tests

    Run privilege-dependent docker operations

    [Teardown]  Cleanup VIC Appliance On Test Server

granted ops-user perms work after upgrade
    [Tags]  vsphere70-not-support
    Install VIC with version to Test Server  v1.4.0  additional-args=--ops-user ${ops_user_name} --ops-password ${ops_user_password} --ops-grant-perms

    Check Original Version
    Upgrade
    Check Upgraded Version

    # Run a govc test to check that access is denied on some resources
    Attempt To Create Resource Pool

    Run Regression Tests

    Run privilege-dependent docker operations

    [Teardown]  Cleanup VIC Appliance On Test Server

Test with VM-Host Affinity
    Log To Console  \nStarting test...
    Install VIC Appliance To Test Server  additional-args=--ops-user ${ops_user_name} --ops-password ${ops_user_password} --ops-grant-perms --affinity-vm-group

    # Run a govc test to check that access is denied on some resources
    Attempt To Create Resource Pool

    Run Regression Tests

    Run privilege-dependent docker operations

    [Teardown]  Cleanup VIC Appliance On Test Server

vic-machine configure grants ops-user perms
    Install VIC Appliance To Test Server

    Reconfigure VCH With Ops User

    # Run a govc test to check that access is denied on some resources
    Attempt To Disable DRS

    Run privilege-dependent docker operations

    [Teardown]  Cleanup VIC Appliance On Test Server
