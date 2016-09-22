*** Settings ***
Documentation  Test 1-5 - Docker Start
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple start
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
Start with attach and interactive
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start -ai ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
Start from image that has no PATH
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull vmware/photon
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it vmware/photon
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
Start non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer
    Should Contain  ${output}  Error: failed to start containers: fakeContainer
Start with no ethernet card
    # Testing that port layer doesn't hang forever if tether fails to initialize (see issue #2327)
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${name}=  Generate Random String  15
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name ${name} busybox date
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc device.remove -vm ${name}-* ethernet-0
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${name}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  unable to wait for process launch status
    Should Not Contain  ${output}  context deadline exceeded
Restart a stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it busybox /bin/ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
