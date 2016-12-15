*** Settings ***
Documentation  Test 5-14 - Remove Container OOB
Resource  ../../resources/Util.robot
#Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup

*** Test Cases ***
Docker run an image from a container that was removed OOB
    ${status}=  Get State Of Github Issue  2928
    Run Keyword If  '${status}' == 'closed'  Fail  Test 5-14-Remove-Container-OOB.robot needs to be updated now that Issue #2928 has been resolved
    Log  Issue \#2928 is blocking implementation  WARN
#    Create a Simple VC Cluster

#    Install VIC Appliance To Test Server

#    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
#    Should Be Equal As Integers  ${rc}  0
#    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -itd busybox /bin/top 
#    Should Be Equal As Integers  ${rc}  0
#    Destroy VM OOB  ${container}
#    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -itd busybox /bin/top 
#    Should Be Equal As Integers  ${rc}  0
