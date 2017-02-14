*** Settings ***
Documentation  Test 1-33 - Docker Service
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker service create 
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service create test-service
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker service ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service ls
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker service ps
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service ps test-service
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker serivce rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service rm test-service
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker service scale
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service scale test-service=3
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker service update
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service update test-service
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker service logs
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} service logs test
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  only supported with experimental daemon

