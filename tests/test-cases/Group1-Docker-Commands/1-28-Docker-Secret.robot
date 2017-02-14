*** Settings ***
Documentation  Test 1-28 - Docker Secret
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${fake-secret}  test

*** Test Cases ***
Docker secret ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} secret ls
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support Docker Swarm

Docker secret create
    Run  echo '${fake-secret}' > secret.file
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} secret create mysecret ./secret.file
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support Docker Swarm

Docker secret inspect
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} secret inspect my_secret
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support Docker Swarm

Docker secret rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} secret rm my_secret
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  vSphere Integrated Containers does not yet support Docker Swarm
	
