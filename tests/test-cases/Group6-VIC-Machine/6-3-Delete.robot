*** Settings ***
Documentation  Test 6-3 - Verify delete clean up all resources
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Initial load
    # Create container VM first
    Log To Console  \nRunning docker pull busybox...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container-id}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container-id}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container-id}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    Set Suite Variable  ${containerID}  ${container-id}
    
Delete VCH and verify
    # Get VCH uuid and container VM uuid, to check if resources are removed correctly
    Run Keyword And Ignore Error  Gather Logs From Test Server
    ${uuid}=  Run  govc vm.info -json\=true ${vch-name} | jq -r '.VirtualMachines[0].Config.Uuid'
    Run  govc vm.power -on=true ${containerID}
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name}
    Should Contain  ${ret}  is powered on

    # Delete with force
    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name} --force
    Should Contain  ${ret}  Completed successfully

    # Check VM is removed
    ${ret}=  Run  govc vm.info -json=true ${containerID}
    Should Contain  ${ret}  {"VirtualMachines":null}
    ${ret}=  Run  govc vm.info -json=true ${vch-name}
    Should Contain  ${ret}  {"VirtualMachines":null}

    # Check directories is removed
    ${ret}=  Run  govc datastore.ls VIC/${uuid}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls ${vch-name}
    Should Contain  ${ret}   was not found
    ${ret}=  Run  govc datastore.ls VIC/${containerID}
    Should Contain  ${ret}   was not found

    # Check resource pool is removed
    ${ret}=  Run  govc pool.info -json=true host/*/Resources/${vch-name}
	Should Contain  ${ret}  {"ResourcePools":null}
