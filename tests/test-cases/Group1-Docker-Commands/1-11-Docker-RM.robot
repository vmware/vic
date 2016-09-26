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
    ${rc}  ${output}=  Run And Return Rc And Output  govc datastore.ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${container}

Remove a running container
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: You cannot remove a running container. Stop the container before attempting removal or use -f
    
Force remove a running container
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm -f ${container}
    Should Be Equal As Integers  ${rc}  0
    
Remove a fake container
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: fakeContainer

Remove a container deleted out of band
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name testRMOOB busybox
    Should Be Equal As Integers  ${rc}  0
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "testRMOOB*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy ${vch-name}/"testRMOOB*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm testRMOOB
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: No such container: testRMOOB

Remove a container created with unknown executable
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox xxxx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    
