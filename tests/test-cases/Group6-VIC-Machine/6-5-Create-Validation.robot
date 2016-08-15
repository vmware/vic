*** Settings ***
Documentation  Test 6-5 - Verify vic-machine create validation function
Resource  ../../resources/Util.robot

*** Test Cases ***
Suggest resources - Invalid datacenter
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables  ${true}  default  network  'VM Network'
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL}/WOW --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD}
    Should Contain  ${output}  Suggested datacenters:
    Should Contain  ${output}  vic-machine-linux failed


Suggest resources - Invalid target path
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables  ${true}  default  network  'VM Network'
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL}/MUCH/DATACENTER --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD}
    Should Contain  ${output}  Suggested datacenters:
    Should Contain  ${output}  vic-machine-linux failed

Default image datastore
    # This test case is dependent on the ESX environment having only one datastore
    Log To Console  \nRunning vic-machine create - default image datastore
    Set Test Environment Variables  ${true}  default  network  'VM Network'
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Log To Console  ${output}
    Should Contain  ${output}  Using default datastore
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...
    Sleep  10 seconds
    ${status}=  Get State Of Github Issue  1109
    Run Keyword If  '${status}' == 'closed'  Fail  6-5-Create-Validation.robot needs to be updated now that Issue #1109 has been resolved
    Run Regression Tests
    Cleanup VIC Appliance On Test Server
