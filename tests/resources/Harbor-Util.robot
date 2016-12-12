*** Settings ***
Documentation  This resource provides any keywords related to the Harbor private registry appliance

*** Keywords ***
Install Harbor To Test Server
    [Arguments]  ${user}=%{TEST_USERNAME}  ${password}=%{TEST_PASSWORD}  ${host}=%{TEST_URL}  ${datastore}=${TEST_DATASTORE}  ${network}=%{BRIDGE_NETWORK}  ${name}=harbor
    ${out}=  Run  wget https://github.com/vmware/harbor/releases/download/0.4.5/harbor_0.4.5_beta_respin2.ova
    ${out}=  Run  ovftool harbor_0.4.5_beta_respin2.ova harbor_0.4.5_beta_respin2.ovf
    ${out}=  Run  ovftool --datastore=${datastore} --name=${name} --net:"Network 1"="${network}" --diskMode=thin --powerOn --X:waitForIp --X:injectOvfEnv --X:enableHiddenProperties --prop:vami.domain.Harbor=mgmt.local --prop:vami.searchpath.Harbor=mgmt.local --prop:vami.DNS.Harbor=8.8.8.8 --prop:vm.vmname=Harbor harbor_0.4.5_beta_respin2.ovf 'vi://${user}:${password}@${host}'
    ${out}=  Split To Lines  ${out}

    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  Received IP address:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Environment Variable  HARBOR_IP  ${ip}
    \   Exit For Loop If  ${status}
