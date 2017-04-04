# Copyright 2017 VMware, Inc. All Rights Reserved.
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
Documentation   Test 1-38 - Docker Stats
Resource        ../../resources/Util.robot
Suite Setup     Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create test containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --name stresser busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Set Environment Variable  STRESSED  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name stopper busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    Set Environment Variable  STOPPER  ${output}

Stats No Stream
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stats --no-stream
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Get Line  ${output}  1
    ${short}=  Get Container ShortID  %{STRESSED}
    Should Contain  ${output}  ${short}
    ${vals}=  Split String  ${output}
    ${vicMemory}=  Get From List  ${vals}  7
    # only care about the integer value of memory usage
    ${vicMemory}=  Fetch From Left  ${vicMemory}  .
    # get the latest memory value for the "stresser" vm
    ${rc}  ${vmomiMemory}=  Run And Return Rc And Output  govc metric.sample -n 1 -json vm/*${short} mem.active.average | jq -r .Sample[].Value[].Value[0]
    Should Be Equal As Integers  ${rc}  0
    Should Be True  ${vmomiMemory} > 0
    # convert to percent and move decimal
    ${percent}=  Evaluate  (${vmomiMemory}/2048000)*100
    ${diff}=  Evaluate  ${percent}-${vicMemory}
    # due to timing we could see some variation, but shouldn't exceed 5
    Should Be True  ${diff} < 5

Stats No Stream All Containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stats --no-stream -a
    Should Be Equal As Integers  ${rc}  0
    ${stress}=  Get Container ShortID  %{STRESSED}
    ${stop}=  Get Container ShortID  %{STOPPER}
    Should Contain  ${output}  ${stress}
    Should Contain  ${output}  ${stop}

Stats API Memory Validation
    ${rc}  ${apiMem}=  Run And Return Rc And Output  curl -s -k -H "Accept: application/json" -H "Content-Type: application/json" -X GET https://%{VCH-IP}:%{VCH-PORT}/containers/%{STRESSED}/stats?stream=false | jq -r .memory_stats.usage
    Should Be Equal As Integers  ${rc}  0
    ${stress}=  Get Container ShortID  %{STRESSED}
    ${rc}  ${vmomiMem}=  Run And Return Rc And Output  govc metric.sample -n 1 -json vm/*${stress} mem.active.average | jq -r .Sample[].Value[].Value[0]
    Should Be Equal As Integers  ${rc}  0
    ${vmomiMem}=  Evaluate  ${vmomiMem}*1024
    ${diff}=  Evaluate  ${apiMem}-${vmomiMem}
    ${diff}=  Set Variable  abs(${diff})
    Should Be True  ${diff} < 1000

Stats API CPU Validation
    ${rc}  ${apiCPU}=  Run And Return Rc And Output  curl -s -k -H "Accept: application/json" -H "Content-Type: application/json" -X GET https://%{VCH-IP}:%{VCH-PORT}/containers/%{STRESSED}/stats?stream=false | jq -r .cpu_stats.cpu_usage.percpu_usage[0]
    Should Be Equal As Integers  ${rc}  0
    ${stress}=  Get Container ShortID  %{STRESSED}
    ${rc}  ${vmomiCPU}=  Run And Return Rc And Output  govc metric.sample -json vm/*${stress} cpu.usagemhz.average | jq -r .Sample[].Value[0].Value[]
    Should Contain  ${vmomiCPU}  ${apiCPU}