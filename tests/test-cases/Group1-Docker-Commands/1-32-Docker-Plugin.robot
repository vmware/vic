*** Settings ***
Documentation  Test 1-32 - Docker plugin
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Docker plugin install
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin install vieux/sshfs
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin create
    Run  echo '{}' > config.json
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin create test-plugin .
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin enable
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin enable test-plugin
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin disable
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin disable test-plugin
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin inspect
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin inspect test-plugin
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin ls
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin push
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin push test-plugin
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin rm
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin rm test-plugin
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins

Docker plugin set
    ${rc}  ${output}=  Run And Return Rc And Output  docker1.13 %{VCH-PARAMS} plugin set test-plugin test-data
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  does not yet support plugins
