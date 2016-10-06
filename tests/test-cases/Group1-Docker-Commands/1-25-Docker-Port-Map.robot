*** Settings ***
Documentation  Test 1-25 - Docker Port Map
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create container with port mappings
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop webserver
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl ${vch-ip}:10000 --connect-timeout 5
    Should Not Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl ${vch-ip}:10001 --connect-timeout 5
    Should Not Be Equal As Integers  ${rc}  0
    
Create container with conflicting port mapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8083:80 --name webserver2 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8083:80 --name webserver3 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver3
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  port 8083 is not available

Create container with port range
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8081-8088:80 --name webserver5 nginx
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  host port ranges are not supported for port bindings

Create container with host ip
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10.10.10.10:8088:80 --name webserver5 nginx
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  host IP is not supported for port bindings

Create container without specifying host port
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 6379 --name test-redis redis:alpine
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start test-redis
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop test-redis
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error