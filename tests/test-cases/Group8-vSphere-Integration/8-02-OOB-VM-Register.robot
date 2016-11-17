*** Settings ***
Documentation  Test 8-02 OOB VM Register
Resource  ../../resources/Util.robot
#Suite Teardown  Extra Cleanup

*** Keywords ***
Extra Cleanup
    ${out}=  Run Keyword And Ignore Error  Run  govc vm.destroy ${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc pool.destroy host/*/Resources/${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc datastore.rm ${old-vm}
    ${out}=  Run Keyword And Ignore Error  Run  govc host.portgroup.remove ${old-vm}-bridge
    Cleanup VIC Appliance On Test Server

*** Test Cases ***
Verify VIC Still Works When Different VM Is Registered
    ${status}=  Get State Of Github Issue  3201
    Run Keyword If  '${status}' == 'closed'  Fail  Test 8-02-OOB-VM-Register.robot needs to be updated now that Issue #3201 has been resolved
    Log  Issue \#3201 is blocking implementation  WARN
#    Install VIC Appliance To Test Server
#    Set Suite Variable  ${old-vm}  ${vch-name}
#    Install VIC Appliance To Test Server

#    ${out}=  Run  govc vm.power -off ${old-vm}
#    Should Contain  ${out}  OK
#    ${out}=  Run  govc vm.unregister ${old-vm}
#    Should Be Empty  ${out}
#    ${out}=  Run  govc vm.register ${old-vm}/${old-vm}.vmx
#    Should Be Empty  ${out}

#    ${out}=  Run  docker ${params} ps -a
#    Log  ${out}
#    Should Contain  ${out}  CONTAINER ID
#    Should Contain  ${out}  IMAGE
#    Should Contain  ${out}  COMMAND

#    Run Regression Tests

#    ${out}=  Run  govc vm.destroy ${old-vm}