*** Settings ***
Documentation  Test 10-1 - VCH Restart
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
    [Arguments]  ${name}  ${network}=default
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --name ${name} --net ${network} -itd busybox
    Should Be Equal As Integers  ${rc}  0
    ${id}=  Get Line  ${output}  -1
    ${ip}=  Get Container IP  ${id}  ${network}
    [Return]  ${id}  ${ip}


*** Test Cases ***
Created Network And Images Persists As Well As Containers Are Discovered With Correct IPs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create bar
    Should Be Equal As Integers  ${rc}  0
    Comment  Launch first container on bridge network
    ${id1}  ${ip1}=  Launch Container  vch-restart-test1  bridge
    ${id2}  ${ip2}=  Launch Container  vch-restart-test2  bridge

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    Log To Console  Rebooting VCH ...
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -reset=true ${vch-name}
    Should Be Equal As Integers  ${rc}  0
    Log To Console  Waiting for VCH to boot ...
    Wait Until Vm Powers On  ${vch-name}
    Log To Console  VCH Powered On
    Sleep  5
    Log To Console  Getting VCH IP ...
    ${new-vch-ip}=  Get VM IP  ${vch-name}
    Log To Console  New VCH IP is ${new-vch-ip}
    Replace String  ${params}  ${vch-ip}  ${new-vch-ip}

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  nginx
    Should Contain  ${output}  busybox

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bridge
    Should Be Equal As Integers  ${rc}  0
    ${ip}=  Get Container IP  ${id1}  bridge
    Should Be Equal  ${ip}  ${ip1}
    ${ip}=  Get Container IP  ${id2}  bridge
    Should Be Equal  ${ip}  ${ip2}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  "Id"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Exited (0)
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Exited (0)

    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver1 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  port 10000 is not available
