*** Settings ***
Documentation  Test 10-01 - VCH Restart
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Get Container IP
    [Arguments]  ${id}  ${network}=default
    ${rc}  ${ip}=  Run And Return Rc And Output  docker ${params} network inspect ${network} | jq '.[0].Containers."${id}".IPv4Address'
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${ip}

Launch Container
    [Arguments]  ${name}  ${network}=default  ${command}=sh
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --name ${name} --net ${network} -itd busybox ${command}
    Should Be Equal As Integers  ${rc}  0
    ${id}=  Get Line  ${output}  -1
    [Return]  ${id}

Reboot VCH
    Log To Console  Rebooting VCH ...
    ${vm}=  Get VM Name  ${vch-name}
    Power Off VM OOB  ${vm}
    Power On VM OOB  ${vm}
    Log To Console  VCH Powered On

*** Test Cases ***
Created Network And Images Persists As Well As Containers Are Discovered With Correct IPs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create foo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create bar
    Should Be Equal As Integers  ${rc}  0

    ${bridge-exited}=  Launch Container  vch-restart-bridge-exited  bridge  ls
    ${bridge-running}=  Launch Container  vch-restart-bridge-running  bridge
    ${bridge-running-ip}=  Get Container IP  ${bridge-running}  bridge
    ${bar-exited}=  Launch Container  vch-restart-bar-exited  bar  ls
    ${bar-running}=  Launch Container  vch-restart-bar-running  bar
    ${bar-running-ip}=  Get Container IP  ${bar-running}  bar

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    Reboot VCH

    Log To Console  Getting VCH IP ...
    ${new-vch-ip}=  Get VM IP  ${vch-name}
    Log To Console  New VCH IP is ${new-vch-ip}
    Replace String  ${params}  ${vch-ip}  ${new-vch-ip}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  ${params}

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  nginx
    Should Contain  ${output}  busybox

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  foo
    Should Contain  ${output}  bar
    Should Contain  ${output}  bridge
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bridge
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bar
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect foo
    Should Be Equal As Integers  ${rc}  0

    ${ip}=  Get Container IP  ${bridge-running}  bridge
    Should Be Equal  ${ip}  ${bridge-running-ip}
    ${ip}=  Get Container IP  ${bar-running}  bar
    Should Be Equal  ${ip}  ${bar-running-ip}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${bridge-running} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"running\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${bar-running} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"running\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${bridge-exited} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"exited\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect ${bar-exited} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"exited\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${bar-exited}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${bridge-exited}
    Should Be Equal As Integers  ${rc}  0

    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver1 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  port 10000 is not available

    # docker pull should work
    # if this fails, very likely the default gateway on the VCH is not set
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull alpine
    Should Be Equal As Integers  ${rc}  0
