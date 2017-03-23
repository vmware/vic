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

*** Variables ***
${HARBOR_SHORT_VERSION}  0.5.0
${HARBOR_VERSION}  harbor_0.5.0-9e4c90e

*** Keywords ***
Install Harbor To Test Server
    [Arguments]  ${user}=%{TEST_USERNAME}  ${password}=%{TEST_PASSWORD}  ${host}=%{TEST_URL_ARRAY}  ${datastore}=%{TEST_DATASTORE}  ${network}=%{BRIDGE_NETWORK}  ${name}=harbor  ${protocol}=http  ${verify}=false
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  DRONE_BUILD_NUMBER
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  DRONE_BUILD_NUMBER  0

    @{URLs}=  Split String  %{TEST_URL_ARRAY}
    ${len}=  Get Length  ${URLs}
    ${IDX}=  Evaluate  %{DRONE_BUILD_NUMBER} \% ${len}

    Set Environment Variable  TEST_URL  @{URLs}[${IDX}]
    Set Environment Variable  GOVC_URL  %{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}

    Log To Console  Downloading Harbor OVA...
    ${out}=  Run  wget https://github.com/vmware/harbor/releases/download/${HARBOR_SHORT_VERSION}/${HARBOR_VERSION}.ova
    Log To Console  Generating OVF file...
    ${out}=  Run  ovftool ${HARBOR_VERSION}.ova ${HARBOR_VERSION}.ovf
    Log To Console  Installing Harbor into test server...
    ${out}=  Run  ovftool --noSSLVerify --acceptAllEulas --datastore=${datastore} --name=${name} --net:"Network 1"="${network}" --diskMode=thin --powerOn --X:waitForIp --X:injectOvfEnv --X:enableHiddenProperties --prop:vami.domain.Harbor=mgmt.local --prop:vami.searchpath.Harbor=mgmt.local --prop:vami.DNS.Harbor=8.8.8.8 --prop:vm.vmname=Harbor --prop:root_pwd=${password} --prop:harbor_admin_password=${password} --prop:verify_remote_cert=${verify} --prop:protocol=${protocol} ${HARBOR_VERSION}.ovf 'vi://${user}:${password}@${host}'
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
    # Requires to be run as root
    ${out}=  Run  service docker stop
    ${out}=  Run  echo %{HARBOR_IP} > /tmp/harbor
    ${out}=  Run  service docker start

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
    [Arguments]  ${name}  ${public}=${true}
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Click Button  css=button.btn-success:nth-child(2)
    Wait Until Element Is Visible  name=uProjectName
    Wait Until Element Is Enabled  name=uProjectName
    Input Text  uProjectName  ${name}
    Wait Until Element Is Visible  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Wait Until Element Is Enabled  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Sleep  1
    Run Keyword If  ${public}  Select Checkbox  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > add-project > div > form > div > div.col-xs-10.col-md-10 > div:nth-child(2) > label > input
    Click Button  Save
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > div.each-tab-pane > div > div.table-body-container > table  ${name}

Create A New User
    [Arguments]  ${name}  ${email}  ${fullName}  ${password}
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Click Element  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    Click Element  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    
    Wait Until Element Is Visible  username
    Wait Until Element Is Visible  email
    Wait Until Element Is Visible  fullName
    Wait Until Element Is Visible  password
    Wait Until Element Is Visible  confirmPassword
    
    Input Text  username  ${name}
    Input Text  email  ${email}
    Input Text  fullName  ${fullName}
    Input Text  password  ${password}
    Input Text  confirmPassword  ${password}
    
    Sleep  1
    Click Button  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > div > div > form > div:nth-child(7) > div > button

    Wait Until Page Contains  New user added successfully.
    Click Button  css=div.in:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(3) > button
    Sleep  1

Toggle Admin Priviledges For User
    [Arguments]  ${username}
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > ul > li:nth-child(1) > navigation-header > ul > li:nth-child(3) > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > ul > li:nth-child(1) > navigation-header > ul > li:nth-child(3) > a
    Click Link  css=#bs-harbor-navbar-collapse-1 > ul > li:nth-child(1) > navigation-header > ul > li:nth-child(3) > a
    
    Table Should Contain  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > list-user > div > div > div.pane > div.sub-pane > div.table-body-container > table  ${username}
    
    ${rowNum}=  Get Matching Xpath Count  xpath=/html/body/div[1]/div/div/div/list-user/div/div/div[2]/div[1]/div[2]/table/tbody/tr[1]
    :FOR  ${idx}  IN RANGE  1  ${rowNum}
    \   ${status}=  Run Keyword And Return Status  Table Row Should Contain  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > list-user > div > div > div.pane > div.sub-pane > div.table-body-container > table  ${idx}  ${username}
    \   Run Keyword If  ${status}  Set Test Variable  ROW_NUMBER  ${idx}
    \   Exit For Loop If  ${status}
    
    Click Element  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > list-user > div > div > div.pane > div.sub-pane > div.table-body-container > table > tbody > tr:nth-child(${ROW_NUMBER}) > td:nth-child(4) > toggle-admin > button.btn.btn-danger.ng-binding

