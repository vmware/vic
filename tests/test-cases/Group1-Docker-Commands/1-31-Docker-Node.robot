*** Settings ***
Documentation  Test 1-31 - Docker Node
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker node demote
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node demote self
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker node ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node ls
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker node promote
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node promote self
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker node rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node rm self
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker node update
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node update self
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support Docker Swarm

Docker node ps
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node ps
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No such node

Docker node inspect
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} node inspect self
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  No such node

