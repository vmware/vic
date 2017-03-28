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
    [Arguments]  ${user}=%{TEST_USERNAME}  ${password}=%{TEST_PASSWORD}  ${host}=%{TEST_URL_ARRAY}  ${datastore}=%{TEST_DATASTORE}  ${network}=%{BRIDGE_NETWORK}  ${name}=harbor  ${protocol}=http  ${verify}=off  ${datacenter}=%{TEST_DATACENTER}  ${cluster}=%{TEST_RESOURCE}
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
    ${out}=  Run  ovftool --noSSLVerify --acceptAllEulas --datastore=${datastore} --name=${name} --diskMode=thin --powerOn --X:waitForIp --X:injectOvfEnv --X:enableHiddenProperties --net:"Network 1"="${network}" --prop:vm.vmname=Harbor --prop:root_pwd='${password}' --prop:harbor_admin_password='${password}' --prop:db_password='${password}' --prop:auth_mode=db_auth --prop:permit_root_login=true --prop:verify_remote_cert=${verify} --prop:protocol=${protocol} ./bin/harbor_0.5.0-9e4c90e.ova 'vi://${user}:${password}@${host}${datacenter}/host/${cluster}'
    ${out}=  Split To Lines  ${out}
    
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  Received IP address:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Environment Variable  HARBOR_IP  ${ip}
    \   Return From Keyword If  ${status}

    Fail  Harbor failed to install

Restart Docker With Insecure Registry Option
    # Requires you to edit /etc/systemd/system/docker.service.d/overlay.conf or docker.conf to be:
    # ExecStart=/bin/bash -c "usr/bin/docker daemon -H fd:// -s overlay $DOCKER_OPTS --insecure-registry='cat /tmp/harbor'"
    # Requires to be run as root
    ${out}=  Run  service docker stop
    Log  ${out}
    # Since docker couldn't access systemd, just modify the docker env
    ${docker_opts}=  Get Environment Variable  DOCKER_OPTS  ${EMPTY}
    Log to Console  ${docker_opts}
    Set Environment Variable  DOCKER_OPTS  ${docker_opts} --insecure-registry=%{HARBOR_IP}
    ${out}=  Run  service docker status
    Log  ${out}
    ${rc}  ${output}=  Run And Return Rc And Output  service docker start
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Log to Console  %{DOCKER_OPTS}

Log Into Harbor
    [Arguments]  ${user}=%{TEST_USERNAME}  ${pw}=%{TEST_PASSWORD}
    Maximize Browser Window
    Input Text  username  ${user}
    Input Text  uPassword  ${pw}
    Click button  Sign In
    Wait Until Page Contains  Summary
    Wait Until Page Contains  My Projects:
    Wait Until Keyword Succeeds  5x  1  Page Should Contain Element  xpath=//optional-menu/div/a[contains(., '${user}')]

Create A New Project
    [Arguments]  ${name}  ${public}=${true}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
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
    [Arguments]  ${name}  ${email}  ${fullName}  ${password}  ${comments}
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Click Element  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    Click Element  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > ul > li:nth-child(1) > a
    
    Wait Until Element Is Visible  username
    Wait Until Element Is Enabled  username
    Wait Until Element Is Visible  email
    Wait Until Element Is Enabled  email
    Wait Until Element Is Visible  fullName
    Wait Until Element Is Enabled  fullName
    Wait Until Element Is Visible  password
    Wait Until Element Is Enabled  password
    Wait Until Element Is Visible  confirmPassword
    Wait Until Element Is Enabled  confirmPassword
    Wait Until Element Is Visible  comments
    Wait Until Element Is Enabled  comments

    Input Text  username  ${name}
    Input Text  email  ${email}
    Input Text  fullName  ${fullName}
    Input Text  password  ${password}
    Input Text  confirmPassword  ${password}
    Input Text  comments  ${comments}

    Sleep  1
    Click Button  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > div > div > form > div:nth-child(7) > div > button

    Wait Until Page Contains  New user added successfully.
    Click Button  css=div.in:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(3) > button
    Sleep  1

