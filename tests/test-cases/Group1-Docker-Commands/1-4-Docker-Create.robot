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

Create with anonymous volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -v /var/log busybox ls /var/log
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create with named volume
    ${disk-size}=  Run  docker ${params} logs $(docker ${params} start $(docker ${params} create -v test-named-vol:/testdir busybox /bin/df -Ph) && sleep 10) | grep by-label | awk '{print $2}'
    Should Be Equal As Strings  ${disk-size}  975.9M

Create with a directory as a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -v /dir:/dir busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: Bad request error from portlayer: vSphere Integrated Containers does not support mounting directories as a data volume.

Create simple top example
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create fakeimage image
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create fakeimage
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: image library/fakeimage not found

Create fakeImage repository
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create fakeImage
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error parsing reference: "fakeImage" is not a valid repository/tag

Create and start named container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --name busy1 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start busy1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create linked containers that can ping
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create --link busy1:busy1 --name busy2 busybox ping -c2 busy1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start busy2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --follow busy2
    Should Be Equal As Integers  ${rc}  0
    ${status}=  Get State Of Github Issue  1459
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-4-Docker-Create.robot needs to be updated now that Issue #1459 has been resolved
    Log  Issue \#1459 is blocking implementation  WARN
    #Should Contain  ${output}  2 packets transmitted, 2 received

Create a container after the last container is removed
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${cid}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${cid}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${cid}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${cid2}=  Run And Return Rc And Output  docker ${params} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${cid2}  Error

Create a container from an image that has not been pulled yet
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create alpine bash
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create a container with no command specified
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create alpine
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No command specified
