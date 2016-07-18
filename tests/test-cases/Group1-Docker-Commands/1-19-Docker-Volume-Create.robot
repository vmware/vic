*** Settings ***
Documentation  Test 1-19 - Docker Volume Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker volume create
    Log  Not supported yet  WARN
    # Auto-named volumes not supported yet
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --opt VolumeStore=default --opt Capacity=1
    #Should Be Equal As Integers  ${rc}  0

Docker volume create named volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create --name=test --opt VolumeStore=default --opt Capacity=1
    Should Be Equal As Integers  ${rc}  0

Docker volume create remote volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} volume create -d remote --name=test2 --opt VolumeStore=default --opt Capacity=1
    Should Be Equal As Integers  ${rc}  0
    
    
# test when default volume store doesn't exist