Toggle Admin Priviledges For User
    [Arguments]  ${user}
    Wait Until Element Is Visible  //a[@tag='admin_option']
    Wait Until Element Is Enabled  //a[@tag='admin_option']
    Click Link  Admin Options
    Table Should Contain  css=body > div.container-fluid.container-fluid-custom.ng-scope > div > div > div > list-user > div > div > div.pane > div.sub-pane > div.table-body-container > table  ${user}
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()-1]/toggle-admin/button[not(contains(@class, 'ng-hide'))]
    Wait Until Element Is Enabled  xpath=//td[text()='test-user']/../td[last()-1]/toggle-admin/button[not(contains(@class, 'ng-hide'))]
    ${oldPublicity}=  Get Text  xpath=//td[text()='test-user']/../td[last()-1]/toggle-admin/button[not(contains(@class, 'ng-hide'))]
    Click Button  xpath=//td[text()='test-user']/../td[last()-1]/toggle-admin/button[not(contains(@class, 'ng-hide'))]
    Sleep  1
    Log to Console  show result
    ${newPublicity}=  Get Text  xpath=//td[text()='test-user']/../td[last()-1]/toggle-admin/button[not(contains(@class, 'ng-hide'))]
    Should Not Be Equal  ${oldPublicity}  ${newPublicity}
    [return]  ${newPublicity}

Create Self Signed Cert
    ${out}=  Run  openssl req -newkey rsa:4096 -nodes -sha256 -keyout ca.key -x509 -days 365 -out ca.crt

Install Harbor Self Signed Cert
    ${out}=  Run  wget --auth-no-challenge --no-check-certificate --user admin --password %{TEST_PASSWORD} https://%{HARBOR_IP}/api/systeminfo/getcert
    ${out}=  Run  mkdir -p /etc/docker/certs.d/%{HARBOR_IP}
    Move File  getcert  /etc/docker/certs.d/%{HARBOR_IP}/ca.crt
    ${out}=  Run  systemctl deamon-reload
    ${out}=  Run  systemctl restart docker

Delete A User
    [Arguments]  ${user}
    Wait Until Element Is Visible  //a[@tag='admin_option']
    Wait Until Element Is Enabled  //a[@tag='admin_option']
    Click Link  Admin Options
    Wait Until Element Is Visible  xpath=//td[text()='${user}']/../td[last()]/a
    Wait Until Element Is Enabled  xpath=//td[text()='${user}']/../td[last()]/a
    Click Link  xpath=//td[text()='${user}']/../td[last()]/a
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Are you sure to delete the user "${user}" ?
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  css=div.table-body-container > table  ${user} 

Search For A User
    [Arguments]  ${keyword}
    Wait Until Element Is Visible  //a[@tag='admin_option']
    Wait Until Element Is Enabled  //a[@tag='admin_option']
    Click Link  Admin Options
    Wait Until Element Is Visible  txtSearchInput
    Wait Until Element Is Enabled  txtSearchInput
    Input Text  txtSearchInput  ${keyword}
    Wait Until Element Is Visible  css=span.input-group-btn > button
    Wait Until Element Is Enabled  css=span.input-group-btn > button
    Click Button  css=span.input-group-btn > button
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=div.table-body-container > table  ${keyword} 
    # check all result contains the search keyword
    Wait Until Element Is Visible  xpath=//tbody/tr/td[1]
    Wait Until Element Is Enabled  xpath=//tbody/tr/td[1]
    ${rowNum}=  Get Matching Xpath Count  xpath=//tbody/tr/td[1]
    ${names}=  Create List
    ${realRowNum}=  Evaluate  ${rowNum} + 1
    :FOR  ${idx}  IN RANGE  1  ${realRowNum}
    \  ${searchName}=  Get Text  xpath=//tbody/tr[${idx}]/td[1]
    \  Should Match Regexp  ${searchName}  .*${keyword}.*
    \  Append To List  ${names}  ${searchName}
    [return]  ${names}

Change User Information
    [Arguments]  ${email}  ${fullName}  ${comments}
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Click Link  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Visible  xpath=//a[contains(., 'Account Settings')]
    Wait Until Element Is Enabled  xpath=//a[contains(., 'Account Settings')]
    Click Link  xpath=//a[contains(., 'Account Settings')]
    Wait Until Element Is Visible  email
    Wait Until Element Is Enabled  email
    Wait Until Element Is Visible  fullName
    Wait Until Element Is Enabled  fullName
    Wait Until Element Is Visible  comments
    Wait Until Element Is Enabled  comments
    Input Text  email  ${email}
    Input Text  fullName  ${fullName}
    Input Text  comments  ${comments}

    Wait Until Element Is Visible  xpath=//input[@value='Save']
    Wait Until Element Is Enabled  xpath=//input[@value='Save']
    Click Element  xpath=//input[@value='Save']

    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  User profile has been changed successfully.
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Wait Until Page Contains  Summary
    Wait Until Page Contains  My Projects:

