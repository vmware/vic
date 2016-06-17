*** Settings ***
Documentation  Test 1-2 - Docker Pull
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

*** Test Cases ***
Pull nginx
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx

Pull busybox
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  busybox

Pull ubuntu
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  ubuntu

Pull registry
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  registry

Pull swarm
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  swarm
    
Pull redis
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  redis
    
Pull mongo
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  mongo
    
Pull mysql
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  mysql
    
Pull node
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  node
    
Pull postgres
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  postgres

Pull non-default tag
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx:alpine
    
Pull an image based on digest
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx@sha256:7281cf7c854b0dfc7c68a6a4de9a785a973a14f1481bc028e2022bcd6a8d9f64

Pull an image from non-default repo
    #${result}=  Run Process  docker run -d -p 5000:5000 --name registry registry:2  shell=True
    #Log to console  ${result.stdout}
    #Log to console  ${result.stderr}
    #${result}=  Run Process  docker pull nginx  shell=True
    #Log to console  ${result.stdout}
    #Log to console  ${result.stderr}
    #${result}=  Run Process  docker tag nginx localhost:5000/testImage  shell=True
    #Log to console  ${result.stdout}
    #Log to console  ${result.stderr}
    #Wait Until Keyword Succeeds  5x  15 seconds  Pull image  localhost:5000/testImage
    Log  Not quite working yet...  WARN
    Log To Console  Not quite working yet...
    
Pull an image with all tags
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  --all-tags nginx
    
Pull non-existent image
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull fakebadimage
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    ${status}=  Get State Of Github Issue  757
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-2-Docker-Pull.robot needs to be updated now that Issue #757 has been resolved
    #Should contain  ${output}  image library/fakebadimage not found
    
Pull image from non-existent repo
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull fakebadrepo.com:9999/ubuntu
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  no such host

Pull image that already has been pulled
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  alpine
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  alpine