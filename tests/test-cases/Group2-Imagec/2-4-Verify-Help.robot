*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${ret}=  Run  bin/imagec -help
    Should Contain  ${ret}  Usage of bin/imagec: