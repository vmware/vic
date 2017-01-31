*** Settings ***
Documentation  Test 5-12 - Multiple VLAN
Resource  ../../resources/Util.robot
Suite Setup  Multiple VLAN Setup
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}
Test Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Multiple VLAN Setup
    ${esx1}  ${esx2}  ${esx3}  ${vc}  ${vc-ip}=  Create a Simple VC Cluster  multi-vlan-1  cls
    Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${vc}

*** Test Cases ***
Test1
    ${out}=  Run  govc dvs.portgroup.change -vlan 1 bridge
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.change -vlan 1 management
    Should Contain  ${out}  OK

    Install VIC Appliance To Test Server
    Run Regression Tests

Test2
    ${out}=  Run  govc dvs.portgroup.change -vlan 1 bridge
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.change -vlan 2 management
    Should Contain  ${out}  OK

    Install VIC Appliance To Test Server
    Run Regression Tests
