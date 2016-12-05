*** Settings ***
Documentation  Test 1-22 - Docker Volume RM
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple volume rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=test
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=test2
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test2
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm test2
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test2

Volume rm when in use
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v test:/test busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm test
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: volume test in use by

Volume rm invalid volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm test3
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: Get test3: no such volume

Volume rm freed up volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=test4
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test4
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v test4:/test4 busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${containerID}
   Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm test4
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  test4
