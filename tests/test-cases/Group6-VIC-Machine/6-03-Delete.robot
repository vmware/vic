*** Settings ***
Documentation  Test 6-03 - Verify delete clean up all resources
Resource  ../../resources/Util.robot
Test Setup  Install VIC Appliance To Test Server
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Keywords ***
Initial load
    # Create container VM first
    Log To Console  \nRunning docker pull busybox...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${name}=  Generate Random String  15
    ${rc}  ${container-id}=  Run And Return Rc And Output  docker ${params} create --name ${name} busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container-id}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container-id}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    Set Suite Variable  ${containerName}  ${name}

*** Test Cases ***
Delete VCH and verify
    Initial load
    # Get VCH uuid and container VM uuid, to check if resources are removed correctly
    Run Keyword And Ignore Error  Gather Logs From Test Server
    ${uuid}=  Run  govc vm.info -json\=true ${vch-name} | jq -r '.VirtualMachines[0].Config.Uuid'
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name}
    Should Contain  ${ret}  is powered on

    # Delete with force
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name} --force
    Should Contain  ${ret}  Completed successfully
    Should Not Contain  ${ret}  Operation failed: Error caused by file

    # Check VM is removed
    ${ret}=  Run  govc vm.info -json=true ${containerName}-*
    Should Contain  ${ret}  {"VirtualMachines":null}
    ${ret}=  Run  govc vm.info -json=true ${vch-name}
    Should Contain  ${ret}  {"VirtualMachines":null}

    # Check directories is removed
    ${ret}=  Run  govc datastore.ls VIC/${uuid}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls ${vch-name}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls VIC/${containerName}-*
    Should Contain  ${ret}   was not found

    # Check resource pool is removed
    ${ret}=  Run  govc pool.info -json=true host/*/Resources/${vch-name}
	Should Contain  ${ret}  {"ResourcePools":null}
