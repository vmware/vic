*** Settings ***
Documentation  Test 1-24 - Docker Link
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Link and alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name first busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${_}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${_}  Error:
    
    # the name
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox ping -c3 first
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${_}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${_}  Error:

    # the link
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --link first:1st busybox ping -c3 1st
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    
    ${rc}  ${_}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${_}  Error:

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --net-alias 2nd busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    
    ${rc}  ${_}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${_}  Error:

    # the alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox ping -c3 2nd
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${_}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${_}  Error:
