*** Settings ***
Documentation  Test 10-1 - VCH Restart
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  ${true}
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
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -off=true ${vch-name}
    Should Be Equal As Integers  ${rc}  0
    Log To Console  Waiting for VCH to power off ...
    Wait Until VM Powers Off  ${vch-name}
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on=true ${vch-name}
    Should Be Equal As Integers  ${rc}  0
    Log To Console  Waiting for VCH to power on ...
    Wait Until Vm Powers On  ${vch-name}
    Log To Console  VCH Powered On

*** Test Cases ***
Created Network Persists And Containers Are Discovered With Correct IPs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create bar
    Should Be Equal As Integers  ${rc}  0
    ${bridge-exited}=  Launch Container  vch-restart-bridge-exited  bridge  ls
    ${bridge-running}=  Launch Container  vch-restart-bridge-running  bridge
    ${bridge-running-ip}=  Get Container IP  ${bridge-running}  bridge
    ${bar-exited}=  Launch Container  vch-restart-bar-exited  bar  ls
    ${bar-running}=  Launch Container  vch-restart-bar-running  bar
    ${bar-running-ip}=  Get Container IP  ${bar-running}  bar
    Reboot VCH
    Sleep  10
    Log To Console  Getting VCH IP ...
    ${new-vch-ip}=  Get VM IP  ${vch-name}
    Log To Console  New VCH IP is ${new-vch-ip}
    Replace String  ${params}  ${vch-ip}  ${new-vch-ip}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  bar
    Should Contain  ${output}  bridge
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bridge
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bar
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
    ${status}=  Get State Of Github Issue  2448
    Run Keyword If  '${status}' == 'closed'  Fail  Test 10-1-VCH-Restart.robot needs to be updated now that Issue #2448 has been resolved
    Log  Issue \#2448 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect foo
    #Should Be Equal As Integers  ${rc}  0
