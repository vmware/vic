*** Settings ***
Documentation  Test 11-4-Many-Volumes
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker volume create 1000 volumes rapidly
    ${pids}=  Create List

    # Create 1000 volumes rapidly
    :FOR  ${idx}  IN RANGE  0  1000
    \   ${pid}=  Start Process  docker ${params} volume create --name\=multiple${idx} --opt Capacity\=32MB  shell=True
    \   Append To List  ${pids}  ${pid}

    # Wait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0

    Run Regression Tests
