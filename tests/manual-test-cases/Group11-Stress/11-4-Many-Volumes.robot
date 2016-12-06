*** Settings ***
Documentation  Test 11-4-Many-Volumes
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker volume create 1000 volumes rapidly
    ${pids}=  Create List

    Log To Console  \nCreate 1000 volumes rapidly
    :FOR  ${idx}  IN RANGE  0  1000
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} volume create --name\=multiple${idx} --opt Capacity\=32MB  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0

    Run Regression Tests
