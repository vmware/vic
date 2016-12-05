*** Settings ***
Documentation  Test 11-8-Image-Pull-Stress
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Image Pull Stress
    ${pids}=  Create List

    Log To Console  \nRapidly pull images
    :FOR  ${idx}  IN RANGE  0  50
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} pull busybox  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nRapidly pull images
    :FOR  ${idx}  IN RANGE  0  25
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} pull alpine  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nRapidly pull images
    :FOR  ${idx}  IN RANGE  0  25
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} pull ubuntu  shell=True
    \   Append To List  ${pids}  ${pid}

    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0

    Run Regression Tests