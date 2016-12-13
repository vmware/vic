*** Settings ***
Documentation  This resource contains any keywords related to using the Nimbus cluster

*** Variables ***
${ESX_VERSION}  %{ESX_VERSION}  #6.5 RTM
${VC_VERSION}  %{VC_VERSION}   #6.5 RTM

*** Keywords ***
Deploy Nimbus ESXi Server
    [Arguments]  ${user}  ${password}  ${version}=${ESX_VERSION}
    ${name}=  Evaluate  'ESX-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nDeploying Nimbus ESXi server: ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}

    ${out}=  Execute Command  nimbus-esxdeploy ${name} --disk=48000000 --ssd=24000000 --memory=8192 --nics 2 ${version}
    # Make sure the deploy actually worked
    Should Contain  ${out}  To manage this VM use
    # Now grab the IP address and return the name and ip for later use
    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  IP is
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    @{gotIP}=  Split String  ${line}  ${SPACE}
    ${ip}=  Remove String  @{gotIP}[5]  ,

    # Let's set a password so govc doesn't complain
    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${ip}
    ${out}=  Run  govc host.account.update -id root -password e2eFunctionalTest
    Should Be Empty  ${out}
    Disable TLS On ESX Host
    Log To Console  Successfully deployed new ESXi server - ${user}-${name}
    Close connection
    [Return]  ${user}-${name}  ${ip}

Deploy Multiple Nimbus ESXi Servers in Parallel
    [Arguments]  ${user}  ${password}  ${version}=${ESX_VERSION}
    ${name1}=  Evaluate  'ESX-' + str(random.randint(1000,9999))  modules=random
    ${name2}=  Evaluate  'ESX-' + str(random.randint(1000,9999))  modules=random
    ${name3}=  Evaluate  'ESX-' + str(random.randint(1000,9999))  modules=random

    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}

    ${out1}=  Deploy Nimbus ESXi Server Async  ${name1}
    ${out2}=  Deploy Nimbus ESXi Server Async  ${name2}
    ${out3}=  Deploy Nimbus ESXi Server Async  ${name3}

    Wait For Process  ${out1}
    Wait For Process  ${out2}
    Wait For Process  ${out3}

    ${out}=  Execute Command  nimbus-ctl ip ${user}-${name1}

    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  ${user}-${name1}
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    @{gotIP}=  Split String  ${line}  ${SPACE}
    ${ip1}=  Remove String  @{gotIP}[2]

    ${out}=  Execute Command  nimbus-ctl ip ${user}-${name2}

    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  ${user}-${name2}
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    @{gotIP}=  Split String  ${line}  ${SPACE}
    ${ip2}=  Remove String  @{gotIP}[2]

    ${out}=  Execute Command  nimbus-ctl ip ${user}-${name3}

    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  ${user}-${name3}
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    @{gotIP}=  Split String  ${line}  ${SPACE}
    ${ip3}=  Remove String  @{gotIP}[2]

    Log To Console  \nDeploying Nimbus ESXi server: ${gotIP}

    # Let's set a password so govc doesn't complain
    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${ip1}
    ${out}=  Run  govc host.account.update -id root -password e2eFunctionalTest
    Should Be Empty  ${out}
    Disable TLS On ESX Host
    Log To Console  Successfully deployed new ESXi server - ${user}-${name1}
    Log To Console  \nNimbus ESXi server IP: ${ip1}

    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${ip2}
    ${out}=  Run  govc host.account.update -id root -password e2eFunctionalTest
    Should Be Empty  ${out}
    Disable TLS On ESX Host
    Log To Console  Successfully deployed new ESXi server - ${user}-${name2}
    Log To Console  \nNimbus ESXi server IP: ${ip2}

    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${ip3}
    ${out}=  Run  govc host.account.update -id root -password e2eFunctionalTest
    Should Be Empty  ${out}
    Disable TLS On ESX Host
    Log To Console  Successfully deployed new ESXi server - ${user}-${name3}
    Log To Console  \nNimbus ESXi server IP: ${ip3}

    Close connection
    [Return]  ${user}-${name1}  ${ip1}  ${user}-${name2}  ${ip2}  ${user}-${name3}  ${ip3}

Deploy Nimbus vCenter Server
    [Arguments]  ${user}  ${password}  ${version}=${VC_VERSION}
    ${name}=  Evaluate  'VC-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nDeploying Nimbus vCenter server: ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}

    ${out}=  Execute Command  nimbus-vcvadeploy --vcvaBuild ${version} ${name}
    # Make sure the deploy actually worked
    Should Contain  ${out}  Overall Status: Succeeded
    # Now grab the IP address and return the name and ip for later use
    @{out}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{out}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  Cloudvm is running on IP
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    ${ip}=  Fetch From Right  ${line}  ${SPACE}

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${ip}
    Log To Console  Successfully deployed new vCenter server - ${user}-${name}
    Close connection
    [Return]  ${user}-${name}  ${ip}

Deploy Nimbus ESXi Server Async
    [Tags]  secret
    [Arguments]  ${name}  ${version}=${ESX_VERSION}
    Log To Console  \nDeploying Nimbus ESXi server: ${name}
    
    ${out}=  Run Secret SSHPASS command  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  'nimbus-esxdeploy ${name} --disk\=48000000 --ssd\=24000000 --memory\=8192 --nics 2 ${version}'
    [Return]  ${out}

Run Secret SSHPASS command
    [Tags]  secret
    [Arguments]  ${user}  ${password}  ${cmd}

    ${out}=  Start Process  sshpass -p ${password} ssh -o StrictHostKeyChecking\=no ${user}@%{NIMBUS_GW} ${cmd}  shell=True
    [Return]  ${out}

