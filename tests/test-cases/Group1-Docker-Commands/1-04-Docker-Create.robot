*** Settings ***
Documentation  Test 1-04 - Docker Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple creates
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -t -i busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name test1 busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create with anonymous volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v /var/log busybox ls /var/log
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs --follow ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create with named volume
    ${disk-size}=  Run  docker %{VCH-PARAMS} logs $(docker %{VCH-PARAMS} start $(docker %{VCH-PARAMS} create -v test-named-vol:/testdir busybox /bin/df -Ph) && sleep 10) | grep by-label | awk '{print $2}'
    Should Be Equal As Strings  ${disk-size}  975.9M

Create with a directory as a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v /dir:/dir busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: Bad request error from portlayer: vSphere Integrated Containers does not support mounting directories as a data volume.

Create simple top example
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create fakeimage image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create fakeimage
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: image library/fakeimage not found

Create fakeImage repository
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create fakeImage
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error parsing reference: "fakeImage" is not a valid repository/tag

Create and start named container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name busy1 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start busy1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create linked containers that can ping
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull debian
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --link busy1:busy1 --name busy2 debian ping -c2 busy1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start busy2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} wait busy2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs busy2
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  2 packets transmitted, 2 packets received

Create a container after the last container is removed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${cid}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${cid}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${cid}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${cid2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${cid2}  Error

Create a container from an image that has not been pulled yet
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create alpine bash
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Create a container with no command specified
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create alpine
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No command specified

Create a container with custom CPU count
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it --cpuset-cpus 3 busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  govc datastore.ls |grep ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info ${id} |awk '/CPU:/ {print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  3

Create a container with custom amount of memory in GB
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -m 4G busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  govc datastore.ls |grep ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info ${id} |awk '/Memory:/ {print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  4096MB

Create a container with custom amount of memory in MB
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -m 2048M busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  govc datastore.ls |grep ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info ${id} |awk '/Memory:/ {print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  2048MB

Create a container with custom amount of memory in KB
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -m 2097152K busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  govc datastore.ls |grep ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info ${id} |awk '/Memory:/ {print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  2048MB

Create a container with custom amount of memory in Bytes
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -m 2147483648B busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  govc datastore.ls |grep ${id}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info ${id} |awk '/Memory:/ {print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  2048MB
