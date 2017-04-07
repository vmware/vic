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

Attach to fake container
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} attach fakeContainer
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${out}  Error: No such container: fakeContainer

Attach with short input
    ${status}=  Get State Of Github Issue  4367
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-09-Docker-Attach.robot needs to be updated now that Issue #4367 has been resolved
    Log  Issue \#4367 is blocking implementation  WARN
    ${rc}  ${tmp}=  Run And Return Rc And Output  mktemp -d -p /tmp
    Should Be Equal As Integers  ${rc}  0
    ${fifo}=  Catenate  SEPARATOR=/  ${tmp}  fifo
    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo ${fifo}
    Should Be Equal As Integers  ${rc}  0
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i busybox sort
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  20
    \     Start Process  docker %{VCH-PARAMS} start -ai ${containerID} < ${fifo}  shell=True  alias=custom
    \     Run  echo one > ${fifo}
    \     ${ret}=  Wait For Process  custom
    \     Log  ${ret.stderr}
    \     Should Be Equal  ${ret.stdout}  one
    \     Should Be Equal As Integers  ${ret.rc}  0
    \     Should Be Empty  ${ret.stderr}
    # should be removed once 4410 is fixed
    \     Sleep  1
    Run  rm -rf ${tmp}

Attach with short output
    ${status}=  Get State Of Github Issue  4410
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-09-Docker-Attach.robot needs to be updated now that Issue #4410 has been resolved
    Log  Issue \#4410 is blocking implementation  WARN
    Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox echo one
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  20
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -a ${containerID}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  one
    # should be removed once 4410 is fixed
    \     Sleep  1

Attach with short output with tty
    ${status}=  Get State Of Github Issue  4410
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-09-Docker-Attach.robot needs to be updated now that Issue #4410 has been resolved
    Log  Issue \#4410 is blocking implementation  WARN
    Run  docker %{VCH-PARAMS} pull busybox
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -t busybox echo one
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  20
    \     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -a ${containerID}
    \     Should Be Equal As Integers  ${rc}  0
    \     Should Be Equal  ${output}  one
    # should be removed once 4410 is fixed
    \     Sleep  1
