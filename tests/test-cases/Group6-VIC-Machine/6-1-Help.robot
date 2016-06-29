*** Settings ***
Documentation  Test 6-1 - Verify Help
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${ret}=  Run  bin/vic-machine-linux delete -h
    Should Contain  ${ret}  vic-machine-linux delete - Delete VCH