# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  This resource contains all keywords related to creating, deleting, maintaining an instance of Admiral
Library  Selenium2Library  5  5
Suite Teardown  Close All Browsers

*** Keywords ***
Install Admiral
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p 8282:8282 --name admiral vmware/admiral
    Should Be Equal As Integers  0  ${rc}
    Set Environment Variable  ADMIRAL_IP  %{VCH-IP}:8282
    :FOR  ${idx}  IN RANGE  0  10
    \   ${out}=  Run  curl %{ADMIRAL_IP}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  <body class="admiral-default">
    \   Return From Keyword If  ${status}
    \   Sleep  5
    Fail  Install Admiral failed: Admiral endpoint failed to respond to curl

Cleanup Admiral
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f admiral
    Should Be Equal As Integers  0  ${rc}

Login To Admiral
    [Arguments]  ${url}=localhost:8282  ${browser}=chrome
    Open Browser  ${url}  ${browser}
    Maximize Browser Window
    Wait Until Page Contains  Welcome!
    Wait Until Page Contains  This is the place for you to create, provision, manage and monitor containerized applications.
    Wait Until Page Contains  Let's get started now!
    Wait Until Element Is Visible  css=button.admiral-btn.enter-btn
    Wait Until Element Is Enabled  css=button.admiral-btn.enter-btn
    Click Button  css=button.admiral-btn.enter-btn
    Wait Until Element Is Visible  css=div.query-search-input-controls.form-control
    Wait Until Element Is Enabled  css=div.query-search-input-controls.form-control
    Wait Until Element Is Visible  css=a.admiral-btn.addHost-btn
    Wait Until Element Is Enabled  css=a.admiral-btn.addHost-btn

Add Host To Admiral
    [Arguments]  ${address}  ${tags}=${EMPTY}
    Wait Until Element Is Visible  css=a.nav-item.hosts
    Wait Until Element Is Enabled  css=a.nav-item.hosts
    Click Element  css=a.nav-item.hosts

    Wait Until Element Is Visible  css=div.query-search-input-controls.form-control
    Wait Until Element Is Enabled  css=div.query-search-input-controls.form-control
    Wait Until Element Is Visible  css=a.admiral-btn.addHost-btn
    Wait Until Element Is Enabled  css=a.admiral-btn.addHost-btn
    Click Element  css=a.admiral-btn.addHost-btn

    Wait Until Page Contains  Add Host
    Wait Until Element Is Visible  css=div.hostname-holder > input.form-control
    Wait Until Element Is Enabled  css=div.hostname-holder > input.form-control
    Wait Until Element Is Visible  css=button.btn.btn-default.dropdown-toggle
    Wait Until Element Is Enabled  css=button.btn.btn-default.dropdown-toggle

    Input Text  css=div.hostname-holder > input.form-control  ${address}

    Click Element  css=#resourcePool > div.col-sm-9 > div > div > button.btn.btn-default.dropdown-toggle
    Click Element  css=a[data-name=default-placement-zone]

    Click Element  css=#credential > div.col-sm-9 > div > div > button.btn.btn-default.dropdown-toggle
    Click Element  css=a[data-name=default-client-cert]

    Input Text  css=#tags > div.col-sm-9 > div > div > span > input.token-input.tt-input  ${tags}

    Wait Until Element Is Visible  css=button.btn.admiral-btn.saveHost
    Wait Until Element Is Enabled  css=button.btn.admiral-btn.saveHost
    Click Button  css=button.btn.admiral-btn.saveHost

Add Project to Admiral
    [Arguments]  ${name}
    Wait Until Element Is Visible  css=a.nav-item.placements
    Wait Until Element Is Enabled  css=a.nav-item.placements
    Click Element  css=a.nav-item.placements
    
    Wait Until Element Is Visible  css=div.right-context-panel > div.toolbar > div:nth-child(2) > a
    Wait Until Element Is Enabled  css=div.right-context-panel > div.toolbar > div:nth-child(2) > a
    Click Element  css=div.right-context-panel > div.toolbar > div:nth-child(2) > a

    Wait Until Page Contains  Project
    Wait Until Element Is Visible  css=div.right-context-panel > div.content > div > div.list-holder > div.inline-editable-list-holder > div.inline-editable-list > div.toolbar > a.new-item
    Wait Until Element Is Enabled  css=div.right-context-panel > div.content > div > div.list-holder > div.inline-editable-list-holder > div.inline-editable-list > div.toolbar > a.new-item
    Click Element  css=div.right-context-panel > div.content > div > div.list-holder > div.inline-editable-list-holder > div.inline-editable-list > div.toolbar > a.new-item
    Input Text  css=input.name-input  ${name}
    Click Element  css=a.resourceGroupEdit-save

    Table Should Contain  css=div.right-context-panel > div.content > div > div > div > div > table  ${name}

*** Test Cases ***
test
    Login To Admiral
    #Add Host To Admiral  test  test:test
    Add Project to Admiral  test
    