*** Settings ***
Documentation  Test 9-2 - VICAdmin CertAuth
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${true}
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Curl
    [Arguments]  ${path}
    ${output}=  Run  curl -sk --cert /drone/src/github.com/vmware/vic/${vch-name}/cert.pem --key /drone/src/github.com/vmware/vic/${vch-name}/key.pem ${vic-admin}${path}
    Should Not Be Equal As Strings  ''  ${output}
    [Return]  ${output}

*** Test Cases ***
# Display HTML
#     ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  ${EMPTY}
#     Should contain  ${output}  <title>VCH Admin</title>

# Get Portlayer Log
#     ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/port-layer.log
#     Should contain  ${output}  Launching portlayer server

# Get VCH-Init Log
#     ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/init.log
#     Should contain  ${output}  reaping child processes

# Get Docker Personality Log
#     ${output}=  Wait Until Keyword Succeeds  10x  10s  Curl  /logs/docker-personality.log
#     Should contain  ${output}  docker personality

Get Container Logs
    Run Keyword  Set Environment Variable  DOMAIN  ${EMPTY}
    ${rc}  ${output}=  Run And Return Rc and Output  docker ${params} pull busybox
    Log To Console  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container}=  Run And Return Rc and Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Run  curl -sk --cert /drone/src/github.com/vmware/vic/${vch-name}/cert.pem --key /drone/src/github.com/vmware/vic/${vch-name}/key.pem ${vic-admin}/container-logs.tar.gz | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug

Get VICAdmin Log
    ${output}=  Run  curl -sk --cert /drone/src/github.com/vmware/vic/${vch-name}/cert.pem --key /drone/src/github.com/vmware/vic/${vch-name}/key.pem ${vic-admin}/logs/vicadmin.log
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server
