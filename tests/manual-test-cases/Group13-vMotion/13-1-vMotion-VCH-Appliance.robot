*** Settings ***
Documentation  Test 13-1 - vMotion VCH Appliance
Resource  ../../resources/Util.robot
Suite Setup  Create a VSAN Cluster
Suite Teardown  Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  *

*** Test Cases ***
Step 1-5
    ${status}=  Get State Of Github Issue  701
    Run Keyword If  '${status}' == 'closed'  Fail  Test 13-1-vMotion-VCH-Appliance.robot needs to be updated now that Issue #701 has been resolved
    Log  Issue \#701 is blocking implementation  WARN
#    Install VIC Appliance To Test Server  ${false}  default
#    Run Regression Tests
#    ${host}=  Get VM Host Name  ${vch-name}
#    Power Off VM OOB  ${vch-name}
#    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
#    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vch-name}/${vch-name}
#    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vch-name}/${vch-name}
#    Set Suite Variable  ${vch-name}  "${vch-name} (1)"
#    Power On VM OOB  ${vch-name}
#    Run Regression Tests
#    Cleanup VIC Appliance On Test Server

Step 6-9
    Install VIC Appliance To Test Server  ${false}  default
    Run Regression Tests
    ${host}=  Get VM Host Name  ${vch-name}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Set Suite Variable  ${vch-name}  "${vch-name} (1)"
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Step 10-13
    Install VIC Appliance To Test Server  ${false}  default
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container3}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container3}
    Should Be Equal As Integers  ${rc}  0

    ${host}=  Get VM Host Name  ${vch-name}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Set Suite Variable  ${vch-name}  "${vch-name} (1)"

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container1}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs ${container3}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${container3}
    Should Be Equal As Integers  ${rc}  0

    Cleanup VIC Appliance On Test Server