Create Self Signed Cert
    ${out}=  Run  openssl req -newkey rsa:4096 -nodes -sha256 -keyout ca.key -x509 -days 365 -out ca.crt

Install Harbor Self Signed Cert
    ${out}=  Run  wget --auth-no-challenge --no-check-certificate --user admin --password %{TEST_PASSWORD} https://%{HARBOR_IP}/api/systeminfo/getcert
    ${out}=  Run  mkdir -p /etc/docker/certs.d/%{HARBOR_IP}
    Move File  getcert  /etc/docker/certs.d/%{HARBOR_IP}/ca.crt
    ${out}=  Run  systemctl deamon-reload
    ${out}=  Run  systemctl restart docker

Delete A User
    [Arguments]  ${user}=%{TEST_USERNAME}
    Click Link  Admin Options
    Wait Until Element Is Visible  css=span.glyphicon-trash
    Wait Until Element Is Enabled  css=span.glyphicon-trash
    Click Link  xpath=//td[text()='${user}']/../td[last()]/a
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Are you sure to delete the user "${user}" ?
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  css=div.table-body-container > table  ${user} 


Add A User To A Project
    # role should be one of the strings : 'Project Admin'/'Developer'/'Guest'
    [Arguments]  ${user}=%{TEST_USERNAME}  ${project}=%{TEST_PROJECT}  ${role}=%{TEST_USER_ROLE}
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[@tag='users']
    Wait Until Element Is Enabled  xpath=//a[@tag='users']
    Click Link  Users
    Wait Until Element Is Visible  css=button.btn-success
    Wait Until Element Is Enabled  css=button.btn-success
    Click Button  css=button.btn-success
    Wait Until Element Is Visible  uUsername
    Wait Until Element Is Enabled  uUsername
    Input Text  uUsername  ${user}
    Wait Until Element Is Visible  xpath=//span[contains(., '${role}')]/input
    Wait Until Element Is Enabled  xpath=//span[contains(., '${role}')]/input
    Select Checkbox  xpath=//span[contains(., '${role}')]/input
    Wait Until Element Is Visible  btnSave
    Wait Until Element Is Enabled  btnSave
    Click Button  btnSave
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=div.sub-pane > div:nth-child(2) > table  ${user}

Remove A User From A Project
    [Arguments]  ${user}=%{TEST_USERNAME}  ${project}=%{TEST_PROJECT} 
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[@tag='users']
    Wait Until Element Is Enabled  xpath=//a[@tag='users']
    Click Link  Users
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[last()]
    Wait Until Element Is Enabled  xpath=//td[text()='${user}']/../td[last()]/a[last()]
    Click Link  xpath=//td[text()='${user}']/../td[last()]/a[last()]
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Page Should Not Contain  ${user}

Change A User's Role In A Project
    [Arguments]  ${user}=%{TEST_USERNAME}  ${project}=%{TEST_PROJECT}  ${role}=%{TEST_USER_ROLE}
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Wait Until Element Is Visible  css=div.custom-sub-pane > div:nth-child(2) > table 
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[@tag='users']
    Wait Until Element Is Enabled  xpath=//a[@tag='users']
    Click Link  Users
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Enabled  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]/span[@title='Edit']
    Click Link  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Visible  css=select.form-control
    Wait Until Element Is Enabled  css=select.form-control
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Enabled  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]/span[@title='Confirm']
    Select From List By Label  css=select  ${role}
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Wait Until Element Is Enabled  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Click Link  xpath=//td[text()='${user}']/../td[last()]/a[1]
    Sleep  1
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a[1]/span[@title='Edit']
    Wait Until Keyword Succeeds  5x  1  Page Should Contain Element  //td[text()='${user}']/../td[2]/switch-role/ng-switch/span[text()='${role}']

Delete A Project
    [Arguments]  ${project}=%{TEST_PROJECT}
    Click Link  Projects
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']/../../td[last()]/a
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']/../../td[last()]/a
    Click Link  xpath=//td/a[text()='${project}']/../../td[last()]/a
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Are you sure to delete the project "${project}" ?
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  css=div.table-body-container > table  ${project}

Delete Repository From Project
    [Arguments]  ${image}  ${project}=%{TEST_PROJECT} 
    Click Link  Projects
    Wait Until Element Is Visible  css=button.btn-success:nth-child(2)
    Wait Until Element Is Enabled  css=button.btn-success:nth-child(2)
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[contains(., '${project}/${image}')]/../a[last()]
    Wait Until Element Is Enabled  xpath=//a[contains(., '${project}/${image}')]/../a[last()]
    Click Link  xpath=//a[contains(., '${project}/${image}')]/../a[last()]
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Delete repository "${project}/${image}" now?
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  css=div.sub-pane  ${project}/${image}

Toggle Publicity On Project
    [Arguments]  ${project}=%{TEST_PROJECT}
    Click Link  Projects
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    ${oldPublicity}=  Get Text  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Click Button  //td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Sleep  1
    ${newPublicity}=  Get Text  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Wait Until Keyword Succeeds  5x  1  Should Not Be Equal  ${oldPublicity}  ${newPublicity}

