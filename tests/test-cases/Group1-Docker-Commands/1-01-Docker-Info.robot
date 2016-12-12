*** Settings ***
Documentation  Test 1-01 - Docker Info
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Get resource pool CPU and mem values
    [Arguments]  ${info}

    ${cpuline}=  Get Lines Containing String  ${info}  CPUs:
    ${memline}=  Get Lines Containing String  ${info}  Total Memory:
    @{cpuline}=  Split String  ${cpuline}
    Length Should Be  ${cpuline}  2
    @{memline}=  Split String  ${memline}
    Length Should Be  ${memline}  4
    ${cpuval}=  Set Variable  @{cpuline}[1]
    ${memunit}=  Set Variable  @{memline}[3]
    ${memval}=  Set Variable  @{memline}[2]
    # Since govc accepts a mem value only in MB, convert the value if necessary
    ${memval}=  Run Keyword If  '${memunit}' == 'GiB'  Evaluate  int(round(${memval} * 1024))  ELSE  Evaluate  ${memval}

    [Return]  ${cpuval}  ${memval}

Set resource pool CPU and mem values
    [Arguments]  ${cpuval}  ${memval}

    ${rc}  ${output}=  Run And Return Rc And Output  govc pool.change -cpu.limit=${cpuval} %{TEST_RESOURCE}/%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc pool.change -mem.limit=${memval} %{TEST_RESOURCE}/%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

*** Test Cases ***
Basic Info
    Log To Console  \nRunning docker info command...
    ${output}=  Run  docker %{VCH-PARAMS} info
    Log  ${output}
    Should Contain  ${output}  vSphere

Debug Info
    ${status}=  Get State Of Github Issue  780
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-1-Docker-Info.robot needs to be updated now that Issue #780 has been resolved
    #Log To Console  \nRunning docker -D info command...
    #${output}=  Run  docker %{VCH-PARAMS} -D info
    #Log  ${output}
    #Should Contain  ${output}  Debug mode

Correct container count
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  Containers: 0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${cid}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${cid}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  Containers: 1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${cid}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  Containers: 1
    Should Contain  ${output}  Running: 1

Check modified resource pool CPU and memory values
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${oldcpuval}  ${oldmemval}=  Get resource pool CPU and mem values  ${output}

    ${newcpuval}=  Evaluate  ${oldcpuval} - 1
    ${newmemval}=  Evaluate  1000
    Set resource pool CPU and mem values  ${newcpuval}  ${newmemval}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${cpuval}  ${memval}=  Get resource pool CPU and mem values  ${output}
    Should Be Equal As Integers  ${cpuval}  ${newcpuval}
    Should Be Equal As Integers  ${memval}  ${newmemval}

    Set resource pool CPU and mem values  ${oldcpuval}  ${oldmemval}
