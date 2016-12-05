*** Settings ***
Documentation  Test 5-9 - Private Registry
Resource  ../../resources/Util.robot
Suite Setup  Private Registry Setup
Suite Teardown  Private Registry Cleanup

*** Keywords ***
Private Registry Setup
    ${dockerHost}=  Get Environment Variable  DOCKER_HOST  ${SPACE}
    Remove Environment Variable  DOCKER_HOST
    ${rc}  ${output}=  Run And Return Rc And Output  docker run -d -p 5000:5000 --name registry registry
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker tag busybox localhost:5000/busybox:latest
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker push localhost:5000/busybox
    Should Be Equal As Integers  ${rc}  0
    Set Environment Variable  DOCKER_HOST  ${dockerHost}

Private Registry Cleanup
    ${dockerHost}=  Get Environment Variable  DOCKER_HOST  ${SPACE}
    Remove Environment Variable  DOCKER_HOST
    ${rc}  ${output}=  Run And Return Rc And Output  docker rm -f registry
    Should Be Equal As Integers  ${rc}  0
    Set Environment Variable  DOCKER_HOST  ${dockerHost}

Pull image
    [Arguments]  ${image}
    Log To Console  \nRunning docker pull ${image}...
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Digest:
    Should Contain  ${output}  Status:
    Should Not Contain  ${output}  No such image:

*** Test Cases ***
Pull an image from non-default repo
    Install VIC Appliance To Test Server  vol=default --insecure-registry 172.17.0.1:5000
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  172.17.0.1:5000/busybox
    Cleanup VIC Appliance On Test Server

Pull image from non-whitelisted repo
    Install VIC Appliance To Test Server  vol=default
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull 172.17.0.1:5000/busybox
    Should Contain  ${output}  Error response from daemon: Head https://172.17.0.1:5000/v2/: http: server gave HTTP response to HTTPS client
