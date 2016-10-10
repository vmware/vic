*** Settings ***
Documentation  Test 1-17 - Docker Network Connect
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
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

Connect containers to multiple networks overlapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross1-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross1-network2
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross1-network --name cross1-container busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect cross1-network2 ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross1-network --name cross1-container2 busybox ping -c2 cross1-container
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect cross1-network2 ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow cross1-container2
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  2 packets transmitted, 2 packets received

Connect containers to multiple networks non-overlapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross2-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross2-network2
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross2-network --name cross2-container busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross2-network2 --name cross2-container2 busybox ping -c2 cross2-container
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow cross2-container2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  2 packets transmitted, 2 packets received

#Connect containers to multiple networks non-overlapping with a bridge container
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross3-network
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create cross3-network2
#    Should Be Equal As Integers  ${rc}  0
#
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross3-network --name cross3-container busybox /bin/top
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
#    Should Be Equal As Integers  ${rc}  0
#
#    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross3-network2 --name cross3-container2 busybox /bin/top
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
#    Should Be Equal As Integers  ${rc}  0
#
#    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create --net cross3-network --name cross3-container3 busybox /bin/sh -c "ping -c2 cross3-container && ping -c2 cross3-container2"
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network connect cross3-network2 ${containerID}
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
#    Should Be Equal As Integers  ${rc}  0
#    Wait Until Container Stops  ${containerID}
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow cross3-container3
#    Should Be Equal As Integers  ${rc}  0
#    Should Contain X Times  ${output}  2 packets transmitted, 2 packets received  2
