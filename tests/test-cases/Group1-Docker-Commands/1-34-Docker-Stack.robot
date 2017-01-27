*** Settings ***
Documentation  Test 1-34 - Docker Stack
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
#Docker stack deploy
    #${rc}  ${output}=  Run And Return Rc And Output  wget #https://raw.githubusercontent.com/vfarcic/docker-flow-proxy/master/docker-compose-stack.yml
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} stack deploy -c #./docker-compose-stack.yml proxy
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  does not yet support Docker Swarm

Docker stack ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} stack ls
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker stack ps
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} stack ps test-stack
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker stack rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} stack rm test-stack
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker stack services
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} stack services test-stack
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm
