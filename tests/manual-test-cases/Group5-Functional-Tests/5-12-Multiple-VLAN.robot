*** Settings ***
Documentation  Test 5-12 - Multiple VLAN
Resource  ../../resources/Util.robot
Suite Setup  Create a Simple VC Cluster  multi-vlan-1  cls
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup
Test Teardown  Cleanup VIC Appliance On Test Server

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