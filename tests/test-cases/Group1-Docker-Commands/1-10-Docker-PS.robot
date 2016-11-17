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

Create several containers
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+1}

Docker ps all containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    Should Contain  ${output}  dmesg
    Should Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+3}

Docker ps start/stop status
    ${output}=  Run  docker ${params} ps -a
    Log to Console  ${output}
    ${output}=  Run  docker ${params} ps
    Should Contain  ${output}  Up
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  3
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name test-status busybox /bin/top
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start test-status; sleep 5
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop test-status; sleep 5
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Run  docker ${params} ps -a | grep test-status
    Should Contain  ${output}  Exited

Docker ps powerOn container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name jojo busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}

    Power On VM OOB  jojo*

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+1}

    Run  sleep 5
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps | grep jojo
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Up

Docker ps powerOff container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name koko busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start koko
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}

    Power Off VM OOB  koko*

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len-1}

    Run  sleep 5
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a | grep koko
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Exited

Docker ps ports output
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create -p 8000:80 -p 8443:443 nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  :8000->80/tcp
    Should Contain  ${output}  :8443->443/tcp

    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} run -d -p 6379 redis:alpine
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ->6379/tcp

Docker ps Remove container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create --name lolo busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start lolo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop lolo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy ${vch-name}/"lolo*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "lolo*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Wait Until VM Is Destroyed  "lolo*"
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len-1}

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
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers 
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -aq
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  CONTAINER ID
    Should Not Contain  ${output}  /bin/top
    Should Not Contain  ${output}  dmesg
    Should Not Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+3}

Docker ps with filter
    ${status}=  Get State Of Github Issue  1676
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #1676 has been resolved
    Log  Issue \#1676 is blocking implementation  WARN
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -f status=created
    #Should Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  ls
    #${output}=  Split To Lines  ${output}
    #Length Should Be  ${output}  2
