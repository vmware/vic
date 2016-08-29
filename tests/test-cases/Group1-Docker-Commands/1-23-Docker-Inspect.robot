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

Docker inspect container check cmd and image name
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/bash
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${config}=  Get From Dictionary  ${output[0]}  Config
    ${image}=  Get From Dictionary  ${config}  Image
    Should Contain  ${image}  busybox
    ${cmd}=  Get From Dictionary  ${config}  Cmd
    Should Be Equal As Strings  ${cmd}  [u'/bin/bash']

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