Deploy Nimbus vCenter Server Async
    [Tags]  secret
    [Arguments]  ${name}  ${version}=${VC_VERSION}
    Log To Console  \nDeploying Nimbus VC server: ${name}

    ${out}=  Run Secret SSHPASS command  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  'nimbus-vcvadeploy --vcvaBuild ${version} ${name}'
    [Return]  ${out}

Deploy Nimbus Testbed
    [Arguments]  ${user}  ${password}  ${testbed}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}
    ${out}=  Execute Command  nimbus-testbeddeploy ${testbed}
    [Return]  ${out}

Kill Nimbus Server
    [Arguments]  ${user}  ${password}  ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}
    ${out}=  Execute Command  nimbus-ctl kill '${name}'
    Close connection

Nimbus Cleanup
    Run Keyword And Continue On Failure  Gather Logs From Test Server
    Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  *

Gather Host IPs
    ${out}=  Run  govc ls host/cls
    ${out}=  Split To Lines  ${out}
    ${idx}=  Set Variable  1
    :FOR  ${line}  IN  @{out}
    \   Continue For Loop If  '${line}' == '/vcqaDC/host/cls/Resources'
    \   ${ip}=  Fetch From Right  ${line}  /
    \   Set Suite Variable  ${esx${idx}-ip}  ${ip}
    \   ${idx}=  Evaluate  ${idx}+1

vMotion A VM
    [Arguments]  ${vm}
    ${host}=  Get VM Host Name  ${vm}
    ${status}=  Run Keyword And Return Status  Should Contain  ${host}  ${esx1-ip}
    Run Keyword If  ${status}  Run  govc vm.migrate -host cls/${esx2-ip} -pool cls/Resources ${vm}
    Run Keyword Unless  ${status}  Run  govc vm.migrate -host cls/${esx1-ip} -pool cls/Resources ${vm}

Create a VSAN Cluster
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --vcvaBuild ${VC_VERSION} --esxPxeDir ${ESX_VERSION} --esxBuild ${ESX_VERSION} --testbedName vcqa-vsan-simple-pxeBoot-vcva --runName vic-vmotion
    ${out}=  Split To Lines  ${out}
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  .vcva-${VC_VERSION}' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Suite Variable  ${vc-ip}  ${ip}
    \   Exit For Loop If  ${status}

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Log To Console  Create a distributed switch
    ${out}=  Run  govc dvs.create -dc=vcqaDC test-ds
    Should Contain  ${out}  OK

    Log To Console  Create three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=vcqaDC -dvs=test-ds bridge
    Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    ${out}=  Run  govc dvs.add -dvs=test-ds -pnic=vmnic1 /vcqaDC/host/cls
    Should Contain  ${out}  OK

    Log To Console  Enable DRS and VSAN on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  vsanDatastore
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    Gather Host IPs

Create a Simple VC Cluster
    [Arguments]  ${datacenter}=ha-datacenter  ${cluster}=cls
    Log To Console  \nStarting simple VC cluster deploy...
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX1}  ${esx1}
    ${esx2}  ${esx2-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX2}  ${esx2}
    ${esx3}  ${esx3-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX3}  ${esx3}

    ${vc}  ${vc-ip}=  Deploy Nimbus vCenter Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${VC}  ${vc}

    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ${datacenter}
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create ${cluster}
    Should Be Empty  ${out}

    Log To Console  Add ESX host to the VC
    ${out}=  Run  govc cluster.add -hostname=${esx1-ip} -username=root -dc=${datacenter} -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx2-ip} -username=root -dc=${datacenter} -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx3-ip} -username=root -dc=${datacenter} -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create a distributed switch
    ${out}=  Run  govc dvs.create -dc=${datacenter} test-ds
    Should Contain  ${out}  OK

    Log To Console  Create three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=test-ds management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=test-ds vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=test-ds bridge
    Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    ${out}=  Run  govc dvs.add -dvs=test-ds -pnic=vmnic1 /${datacenter}/host/${cluster}
    Should Contain  ${out}  OK

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /${datacenter}/host/${cluster}
    Should Be Empty  ${out}

    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_RESOURCE  ${cluster}
    Set Environment Variable  TEST_TIMEOUT  30m
    [Return]  ${vc-ip}

Create A Distributed Switch
    [Arguments]  ${datacenter}  ${dvs}=test-ds
    Log To Console  \nCreate a distributed switch
    ${out}=  Run  govc dvs.create -product-version 5.5.0 -dc=${datacenter} ${dvs}
    Should Contain  ${out}  OK

Create Three Distributed Port Groups
    [Arguments]  ${datacenter}  ${dvs}=test-ds
    Log To Console  \nCreate three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=${dvs} management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=${dvs} vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=${datacenter} -dvs=${dvs} bridge
    Should Contain  ${out}  OK

Add Host To Distributed Switch
    [Arguments]  ${host}  ${dvs}=test-ds
    Log To Console  \nAdd host(s) to the distributed switch
    ${out}=  Run  govc dvs.add -dvs=${dvs} -pnic=vmnic1 ${host}
    Should Contain  ${out}  OK

Disable TLS On ESX Host
    Log To Console  \nDisable TLS on the host
    ${ver}=  Get Vsphere Version
    ${out}=  Run Keyword If  '${ver}' != '5.5.0'  Run  govc host.option.set UserVars.ESXiVPsDisabledProtocols sslv3,tlsv1,tlsv1.1
    Run Keyword If  '${ver}' != '5.5.0'  Should Be Empty  ${out}

Get Vsphere Version
    ${out}=  Run  govc about
    ${out}=  Split To Lines  ${out}
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  Version:
    \   Run Keyword And Return If  ${status}  Fetch From Right  ${line}  ${SPACE}
