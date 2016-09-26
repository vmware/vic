*** Settings ***
Documentation  Test 1-7 - Docker Stop
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Trap Signal Command
    # Container command runs an infinite loop, trapping and logging the given signal name
    [Arguments]  ${sig}
    [Return]  busybox sh -c "trap 'echo StopSignal${sig}' ${sig}; echo READY; while true; do sleep 1; done"

Assert Ready
    # Assert the docker stop signal trap has been set
    [Arguments]  ${id}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  READY

Assert Stop Signal
    # Assert the docker stop signal was trapped by checking the container output log file
    [Arguments]  ${id}  ${sig}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  StopSignal${sig}

Assert Kill Signal
    # Assert SIGKILL was sent or not by checking the tether debug log file
    [Arguments]  ${id}  ${expect}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.info -json ${vch-name}/*-${id} | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal  ${output}  poweredOff
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.info -json *-${id} | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal  ${output}  poweredOff
    ${rc}  ${dir}=  Run And Return Rc And Output  govc datastore.ls *-${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc datastore.download ${dir}/${id}.debug -
    Should Be Equal As Integers  ${rc}  0
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
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Ready  ${container}
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Kill Signal  ${container}  False

Stop a container with SIGKILL using specific stop signal
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  USR1
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d --stop-signal USR1 ${trap}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Ready  ${container}
    ${rc}=  Run And Return Rc  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Stop Signal  ${container}  USR1
    Assert Kill Signal  ${container}  True

Stop a container with SIGKILL using specific grace period
    ${rc}=  Run And Return Rc  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${trap}=  Trap Signal Command  HUP
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --stop-signal HUP ${trap}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Assert Ready  ${container}
    ${rc}=  Run And Return Rc  docker ${params} stop -t 2 ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Stop Signal  ${container}  HUP
    Assert Kill Signal  ${container}  True

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
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -on=true ${vch-name}/${name}-*
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -on=true ${name}-*
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Assert Kill Signal  ${container}  False

Restart a stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it busybox /bin/ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    Wait Until VM Powers Off  *-${output}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
