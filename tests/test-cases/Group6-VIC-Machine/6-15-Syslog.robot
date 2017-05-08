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
Documentation  Test 6-15 - Verify remote syslog
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword  Cleanup VIC Appliance On Test Server

*** Variables ***
${SYSLOG_FILE}  /var/log/daemon.log

*** Keywords ***
Get Remote PID
    [Arguments]  ${proc}
    ${pid}=  Execute Command  ps -C ${proc} -o pid=
    ${pid}=  Strip String  ${pid}
    [Return]  ${pid}

*** Test Cases ***
Add remote syslog to VCH
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --syslog-address tcp://%{SYSLOG_SERVER}:514 ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully

    ${output}=  Run  bin/vic-machine-linux debug --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Should Contain  ${output}  Completed successfully

    Get Docker Params  ${output}  ${true}

    ${vch-conn}=  Open Connection  %{VCH-IP}
    Login  root  password

    @{procs}=  Create List  port-layer-server  docker-engine-server  vic-init  vicadmin
    &{proc-pids}=  Create Dictionary

    :FOR  ${proc}  IN  @{procs}
    \     ${pid}=  Get Remote PID  ${proc}
    \     Set To Dictionary  ${proc-pids}  ${proc}  ${pid}

    ${syslog-conn}=  Open Connection  %{SYSLOG_SERVER}
    Login  %{SYSLOG_USER}  %{SYSLOG_PASSWD}

    ${out}=  Execute Command  cat ${SYSLOG_FILE}
    ${keys}=  Get Dictionary Keys  ${proc-pids}
    :FOR  ${proc}  IN  ${keys}
    \     ${pid}=  Get From Dictionary  ${proc-pids}  ${proc}
    \     Should Contain  ${out}  %{VCH-IP} ${proc}[${pid}]:

    ${out}  ${rc}=  Run And Return Rc And Output  docker ${VCH-PARAMS} ps -a
    Should Be Equal As Integers  ${rc}  0

    ${out}=  Execute Command  cat ${SYSLOG_FILE}
    ${pid}=  Get From Dictionary  ${proc-pids}  docker-engine-server
    Should Contain  ${out}  docker-engine-server[${pid}]: Calling GET /v1.25/containers/json?all=1
