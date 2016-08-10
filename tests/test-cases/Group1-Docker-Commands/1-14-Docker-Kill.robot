*** Settings ***
Documentation  Test 1-14 - Docker Kill
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Trap Signal Command
    # Container command runs an infinite loop, trapping and logging the given signal name
    [Arguments]  ${sig}
    [Return]  busybox sh -c "trap 'echo KillSignal${sig}' ${sig}; while true; do date && sleep 1; done"

Assert Kill Signal
    # Assert the docker kill signal was trapped by checking the container output log file
    [Arguments]  ${id}  ${sig}
    ${rc}=  Run And Return Rc  govc datastore.download ${id}/${id}.log ${TEMPDIR}/${id}.log
    Should Be Equal As Integers  ${rc}  0
    ${output}=  OperatingSystem.Get File  ${TEMPDIR}/${id}.log
    Remove File  ${TEMPDIR}/${id}.log
    Should Contain  ${output}  KillSignal${sig}

Inspect State Running
    [Arguments]  ${id}  ${expected}
    ${rc}  ${state}=  Run And Return Rc And Output  docker ${params} inspect --format="{{ .State.Running }}" ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${state}  ${expected}

*** Test Cases ***
Signal a container with default kill signal
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} create busybox sleep 300
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} start ${id}
    Should Be Equal As Integers  ${rc}  0
    Inspect State Running  ${id}  true
    ${rc}=  Run And Return Rc  docker ${params} kill ${id}
    Should Be Equal As Integers  ${rc}  0
    # Wait for container VM to stop/powerOff
    Wait Until Keyword Succeeds  20x  500 milliseconds  Inspect State Running  ${id}  false
    # Cannot send signal to a powered off container VM
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} kill ${id}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Cannot kill container ${id}

Signal a container with SIGHUP
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  HUP
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} run -d ${trap}
    Should Be Equal As Integers  ${rc}  0
    # Expect failure with unknown signal name
    ${rc}=  Run And Return Rc  docker ${params} kill -s NOPE ${id}
    Should Be Equal As Integers  ${rc}  1
    ${rc}=  Run And Return Rc  docker ${params} kill -s HUP ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  5x  1 seconds  Assert Kill Signal  ${id}  HUP
    Inspect State Running  ${id}  true
    ${rc}=  Run And Return Rc  docker ${params} kill -s TERM ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  500 milliseconds  Inspect State Running  ${id}  false

Signal a non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} kill fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No such container: fakeContainer