Change User Password
    [Arguments]  ${password}  ${newPassword}
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Click Link  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Visible  xpath=//a[contains(., 'Account Settings')]
    Wait Until Element Is Enabled  xpath=//a[contains(., 'Account Settings')]
    Click Link  xpath=//a[contains(., 'Account Settings')]
    Wait Until Element Is Visible  toggleChangePassword
    Wait Until Element Is Enabled  toggleChangePassword
    Click Link  toggleChangePassword
    Wait Until Element Is Visible  oldPassword
    Wait Until Element Is Enabled  oldPassword
    Wait Until Element Is Visible  password
    Wait Until Element Is Enabled  password
    Wait Until Element Is Visible  confirmPassword
    Wait Until Element Is Enabled  confirmPassword
    Input Text  oldPassword  ${password}
    Input Text  password  ${newPassword}
    Input Text  confirmPassword  ${newPassword}

    Wait Until Element Is Visible  xpath=//input[@value='Save']
    Wait Until Element Is Enabled  xpath=//input[@value='Save']
    Click Element  xpath=//input[@value='Save']
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Password has been changed successfully.
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Wait Until Page Contains  Summary
    Wait Until Page Contains  My Projects:

Logout Harbor
    Wait Until Element Is Visible  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Enabled  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Click Link  css=#bs-harbor-navbar-collapse-1 > optional-menu > div > a
    Wait Until Element Is Visible  xpath=//a[contains(., 'Log Out')]
    Wait Until Element Is Enabled  xpath=//a[contains(., 'Log Out')]
    Click Link  xpath=//a[contains(., 'Log Out')]
    Wait Until Keyword Succeeds  5x  1  Page Should Contain Element  xpath=//h4[text()='Login Now']

Sign up
    [Arguments]  ${name}  ${email}  ${fullName}  ${password}  ${comments}
    Wait Until Element Is Visible  xpath=//button[text()='Sign Up']
    Wait Until Element Is Enabled  xpath=//button[text()='Sign Up']
    Click Button  xpath=//button[text()='Sign Up']
    Wait Until Keyword Succeeds  5x  1  Page Should Contain Element  xpath=//button[text()='Sign Up']
    Wait Until Element Is Visible  username
    Wait Until Element Is Enabled  username
    Wait Until Element Is Visible  email
    Wait Until Element Is Enabled  email
    Wait Until Element Is Visible  fullName
    Wait Until Element Is Enabled  fullName
    Wait Until Element Is Visible  password
    Wait Until Element Is Enabled  password
    Wait Until Element Is Visible  confirmPassword
    Wait Until Element Is Enabled  confirmPassword
    Wait Until Element Is Visible  comments
    Wait Until Element Is Enabled  comments

    Input Text  username  ${name}
    Input Text  email  ${email}
    Input Text  fullName  ${fullName}
    Input Text  password  ${password}
    Input Text  confirmPassword  ${password}
    Input Text  comments  ${comments}
    
    Wait Until Element Is Visible  xpath=//button[text()='Sign Up']
    Wait Until Element Is Enabled  xpath=//button[text()='Sign Up']
    Click Button  xpath=//button[text()='Sign Up']
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Signed up successfully.
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)

    Wait Until Keyword Succeeds  5x  1  Page Should Contain Element  xpath=//h4[text()='Login Now']

Add A User To A Project
    # role should be one of the strings : 'Project Admin'/'Developer'/'Guest'
    [Arguments]  ${user}  ${project}  ${role}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[@tag='users']
    Wait Until Element Is Enabled  xpath=//a[@tag='users']
    Click Link  Users
    Wait Until Element Is Visible  css=button.btn-success
    Wait Until Element Is Enabled  css=button.btn-success
    Click Button  css=button.btn-success
    Wait Until Element Is Visible  addUsername
    Wait Until Element Is Enabled  addUsername
    Input Text  addUsername  ${user}
    Wait Until Element Is Visible  xpath=//span[contains(., '${role}')]/input
    Wait Until Element Is Enabled  xpath=//span[contains(., '${role}')]/input
    Select Checkbox  xpath=//span[contains(., '${role}')]/input
    Wait Until Element Is Visible  btnSave
    Wait Until Element Is Enabled  btnSave
    Click Button  btnSave
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=div.sub-pane > div:nth-child(2) > table  ${user}

Remove A User From A Project
    [Arguments]  ${user}  ${project}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']
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
    [Arguments]  ${user}  ${project}  ${role}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
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
    [Arguments]  ${project}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
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

Search For A Project
    # search for the project contains the keyword, and return all result as a list
    [Arguments]  ${keyword}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Wait Until Element Is Visible  xpath=//input[@ng-model='vm.projectName']
    Wait Until Element Is Enabled  xpath=//input[@ng-model='vm.projectName']
    Input Text  xpath=//input[@ng-model='vm.projectName']  ${keyword}
    Wait Until Element Is Visible  css=span.input-group-btn > button
    Wait Until Element Is Enabled  css=span.input-group-btn > button
    Click Button  css=span.input-group-btn > button
    Sleep  1
    Wait Until Keyword Succeeds  5x  1  Table Should Contain  css=div.table-body-container > table  ${keyword} 
    # check all result contains the search keyword
    Wait Until Element Is Visible  xpath=//tbody/tr/td[1]
    Wait Until Element Is Enabled  xpath=//tbody/tr/td[1]
    ${rowNum}=  Get Matching Xpath Count  xpath=//tbody/tr/td[1]
    ${names}=  Create List
    ${realRowNum}=  Evaluate  ${rowNum} + 2
    :FOR  ${idx}  IN RANGE  2  ${realRowNum}
    \  ${searchName}=  Get Text  xpath=//tbody/tr[${idx}]/td[1]
    \  Should Match Regexp  ${searchName}  .*${keyword}.*
    \  Append To List  ${names}  ${searchName}
    [return]  ${names}

Delete Repository From Project
    [Arguments]  ${image}  ${project}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']
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

Delete Image From Project
    [Arguments]  ${image}  ${tag}  ${project}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Table Should Contain  css=div.custom-sub-pane > div:nth-child(2) > table  ${project}
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']
    Click Link  xpath=//td/a[text()='${project}']
    Wait Until Element Is Visible  xpath=//a[contains(., '${project}/${image}')]
    Wait Until Element Is Enabled  xpath=//a[contains(., '${project}/${image}')]
    Click Link  xpath=//a[contains(., '${project}/${image}')]
    Wait Until Element Is Visible  xpath=//a[contains(., '${project}/${image}')]/../../../div[2]/div/table/tbody/tr/td[text()='${tag}']/../td[last()]/a
    Wait Until Element Is Enabled  xpath=//a[contains(., '${project}/${image}')]/../../../div[2]/div/table/tbody/tr/td[text()='${tag}']/../td[last()]/a
    Click Link  xpath=//a[contains(., '${project}/${image}')]/../../../div[2]/div/table/tbody/tr/td[text()='${tag}']/../td[last()]/a
    Wait Until Element Is Visible  css=div.modal.fade.in > div > div > div:nth-child(2)
    Wait Until Element Contains  css=div.modal.fade.in > div > div > div:nth-child(2)  Delete tag "${tag}" now?
    Wait Until Element Is Enabled  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)
    Click Button  css=div.modal.fade.in > div > div > div:nth-child(3) > button:nth-child(1)    
    Sleep  1
    # if it is last image in this repo, the repo will be deleted
    ${imageNum}=  Get Text  xpath=//a[contains(., '${project}/${image}')]/span[2]
    Run Keyword If  '${imageNum}'=='1'  Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  css=div.sub-pane  ${tag}
    ...  Else  Wait Until Keyword Succeeds  5x  1  Element Should Not Contain  xpath=//a[contains(., '${project}/${image}')]/../../../div[2]/div/table/tbody  ${tag}    
    
Toggle Publicity On Project
    [Arguments]  ${project}
    Wait Until Element Is Visible  //a[@tag='project']
    Wait Until Element Is Enabled  //a[@tag='project']
    Click Link  Projects
    Wait Until Element Is Visible  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Wait Until Element Is Enabled  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    ${oldPublicity}=  Get Text  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Click Button  //td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Sleep  1
    ${newPublicity}=  Get Text  xpath=//td/a[text()='${project}']/../../td[last()-1]/publicity-button/button
    Should Not Be Equal  ${oldPublicity}  ${newPublicity}
    [return]  ${newPublicity}

Go To HomePage
    Wait Until Element Is Visible  css=a.navbar-brand
    Wait Until Element Is Enabled  css=a.navbar-brand
    Click Link  css=a.navbar-brand

    Wait Until Page Contains  Summary
    Wait Until Page Contains  My Projects: