*** Settings ***
Documentation  Test 1-22 - Docker Volume RM
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple volume rm
    ${status}=  Get State Of Github Issue  1720
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-22-Docker-Volume-RM.robot needs to be updated now that Issue #1720 has been resolved
    Log  Issue \#1720 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test
    #Should Be Equal As Integers  ${rc}  0
    #Should Be Equal As Strings  ${output}  test
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test2
    #Should Be Equal As Integers  ${rc}  0
    #Should Be Equal As Strings  ${output}  test2
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume rm test2
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  test2
    
#Volume rm when in use
#    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create -v test:/test busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume rm test
#    Should Be Equal As Integers  ${rc}  1
#    Should Contain  ${output}  Error response from daemon: Conflict: remove test2: volume is in use - ${containerID}

#Volume rm invalid volume
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume rm test3
#    Should Be Equal As Integers  ${rc}  1
#    Should Contain  ${output}  Error response from daemon: get test3: no such volume
    
#Volume rm freed up volume
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test4
#    Should Be Equal As Integers  ${rc}  0
#    Should Be Equal As Strings  ${output}  test4
#    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create -v test4:/test4 busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${containerID}
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume rm test4
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume ls
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  test4