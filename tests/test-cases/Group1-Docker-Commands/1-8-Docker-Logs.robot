*** Settings ***
Documentation  Test 1-8 - Docker Logs
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker logs with tail
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;'
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
	#Run  Sleep 6
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=all ${containerID}
	#${linecount}=  Get Line Count  ${output}
    #Should Be Equal As Integers  ${rc}  0
	#Should Be Equal As Integers  ${linecount}  5
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=2 ${containerID}
	#${linecount}=  Get Line Count  ${output}
    #Should Be Equal As Integers  ${rc}  0
	#Should Be Equal As Integers  ${linecount}  2
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --tail=0 ${containerID}
    #Should Be Equal As Integers  ${rc}  0
	#${linecount}=  Get Line Count  ${output}
	#Should Be Equal As Integers  ${linecount}  0
    
Docker logs with follow
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;'
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
	#Run  Sleep 2
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow ${containerID}
    #Should Be Equal As Integers  ${rc}  0
	#${linecount}=  Get Line Count  ${output}
	#${lastline}=  Get Line  ${output}  ${linecount}-1
	#Should Contain  ${lastline}  line 5

Docker logs with follow and tail
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;'
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
	#Run  Sleep 2.5, wait for first 2 lines to output
	#${rc}  ${output}=  Run and Return Rc And Output  docker ${params} logs --tail=1 --follow ${containerID}
	#Should Be Equal As Integers  ${rc}  0
	#${linecount}=  Get Line Count  ${Output}
	#Should Be Equal As Integers  ${linecount}  3
    
Docker logs with timestamps and since certain time
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;'
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
	Run  Sleep 6, wait for container to finish
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --since=1s ${containerID}
	Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support timestampped logs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --timestamps ${containerID}
	Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support timestampped logs
    
Docker logs non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such container: fakeContainer