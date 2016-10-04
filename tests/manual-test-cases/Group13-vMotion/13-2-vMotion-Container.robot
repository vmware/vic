*** Settings ***
Documentation  Test 13-2 - vMotion Container
Resource  ../../resources/Nimbus-Util.robot
Suite Setup  Create a VSAN Cluster
Suite Teardown  Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  *

*** Test Cases ***
Test
    Install VIC Appliance To Test Server  ${false}  default
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container2}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container3}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container3}
    Should Be Equal As Integers  ${rc}  0

    vMotion A VM  ${vch-name}/*-${container1}
    vMotion A VM  ${vch-name}/*-${container2}
    vMotion A VM  ${vch-name}/*-${container3}
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container1}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${container3}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container3}
    Should Be Equal As Integers  ${rc}  0

    Cleanup VIC Appliance On Test Server