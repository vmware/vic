*** Settings ***
Documentation  Test 11-03 - Upgrade-InsecureRegistry
Suite Setup  Install VIC with version to Test Server  7315  --insecure-registry 10.10.10.10:1234
Suite Teardown  Clean up VIC Appliance And Local Binary
Resource  ../../resources/Util.robot

*** Keywords ***
Get host and path from guestinfo
    ${host}=  Run  govc vm.info -e -json %{VCH-NAME} | jq -c '.VirtualMachines[0].Config.ExtraConfig[] | select(.Key | . and contains("guestinfo.vice./registry/insecure_registries|0/Host")) .Value'
    ${path}=  Run  govc vm.info -e -json %{VCH-NAME} | jq -c '.VirtualMachines[0].Config.ExtraConfig[] | select(.Key | . and contains("guestinfo.vice./registry/insecure_registries|0/Path")) .Value'
    [Return]  ${host}  ${path}

*** Test Cases ***
Upgrade VCH with InsecureRegistry
    ${oldHost}  ${oldPath}=  Get host and path from guestinfo
    Log  ${oldHost} ${oldPath}
    Should Be Equal As Strings  ${oldHost}  "<nil>"
    Should Be Equal As Strings  ${oldPath}  "10.10.10.10:1234"
    Upgrade
    Check Upgraded Version
    ${newHost}  ${newPath}=  Get host and path from guestinfo
    Log  ${newHost} ${newPath}
    Should Be Equal As Strings  ${newPath}  "<nil>"
    Should Be Equal As Strings  ${newHost}  "10.10.10.10:1234"
