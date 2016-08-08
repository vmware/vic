*** Settings ***
Documentation  Test 9-1 - VICAdmin DisplayHTML
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Display HTML
    ${rc}  ${output}=  Run And Return Rc And Output  curl -k http://${vch-ip}:2378
    Should contain  ${output}  <title>VCH Admin</title>

Get Portlayer Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -k http://${vch-ip}:2378/logs/port-layer.log
	Should contain  ${output}  Launching portlayer server

Get VCH-Init Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -k http://${vch-ip}:2378/logs/init.log
	Should contain  ${output}  reaping child processes

Get Docker Personality Log
    ${rc}  ${output}=  Run And Return Rc And Output  curl -k http://${vch-ip}:2378/logs/docker-personality.log
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
	${rc}=  Run And Return Rc  curl -k https://${vch-ip}:2378/logs/container-logs.zip > /tmp/container-logs.zip
    Should Be Equal As Integers  ${rc}  0
	File Should Exist  /tmp/container-logs.zip
    File Should Not Be Empty  /tmp/container-logs.zip
    ${rc}  ${output}=  Run And Return Rc And Output  zipinfo /tmp/container-logs.zip
	Should Contain  ${output}  ${container}/vmware.log
	Remove File  /tmp/container-logs.zip 
