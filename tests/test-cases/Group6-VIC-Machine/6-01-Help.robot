*** Settings ***
Documentation  Test 6-01 - Verify Help
Resource  ../../resources/Util.robot

*** Test Cases ***
Delete help basic
    ${ret}=  Run  bin/vic-machine-linux delete -h
    Should Contain  ${ret}  vic-machine-linux delete - Delete VCH
