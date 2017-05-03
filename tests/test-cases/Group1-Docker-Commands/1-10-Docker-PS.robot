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
Documentation  Test 1-10 - Docker PS
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Assert VM Power State
    [Arguments]  ${name}  ${state}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.info -json %{VCH-NAME}/${name}-* | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal  ${output}  ${state}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.info -json ${name}-* | jq -r .VirtualMachines[].Runtime.PowerState
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal  ${output}  ${state}

Create several containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox ls
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container2}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container3}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox dmesg
    Should Be Equal As Integers  ${rc}  0
    ${container2shortID}=  Get container shortID  ${container2}
    Wait Until VM Powers Off  *-${container2shortID}*

Assert Number Of Containers
    [Arguments]  ${num}  ${type}=-q
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps ${type}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${num}

Check Length Of PS
    [Arguments]  ${len}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len}

*** Test Cases ***
Empty docker ps command
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers
    ${status}=  Get State Of Github Issue  4997
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-10-Docker-PS.robot needs to be updated now that Issue #4997 has been resolved
    #Wait Until Keyword Succeeds  5x  5 seconds  Check Length of PS  ${len+1}

Docker ps all containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    Should Contain  ${output}  dmesg
    Should Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+3}

Docker ps powerOn container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name jojo busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}

    Power On VM OOB  jojo*

    Wait Until Keyword Succeeds  10x  6s  Assert Number Of Containers  ${len+1}

Docker ps powerOff container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name koko busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start koko
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -q
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}

    Power Off VM OOB  koko*

    Wait Until Keyword Succeeds  10x  6s  Assert Number Of Containers  ${len-1}

Docker ps ports output
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -p 8000:80 -p 8443:443 nginx
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  :8000->80/tcp
    Should Contain  ${output}  :8443->443/tcp

    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p 6379 redis:alpine
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ->6379/tcp

Docker ps Remove container OOB
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name lolo busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start lolo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop lolo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    # Remove container VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy %{VCH-NAME}/"lolo*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "lolo*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Wait Until VM Is Destroyed  "lolo*"
    Wait Until Keyword Succeeds  10x  6s  Assert Number Of Containers  ${len-1}  -aq
    ${rc}  ${output}=  Run Keyword If  '%{DATASTORE_TYPE}' == 'VSAN'  Run And Return Rc And Output  govc datastore.ls | grep "lolo*" | xargs -n1 govc datastore.rm
    Run Keyword If  '%{DATASTORE_TYPE}' == 'VSAN'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{DATASTORE_TYPE}' == 'Non_VSAN'  Run And Return Rc And Output  govc datastore.ls | grep ${container} | xargs -n1 govc datastore.rm
    Run Keyword If  '%{DATASTORE_TYPE}' == 'Non_VSAN'  Should Be Equal As Integers  ${rc}  0

Docker ps last container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -l
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  redis
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

Docker ps two containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -n=2
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  redis
    Should Contain  ${output}  nginx
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  3

Docker ps last container with size
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  SIZE
    Should Contain  ${output}  redis
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

Docker ps all containers with only IDs
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq
    ${output}=  Split To Lines  ${output}
    ${len}=  Get Length  ${output}
    Create several containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  CONTAINER ID
    Should Not Contain  ${output}  /bin/top
    Should Not Contain  ${output}  dmesg
    Should Not Contain  ${output}  ls
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  ${len+3}

Docker ps with status filter
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -f status=created
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  nginx
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  5

Docker ps with label and name filter
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name abe --label prod busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f label=prod
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  busybox
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f name=abe
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  busybox
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

Docker ps with volume filter
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v foo:/dir --name fooContainer busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=foo
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  fooContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=foo -f volume=bar
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  fooContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=fo
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

Docker ps with network filter
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create fooNet
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --net=fooNet --name fooNetContainer busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f network=fooNet
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f network=fooNet -f network=barNet
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f network=fo
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

Docker ps with volume and network filters
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=foo -f network=bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooContainer
    Should Not Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f network=bar -f volume=foo
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooContainer
    Should Not Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=foo -f volume=buz -f network=bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooContainer
    Should Not Contain  ${output}  fooNetContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -v buz:/dir --net=fooNet --name buzFooContainer busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a -f volume=buz -f network=fooNet
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  fooContainer
    Should Not Contain  ${output}  fooNetContainer
    Should Contain  ${output}  buzFooContainer
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  2
