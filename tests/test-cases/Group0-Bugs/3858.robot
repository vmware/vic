*** Settings ***
Documentation  Test 3858
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Check kernel modules
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -t ubuntu find /lib/modules -name "*.ko" -printf "FOUND\n" -quit
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  FOUND