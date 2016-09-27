*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${ret}=  Run  ${bin-dir}/imagec -help
    Should Contain  ${ret}  Usage of
    Should Contain  ${ret}  bin/imagec:
    Should Contain  ${ret}  -version
    Should Contain  ${ret}  Show version info