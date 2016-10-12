*** Settings ***
Documentation  Test 5-9 - Private Registry
Resource  ../../resources/Util.robot
Suite Setup  Private Registry Setup
Suite Teardown  Private Registry Cleanup

*** Keywords ***
Private Registry Setup
    Install VIC Appliance To Test Server  certs=${true}  vol=default
    ${rc}  ${output}=  Run And Return Rc And Output  docker run -d -p 5000:5000 --name registry registry
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker tag busybox localhost:5000/busybox:latest
    Should Be Equal As Integers  ${rc}  0  
    ${rc}  ${output}=  Run And Return Rc And Output  docker push localhost:5000/busybox
    Should Be Equal As Integers  ${rc}  0
    
Private Registry Cleanup
    Cleanup VIC Appliance On Test Server
    ${rc}  ${output}=  Run And Return Rc And Output  docker rm -f registry
    Should Be Equal As Integers  ${rc}  0

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
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  172.17.0.1:5000/busybox
    
Standalone imagec pull from private registry
    ${result}=  Run Process  ${bin-dir}/imagec -standalone -insecure-allow-http -reference localhost:5000/busybox  shell=True  cwd=/
    Log  ${result.stdout}
    Log  ${result.stderr}
    Should Contain  ${result.stdout}  "status":"Download complete"
    Should Contain  ${result.stdout}  "status":"Digest: sha256:
