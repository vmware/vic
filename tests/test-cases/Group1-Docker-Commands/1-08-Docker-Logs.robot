*** Settings ***
Documentation  Test 1-08 - Docker Logs
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Grep Logs And Count Lines
    [Arguments]  ${id}  ${match}  ${total}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${id}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${match}
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${linecount}  ${total}

*** Test Cases ***
Docker logs with tail
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} create busybox sh -c 'seq 1 5000'
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${id}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Grep Logs And Count Lines  ${id}  2500  5000
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=all ${id}
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Integers  ${linecount}  5000
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=200 ${id}
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Integers  ${linecount}  200
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=0 ${id}
    Should Be Equal As Integers  ${rc}  0
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${linecount}  0
    
Docker logs with follow
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} create busybox sh -c 'for i in $(seq 1 5) ; do sleep 1 && echo line $i; done'
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow ${id}
    Should Be Equal As Integers  ${rc}  0
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${linecount}  5
    ${lastline}=  Get Line  ${output}  4
    Should Contain  ${lastline}  line 5
    # Container is stopped at this point, verify that --follow does not block.
    ${rc}  ${output2}=  Run And Return Rc And Output  docker ${params} logs --follow ${id}
    Should Be Equal  ${output}  ${output2}

Docker logs with follow and tail
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} create busybox sh -c 'trap "seq 11 20; exit" HUP; seq 1 10; while true; do sleep 1; done'
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${id}
    Should Be Equal As Integers  ${rc}  0
    # Wait for the first 10 lines to be logged
    Wait Until Keyword Succeeds  20x  200 milliseconds  Grep Logs And Count Lines  ${id}  5  10
    # kill -HUP will create another 5 lines of log output
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} kill -s HUP ${id}
    # --tail=5 to skip the first 5 lines and --follow to wait for the rest
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=5 --follow ${id}
    Should Be Equal As Integers  ${rc}  0
    ${linecount}=  Get Line Count  ${output}
    Should Be Equal As Integers  ${linecount}  15

Docker logs with timestamps and since certain time
    ${status}=  Get State Of Github Issue  1738
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#1738 is blocking implementation  WARN
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;'
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
	Run  Sleep 6, wait for container to finish
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --since=1s ${containerID}
	Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support '--since'
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --timestamps ${containerID}
	Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support '--timestamps'

Docker logs with no flags
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker ${params} run -d busybox sh -c "seq 1 128 | xargs -n1 echo"
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  20x  200 milliseconds  Grep Logs And Count Lines  ${id}  42  128

Docker logs non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such container: fakeContainer
