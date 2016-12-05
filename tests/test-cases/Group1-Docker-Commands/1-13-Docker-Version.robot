*** Settings ***
Documentation  Test 1-13 - Docker Version
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple Docker Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Client:
    Should Contain  ${output}  Server:
    Should Contain  ${output}  Version:
    Should Contain  ${output}  Built:

Docker Version Format Client Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Client.Version}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

Docker Version Format Client API Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Client.APIVersion}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

Docker Version Format Client Go Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Client.GoVersion}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

Docker Version Format Server Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Server.Version}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

Docker Version Format Server API Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Server.APIVersion}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

Docker Version Format Server Go Version
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.Server.GoVersion}}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}
    
Docker Version Format Bad Field
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} version --format '{{.fakeItem}}'
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  fakeItem is not a field of struct type types.VersionResponse