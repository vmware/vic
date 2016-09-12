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

    # ${status}=  Get State Of Github Issue  1459
    # Run Keyword If  '${status}' == 'closed'  Fail  Test 1-24-Docker-Link.robot needs to be updated now that Issue #1459 has been resolved
    # Log  Issue \#1459 is blocking implementation  WARN
    # the name
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi busybox ping -c1 first
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error

    # the link
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi --link first:1st busybox ping -c1 1st
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error

    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it -d --net jedi --net-alias 2nd busybox
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error

    # the alias
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --net jedi busybox ping -c1 2nd
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error