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
Documentation  Test 10-01 - VCH Restart
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Default Tags

*** Keywords ***
Get Container IP
    [Arguments]  ${id}  ${network}=default
    ${rc}  ${ip}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network inspect ${network} | jq '.[0].Containers."${id}".IPv4Address'
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${ip}

Launch Container
    [Arguments]  ${name}  ${network}=default  ${command}=sh
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${name} --net ${network} -itd ${busybox} ${command}
    Should Be Equal As Integers  ${rc}  0
    ${id}=  Get Line  ${output}  -1
    [Return]  ${id}

*** Test Cases ***
Created Network And Images Persists As Well As Containers Are Discovered With Correct IPs
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${nginx}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create foo
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create bar
    Should Be Equal As Integers  ${rc}  0

    ${bridge-exited}=  Launch Container  vch-restart-bridge-exited  bridge  ls
    ${bridge-running}=  Launch Container  vch-restart-bridge-running  bridge
    ${bridge-running-ip}=  Get Container IP  ${bridge-running}  bridge
    ${bar-exited}=  Launch Container  vch-restart-bar-exited  bar  ls
    ${bar-running}=  Launch Container  vch-restart-bar-running  bar
    ${bar-running-ip}=  Get Container IP  ${bar-running}  bar

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -p 10000:80 -p 10001:80 --name webserver ${nginx}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10001

    Reboot VM  %{VCH-NAME}

    Log To Console  Getting VCH IP ...
    ${new-vch-ip}=  Get VM IP  %{VCH-NAME}
    Log To Console  New VCH IP is ${new-vch-ip}
    Replace String  %{VCH-PARAMS}  %{VCH-IP}  ${new-vch-ip}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  nginx
    Should Contain  ${output}  busybox

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  foo
    Should Contain  ${output}  bar
    Should Contain  ${output}  bridge
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network inspect bridge
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network inspect bar
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network inspect foo
    Should Be Equal As Integers  ${rc}  0

    ${ip}=  Get Container IP  ${bridge-running}  bridge
    Should Be Equal  ${ip}  ${bridge-running-ip}
    ${ip}=  Get Container IP  ${bar-running}  bar
    Should Be Equal  ${ip}  ${bar-running-ip}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${bridge-running} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"running\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${bar-running} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"running\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${bridge-exited} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"exited\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${bar-exited} | jq '.[0].State.Status'
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal  ${output}  \"exited\"
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${bar-exited}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${bridge-exited}
    Should Be Equal As Integers  ${rc}  0

    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -p 10000:80 -p 10001:80 --name webserver1 ${nginx}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start webserver1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  port 10000 is not available

    # docker pull should work
    # if this fails, very likely the default gateway on the VCH is not set
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${alpine}
    Should Be Equal As Integers  ${rc}  0

Create VCH attach disk and reboot
    ${rc}=  Run And Return Rc  govc vm.disk.create -vm=%{VCH-NAME} -name=%{VCH-NAME}/deleteme -size "16M"
    Should Be Equal As Integers  ${rc}  0

    Reboot VM  %{VCH-NAME}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}
    ${rc}=  Run And Return Rc  govc device.ls -vm=%{VCH-NAME} | grep disk
    Should Be Equal As Integers  ${rc}  1

Docker inspect mount and cmd data after reboot
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=named-volume
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=mount-data-test -v /mnt/test -v named-volume:/mnt/named busybox /bin/ls -la /
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Mounts}}' mount-data-test
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  /mnt/test
    Should Contain  ${out}  /mnt/named

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Config.Cmd}}' mount-data-test
    Should Be Equal As Integers  ${rc}  0
    Should Contain X Times  ${out}  /bin/ls  1
    Should Contain X Times  ${out}  -la  1
    Should Contain X Times  ${out}  ${SPACE}/  1

    Reboot VM  %{VCH-NAME}

    # wait for docker info to succeed
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Mounts}}' mount-data-test
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  /mnt/test
    Should Contain  ${out}  /mnt/named

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Config.Cmd}}' mount-data-test
    Should Be Equal As Integers  ${rc}  0
    Should Contain X Times  ${out}  /bin/ls  1
    Should Contain X Times  ${out}  -la  1
    Should Contain X Times  ${out}  ${SPACE}/  1
