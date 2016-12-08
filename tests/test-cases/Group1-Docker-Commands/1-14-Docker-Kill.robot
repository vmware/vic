*** Settings ***
Documentation  Test 1-14 - Docker Kill
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Trap Signal Command
    # Container command runs an infinite loop, trapping and logging the given signal name
    [Arguments]  ${sig}
    [Return]  busybox sh -c "trap 'echo KillSignal${sig}' ${sig}; echo READY; while true; do date && sleep 1; done"

Assert Container Output
    [Arguments]  ${id}  ${match}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${match}

Check That Container Was Killed
    [Arguments]  ${container}
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} ${container}
    Log  ${out}
    Should Contain  ${out}  false
    Should Be Equal As Integers  ${rc}  0

*** Test Cases ***
Signal a container with default kill signal
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  HUP
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Container Output  ${id}  READY
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill ${id}
    Should Be Equal As Integers  ${rc}  0
    # Wait for container VM to stop/powerOff
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs --follow ${id}
    # Cannot send signal to a powered off container VM
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill ${id}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Cannot kill container ${id}

Signal a container with SIGHUP
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  HUP
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Container Output  ${id}  READY
    # Expect failure with unknown signal name
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill -s NOPE ${id}
    Should Be Equal As Integers  ${rc}  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill -s HUP ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Container Output  ${id}  KillSignalHUP
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill -s TERM ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs --follow ${id}

Signal a non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No such container: fakeContainer

Signal a tough to kill container - nginx
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} pull nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} kill ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That Container Was Killed  ${id}
