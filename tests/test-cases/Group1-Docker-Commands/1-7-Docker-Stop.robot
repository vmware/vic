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

Assert Kill Signal
    # Assert SIGKILL was sent or not by checking the tether debug log file
    [Arguments]  ${id}  ${expect}
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -json -vm.path "[%{TEST_DATASTORE}] ${id}/${id}.vmx" | jq -r .VirtualMachines[].Runtime.PowerState
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  poweredOff
    ${rc}=  Run And Return Rc  govc datastore.download ${id}/${id}.debug ${TEMPDIR}/${id}.debug
    Should Be Equal As Integers  ${rc}  0
    ${output}=  OperatingSystem.Get File  ${TEMPDIR}/${id}.debug
    Run Keyword If  ${expect}  Should Contain  ${output}  sending signal KILL
    Run Keyword Unless  ${expect}  Should Not Contain  ${output}  sending signal KILL

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
    Assert Kill Signal  ${container}  False

Stop a container with SIGKILL using default grace period
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  HUP
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Kill Signal  ${container}  False

Stop a container with SIGKILL using specific stop signal
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  USR1
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d --stop-signal USR1 ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Stop Signal  ${container}  USR1
    Assert Kill Signal  ${container}  True

Stop a container with SIGKILL using specific grace period
    ${status}=  Get State Of Github Issue  1924
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-7-Docker-Stop.robot needs to be updated now that Issue #1924 has been resolved
    Log  Issue \#1924 is blocking implementation  WARN
    #${rc}=  Run And Return Rc  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${trap}=  Trap Signal Command  HUP
    #${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d --stop-signal HUP ${trap}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}=  Run And Return Rc  docker ${params} stop -t 2 ${container}
    #Should Be Equal As Integers  ${rc}  0
    #Assert Stop Signal  ${container}  HUP
    #Assert Kill Signal  ${container}  True

Stop a non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer

Attempt to stop a container that has been started out of band
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${name}=  Generate Random String  15
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name ${name} busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on=true ${name}-*
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Kill Signal  ${container}  False
