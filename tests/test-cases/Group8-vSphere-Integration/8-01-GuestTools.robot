*** Settings ***
Documentation  Test 8-01 - Verify VM guest tools integration
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Assert VM Power State
    [Arguments]  ${name}  ${state}
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -json ${name}-* | jq -r .VirtualMachines[].Runtime.PowerState
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  ${state}

*** Test Cases ***
Verify VCH VM guest IP is reported
    ${ip}=  Run  govc vm.ip %{VCH-NAME}
    # VCH ip should be the same as docker host param
    Should Contain  %{VCH-PARAMS}  ${ip}

Verify container VM guest IP is reported
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${name}=  Generate Random String  15
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${name} -d busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.ip ${name}-*
    Should Be Equal As Integers  ${rc}  0

Stop container VM using guest shutdown
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${name}=  Generate Random String  15
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${name} -d busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -s ${name}-*
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  500 milliseconds  Assert VM Power State  ${name}  poweredOff

Signal container VM using vix command
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${name}=  Generate Random String  15
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${name} -d busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Run  govc vm.ip ${name}-*
    # Invalid command
    ${rc}=  Run And Return Rc  govc guest.start -vm ${name}-* -l ${id} hello world
    Should Be Equal As Integers  ${rc}  1
    # Invalid id (via auth user)
    ${rc}=  Run And Return Rc  govc guest.start -vm ${name}-* kill USR1
    Should Be Equal As Integers  ${rc}  1
    # OK
    ${rc}  ${output}=  Run And Return Rc And Output  govc guest.start -vm ${name}-* -l ${id} kill USR1
    Should Be Equal As Integers  ${rc}  0
