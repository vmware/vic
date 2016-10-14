*** Settings ***
Documentation  Test 9-2 - VICAdmin CertAuth
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${true}
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Curl
    [Arguments]  ${path}
    ${output}=  Run  curl -sk --cert %{DOCKER_CERT_PATH}/cert.pem --key %{DOCKER_CERT_PATH}/key.pem ${vic-admin}${path}
    Should Not Be Equal As Strings  ''  ${output}
    [Return]  ${output}

*** Test Cases ***
Display HTML    
     ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  ${EMPTY}
     Should contain  ${output}  <title>VIC: ${vch-name}</title>

Get Portlayer Log
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/port-layer.log
    Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/init.log
    Should contain  ${output}  reaping child processes

Get Docker Personality Log
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/docker-personality.log
    Should contain  ${output}  docker personality

Get VICAdmin Log
    ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/vicadmin.log
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server

Fail to Get VICAdmin Log without cert
    ${output}=  Run  curl -sk ${vic-admin}/logs/vicadmin.log
    Log  ${output}
    Should Not contain  ${output}  Launching vicadmin pprof server
    
Fail to Display HTML without cert
    ${output}=  Run  curl -sk ${vic-admin}
    Log  ${output}
    Should Not contain  ${output}  <title>VCH ${vch-name}</title>

Fail to get Portlayer Log without cert
    ${output}=  Run  curl -sk ${vic-admin}/logs/port-layer.log
    Log  ${output}
    Should Not contain  ${output}  Launching portlayer server

Fail to get Docker Personality Log without cert
    ${output}=  Run  curl -sk ${vic-admin}/logs/docker-personality.log
    Log  ${output}
    Should Not contain  ${output}  docker personality

Fail to get VCH init logs without cert
    ${output}=  Run  curl -sk ${vic-admin}/logs/init.log
    Log  ${output}
    Should Not contain  ${output}  reaping child processes
