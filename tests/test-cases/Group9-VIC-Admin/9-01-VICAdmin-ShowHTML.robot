*** Settings ***
Documentation  Test 9-01 - VICAdmin ShowHTML
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Test Cases ***
Get Login Page
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/authentication
    Should contain  ${output}  <title>VCH Admin</title>

While Logged Out Fail To Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}
    Should not contain  ${output}  <title>VIC: %{VCH-NAME}</title>
    Should Contain  ${output}  <a href="/authentication">See Other</a>.

While Logged Out Fail To Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/port-layer.log
    Should Not Contain  ${output}  Launching portlayer server
    Should Contain  ${output}  <a href="/authentication">See Other</a>.

While Logged Out Fail To Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/init.log
    Should not contain  ${output}  reaping child processes
    Should Contain  ${output}  <a href="/authentication">See Other</a>.

While Logged Out Fail To Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/docker-personality.log
    Should not contain  ${output}  docker personality
    Should Contain  ${output}  <a href="/authentication">See Other</a>.

While Logged Out Fail To Get Container Logs
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container}=  Run And Return Rc and Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk %{VIC-ADMIN}/container-logs.tar.gz | tar tvzf -
    Should not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  gzip: stdin: not in gzip format
    Log  ${output}
    Should not Contain  ${output}  ${container}/vmware.log
    Should not Contain  ${output}  ${container}/tether.debug

While Logged Out Fail To Get VICAdmin Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/vicadmin.log
    Log  ${output}
    Should not contain  ${output}  Launching vicadmin pprof server
    Should Contain  ${output}  <a href="/authentication">See Other</a>.

Login
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/authentication -XPOST -F username=%{GOVC_USERNAME} -F password=%{GOVC_PASSWORD} -D /tmp/cookies-%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN} -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  <title>VIC: %{VCH-NAME}</title>

Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/port-layer.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/init.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  reaping child processes

Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/docker-personality.log -b /tmp/cookies-%{VCH-NAME}
    Should contain  ${output}  docker personality

Get Container Logs
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${container}=  Run And Return Rc and Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${container}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk %{VIC-ADMIN}/container-logs.tar.gz -b /tmp/cookies-%{VCH-NAME} | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug

Get VICAdmin Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk %{VIC-ADMIN}/logs/vicadmin.log -b /tmp/cookies-%{VCH-NAME}
    Log  ${output}
    Should contain  ${output}  Launching vicadmin pprof server
