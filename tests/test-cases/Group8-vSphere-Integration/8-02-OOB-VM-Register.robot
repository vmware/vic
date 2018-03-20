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
Documentation  Test 8-02 OOB VM Register
Resource  ../../resources/Util.robot
Suite Teardown  Extra Cleanup

*** Keywords ***
Extra Cleanup
    ${out}=  Run Keyword And Ignore Error  Run  govc vm.destroy ${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc pool.destroy host/*/Resources/${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc datastore.rm ${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc host.portgroup.remove ${old-vm}-bridge
    Cleanup VIC Appliance On Test Server

*** Test Cases ***
Verify VIC Still Works When Different VM Is Registered
    Install VIC Appliance To Test Server
    Set Suite Variable  ${old-vm}  %{VCH-NAME}
    Add VCH to Removal Exception List  ${old-vm}

    # we need to be sure that we do not treat the original vch as a dangling entity during the second install
    Install VIC Appliance To Test Server
    Remove VCH from Removal Exception List  ${old-vm}

    ${out}=  Run  govc ls vm/${old-vm}
    Should Contain  ${out}  ${old-vm}/${old-vm}
    ${out}=  Run  govc vm.power -off ${old-vm}
    Should Contain  ${out}  OK
    ${out}=  Run  govc vm.unregister ${old-vm}
    Should Be Empty  ${out}

    # At this point the vm is unregsitered and we will need to reregister this vm...
    # we need to put it into the original inventory folder. so we need to fetch that
    # path. We also want to be explicit about the resource pool.
    ${old-vm-folder}=  Run  govc find / -name ${old-vm} -type f
    ${old-vm-pool}=  Run  govc find / -name ${old-vm} -type p
    ${out}=  Run  govc vm.register -pool ${old-vm-pool} -folder ${old-vm-folder} ${old-vm}/${old-vm}.vmx
    Should Be Empty  ${out}

    ${out}=  Run  docker %{VCH-PARAMS} ps -a
    Log  ${out}
    Should Contain  ${out}  CONTAINER ID
    Should Contain  ${out}  IMAGE
    Should Contain  ${out}  COMMAND

    Run Regression Tests

    ${out}=  Run  govc vm.destroy ${old-vm}
