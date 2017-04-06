# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 17-1 - TTY Tests
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Make sure container starts
    :FOR  ${idx}  IN RANGE  0  30
    \   ${out}=  Run  docker %{VCH-PARAMS} ps
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  /bin/top
    \   Exit For Loop If  ${status}
    \   Sleep  1

*** Test Cases ***
Docker run -it date
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it busybox date
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  UTC

Docker run -it df
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it busybox df
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  Filesystem

Docker run -it command that doesn't stop
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq | xargs -n1 docker %{VCH-PARAMS} rm -f
    ${result}=  Start Process  docker %{VCH-PARAMS} run -itd busybox /bin/top  shell=True  alias=top

    Make sure container starts
    ${containerID}=  Run  docker %{VCH-PARAMS} ps -q
    ${out}=  Run  docker %{VCH-PARAMS} logs ${containerID}
    Should Contain  ${out}  Mem:
    Should Contain  ${out}  CPU:
    Should Contain  ${out}  Load average:

Docker run with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -i busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Docker run with -it
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Hello world with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -i hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/

Hello world with -it
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/

Start with attach and interactive
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -ai ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error:

Stop a container stuck in starting state
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${name}=  Generate Random String  15
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it --name ${name} busybox /bin/date
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc device.disconnect -vm */${name}* serialport-9000
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${name}
    Should Contain  ${output}  Error:
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${name}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -json */${name}* | jq -r .VirtualMachines[].Runtime.PowerState
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  poweredOff

Attach with custom detach keys
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox /bin/top
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
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it busybox /bin/top
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

Exec Echo -it
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  5
    \   ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -it ${d} /bin/echo "I find your lack of faith disturbing."
    \   Should Be Equal As Integers  ${rc}  0
    \   Should Be Equal As Strings  ${output}  I find your lack of faith disturbing.

Exec Sort -it
    ${rc}  ${tmp}=  Run And Return Rc And Output  mktemp -d -p /tmp
    Should Be Equal As Integers  ${rc}  0
    ${fifo}=  Catenate  SEPARATOR=/  ${tmp}  fifo
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo ${fifo}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  5
    \     Start Process  docker %{VCH-PARAMS} exec -it ${output} /bin/sort < ${fifo}  shell=True  alias=custom
    \     Run  echo one > ${fifo}
    \     ${ret}=  Wait For Process  custom
    \     Log  ${ret.stderr}
    \     Should Be Equal  ${ret.stdout}  cannot enable tty mode on non tty input
    \     Should Be Equal As Integers  ${ret.rc}  0
    \     Should Be Empty  ${ret.stderr}
    Run  rm -rf ${tmp}
