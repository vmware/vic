*** Settings ***
Documentation  Test 11-2-Container-Boot-Storm
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Container Boot Storm
    ${containers}=  Create List
    ${pids}=  Create List
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox

    Log To Console  \nCreate 100 containers
    :FOR  ${idx}  IN RANGE  0  100
    \   ${id}=  Run  docker %{VCH-PARAMS} create busybox /bin/top
    \   Append To List  ${containers}  ${id}
    
    Log To Console  \nRapidly start each container
    :FOR  ${id}  IN  @{containers}
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} start ${id}  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0

    Run Regression Tests