*** Settings ***
Documentation  Test 3-1 - Docker Compose LEMP
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server  ${false}
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  vic.mysql:\n${SPACE}container_name: vic.mysql\n${SPACE}image: mysql\n${SPACE}environment:\n${SPACE}- MYSQL_ROOT_PASSWORD=root\n${SPACE}- MYSQL_DATABASE=vmware1\n\nvic.php:\n${SPACE}container_name: vic.php\n${SPACE}image: php:7.0-fpm\n${SPACE}volumes:\n${SPACE}- src:/var/www\n\nvic.nginx:\n${SPACE}container_name: vic.nginx\n${SPACE}image: nginx\n${SPACE}links:\n${SPACE}- vic.mysql\n${SPACE}- vic.php\n${SPACE}volumes:\n${SPACE}- src:/var/www\n${SPACE}ports:\n${SPACE}- "80:80"

*** Test Cases ***
Compose LEMP Server
    ${status}=  Get State Of Github Issue  2357
    Run Keyword If  '${status}' == 'closed'  Fail  Test 3-1-Docker-Compose-LEMP.robot needs to be updated now that Issue #2357 has been resolved
    Log  Issue \#2357 is blocking implementation  WARN
    #Run  echo '${yml}' > lemp-compose.yml
    #${rc}  ${output}=  Run And Return Rc And Output  DOCKER_HOST=${vch-ip}:2375 docker-compose --file lemp-compose.yml up
    #Log  ${output}
    #Should Be Equal As Integers  ${rc}  0