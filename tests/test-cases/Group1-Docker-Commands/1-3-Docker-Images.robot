*** Settings ***
Documentation  Test 1-3 - Docker Images 
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Images
    Log To Console  \nRunnning docker images command...
    ${output}=  Run  docker ${params} images
    Log  ${output}
    Should contain  ${output}  IMAGE ID