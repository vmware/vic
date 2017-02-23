*** Settings ***
Documentation  Test 3-03 - Docker Compose Basic
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  version: "3"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: python:2.7\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5000:5000"\n${SPACE}${SPACE}depends_on:\n${SPACE}${SPACE}- redis\n${SPACE}redis:\n${SPACE}${SPACE}image: redis\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5001:5001"

*** Test Cases ***
Compose basic
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300
    # must set CURL_CA_BUNDLE to work around Compose bug https://github.com/docker/compose/issues/3365
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}

    Run  echo '${yml}' > basic-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create vic_default
    Log  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file basic-compose.yml create
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file basic-compose.yml start
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file basic-compose.yml logs
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file basic-compose.yml stop
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
