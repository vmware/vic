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
    \   ${out}=  Run  docker ${params} ps
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  /bin/top
    \   Exit For Loop If  ${status}
    \   Sleep  1

*** Test Cases ***
Docker run -it date
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} run -it busybox date
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  UTC

Docker run -it df
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} run -it busybox df
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  Filesystem

Docker run -it command that doesn't stop
    ${rc}  ${out}=  Run And Return Rc And Output  docker ${params} ps -aq | xargs -n1 docker ${params} rm -f
    ${result}=  Start Process  docker ${params} run -itd busybox /bin/top  shell=True  alias=top

    Make sure container starts
    ${containerID}=  Run  docker ${params} ps -q
    ${out}=  Run  docker ${params} logs ${containerID}
    Should Contain  ${out}  Mem:
    Should Contain  ${out}  CPU:
    Should Contain  ${out}  Load average:

Docker run with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -i busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Docker run with -it
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it busybox /bin/ash -c "dmesg;echo END_OF_THE_TEST"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

Hello world with -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -i hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/

Hello world with -it
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run -it hello-world
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://docs.docker.com/engine/userguide/
