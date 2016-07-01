*** Settings ***
Documentation  Test 1-18 - Docker Network RM
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic network remove
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network rm test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test-network
    
Multiple network remove
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create test-network2
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create test-network3
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network rm test-network2 ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test-network2
    ${status}=  Get State Of Github Issue  1318
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-18-Docker-Network-RM.robot needs to be updated now that Issue #1318 has been resolved
    Log  Issue \#1318 is blocking implementation  WARN
    #Should Not Contain  ${output}  test-network3

Remove already removed network
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network rm test-network
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: network test-network not found
    
Remove network with running container
    ${status}=  Get State Of Github Issue  1235
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-18-Docker-Network-RM.robot needs to be updated now that Issue #1235 has been resolved
    Log  Issue \#1235 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create test-network
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect test-network ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network rm test-network
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error response from daemon: network test-network has active endpoints
    
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network rm test-network
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  test-network