*** Settings ***
Documentation  Test 1-27 - Docker Login
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  ${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker login and pull from docker.io
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull victest/busybox
    Should Be Equal As Integers  ${rc}  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull victest/public-hello-world
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull victest/busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logout
    Should Be Equal As Integers  ${rc}  0
	
	
