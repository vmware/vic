*** Settings ***
Documentation  Test has been killed, cleanup VCH
Resource  Util.robot

*** Test Cases ***
ForceCleanup
    ${ret}=  Run  bin/vic-machine-linux delete --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name=%{TEST_VCH_NAME} --force
    Log To Console  vic-machine delete return code: ${ret}