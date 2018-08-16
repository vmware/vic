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
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Keywords ***
Remove Container OOB Setup
    [Timeout]    110 minutes
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${esx1}  ${esx2}  ${esx3}  ${vc}  ${esx1-ip}  ${esx2-ip}  ${esx3-ip}  ${vc-ip}=  Create a Simple VC Cluster
    Set Suite Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  %{NIMBUS_USER}-${vc}

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

