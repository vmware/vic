*** Settings ***
Documentation  Test 1-7 - Docker Stop
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Trap Signal Command
    # Container command runs an infinite loop, trapping and logging the given signal name
    [Arguments]  ${sig}
    [Return]  busybox sh -c "trap 'echo StopSignal${sig}' ${sig}; while true; do sleep 1; done"

Assert Stop Signal
    # Assert the docker stop signal was trapped by checking the container output log file
    [Arguments]  ${id}  ${sig}
    ${rc}=  Run And Return Rc  govc datastore.download ${id}/${id}.log ${TEMPDIR}/${id}.log
    Should Be Equal As Integers  ${rc}  0
    ${output}=  OperatingSystem.Get File  ${TEMPDIR}/${id}.log
    Should Contain  ${output}  StopSignal${sig}

*** Test Cases ***
Stop an already stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0

Basic docker container stop
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox sleep 30
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0

Stop a container with SIGKILL using default grace period
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  TERM
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0

Stop a container with SIGKILL using specific stop signal
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  USR1
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d --stop-signal USR1 ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Stop Signal  ${container}  USR1

Stop a container with SIGKILL using specific grace period
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox sleep 30
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${before}=  Get Current Date
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop -t 2 ${container}
    ${after}=  Get Current Date
    Should Be Equal As Integers  ${rc}  0
    ${result}=  Subtract Date From Date  ${after}  ${before}

    ${status}=  Get State Of Github Issue  1321
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-7-Docker-Stop.robot needs to be updated now that Issue #1321 has been resolved
    Log  Issue \#1321 is blocking implementation  WARN
    #Should Be True  ${1} < ${result} < ${3}

Stop a non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer

Attempt to stop a container that has been started out of band
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on=true ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
