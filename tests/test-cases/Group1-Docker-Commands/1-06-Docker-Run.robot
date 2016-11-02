*** Settings ***
Documentation  Test 1-06 - Docker Run
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker run
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox dmesg
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Docker run with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -i busybox dmesg
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Docker run with -it
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${result}=  Start Process  docker ${params} run -i busybox /bin/top < /tmp/fifo  shell=True  alias=top
    Sleep  5
    ${rc2}  ${output2}=  Run And Return Rc And Output  echo q > /tmp/fifo
    ${result2}=  Wait for process  top
    Log  ${result2.stdout}
    Log  ${result2.stderr}

Simple docker run with app that doesn't exit
    Log To Console  Not sure how to implement this just yet...
    #${result}=  Run Process  docker ${params} run busybox /bin/top  shell=True  timeout=5s  on_timeout=terminate
    #Log  ${result.stdout}
    #Should Be Equal As Integers  ${result.rc}  0

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

Docker run df command
    ${status}=  Get State Of Github Issue  1582
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-6-Docker-Run.robot needs to be updated now that Issue #1582 has been resolved
    Log  Issue \#1582 is blocking implementation  WARN
#   ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it busybox /bin/df
#   Should Be Equal As Integers  ${rc}  0
#   Should Contain  ${output}  Filesystem

Docker run -d unspecified host port
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -d -p 6379 redis:alpine
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Docker run check exit codes
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox true
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run busybox false
    Should Be Equal As Integers  ${rc}  1

Docker run date
    ${status}=  Get State Of Github Issue  1582
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-6-Docker-Run.robot needs to be updated now that Issue #1582 has been resolved
    Log  Issue \#1582 is blocking implementation  WARN
#   ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it busybox date
#   Should Be Equal As Integers  ${rc}  0
#   Should Contain  ${output}  UTC

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
