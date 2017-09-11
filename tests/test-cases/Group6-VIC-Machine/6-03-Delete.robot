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
Documentation  Test 6-03 - Verify delete clean up all resources
Resource  ../../resources/Util.robot
Test Setup  Install VIC Appliance To Test Server
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server
Test Timeout  20 minutes

*** Keywords ***
Initial load
    # Create container VM first
    Log To Console  \nRunning docker pull ${busybox}...
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${name}=  Generate Random String  15
    ${rc}  ${container-id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name ${name} ${busybox} /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container-id}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container-id}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    Set Suite Variable  ${containerName}  ${name}

*** Test Cases ***
Delete VCH and verify
    Initial load
    # Get VCH uuid and container VM uuid, to check if resources are removed correctly
    Run Keyword And Ignore Error  Gather Logs From Test Server
    ${uuid}=  Run  govc vm.info -json\=true %{VCH-NAME} | jq -r '.VirtualMachines[0].Config.Uuid'
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME}
    Should Contain  ${ret}  is powered on

    # Delete with force
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME} --force
    Should Contain  ${ret}  Completed successfully
    Should Not Contain  ${ret}  Operation failed: Error caused by file

    # Check VM is removed
    ${ret}=  Run  govc vm.info -json=true ${containerName}-*
    Should Contain  ${ret}  {"VirtualMachines":null}
    ${ret}=  Run  govc vm.info -json=true %{VCH-NAME}
    Should Contain  ${ret}  {"VirtualMachines":null}

    # Check directories is removed
    ${ret}=  Run  govc datastore.ls VIC/${uuid}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls %{VCH-NAME}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls VIC/${containerName}-*
    Should Contain  ${ret}   was not found

    # Check resource pool is removed
    ${ret}=  Run  govc pool.info -json=true host/*/Resources/%{VCH-NAME}
    Should Contain  ${ret}  {"ResourcePools":null}

Attach Disks and Delete VCH
    # VCH should delete normally during commit/pull/cp/push operations
    ${output}=  Install VIC Appliance To Test Server
    Should Contain    ${output}    Installer completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${ubuntu}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls %{VCH-NAME}/VIC/
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # iterate through found images and attach them to the appliance VM
    ${vm-ip}=  Run  govc vm.ip %{VCH-NAME}
    ${imagedir}=  Run  govc datastore.ls %{VCH-NAME}/VIC/
    ${images}=  Run  govc datastore.ls %{VCH-NAME}/VIC/${imagedir}/images/ | tr '${\n}' ' '
    ${rc}  ${output}=  Run And Return Rc And Output  (set -e; for x in ${images}; do echo $x; govc vm.disk.attach -disk=%{VCH-NAME}/VIC/${imagedir}/images/$x/$x.vmdk -vm.ip=${vm-ip}; done)
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME}
    Log  ${output}

    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully

    ${rc}=  Run And Return Rc  govc datastore.ls -dc=%{TEST_DATACENTER} %{VCH-NAME}/VIC/
    Should Be Equal As Integers  ${rc}  1
