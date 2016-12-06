*** Settings ***
Documentation  Test 3-01 - Docker Compose LEMP
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Compose LEMP Server
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0

    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300
    # must set CURL_CA_BUNDLE to work around Compose bug https://github.com/docker/compose/issues/3365
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}
	
    ${vch_ip}=  Get Environment Variable  VCH_IP  %{VCH-IP}
    Log To Console  \nThe VCH IP is %{VCH-IP}

    Run  cat %{GOPATH}/src/github.com/vmware/vic/demos/compose/webserving-app/docker-compose.yml | sed -e "s/192.168.60.130/${vch_ip}/g" > lemp-compose.yml
    Run  cat lemp-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file lemp-compose.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
