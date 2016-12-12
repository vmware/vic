*** Settings ***
Documentation  Test 1-18 - Docker Network RM
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic network remove
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network rm test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test-network
    
Multiple network remove
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network2
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network3
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network rm test-network2 ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test-network2
    Should Not Contain  ${output}  test-network3

Remove already removed network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network rm test-network
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: network test-network not found

Remove network with running container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network connect test-network ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network rm test-network
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: test-network has active endpoints
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network rm test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test-network
