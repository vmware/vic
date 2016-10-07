*** Settings ***
Documentation  Test 6-11 - Verify enable of ssh in the appliance
Resource  ../../resources/Util.robot
Test Setup  Install VIC Appliance To Test Server
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Enable SSH and verify
    # generate a key to use for the Test
    ${rc}=  Run And Return Rc  ssh-keygen -t rsa -N "" -f ${vch-name}.key
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  chmod 600 ${vch-name}.key
    Should Be Equal As Integers  ${rc}  0

    ${rc}=  Run And Return Rc  bin/vic-machine-linux debug --target %{TEST_URL} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name} --enable-ssh --authorized-key=${vch-name}.key.pub
    Should Be Equal As Integers  ${rc}  0

    # check the ssh
    ${rc}=  Run And Return Rc  ssh -vv -o StrictHostKeyChecking=no -i ${vch-name}.key root@${vch-ip} /bin/true
    Should Be Equal As Integers  ${rc}  0

    # delete the keys
    Remove Files  ${vch-name}.key  ${vch-name}.key.pub 
