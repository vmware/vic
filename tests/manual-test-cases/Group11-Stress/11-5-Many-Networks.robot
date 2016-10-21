*** Settings ***
Documentation  Test 11-5-Many-Networks
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Many Networks
    Log To Console  Create 1000 networks
    :FOR  ${idx}  IN RANGE  0  1000
    \   Log To Console  \nCreate network ${idx}
    \   ${rc}=  Run And Return Rc  docker ${params} network create net${idx}
    \   Should Be Equal As Integers  ${rc}  0

    ${out}=  Run  docker ${params} pull busybox
    ${container}=  Run  docker ${params} create --net=net999 busybox ping -C2 google.com 
    ${out}=  Run  docker ${params} start ${container} 
    Should Contain  ${out}  2 packets transmitted, 2 received

    Run Regression Tests