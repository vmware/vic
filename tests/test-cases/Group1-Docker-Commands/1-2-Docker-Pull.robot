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
    Should Not Contain  ${output}  No such image:

*** Test Cases ***
Pull nginx
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx

Pull busybox
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  busybox

Pull ubuntu
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  ubuntu

Pull non-default tag
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx:alpine

Pull an image based on digest
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  nginx@sha256:7281cf7c854b0dfc7c68a6a4de9a785a973a14f1481bc028e2022bcd6a8d9f64

Pull an image with the full docker registry URL
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  registry.hub.docker.com/library/hello-world

Pull an image with all tags
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  --all-tags nginx

Pull non-existent image
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull fakebadimage
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    Should contain  ${output}  image library/fakebadimage not found

Pull image from non-existent repo
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull fakebadrepo.com:9999/ubuntu
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  no such host

Pull image with a tag that doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox:faketag
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Tag faketag not found in repository library/busybox

Pull image that already has been pulled
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  alpine
    Wait Until Keyword Succeeds  5x  15 seconds  Pull image  alpine

Pull the same image concurrently
     ${pids}=  Create List

     # Create 5 processes to pull the same image at once
     :FOR  ${idx}  IN RANGE  0  5
     \   ${pid}=  Start Process  docker ${params} pull redis  shell=True
     \   Append To List  ${pids}  ${pid}

     # Wait for them to finish and check their output
     :FOR  ${pid}  IN  @{pids}
     \   ${res}=  Wait For Process  ${pid}
     \   Should Be Equal As Integers  ${res.rc}  0
     \   Should Contain  ${res.stdout}  Downloaded newer image for library/redis:latest

Pull two images that share layers concurrently
     ${pid1}=  Start Process  docker ${params} pull golang:1.7  shell=True
     ${pid2}=  Start Process  docker ${params} pull golang:1.6  shell=True

    # Wait for them to finish and check their output
    ${res1}=  Wait For Process  ${pid1}
    ${res2}=  Wait For Process  ${pid2}
    Should Be Equal As Integers  ${res1.rc}  0
    Should Be Equal As Integers  ${res2.rc}  0
    Should Contain  ${res1.stdout}  Downloaded newer image for library/golang:1.7
    Should Contain  ${res2.stdout}  Downloaded newer image for library/golang:1.6
