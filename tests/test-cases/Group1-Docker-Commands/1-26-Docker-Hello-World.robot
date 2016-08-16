*** Settings ***
Documentation  Test 1-26 - Docker Hello World
Resource  ../../resources/Util.robot
#Suite Setup  Install VIC Appliance To Test Server
#Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Hello world
    ${status}=  Get State Of Github Issue  384
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-26-Docker-Hello-World.robot needs to be updated now that Issue #384 has been resolved
    Log  Issue \#384 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run hello-world
    #Log  ${output}
    #Should Be Equal As Integers  ${rc}  0