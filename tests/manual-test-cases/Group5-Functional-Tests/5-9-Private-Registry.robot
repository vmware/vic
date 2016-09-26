*** Settings ***
Documentation  Test 5-9 - Private Registry
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Pull image
    [Arguments]  ${image}
    Log To Console  \nRunning docker pull ${image}...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull ${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Digest:
    Should Contain  ${output}  Status:
    Should Not Contain  ${output}  No such image:

*** Test Cases ***
Pull an image from non-default repo
    ${output}=  Run  docker run -d -p 5000:5000 --name registry registry
    Log  ${output}
    ${output}=  Run  docker pull busybox
    Log  ${output}
    ${output}=  Run  docker tag busybox localhost:5000/busybox:latest
    Log  ${output}
    ${output}=  Run  docker push localhost:5000/busybox
    Log  ${output}
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  172.17.0.1:5000/busybox