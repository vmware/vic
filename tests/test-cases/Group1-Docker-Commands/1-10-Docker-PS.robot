*** Settings ***
Documentation  Test 1-10 - Docker PS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Empty docker ps command
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  CONTAINER ID
    Should Contain  ${output}  IMAGE
    Should Contain  ${output}  COMMAND
    Should Contain  ${output}  CREATED
    Should Contain  ${output}  STATUS
    Should Contain  ${output}  PORTS
    Should Contain  ${output}  NAMES
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

Docker ps only running containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container2}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container3}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    Sleep  5 seconds
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2
    
Docker ps all containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    Should Contain  ${output}  dmesg
    Should Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  4
    
Docker ps last container
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -l
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2
    
Docker ps two containers
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -n=2
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  dmesg
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  3
    
Docker ps last container with size
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  SIZE
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2
    
Docker ps all containers with only IDs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  CONTAINER ID    
    Should Not Contain  ${output}  /bin/top
    Should Not Contain  ${output}  dmesg
    Should Not Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  3
    
Docker ps with filter
    ${status}=  Get State Of Github Issue  1676
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1676 has been resolved
    Log  Issue \#1676 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -f status=created
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2