*** Settings ***
Documentation  Test 1-03 - Docker Images
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple images
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull alpine
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull alpine:3.2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull alpine:3.1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain X Times  ${output}  alpine  3

All images
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images -a
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain X Times  ${output}  alpine  3

Quiet images
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images -q
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Not Contain  ${output}  alpine
    @{lines}=  Split To Lines  ${output}
    Length Should Be  ${lines}  3
    Length Should Be  @{lines}[1]  12

No-trunc images
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images --no-trunc
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain X Times  ${output}  alpine  3
    @{lines}=  Split To Lines  ${output}
    @{line}=  Split String  @{lines}[2]
    Length Should Be  @{line}[2]  64

Specific images
    ${status}=  Get State Of Github Issue  1035
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-3-Docker-Images.robot needs to be updated now that Issue #1035 has been resolved
#    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images alpine:3.1
#    Should Be Equal As Integers  ${rc}  0
#    Should Not Contain  ${output}  Error
#    Should Contain  ${output}  3.1
#    Should Contain X Times  ${output}  alpine  1