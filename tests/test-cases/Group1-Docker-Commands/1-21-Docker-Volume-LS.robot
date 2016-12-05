*** Settings ***
Documentation  Test 1-21 - Docker Volume LS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple volume ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=test
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  test
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  vsphere
    Should Contain  ${output}  test
    Should Contain  ${output}  DRIVER
    Should Contain  ${output}  VOLUME NAME
    
Volume ls quiet
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -q
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  vsphere
    Should Contain  ${output}  test
    Should Not Contain  ${output}  DRIVER
    Should Not Contain  ${output}  VOLUME NAME
    
Volume ls dangling volumes
    ${status}=  Get State Of Github Issue  1718
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-21-Docker-Volume-LS.robot needs to be updated now that Issue #1718 has been resolved
    Log  Issue \#1718 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=true
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  vsphere
    #Should Contain  ${output}  test
    #Should Contain  ${output}  DRIVER
    #Should Contain  ${output}  VOLUME NAME
    
Volume ls invalid filter
    ${status}=  Get State Of Github Issue  1718
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-21-Docker-Volume-LS.robot needs to be updated now that Issue #1718 has been resolved
    Log  Issue \#1718 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f name=test
    #Should Be Equal As Integers  ${rc}  1
    #Should Contain  ${output}  Error response from daemon: Invalid filter 'name'
    
Volume ls no dangling volumes
    ${status}=  Get State Of Github Issue  1718
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-21-Docker-Volume-LS.robot needs to be updated now that Issue #1718 has been resolved
    Log  Issue \#1718 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v test:/test busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=true
    #Should Be Equal As Integers  ${rc}  0
    #Log  ${output}