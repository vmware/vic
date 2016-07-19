*** Settings ***
Documentation  Test 1-8 - Docker Logs
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic docker logs
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    
Docker logs with timestamps
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${time}=  Run And Return Rc And Output  docker ${params} logs -t ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    
Docker logs with timestamps and since certain time
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${time}=  Run And Return Rc And Output  docker ${params} logs -t ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${time}=  Split To Lines  ${time}
    #${len}=  Get Length  ${time}
    #${time}=  Fetch From Left  @{time}[${len-5}]  ${SPACE}
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs --since=@{time} ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    
Docker logs last three
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${time}=  Run And Return Rc And Output  docker ${params} logs --tail="3" ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    
Docker logs live output
    ${status}=  Get State Of Github Issue  366
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-8-Docker-Logs.robot needs to be updated now that Issue #366 has been resolved
    Log  Issue \#366 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${containerID}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs -f ${containerID}
    #Should Be Equal As Integers  ${rc}  0
    
Docker logs non-existent container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such container: fakeContainer