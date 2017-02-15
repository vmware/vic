*** Settings ***
Documentation  Test 1-15 - Docker Network Create
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic network create
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  test-network

Create already created network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create test-network
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  already exists 

Create overlay network
    ${status}=  Get State Of Github Issue  1222
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-15-Docker-Network-Create.robot needs to be updated now that Issue #1222 has been resolved
    Log  Issue \#1222 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create -d overlay test-network2
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error response from daemon: failed to parse pool request for address space "GlobalDefault" pool "" subpool "": cannot find address space GlobalDefault (most likely the backing datastore is not configured)