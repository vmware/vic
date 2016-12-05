*** Settings ***
Documentation  Test 1-09 - Docker Attach
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Basic attach
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    Start Process  docker %{VCH-PARAMS} attach ${containerID} < /tmp/fifo  shell=True  alias=custom
    Sleep  3
    Run  echo q > /tmp/fifo
    ${ret}=  Wait For Process  custom
    Should Be Equal As Integers  ${ret.rc}  0
    Should Be Empty  ${ret.stdout}
    Should Be Empty  ${ret.stderr}

Attach to stopped container
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} attach ${out}
    Should Be Equal As Integers  ${rc}  1
    Should Be Equal  ${out}  You cannot attach to a stopped container, start it first

Attach with custom detach keys
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    Start Process  docker %{VCH-PARAMS} attach --detach-keys\=a ${containerID} < /tmp/fifo  shell=True  alias=custom
    Sleep  3
    Run  echo a > /tmp/fifo
    ${ret}=  Wait For Process  custom
    Should Be Equal As Integers  ${ret.rc}  0
    Should Be Empty  ${ret.stdout}
    Should Be Empty  ${ret.stderr}

Reattach to container
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerID}
    Should Be Equal As Integers  ${rc}  0
    Start Process  docker %{VCH-PARAMS} attach --detach-keys\=a ${containerID} < /tmp/fifo  shell=True  alias=custom
    Sleep  3
    Run  echo a > /tmp/fifo
    ${ret}=  Wait For Process  custom
    Should Be Equal As Integers  ${ret.rc}  0
    Should Be Empty  ${ret.stdout}
    Should Be Empty  ${ret.stderr}
    Start Process  docker %{VCH-PARAMS} attach --detach-keys\=a ${containerID} < /tmp/fifo  shell=True  alias=custom2
    Sleep  3
    Run  echo a > /tmp/fifo
    ${ret}=  Wait For Process  custom2
    Should Be Equal As Integers  ${ret.rc}  0
    Should Be Empty  ${ret.stdout}
    Should Be Empty  ${ret.stderr}

Attach to fake container
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} attach fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${out}  Error: No such container: fakeContainer