*** Settings ***
Documentation  Test 1-10 - Docker PS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Assert VM Power State
    [Arguments]  ${name}  ${state}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.info -json ${vch-name}/${name}-* | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal  ${output}  ${state}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.info -json ${name}-* | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal  ${output}  ${state}

*** Test Cases ***
Empty docker ps command
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  CONTAINER ID
    Should Contain  ${output}  IMAGE
    Should Contain  ${output}  COMMAND
    Should Contain  ${output}  CREATED
    Should Contain  ${output}  STATUS
    Should Contain  ${output}  PORTS
    Should Contain  ${output}  NAMES
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

Docker ps only running containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container2}=  Run And Return Rc And Output  docker ${params} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container3}=  Run And Return Rc And Output  docker ${params} create busybox dmesg
    Should Be Equal As Integers  ${rc}  0
    Wait Until VM Powers Off  *-${container2}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

Docker ps all containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    Should Contain  ${output}  dmesg
    Should Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  4

Docker ps powerOn container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name jojo busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1
    # powerOn container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -on ${vch-name}/"jojo*"
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -on "jojo*"
    # a complete powerOn can take some time with reconfigures, so let's ensure state before we proceed
    Wait Until Keyword Succeeds  20x  500 milliseconds  Assert VM Power State  jojo  poweredOn
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

Docker ps powerOff container OOB
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2
    # PowerOff VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -off ${vch-name}/"jojo*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -off "jojo*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Wait Until VM Powers Off  "jojo*"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

Docker ps Remove container OOB
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  4
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy ${vch-name}/"jojo*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "jojo*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Wait Until VM Is Destroyed  "jojo*"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  3

Docker ps last container
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -l
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2

Docker ps two containers
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -n=2
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  dmesg
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  3

Docker ps last container with size
    ${status}=  Get State Of Github Issue  1545
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1545 has been resolved
    Log  Issue \#1545 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -ls
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  SIZE
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2

Docker ps all containers with only IDs
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  CONTAINER ID
    Should Not Contain  ${output}  /bin/top
    Should Not Contain  ${output}  dmesg
    Should Not Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  3

Docker ps with filter
    ${status}=  Get State Of Github Issue  1676
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1676 has been resolved
    Log  Issue \#1676 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -f status=created
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2