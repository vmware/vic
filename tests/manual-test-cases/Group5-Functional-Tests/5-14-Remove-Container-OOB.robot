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
Documentation  Test 5-14 - Remove Container OOB
Resource  ../../resources/Util.robot
Suite Setup  Nimbus Suite Setup  Remove Container OOB Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Pod Cleanup  ${nimbus_pod}  ${testbedname}

*** Keywords ***
Remove Container OOB Setup
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

*** Test Cases ***
Docker run a container and verify cannot destroy the VM Out of Band
    Install VIC Appliance To Test Server

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name removeOOB busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  govc vm.destroy removeOOB*
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  govc: ServerFaultCode: The method is disabled by 'VIC'

Docker run a container destroy VM Out Of Band verify container is cleaned up
    Install VIC Appliance To Test Server

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0

    # Create an anchor container to prevent VC from removing all the images and scractch disks
    # when the main container VM is destroyed out of band. This would leave the VCH in an unusable state.
    # VC uses using reference counting to keep disks around. Having a second active container while the
    # target is being destroyed keeps the reference count from going to 0
    ${rc}  ${anchor}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name anchorVmOOB busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name removeVmOOB busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  govc object.method -name Destroy_Task -enable=true removeVmOOB*
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  govc vm.destroy removeVmOOB*
    Should Be Equal As Integers  ${rc}  0

    # Wait for main contaier to stop
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} wait ${container}

    # Make sure the container has been removed from the list of active containers
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${out}  removeVmOOB
    Should Not Contain  ${out}  ${container}

    # Clenup anchor container
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill ${anchor}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${anchor}
    Should Be Equal As Integers  ${rc}  0

