*** Settings ***
Documentation  Test 1-06 - Docker Run
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
Simple docker run
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox dmesg
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Simple docker run with app that doesn't exit
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq | xargs -n1 docker ${params} rm -f
    ${result}=  Start Process  docker ${params} run busybox /bin/top  shell=True  alias=top

    Make sure container starts
    ${containerID}=  Run  docker ${params} ps -q
    ${out}=  Run  docker ${params} logs ${containerID}
    Should Contain  ${out}  Mem:
    Should Contain  ${out}  CPU:
    Should Contain  ${out}  Load average:
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq | xargs -n1 docker ${params} rm -f

Docker run with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -i busybox dmesg
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Docker run fake command
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox fakeCommand
    Should Be True  ${rc} > 0
    Should Contain  ${output}  docker: Error response from daemon:
    Should Contain  ${output}  fakeCommand: no such executable in PATH.

Docker run fake image
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run fakeImage /bin/bash
    Should Be True  ${rc} > 0
    Should Contain  ${output}  docker: Error parsing reference:
    Should Contain  ${output}  "fakeImage" is not a valid repository/tag.

Docker run named container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -d --name busy3 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

Docker run linked containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull debian
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --link busy3:busy3 debian ping -c2 busy3
    Should Be Equal As Integers  ${rc}  0

Docker run -d unspecified host port
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -d -p 6379 redis:alpine
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Docker run check exit codes
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox true
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox false
    Should Be Equal As Integers  ${rc}  1

Docker run ps password check
    [Tags]  secret
    ${status}=  Get State Of Github Issue  2894
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-6-Docker-Run.robot needs to be updated now that Issue #2894 has been resolved
    Log  Issue \#2894 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox ps auxww
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ps auxww
    #Should Not Contain  ${output}  %{TEST_USERNAME}
    #Should Not Contain  ${output}  %{TEST_PASSWORD}