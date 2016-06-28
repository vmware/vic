*** Settings ***
Documentation  Test 1-19 - Docker Volume Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker volume create
    Log  Not really implemented  WARN
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --opt VolumeStore=datastore1 --opt Capacity=1
    #Should Be Equal As Integers  ${rc}  1

Docker volume create named volume
    Log  Not really implemented  WARN
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test --opt VolumeStore=datastore1 --opt Capacity=1
    #Should Be Equal As Integers  ${rc}  1

Docker volume create remote volume
    Log  Not really implemented  WARN
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create -d remote --name=test --opt VolumeStore=datastore1 --opt Capacity=1
    #Should Be Equal As Integers  ${rc}  1