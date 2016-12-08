*** Settings ***
Documentation  Test 6-02 - Verify default parameters
Resource  ../../resources/Util.robot

*** Test Cases ***
Delete with defaults
    Set Test Environment Variables

    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Should Contain  ${ret}  vic-machine-linux delete failed:  resource pool
    Should Contain  ${ret}  /Resources/virtual-container-host' not found
