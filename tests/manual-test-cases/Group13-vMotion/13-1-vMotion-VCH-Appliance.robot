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
#    Install VIC Appliance To Test Server
#    Run Regression Tests
#    ${host}=  Get VM Host Name  %{VCH-NAME}/%{VCH-NAME}
#    Power Off VM OOB  %{VCH-NAME}
#    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
#    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
#    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
#    Set Environment Variable  VCH-NAME  "%{VCH-NAME} (1)"
#    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on %{VCH-NAME}
#    Should Be Equal As Integers  ${rc}  0

#    Log To Console  Waiting for VM to power on ...
#    :FOR  ${idx}  IN RANGE  0  30
#    \   ${ret}=  Run  govc vm.info %{VCH-NAME}
#    \   ${status}=  Run Keyword And Return Status  Should Contain  ${ret}  poweredOn
#    \   Exit For Loop If  ${status}
#    \   Sleep  1

#    Log To Console  Getting VCH IP ...
#    ${new-vch-ip}=  Get VM IP  %{VCH-NAME}
#    Log To Console  New VCH IP is ${new-vch-ip}
#    Replace String  %{VCH-PARAMS}  %{VCH-IP}  ${new-vch-ip}

#    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

#    Run Regression Tests
    #TODO
    #This does not work currently, as the VM has been migrated out of the vApp
    #Cleanup VIC Appliance On Test Server

Step 6-9
    Install VIC Appliance To Test Server
    Run Regression Tests
    ${host}=  Get VM Host Name  %{VCH-NAME}/%{VCH-NAME}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
    Set Environment Variable  VCH-NAME  "%{VCH-NAME} (1)"
    Run Regression Tests
    #Cleanup VIC Appliance On Test Server

Step 10-13
    Install VIC Appliance To Test Server
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

    ${host}=  Get VM Host Name  %{VCH-NAME}/%{VCH-NAME}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources %{VCH-NAME}/%{VCH-NAME}
    Set Environment Variable  VCH-NAME  "%{VCH-NAME} (1)"

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

    #Cleanup VIC Appliance On Test Server
