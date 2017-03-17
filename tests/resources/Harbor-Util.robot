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
Documentation  This resource provides any keywords related to the Harbor private registry appliance
Library  Selenium2Library

*** Variables ***
${HARBOR_SHORT_VERSION}  0.5.0
${HARBOR_VERSION}  harbor_0.5.0-9e4c90e

*** Keywords ***
Install Harbor To Test Server
    [Arguments]  ${user}=%{TEST_USERNAME}  ${password}=%{TEST_PASSWORD}  ${host}=%{TEST_URL}  ${datastore}=%{TEST_DATASTORE}  ${network}=%{BRIDGE_NETWORK}  ${name}=harbor  ${protocol}=http
    ${out}=  Run  wget https://github.com/vmware/harbor/releases/download/${HARBOR_SHORT_VERSION}/${HARBOR_VERSION}.ova
    ${out}=  Run  ovftool ${HARBOR_VERSION}.ova ${HARBOR_VERSION}.ovf
    ${out}=  Run  ovftool --acceptAllEulas --datastore=${datastore} --name=${name} --net:"Network 1"="${network}" --diskMode=thin --powerOn --X:waitForIp --X:injectOvfEnv --X:enableHiddenProperties --prop:vami.domain.Harbor=mgmt.local --prop:vami.searchpath.Harbor=mgmt.local --prop:vami.DNS.Harbor=8.8.8.8 --prop:vm.vmname=Harbor --prop:root_pwd=${password} --prop:harbor_admin_password=${password} --prop:verify_remote_cert=false --prop:protocol=${protocol} ${HARBOR_VERSION}.ovf 'vi://${user}:${password}@${host}'
    ${out}=  Split To Lines  ${out}

    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  Received IP address:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Environment Variable  HARBOR_IP  ${ip}
    \   Return From Keyword If  ${status}

    Fail  Harbor failed to install

Restart Docker With Insecure Registry Option
    # Requires you to edit /etc/systemd/system/docker.service.d/overlay.conf or docker.conf to be:
    # ExecStart=/bin/bash -c "usr/bin/docker daemon -H fd:// -s overlay $DOCKER_OPTS --insecure-registry=`cat /tmp/harbor`"
    ${out}=  Run  sudo service docker stop
    ${out}=  Run  echo %{HARBOR_IP} > /tmp/harbor
    ${out}=  Run  sudo service docker start

Log Into Harbor
    [Arguments]  ${user}=%{TEST_USERNAME}  ${pw}=%{TEST_PASSWORD}
    Open Browser  http://%{HARBOR_IP}/
    Maximize Browser Window
    Input Text  username  ${user}
    Input Text  uPassword  ${pw}
    Click button  Sign In
    Wait Until Page Contains  Summary
    Wait Until Page Contains  My Projects:

Create A New Project
    [Arguments]  ${name}
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Click Button  css=button.btn-success:nth-child(2)
    Input Text  uProjectName  ${name}
    Wait Until Element Is Visible  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Wait Until Element Is Enabled  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Sleep  1
    Select Checkbox  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Click Button  Save
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > div.each-tab-pane > div > div.table-body-container > table  ${name}