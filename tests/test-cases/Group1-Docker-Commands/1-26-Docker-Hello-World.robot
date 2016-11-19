*** Settings ***
Documentation  Test 1-26 - Docker Hello World
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Hello world
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/

Hello world with -t
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -t hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/