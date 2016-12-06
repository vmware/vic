*** Settings ***
Documentation  Test 1-06 - Docker Run
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Make sure container starts
    :FOR  ${idx}  IN RANGE  0  30
    \   ${out}=  Run  docker %{VCH-PARAMS} ps
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  /bin/top
    \   Exit For Loop If  ${status}
    \   Sleep  1

*** Test Cases ***
Simple docker run
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Docker run with -t
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -t busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Simple docker run with app that doesn't exit
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq | xargs -n1 docker %{VCH-PARAMS} rm -f
    ${result}=  Start Process  docker %{VCH-PARAMS} run busybox /bin/top  shell=True  alias=top

    Make sure container starts
    ${containerID}=  Run  docker %{VCH-PARAMS} ps -q
    ${out}=  Run  docker %{VCH-PARAMS} logs ${containerID}
    Should Contain  ${out}  Mem:
    Should Contain  ${out}  CPU:
    Should Contain  ${out}  Load average:
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq | xargs -n1 docker %{VCH-PARAMS} rm -f

Docker run fake command
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run busybox fakeCommand
    Should Be True  ${rc} > 0
    Should Contain  ${output}  docker: Error response from daemon:
    Should Contain  ${output}  fakeCommand: no such executable in PATH.

Docker run fake image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run fakeImage /bin/bash
    Should Be True  ${rc} > 0
    Should Contain  ${output}  docker: Error parsing reference:
    Should Contain  ${output}  "fakeImage" is not a valid repository/tag.

Docker run named container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --name busy3 busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

Docker run linked containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull debian
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --link busy3:busy3 debian ping -c2 busy3
    Should Be Equal As Integers  ${rc}  0

Docker run -d unspecified host port
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p 6379 redis:alpine
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Docker run check exit codes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run busybox true
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run busybox false
    Should Be Equal As Integers  ${rc}  1

Docker run ps password check
    [Tags]  secret
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run busybox ps auxww
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ps auxww
    ${output}=  Split To Lines  ${output}
    :FOR  ${line}  IN  @{output}
    \   ${line}=  Strip String  ${line}
    \   ${command}=  Split String  ${line}  max_split=3
    \   ${len}=  Get Length  ${command}
    \   Continue For Loop If  ${len} <= 4
    \   Should Not Contain  @{command}[4]  %{TEST_USERNAME}
    \   Should Not Contain  @{command}[4]  %{TEST_PASSWORD}