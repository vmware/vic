*** Settings ***
Documentation  Test 7-1 - Regression
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags  regression

*** Test Cases ***
Regression test
    Run Regression Tests