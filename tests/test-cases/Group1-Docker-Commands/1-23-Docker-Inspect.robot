*** Settings ***
Documentation  Test 1-23 - Docker Inspect
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker inspect of image
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying type
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect --type=image busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying incorrect type    
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect --type=container busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such container: busybox
    
Simple docker inspect of container
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0    
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container specifying type
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect --type=container ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container specifying incorrect type
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect --type=image ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such image: ${container}
    
Docker inspect invalid object
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect fake
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such image or container: fake