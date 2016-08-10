*** Settings ***
Documentation  Test 1-11 - Docker RM
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic docker remove container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc ls vm
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${container}

Remove a stopped container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Container Stops  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc ls vm
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${container}
    ${status}=  Get State Of Github Issue  1313
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-11-Docker-RM.robot needs to be updated now that Issue #1313 has been resolved
    Log  Issue \#1313 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  ${container}

Remove a running container
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: You cannot remove a running container. Stop the container before attempting removal or use -f
    
Force remove a running container
    ${status}=  Get State Of Github Issue  1312
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-11-Docker-RM.robot needs to be updated now that Issue #1312 has been resolved
    Log  Issue \#1312 is blocking implementation  WARN
    #${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm -f ${container}
    #Should Be Equal As Integers  ${rc}  0
    
Remove a fake container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer

Remove a container deleted out of band
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name test busybox
    Should Be Equal As Integers  ${rc}  0
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.destroy "test*"
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm test
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: test
    
