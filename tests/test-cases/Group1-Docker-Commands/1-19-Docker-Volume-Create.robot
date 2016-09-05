*** Settings ***
Documentation  Test 1-19 - Docker Volume Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker volume create
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create
    Should Be Equal As Integers  ${rc}  0

Docker volume create named volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test

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
    
Docker volume create with zero capacity
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test5 --opt Capacity=0
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error
    
Docker volume create with negative one capacity
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test6 --opt Capacity=-1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error    
    
Docker volume create with capacity too big
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test7 --opt Capacity=9223372036854775808
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error
    
Docker volume create with capacity exceeding int size
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test8 --opt Capacity=9999999999999999999
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error
    
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