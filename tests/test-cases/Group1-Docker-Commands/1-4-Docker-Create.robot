*** Settings ***
Documentation  Test 1-4 - Docker Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple creates
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -t -i busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name test1 busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create with volume
    ${status}=  Get State Of Github Issue  806
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-4-Docker-Create.robot needs to be updated now that Issue #806 has been resolved
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -v /var/log:/var/log busybox ls /var/log
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${output}
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error

Create simple top example
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create fakeimage image
    ${status}=  Get State Of Github Issue  1036
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-4-Docker-Create.robot needs to be updated now that Issue #1036 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create fakeimage
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error: image library/fakeimage not found

Create fakeImage repository
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create fakeImage
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error parsing reference: "fakeImage" is not a valid repository/tag

Create and start named container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name busy1 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create linked containers that can ping
    ${status}=  Get State Of Github Issue  430
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-4-Docker-Create.robot needs to be updated now that Issue #430 has been resolved
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --link busy1:busy1 --name busy2 busybox ping -c2 busy1
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${output}
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    Should Contain  ${output}  2 packets transmitted, 2 received