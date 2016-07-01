*** Settings ***
Documentation  Test 1-12 - Docker RMI
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic docker remove image
    ${status}=  Get State Of Github Issue  437
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-12-Docker-RMI.robot needs to be updated now that Issue #437 has been resolved
    Log  Issue \#437 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  busybox
    
#Remove image with a removed container
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${output}
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  busybox
        
#Remove image with a container    
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi busybox
#    Should Be Equal As Integers  ${rc}  1
#    Should Contain  ${output}  Failed to remove image (busybox): Error response from daemon: conflict: unable to remove repository reference "busybox" (must force) - container
#    Should Contain  ${output}  is using its referenced image
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
#    Should Be Equal As Integers  ${rc}  0
#    Should Contain  ${output}  busybox
    
#Force remove image with a container
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi -f busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  busybox
    
#Remove a fake image
#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi fakeImage
#    Should Be Equal As Integers  ${rc}  1
#    Should Contain  ${output}  Failed to remove image (fakeImage): Error response from daemon: No such image: fakeImage:latest