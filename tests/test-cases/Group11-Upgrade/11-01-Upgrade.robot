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
Documentation  Test 11-01 - Upgrade
Resource  ../../resources/Util.robot
Suite Setup  Install VIC with version to Test Server  7315
Suite Teardown  Clean up VIC Appliance And Local Binary
Default Tags

*** Keywords ***
Install VIC with version to Test Server
    [Arguments]  ${version}=7315
    Log To Console  \nDownloading vic ${version} from bintray...
    ${rc}  ${output}=  Run And Return Rc And Output  wget https://bintray.com/vmware/vic-repo/download_file?file_path=vic_${version}.tar.gz -O vic.tar.gz
    ${rc}  ${output}=  Run And Return Rc And Output  tar zxvf vic.tar.gz
	Set Environment Variable  TEST_TIMEOUT  20m0s
	Install VIC Appliance To Test Server  vic-machine=./vic/vic-machine-linux  appliance-iso=./vic/appliance.iso  bootstrap-iso=./vic/bootstrap.iso  certs=${false}
    Set Environment Variable  VIC-ADMIN  %{VCH-IP}:2378

Clean up VIC Appliance And Local Binary
    Cleanup VIC Appliance On Test Server
    Run  rm -rf vic.tar.gz vic

Launch Container
    [Arguments]  ${name}  ${network}=default
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${name} --net ${network} -itd busybox
    Should Be Equal As Integers  ${rc}  0
    ${id}=  Get Line  ${output}  -1
    ${ip}=  Get Container IP  %{VCH-PARAMS}  ${id}  ${network}
    [Return]  ${id}  ${ip}

*** Test Cases ***
Upgrade Present in vic-machine
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux
    Should Contain  ${output}  upgrade

Upgrade VCH with containers
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create bar
    Should Be Equal As Integers  ${rc}  0
    Comment  Launch container on bridge network
    ${id1}  ${ip1}=  Launch Container  vch-restart-test1  bridge
    ${id2}  ${ip2}=  Launch Container  vch-restart-test2  bridge

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10001

    Log To Console  \nUpgrading VCH with 1s timeout ...
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux upgrade --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout 1s
    Should Contain  ${output}  Upgrading VCH exceeded time limit
    Should Not Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc snapshot.tree -vm=%{VCH-NAME}/%{VCH-NAME}
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Not Contain  ${output}  upgrade
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc snapshot.tree -vm=%{VCH-NAME}
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Not Contain  ${output}  upgrade

    Log To Console  \nUpgrading VCH...
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux upgrade --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    Should Not Contain  ${output}  Rolling back upgrade
    Should Be Equal As Integers  ${rc}  0

    # check version
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux version
    @{vers}=  Split String  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  Completed successfully
    Should Contain  ${output}  @{vers}[2]
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Get Docker Params  ${output}  ${true}

    Log To Console  \nTesting rollback...
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux upgrade --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --rollback
    Should Contain  ${output}  Completed successfully
    Should Be Equal As Integers  ${rc}  0

    # check version again, but old version this time
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux version
    @{vers}=  Split String  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  Completed successfully
    Should Contain  ${output}  @{vers}[1]
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Get Docker Params  ${output}  ${true}

    # upgrade again to latest version
    Log To Console  \nUpgrading VCH...
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux upgrade --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    Should Not Contain  ${output}  Rolling back upgrade
    Should Be Equal As Integers  ${rc}  0

    # check version
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux version
    @{vers}=  Split String  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  Completed successfully
    Should Contain  ${output}  @{vers}[2]
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Get Docker Params  ${output}  ${true}

    # wait for docker info to succeed
    Log To Console  Verify Containers...
    Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  %{VCH-PARAMS}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network inspect bridge
    Should Be Equal As Integers  ${rc}  0
    ${ip}=  Get Container IP  %{VCH-PARAMS}  ${id1}  bridge
    Should Be Equal  ${ip}  ${ip1}
    ${ip}=  Get Container IP  %{VCH-PARAMS}  ${id2}  bridge
    Should Be Equal  ${ip}  ${ip2}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  "Id"
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Exited (143)
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start vch-restart-test1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Exited (0)

    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  %{VCH-IP}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -it -p 10000:80 -p 10001:80 --name webserver1 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start webserver1
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  port 10000 is not available

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq | xargs -n1 docker %{VCH-PARAMS} stop
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq | xargs -n1 docker %{VCH-PARAMS} rm
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    Log To Console  Regression Tests...
    Run Regression Tests
