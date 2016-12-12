*** Settings ***
Documentation  This resource provides keywords to interact with Github

*** Keywords ***
Get State Of Github Issue
    [Arguments]  ${num}
    [Tags]  secret
    :FOR  ${idx}  IN RANGE  0  5
    \   ${status}  ${result}=  Run Keyword And Ignore Error  Get  https://api.github.com/repos/vmware/vic/issues/${num}?access_token\=%{GITHUB_AUTOMATION_API_KEY}
    \   Exit For Loop If  '${status}'
    \   Sleep  1
    Should Be Equal  ${result.status_code}  ${200}
    ${status}=  Get From Dictionary  ${result.json()}  state
    [Return]  ${status}
