*** Settings ***
Documentation  Test 17-1 - TTY Tests
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Make sure container starts
    :FOR  ${idx}  IN RANGE  0  30
    \   ${out}=  Run  docker ${params} ps
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  /bin/top
    \   Exit For Loop If  ${status}
    \   Sleep  1

*** Test Cases ***
Docker run -it date
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} run -it busybox date
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  UTC

Docker run -it df 
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} run -it busybox df
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  Filesystem

Docker run -it command that doesn't stop
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} ps -aq | xargs -n1 docker ${params} rm -f
    ${result}=  Start Process  docker ${params} run -itd busybox /bin/top  shell=True  alias=top

    Make sure container starts
    ${containerID}=  Run  docker ${params} ps -q
    ${out}=  Run  docker ${params} logs ${containerID}
    Should Contain  ${out}  Mem:
    Should Contain  ${out}  CPU:
    Should Contain  ${out}  Load average: