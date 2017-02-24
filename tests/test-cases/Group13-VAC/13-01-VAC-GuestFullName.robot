*** Settings ***
Documentation  Test 13-01 - Guest Full Name
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Check VCH VM Guest Operating System
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc vm.info %{VCH-NAME}/%{VCH-NAME} | grep 'Guest name'
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.info %{VCH-NAME} | grep 'Guest name'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Photon - VCH

Create a test container and check Guest Operating System
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name test busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc vm.info %{VCH-NAME}/test-${id} | grep 'Guest name'
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.info test-${id} | grep 'Guest name'
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Photon - Container
