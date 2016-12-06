*** Settings ***
Documentation  Test 7-01 - Regression
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags  regression

*** Test Cases ***
Regression test
    Run Regression Tests

    ${out}=  Run  govc vm.info %{VCH-NAME}
    ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc vm.info %{VCH-NAME}/%{VCH-NAME}
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Set Test Variable  ${out}  ${ret}
    ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc vm.info %{VCH-NAME}
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Test Variable  ${out}  ${ret}
    Should Contain  ${out}  Other 3.x or later Linux (64-bit)