*** Settings ***
Documentation  Test 3-1 - Force Install
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Force Re-install same VCH appliance
    ${output}=  Run  bin/vic-machine-linux create --name ${vch-name} --key ${vch-name}-key.pem --cert ${vch-name}-cert.pem --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-datastore=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --generate-cert=false --password=%{TEST_PASSWORD} --force=true --bridge-network=network --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  Installer completed successfully...