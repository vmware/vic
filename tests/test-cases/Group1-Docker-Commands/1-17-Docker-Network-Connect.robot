*** Settings ***
Documentation  Test 1-17 - Docker Network Connect
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Connect container to a new network
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox ip -4 addr show eth0
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect test-network ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${ips}=  Get Lines Containing String  ${output}  inet
    @{lines}=  Split To Lines  ${ips}
    Length Should Be  ${lines}  2

Connect to non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect test-network fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: container fakeContainer not found

Connect to non-existent network
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name connectTest3 busybox ifconfig
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect fakeNetwork connectTest3
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: network fakeNetwork not found
