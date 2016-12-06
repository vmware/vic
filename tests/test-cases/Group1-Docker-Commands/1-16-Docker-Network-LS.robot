*** Settings ***
Documentation  Test 1-16 - Docker Network LS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic network ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  bridge

Docker network ls -q
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls -q
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  NAME
    Should Not Contain  ${output}  DRIVER
    Should Not Contain  ${output}  bridge

Docker network ls -f
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls -f name=bridge
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  bridge
    @{lines}=  Split To Lines  ${output}
    Length Should Be  ${lines}  2

Docker network ls --no-trunc
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls --no-trunc
    Should Be Equal As Integers  ${rc}  0
    @{lines}=  Split To Lines  ${output}
    @{line}=  Split String  @{lines}[1]
    Length Should Be  @{line}[0]  64

Docker network ls -f fake network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls -f name=fakeName
    Should Be Equal As Integers  ${rc}  0
    @{lines}=  Split To Lines  ${output}
    Length Should Be  ${lines}  1
    Should Contain  @{lines}[0]  NAME
