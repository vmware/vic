*** Settings ***
Documentation  Test 1-7 - Docker Stop
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Stop an already stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox sleep 30
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
    ${before}=  Get Current Date
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    ${after}=  Get Current Date
    Should Be Equal As Integers  ${rc}  0
    ${result}=  Subtract Date From Date  ${after}  ${before}

    ${status}=  Get State Of Github Issue  1320
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-7-Docker-Stop.robot needs to be updated now that Issue #1320 has been resolved
    Log  Issue \#1320 is blocking implementation  WARN
    #Should Be True  ${9} < ${result} < ${11}

Stop a container with a time limit
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
    ${status}=  Get State Of Github Issue  1398
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-7-Docker-Stop.robot needs to be updated now that Issue #1398 has been resolved
    Log  Issue \#1398 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on=true ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    #Should Be Equal As Integers  ${rc}  0