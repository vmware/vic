*** Settings ***
Documentation  Test 6-06 - Verify vic-machine create image store, volume store and container store functions
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Image Store Delete - Image store not found
    Log To Console  \nRunning vic-machine create - custom image store path
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE}/images --password=%{TEST_PASSWORD} --force --kv
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Log To Console  \nDeleting image stores...
    ${out}=  Run  govc datastore.rm -ds=%{TEST_DATASTORE} images

    Log To Console  \nRunning vic-machine delete
    Cleanup VIC Appliance On Test Server
