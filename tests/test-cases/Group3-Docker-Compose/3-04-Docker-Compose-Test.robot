*** Settings ***
Documentation  Test 3-04 - Docker Compose Test
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${true}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: python:2.7\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5000:5000"\n${SPACE}${SPACE}depends_on:\n${SPACE}${SPACE}- redis\n${SPACE}redis:\n${SPACE}${SPACE}image: redis\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5001:5001"
${link-yml}  version: "2"\nservices:\n${SPACE}redis1:\n${SPACE}${SPACE}image: redis:alpine\n${SPACE}${SPACE}container_name: redis1\n${SPACE}${SPACE}ports: ["6379"]\n${SPACE}web1:\n${SPACE}${SPACE}image: busybox\n${SPACE}${SPACE}container_name: a.b.c\n${SPACE}${SPACE}links:\n${SPACE}${SPACE}- redis1:aaa\n${SPACE}${SPACE}command: ["ping", "aaa"]
${hello-yml}  version: "2"\nservices:\n${SPACE}top:\n${SPACE}${SPACE}image: busybox\n${SPACE}${SPACE}container_name: top\n${SPACE}${SPACE}command: ["echo", "hello, world"]

*** Test Cases ***
Compose up in foreground (attach path)
    ${docker_tls_verify}=  Get Environment Variable  DOCKER_TLS_VERIFY
    Log  DOCKER_TLS_VERIFY = ${docker_tls_verify}
    
    Run  echo '${hello-yml}' > hello-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f hello-compose.yml pull
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}

    # Bring up the compose app and wait till they're up and running
    Start Process  docker-compose %{COMPOSE-PARAMS} -f hello-compose.yml up  shell=True  alias=hello
    Sleep  10
    ${ret}=  Wait For Process  hello
    Log To Console  ${ret.stdout}
    Log  ${ret.stdout}
    Should Contain  ${ret.stdout}  hello, world
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --f hello-compose.yml logs
    Log  ${output}

    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f hello-compose.yml down
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
