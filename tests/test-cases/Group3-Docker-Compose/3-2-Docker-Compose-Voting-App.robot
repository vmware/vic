*** Settings ***
Documentation  Test 3-2 - Docker Compose Voting App
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Compose Voting App
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0

    ${out}=  Run Process  docker-compose  ${params}  up  -d  shell=True  cwd=${CURDIR}/../../../demos/compose/voting-app   env:COMPOSE_HTTP_TIMEOUT=300
    Log  ${out.rc}
    Log  ${out.stdout}
    Log  ${out.stderr}
    Should Be Equal As Integers  ${out.rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} inspect -f {{.State.Running}} vote
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} inspect -f {{.State.Running}} result
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0
