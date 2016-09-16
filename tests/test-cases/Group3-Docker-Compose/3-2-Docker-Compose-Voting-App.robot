*** Settings ***
Documentation  Test 3-2 - Docker Compose Voting App
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server  ${false}
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Compose Voting App
    ${status}=  Get State Of Github Issue  2358
    Run Keyword If  '${status}' == 'closed'  Fail  Test 3-2-Docker-Compose-Voting-App.robot needs to be updated now that Issue #2358 has been resolved
    Log  Issue \#2358 is blocking implementation  WARN
    #Run  git clone https://github.com/docker/example-voting-app test
    #${rc}  ${output}=  Run And Return Rc And Output  cd test; DOCKER_HOST=${vch-ip}:2375 docker-compose up
    #Log  ${output}
    #Should Be Equal As Integers  ${rc}  0