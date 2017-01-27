*** Settings ***
Documentation  Test 1-29 - Docker Checkpoint
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker checkpoint create
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} create --name=test-busybox busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} checkpoint create test-busybox new-checkpoint
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet implement checkpointing

Docker checkpoint ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} checkpoint ls test-busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet implement checkpointing

Docker checkpoint rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} checkpoint rm test-busybox new-checkpoint
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No such container

