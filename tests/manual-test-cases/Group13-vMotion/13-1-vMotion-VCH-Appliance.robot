*** Settings ***
Documentation  Test 13-1 - vMotion VCH Appliance
Resource  ../../resources/Util.robot
Suite Setup  Create a VSAN Cluster
Suite Teardown  Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  *

*** Keywords ***
Gather Host IPs
    ${out}=  Run  govc ls host/cls
    ${out}=  Split To Lines  ${out}
    ${idx}=  Set Variable  1
    :FOR  ${line}  IN  @{out}
    \   Continue For Loop If  '${line}' == '/vcqaDC/host/cls/Resources'
    \   ${ip}=  Fetch From Right  ${line}  /
    \   Set Suite Variable  ${esx${idx}-ip}  ${ip}
    \   ${idx}=  Evaluate  ${idx}+1

Create a VSAN Cluster
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --vcvaBuild 3634791 --esxPxeDir 3620759 --esxBuild 3620759 --testbedName vcqa-vsan-simple-pxeBoot-vcva --runName vic-vmotion
    ${out}=  Split To Lines  ${out}
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  .vcva-3634791' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Suite Variable  ${vc-ip}  ${ip}
    \   Exit For Loop If  ${status}

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Log To Console  Create a distributed switch
    ${out}=  Run  govc dvs.create -dc=vcqaDC test-ds
    Should Contain  ${out}  OK

    Log To Console  Create three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds bridge
    Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    ${out}=  Run  govc dvs.add -dvs=test-ds -pnic=vmnic1 /vcqaDC/host/cls
    Should Contain  ${out}  OK

    Log To Console  Enable DRS and VSAN on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  EXTERNAL_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  vsanDatastore
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  10m

    Gather Host IPs

*** Test Cases ***
Step 1-5
    ${status}=  Get State Of Github Issue  701
    Run Keyword If  '${status}' == 'closed'  Fail  Test 13-1-vMotion-VCH-Appliance.robot needs to be updated now that Issue #701 has been resolved
    Log  Issue \#701 is blocking implementation  WARN
#    Install VIC Appliance To Test Server  ${false}  default
#    Run Regression Tests
#    ${host}=  Get VM Host Name  ${vch-name}
#    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -off ${vch-name}/${vch-name}
#    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
#    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vch-name}/${vch-name}
#    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vch-name}/${vch-name}
#    Set Suite Variable  ${vch-name}  "${vch-name} (1)"
#    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.power -on ${vch-name}
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container2}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${container3}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container3}
    Should Be Equal As Integers  ${rc}  0

    ${host}=  Get VM Host Name  ${vch-name}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vch-name}/${vch-name}
    Set Suite Variable  ${vch-name}  "${vch-name} (1)"

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container1}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container2}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} logs ${container3}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container3}
    Should Be Equal As Integers  ${rc}  0

    Cleanup VIC Appliance On Test Server