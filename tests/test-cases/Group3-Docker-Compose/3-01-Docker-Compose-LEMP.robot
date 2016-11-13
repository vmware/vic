*** Settings ***
Documentation  Test 3-01 - Docker Compose LEMP
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
# ${yml}  vic.mysql:\n${SPACE}container_name: vic.mysql\n${SPACE}image: mysql\n${SPACE}environment:\n${SPACE}- MYSQL_ROOT_PASSWORD=root\n${SPACE}- MYSQL_DATABASE=vmware1\n\nvic.php:\n${SPACE}container_name: vic.php\n${SPACE}image: php:7.0-fpm\n${SPACE}volumes:\n${SPACE}- src:/var/www\n\nvic.nginx:\n${SPACE}container_name: vic.nginx\n${SPACE}image: nginx\n${SPACE}links:\n${SPACE}- vic.mysql\n${SPACE}- vic.php\n${SPACE}volumes:\n${SPACE}- src:/var/www\n${SPACE}ports:\n${SPACE}- "80:80"

*** Test Cases ***
Compose LEMP Server
    # ${status}=  Get State Of Github Issue  2357
    # Run Keyword If  '${status}' == 'closed'  Fail  Test 3-1-Docker-Compose-LEMP.robot needs to be updated now that Issue #2357 has been resolved
    # Log  Issue \#2357 is blocking implementation  WARN
	
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0

    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300
    # must set CURL_CA_BUNDLE to work around Compose bug https://github.com/docker/compose/issues/3365
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}
	
    ${vch_ip}=  Get Environment Variable  VCH_IP  ${vch-ip}
    Log To Console  \nThe VCH IP is ${vch_ip}

    # Run  echo '${yml}' > lemp-compose.yml
    Run  cat %{GOPATH}/src/github.com/vmware/vic/demos/compose/webserving-app/docker-compose.yml | sed -e "s/192.168.60.130/${vch_ip}/g" > lemp-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose ${params} --file lemp-compose.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
	
    ${rc}  ${output}=  Run And Return Rc And Output  wget ${vch_ip}:8080
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
