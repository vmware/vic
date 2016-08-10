*** Settings ***
Documentation  Test 3-3 - Docker Compose Basic
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server  ${false}
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: python:2.7\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5000:5000"\n${SPACE}${SPACE}volumes:\n${SPACE}${SPACE}- ".:/code"\n${SPACE}${SPACE}depends_on:\n${SPACE}${SPACE}- redis\n${SPACE}redis:\n${SPACE}${SPACE}image: redis\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5001:5001"

*** Test Cases ***
Compose basic
    ${status}=  Get State Of Github Issue  1734
    Run Keyword If  '${status}' == 'closed'  Fail  Test 3-3-Docker-Compose-Basic.robot needs to be updated now that Issue #1734 has been resolved
    Log  Issue \#1734 is blocking implementation  WARN
    #Run  echo '${yml}' > basic-compose.yml
    #${rc}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker network create vic_default
    #Log  ${output}
    #${rc}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker-compose --file basic-compose.yml create
    #Log  ${output}
    #${rc1}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker-compose --file basic-compose.yml start
    #Log  ${output}
    #${rc2}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker-compose --file basic-compose.yml logs
    #Log  ${output}
    #${rc2}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker-compose --file basic-compose.yml stop
    #Log  ${output}
    #Should Be Equal As Integers  ${rc}  0
    #Should Be Equal As Integers  ${rc1}  0
    #Should Be Equal As Integers  ${rc2}  0