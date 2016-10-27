*** Settings ***
Documentation  Test 1-24 - Docker Link
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Link and alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create jedi
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull debian
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error    

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it -d --net jedi --name first busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi debian ping -c1 first
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # cannot reach first from another network
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run debian ping -c1 first
    Should Not Be Equal As Integers  ${rc}  0

    # the link
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi --link first:1st debian ping -c1 1st
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # cannot reach first using c1 from another container
    # first run a container that has the alias "c1" for the "first" container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -itd --net jedi --link first:1st busybox 
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    # check if we can use alias "c1" from another container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi debian ping -c1 1st
    Should Not Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it -d --net jedi --net-alias 2nd busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # the alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi debian ping -c1 2nd
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
