*** Settings ***
Documentation  Test 1-30 - Docker Deploy
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker deploy
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} deploy %{GOPATH}/src/github.com/vmware/vic/demos/compose/voting-app/votingapp.dab
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  only supported with experimental daemon

