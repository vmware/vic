*** Settings ***
Documentation  Test 11-6-Stop-Storm
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Stop Storm
    ${containers}=  Create List
    ${pids}=  Create List
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox

    Log To Console  \nCreate 100 containers
    :FOR  ${idx}  IN RANGE  0  100
    \   Log To Console  \nCreating container ${idx}
    \   ${id}=  Run  docker %{VCH-PARAMS} run busybox /bin/top
    \   Append To List  ${containers}  ${id}

    Log To Console  \nRapidly stop each container
    :FOR  ${id}  IN  @{containers}
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} stop ${id}  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0

    Run Regression Tests