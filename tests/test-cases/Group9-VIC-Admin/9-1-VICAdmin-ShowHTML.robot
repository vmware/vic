*** Settings ***
Documentation  Test 9-1 - VICAdmin ShowHTML
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Test Cases ***
Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}
    Should contain  ${output}  <title>VCH Admin</title>

Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/port-layer.log
    Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/init.log
    Should contain  ${output}  reaping child processes

Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/docker-personality.log
    Should contain  ${output}  docker personality

Get Container Logs
    ${rc}  ${output}=  Run And Return Rc and Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container}=  Run And Return Rc and Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk ${vic-admin}/container-logs.tar.gz | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug

Get VICAdmin Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/vicadmin.log
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server
