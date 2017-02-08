*** Settings ***
Documentation  Test 1-21 - Docker Volume LS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple volume ls
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=testVol
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${output}  testVol
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  vsphere
    Should Contain  ${output}  testVol
    Should Contain  ${output}  DRIVER
    Should Contain  ${output}  VOLUME NAME
    
Volume ls quiet
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -q
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  vsphere
    Should Contain  ${output}  testVol
    Should Not Contain  ${output}  DRIVER
    Should Not Contain  ${output}  VOLUME NAME

Volume ls invalid filter
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f bogusfilter=test
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: Invalid filter 'bogusfilter'

Volume ls filter by dangling
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=danglingVol
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v testVol:/test busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=true
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=false
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  testVol
    Should Not Contain  ${output}  danglingVol

Volume ls filter by name
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f name=dang
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol

Volume ls filter by driver
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f driver=vsphere
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  danglingVol
    Should Contain  ${output}  testVol

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f driver=vsph
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol

Volume ls filter by label
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=labelVol --label=labeled
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f label=labeled
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  labelVol
    Should Not Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol

Volume ls multiple filters
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=true -f name=dang
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls -f dangling=false -f name=dang
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  danglingVol
    Should Not Contain  ${output}  testVol
