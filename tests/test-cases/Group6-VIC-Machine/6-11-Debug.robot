*** Settings ***
Documentation  Test 6-11 - Verify enable of ssh in the appliance
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Enable SSH and verify
    # generate a key to use for the Test
    ${rc}=  Run And Return Rc  ssh-keygen -t rsa -N "" -f %{VCH-NAME}.key
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  chmod 600 %{VCH-NAME}.key
    Should Be Equal As Integers  ${rc}  0

    ${rc}=  Run And Return Rc  bin/vic-machine-linux debug --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name %{VCH-NAME} --enable-ssh --authorized-key=%{VCH-NAME}.key.pub
    Should Be Equal As Integers  ${rc}  0

    # check the ssh
    ${rc}=  Run And Return Rc  ssh -vv -o StrictHostKeyChecking=no -i %{VCH-NAME}.key root@%{VCH-IP} /bin/true
    Should Be Equal As Integers  ${rc}  0

    # delete the keys
    Remove Files  %{VCH-NAME}.key  %{VCH-NAME}.key.pub
