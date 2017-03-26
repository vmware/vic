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

Attach with short input
       ${status}=  Get State Of Github Issue  4166
       Run Keyword If  '${status}' == 'closed'  Fail  Test 1-09-Docker-Attach.robot needs to be updated now that Issue #4166 has been resolved
       Log  Issue \#4166 is blocking implementation  WARN
