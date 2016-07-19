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

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it -d --net jedi --name first busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    
    # the name
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it --net jedi busybox ping -c64 first
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # the link
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it --net jedi --link first:1st busybox ping -c64 1st
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it -d --net jedi --net-alias 2nd busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    
    # the alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it --net jedi busybox ping -c64 2nd
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
