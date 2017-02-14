*** Settings ***
Documentation  Test 1-35 - Docker Swarm
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker swarm init
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm init
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker swarm join
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm join 127.0.0.1:2375
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker swarm join-token
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm join-token worker
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm join-token manager
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker swarm leave
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm leave
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker swarm unlock-key
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm unlock-key
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker swarm update
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} swarm update --autolock
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

