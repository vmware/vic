*** Settings ***
Documentation  Test 6-5 - Verify vic-machine create validation function
Resource  ../../resources/Util.robot

*** Test Cases ***
Suggest resources - Invalid datacenter
    Log To Console  \nRunning vic-machine create - defaults
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
    Log To Console  \nRunning vic-machine create - defaults
    Set Test Environment Variables  ${true}  default  network  'VM Network'
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL}/MUCH/DATACENTER --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD}
    Should Contain  ${output}  Suggested datacenters:
    Should Contain  ${output}  vic-machine-linux failed
