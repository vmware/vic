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
    ${status}=  Get State Of Github Issue  1562
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1562 has been resolved
    Log  Issue \#1562 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error response from daemon: A volume named test already exists. Choose a different volume name.
    
Docker volume create volume with bad driver
    ${status}=  Get State Of Github Issue  1564
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1564 has been resolved
    Log  Issue \#1564 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create -d fakeDriver --name=test2
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error looking up volume plugin fakeDriver: plugin not found
    
Docker volume create with bad volumestore
    ${status}=  Get State Of Github Issue  1561
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1561 has been resolved
    Log  Issue \#1561 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test3 --opt VolumeStore=fakeStore
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error looking up volume store fakeStore: datastore not found

Docker volume create with specific capacity
    ${status}=  Get State Of Github Issue  1565
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1565 has been resolved
    Log  Issue \#1565 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test4 --opt Capacity=100
    #Should Be Equal As Integers  ${rc}  0
    #Should Be Equal As Strings  ${output}  test4
    #${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls -json=true test/VIC/volumes/test4
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  "FileSize":100
    
Docker volume create with zero capacity
    ${status}=  Get State Of Github Issue  1562
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1562 has been resolved
    Log  Issue \#1562 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test5 --opt Capacity=0
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error
    
Docker volume create with negative one capacity
    ${status}=  Get State Of Github Issue  1562
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1562 has been resolved
    Log  Issue \#1562 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test6 --opt Capacity=-1
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error    
    
Docker volume create with capacity too big
    ${status}=  Get State Of Github Issue  1562
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1562 has been resolved
    Log  Issue \#1562 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test7 --opt Capacity=2147483647
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error
    
Docker volume create with capacity exceeding int size
    ${status}=  Get State Of Github Issue  1562
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1562 has been resolved
    Log  Issue \#1562 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test8 --opt Capacity=9999999999
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error
    
Docker volume create with possibly invalid name
    ${status}=  Get State Of Github Issue  1563
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-19-Docker-Volume-Create.robot needs to be updated now that Issue #1563 has been resolved
    Log  Issue \#1563 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test???
    #Should Be Equal As Integers  ${rc}  1
    #Should Be Equal As Strings  ${output}  Error response from daemon: create test???: "test???" includes invalid characters for a local volume name, only "[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed
    
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