*** Settings ***
Documentation  Test 14-1 - Longevity
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Longevity
    # Each loop should take between 1 and 2 hours
    :FOR  ${idx}  IN RANGE  0  48
    \   ${rand}=  Evaluate  random.randint(10, 50)  modules=random
    \   Log To Console  \nLoop: ${idx}
    \   Install VIC Appliance To Test Server
    \   Repeat Keyword  ${rand} times  Run Regression Tests
    \   Cleanup VIC Appliance On Test Server
    \   ${rand}=  Evaluate  random.randint(10, 50)  modules=random
    \   Install VIC Appliance To Test Server  certs=${true}
    \   Repeat Keyword  ${rand} times  Run Regression Tests
    \   Cleanup VIC Appliance On Test Server