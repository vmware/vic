*** Settings ***
Documentation  Test 1-19 - Docker Volume Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker volume create
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create
    Should Be Equal As Integers  ${rc}  0
    Set Suite Variable  ${ContainerID}  unnamedSpecVol
    Run  docker ${params} run --name ${ContainerID} -d -v ${output}:/mydata busybox /bin/df -Ph
    ${ContainerRC}=  Run  docker ${params} wait ${ContainerID}
    Should Be Equal As Integers  ${ContainerRC}  0
    ${disk-size}=  Run  docker ${params} logs ${ContainerID} | grep by-label | awk '{print $2}'
    Should Be Equal As Strings  ${disk-size}  975.9M

Docker volume create named volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test
    Set Suite Variable  ${ContainerID}  specVol
    Run  docker ${params} run --name ${ContainerID} -d -v ${output}:/mydata busybox /bin/df -Ph
    ${ContainerRC}=  Run  docker ${params} wait ${ContainerID}
    Should Be Equal As Integers  ${ContainerRC}  0
    ${disk-size}=  Run  docker ${params} logs ${ContainerID} | grep by-label | awk '{print $2}'
    Should Be Equal As Strings  ${disk-size}  975.9M

Docker volume create image volume
    Set Suite Variable  ${ContainerID}  imageVol
    Run  docker ${params} run --name ${ContainerID} -d mongo /bin/df -Ph
    ${ContainerRC}=  Run  docker ${params} wait ${ContainerID}
    Should Be Equal As Integers  ${ContainerRC}  0
    ${disk-size}=  Run  docker ${params} logs ${ContainerID} | grep by-label | awk '{print $2}'
    Should Contain  ${disk-size}  976M

Docker volume create anonymous volume
    Set Suite Variable  ${ContainerID}  anonVol
    Run  docker ${params} run --name ${ContainerID} -d -v /mydata busybox /bin/df -Ph
    ${ContainerRC}=  Run  docker ${params} wait ${ContainerID}
    Should Be Equal As Integers  ${ContainerRC}  0
    ${disk-size}=  Run  docker ${params} logs ${ContainerID} | grep by-label | awk '{print $2}'
    Should Contain  ${disk-size}  975.9M

Docker volume create already named volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: A volume named test already exists. Choose a different volume name.

Docker volume create volume with bad driver
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create -d fakeDriver --name=test2
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error looking up volume plugin fakeDriver: plugin not found

Docker volume create with bad volumestore
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test3 --opt VolumeStore=fakeStore
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No volume store named (fakeStore) exists

Docker volume create with specific capacity
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test4 --opt Capacity=100000
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test4
    Set Suite Variable  ${ContainerID}  capacityVol
    Run  docker ${params} run --name ${ContainerID} -d -v ${output}:/mydata busybox /bin/df -Ph
    ${ContainerRC}=  Run  docker ${params} wait ${ContainerID}
    Should Be Equal As Integers  ${ContainerRC}  0
    ${disk-size}=  Run  docker ${params} logs ${ContainerID} | grep by-label | awk '{print $2}'
    Should Be Equal As Strings  ${disk-size}  96.0G

Docker volume create with zero capacity
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test5 --opt Capacity=0
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: bad driver value - Invalid size: 0

Docker volume create with negative one capacity
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test6 --opt Capacity=-1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: bad driver value - Invalid size: -1

Docker volume create with capacity too big
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test7 --opt Capacity=9223372036854775808
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: bad driver value - Capacity value too large: 9223372036854775808

Docker volume create with capacity exceeding int size
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test8 --opt Capacity=9999999999999999999
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: bad driver value - Capacity value too large: 9999999999999999999

Docker volume create with possibly invalid name
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test???
    Should Be Equal As Integers  ${rc}  1
    Should Be Equal As Strings  ${output}  Error response from daemon: volume name "test???" includes invalid characters, only "[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed

Docker volume create 100 volumes rapidly
    ${pids}=  Create List

    # Create 100 volumes rapidly
    :FOR  ${idx}  IN RANGE  0  100
    \   ${pid}=  Start Process  docker ${params} volume create --name\=multiple${idx} --opt Capacity\=512MB  shell=True
    \   Append To List  ${pids}  ${pid}

    # Wait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   ${res}=  Wait For Process  ${pid}
    \   Log  ${res.stdout} ${res.stderr}
    \   Should Be Equal As Integers  ${res.rc}  0
