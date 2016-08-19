*** Settings ***
Documentation  Test 6-4 - Verify vic-machine create basic use cases
Resource  ../../resources/Util.robot
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create VCH - defaults
    Log To Console  \nRunning vic-machine create - defaults
    Set Test Environment Variables  ${true}  default
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - defaults with --no-tls
    Log To Console  \nRunning vic-machine create - defaults with --no-tls
    Set Test Environment Variables  ${true}  default
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Install VIC Appliance To Test Server
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - target URL
    Log To Console  \nRunning vic-machine create - target URL
    Set Test Environment Variables  ${true}  default
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --image-store=%{TEST_DATASTORE}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - full params
    Log To Console  \nRunning vic-machine create
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test Environment Variables  ${true}  default
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:default
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - custom image store directory
    Log To Console  \nRunning vic-machine create
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test Environment Variables  ${true}  default
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-store %{TEST_DATASTORE}/vic-machine-test-images --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=%{TEST_PASSWORD} --force=true --bridge-network=network --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}

    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...
    ${output}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls
    Should Contain  ${output}  vic-machine-test-images

    Run Regression Tests
    Cleanup VIC Appliance On Test Server
    ${output}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls
    Should Not Contain  ${output}  vic-machine-test-images
