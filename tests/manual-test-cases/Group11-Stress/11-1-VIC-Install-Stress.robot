*** Settings ***
Documentation  Test 11-1-VIC-Install-Stress
Resource  ../../resources/Util.robot
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
VIC Install Stress
    :FOR  ${idx}  IN RANGE  0  100
    \   Log To Console  \nLoop ${idx+1}
    \   Install VIC Appliance To Test Server  vol=default %{STATIC_VCH_OPTIONS}
    \   Cleanup VIC Appliance On Test Server
    
    Install VIC Appliance To Test Server  vol=default %{STATIC_VCH_OPTIONS}
    Run Regression Tests