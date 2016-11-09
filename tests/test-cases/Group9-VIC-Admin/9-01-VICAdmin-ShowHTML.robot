*** Settings ***
Documentation  Test 9-01 - VICAdmin ShowHTML
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Test Cases ***
Get Login Page
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/authentication
    Should contain  ${output}  <title>VCH Admin</title>

While Logged Out Fail To Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}
    Should not contain  ${output}  <title>VIC: ${vch-name}</title>
    Should Contain  ${output}  <a href="/authentication">Temporary Redirect</a>.

While Logged Out Fail To Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/port-layer.log
    Should Not Contain  ${output}  Launching portlayer server
    Should Contain  ${output}  <a href="/authentication">Temporary Redirect</a>.

While Logged Out Fail To Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/init.log
    Should not contain  ${output}  reaping child processes
    Should Contain  ${output}  <a href="/authentication">Temporary Redirect</a>.

While Logged Out Fail To Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/docker-personality.log
    Should not contain  ${output}  docker personality
    Should Contain  ${output}  <a href="/authentication">Temporary Redirect</a>.

While Logged Out Fail To Get Container Logs
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
    Should not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  gzip: stdin: not in gzip format
    Log  ${output}
    Should not Contain  ${output}  ${container}/vmware.log
    Should not Contain  ${output}  ${container}/tether.debug

While Logged Out Fail To Get VICAdmin Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/vicadmin.log
    Log  ${output}
    Should not contain  ${output}  Launching vicadmin pprof server
    Should Contain  ${output}  <a href="/authentication">Temporary Redirect</a>.

Login
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/authentication -XPOST -F username=%{GOVC_USERNAME} -F password=%{GOVC_PASSWORD} -D /tmp/cookies-${vch-name}
    Should Be Equal As Integers  ${rc}  0

Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin} -b /tmp/cookies-${vch-name}
    Should contain  ${output}  <title>VIC: ${vch-name}</title>

Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/port-layer.log -b /tmp/cookies-${vch-name}
    Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/init.log -b /tmp/cookies-${vch-name}
    Should contain  ${output}  reaping child processes

Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/docker-personality.log -b /tmp/cookies-${vch-name}
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
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk ${vic-admin}/container-logs.tar.gz -b /tmp/cookies-${vch-name} | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug

Get VICAdmin Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk ${vic-admin}/logs/vicadmin.log -b /tmp/cookies-${vch-name}
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server
