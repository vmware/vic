*** Settings ***
Documentation  Test 6-10 - Verify ls list all VCHs
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
List all VCHs
    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Should Contain  ${ret}  ID
    Should Contain  ${ret}  PATH
    Should Contain  ${ret}  NAME
    Should Not Contain  ${ret}  Error
    @{ret}=  Split To Lines  ${ret}
    ${tLen}=  Get Length  ${ret}
    Should Be True  ${tLen}>3

    # Get VCH ID, PATH and NAME
    @{vch}=  Split String  @{ret}[-1]
    ${vch-id}=  Strip String  @{vch}[0]
    ${vch-path}=  Strip String  @{vch}[1]
    ${vch-name}=  Strip String  @{vch}[2]

    # Run vic-machine inspect
    ${ret}=  Run  bin/vic-machine-linux inspect --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --id ${vch-id}
    Should Contain  ${ret}  Completed successfully
    ${ret}=  Run  bin/vic-machine-linux inspect --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource ${vch-path} --name ${vch-name}
    Should Contain  ${ret}  Completed successfully

List with compute-resource
    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource %{TEST_RESOURCE}
    Should Contain  ${ret}  ID
    Should Contain  ${ret}  PATH
    Should Contain  ${ret}  NAME
    Should Not Contain  ${ret}  Error
    @{ret}=  Split To Lines  ${ret}
    ${tLen}=  Get Length  ${ret}
    Should Be True  ${tLen}>3

    # Get VCH ID, PATH and NAME
    @{vch}=  Split String  @{ret}[-1]
    ${vch-id}=  Strip String  @{vch}[0]
    ${vch-path}=  Strip String  @{vch}[1]
    ${vch-name}=  Strip String  @{vch}[2]

    # Run vic-machine inspect
    ${ret}=  Run  bin/vic-machine-linux inspect --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --id ${vch-id}
    Should Contain  ${ret}  Completed successfully
    ${ret}=  Run  bin/vic-machine-linux inspect --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource ${vch-path} --name ${vch-name}
    Should Contain  ${ret}  Completed successfully

List with trailing slash
    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource %{TEST_RESOURCE}
    Should Contain  ${ret}  ID
    Should Contain  ${ret}  PATH
    Should Contain  ${ret}  NAME
    Should Not Contain  ${ret}  Error
    
List suggest compute resource
    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource fakeComputeResource
    Should Contain  ${ret}  Suggested values for --compute-resource:

List suggest valid datacenter
    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL}/fakeDatacenter --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource %{TEST_RESOURCE}
    Should Contain  ${ret}  Suggesting valid values for datacenter in --target

List with valid datacenter
    ${orig}=  Get Environment Variable  TEST_DATACENTER
    ${dc}=  Run Keyword If  '%{TEST_DATACENTER}' == '${SPACE}'  Get Datacenter Name
    Run Keyword If  '%{TEST_DATACENTER}' == '${SPACE}'  Set Environment Variable  TEST_DATACENTER  ${dc}

    ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL}/%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Set Environment Variable  TEST_DATACENTER  ${orig}
    Should Contain  ${ret}  ID
    Should Contain  ${ret}  PATH
    Should Contain  ${ret}  NAME
    Should Not Contain  ${ret}  Error