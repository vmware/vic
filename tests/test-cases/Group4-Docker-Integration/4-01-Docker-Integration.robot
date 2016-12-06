*** Settings ***
Documentation  Test 4-01 - Docker Integration
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases *** 
Docker Integration Tests
    [Tags]  docker
    Log To Console  \nStarting Docker integration tests... NOT EXECUTING NOW.
    # Currently blocked on issue https://github.com/vmware/vic/issues/2549
    #Set Environment Variable  GOPATH  /go:/go/src/github.com/docker/docker/vendor
    #${ip}=  Remove String  %{VCH-PARAMS}  -H
    #${ip}=  Strip String  ${ip}
    #Run  go get github.com/docker/docker
    #${status}  ${out}=  Run Keyword And Ignore Error  Run Process  DOCKER_HOST\=tcp://${ip} go test  shell=True  cwd=/go/src/github.com/docker/docker/integration-cli  timeout=10min  on_timeout=kill
    #Log  ${out.stdout}
    #Log  ${out.stderr}
    #Should Contain  ${out.stdout}  DockerSuite.Test