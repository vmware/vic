*** Settings ***
Library  OperatingSystem
Library  String
Library  Collections
Library  requests
Library  Process
Library  SSHLibrary  1 minute  prompt=bash-4.1$

*** Variables ***
${bin-dir}  /drone/src/github.com/vmware/vic/bin

*** Keywords ***
Install VIC Appliance To Test Server
    [Arguments]  ${certs}=true
    # Let's try to pro-actively clean up any datastore that doesn't currently have a VM associated with it
    ${datastore}=  Run  govc datastore.ls
    ${datastore}=  Split To Lines  ${datastore}
    :FOR  ${item}  IN  @{datastore}
    \   Continue For Loop If  '${item}' == 'VIC'
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Match Regexp  ${item}  \\w\\w*-\\w\\w*-\\w\\w*-\\w\\w*-\\w\\w*
    \   Continue For Loop If  '${status}' == 'PASS'
    \   ${vms}=  Run  govc ls vm
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Not Contain  ${vms}  ${item}
    \   Run Keyword If  '${status}' == 'PASS'  Run  govc datastore.rm ${item}

    # Now attempt to intall VCH
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  DRONE_BUILD_NUMBER
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  DRONE_BUILD_NUMBER  local
    
    ${name}=  Evaluate  'VCH-%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
    Set Suite Variable  ${vch-name}  ${name}
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run VIC Machine Command  ${certs}
    ${line}=  Get Line  ${output}  -2
    ${ret}=  Fetch From Right  ${line}  ] docker
    ${ret}=  Remove String  ${ret}  info
    ${ret}=  Strip String  ${ret}
    Set Suite Variable  ${params}  ${ret}
    Log To Console  Installer completed successfully: ${vch-name}...

    # Required due to #1109
    Sleep  10 seconds
    ${status}=  Get State Of Github Issue  1109
    Run Keyword If  '${status}' == 'closed'  Fail  Util.robot needs to be updated now that Issue #1109 has been resolved

Run VIC Machine Command
    [Tags]  secret
    [Arguments]  ${certs}
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-datastore=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --generate-cert=${certs} --password=%{TEST_PASSWORD} --force=true --bridge-network=network --compute-resource=%{TEST_RESOURCE}
    ${status}  ${message} =  Run Keyword And Ignore Error  Should Contain  ${output}  Installer completed successfully
    Run Keyword If  "${status}" == "FAIL"  Fail  Installing the VIC appliance failed, wrong credentials or network problems?
    [Return]  ${output}
    
Cleanup VIC Appliance On Test Server
    Run Keyword And Ignore Error  Gather Logs From Test Server
    # Let's attempt to cleanup any container related to the VCH appliance first
    ${list}=  Run  govc ls /ha-datacenter/vm
    ${list}=  Split To Lines  ${list}
    :FOR  ${vm}  IN  @{list}
    \   Continue For Loop If  '${vm}'=='/ha-datacenter/vm/${vch-name}'
    \   ${raw}=  Run  govc vm.info -json=true ${vm}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${raw}  ${vch-name}
    \   ${name}=  Run Keyword If  '${status}'=='PASS'  Run  govc vm.info -json\=true ${vm} | jq -r '.VirtualMachines[0].Name'
    \   ${uuid}=  Run Keyword If  '${status}'=='PASS'  Run  govc vm.info -json\=true ${vm} | jq -r '.VirtualMachines[0].Config.Uuid'
    \   Run Keyword If  '${status}'=='PASS'  Run  govc vm.destroy ${name}
    \   Run Keyword If  '${status}'=='PASS'  Run  govc datastore.rm ${uuid}

    # Then we can try to cleanup the VCH itself
    ${uuid}=  Run  govc vm.info -json\=true ${vch-name} | jq -r '.VirtualMachines[0].Config.Uuid'
    ${output}=  Run  govc vm.destroy ${vch-name}
    ${output}=  Run  govc pool.destroy %{GOVC_RESOURCE_POOL}/${vch-name}
    ${output}=  Run  govc datastore.rm ${vch-name}
    ${output}=  Run  rm -f ${vch-name}-*.pem
    ${output}=  Run  govc datastore.rm VIC/${uuid}

Gather Logs From Test Server
    Variable Should Exist  ${params}
    ${params}=  Strip String  ${params}
    ${status}  ${message}=  Run Keyword And Ignore Error  Should Not Contain  ${params}  --tls
    # Non-certificate case
    ${ip}=  Run Keyword If  '${status}'=='PASS'  Split String  ${params}  :
    Run Keyword If  '${status}'=='PASS'  Run  wget http://@{ip}[0]:2378/container-logs.tar.gz -O ${vch-name}-container-logs.tar.gz
    # Certificate case
    ${ip}=  Run Keyword If  '${status}'=='FAIL'  Split String  ${params}  ${SPACE}
    ${ip}=  Run Keyword If  '${status}'=='FAIL'  Split String  @{ip}[1]  :
    Run Keyword If  '${status}'=='FAIL'  Run  wget --no-check-certificate https://@{ip}[0]:2378/container-logs.tar.gz -O ${vch-name}-container-logs.tar.gz

Get State Of Github Issue
    [Arguments]  ${num}
    [Tags]  secret
    ${result}=  Get  https://api.github.com/repos/vmware/vic/issues/${num}?access_token\=%{GITHUB_AUTOMATION_API_KEY}
    Should Be Equal  ${result.status_code}  ${200}
    ${status}=  Get From Dictionary  ${result.json()}  state
    [Return]  ${status}
    
Get State Of Drone Build
    [Arguments]  ${num}
    ${out}=  Run  drone build info vmware/vic ${num}
    ${lines}=  Split To Lines  ${out}
    [Return]  @{lines}[2]
    
Get Image IDs
    [Arguments]  ${dir}
    ${result}=  Run Process  cat manifest.json | jq -r ".history[].v1Compatibility|fromjson.id"  shell=True  cwd=${dir}
    ${ids}=  Split To Lines  ${result.stdout}
    [Return]  ${ids}
    
Get Checksums
    [Arguments]  ${dir}
    ${result}=  Run Process  cat manifest.json | jq -r ".fsLayers[].blobSum"  shell=True  cwd=${dir}
    ${out}=  Split To Lines  ${result.stdout}
    ${checkSums}=  Create List
    :FOR  ${str}  IN  @{out}
    \   ${sha}  ${sum}=  Split String From Right  ${str}  :
    \   Append To List  ${checkSums}  ${sum}
    [Return]  ${checkSums}
    
Verify Checksums
    [Arguments]  ${dir}
    ${ids}=  Get Image IDs  ${dir}
    ${sums}=  Get Checksums  ${dir}
    ${idx}=  Set Variable  0
    :FOR  ${id}  IN  @{ids}
    \   ${imageSum}=  Run Process  sha256sum ${id}/${id}.tar  shell=True  cwd=${dir}
    \   ${imageSum}=  Split String  ${imageSum.stdout}
    \   Should Be Equal  @{sums}[${idx}]  @{imageSum}[0]
    \   ${idx}=  Evaluate  ${idx}+1

Deploy Nimbus ESXi Server
    [Arguments]  ${user}  ${password}
    ${name}=  Evaluate  'ESX-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nDeploying Nimbus ESXi server: ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}

    ${out}=  Execute Command  nimbus-esxdeploy ${name} 3620759
    # Make sure the deploy actually worked
    ${success}=  Get Line  ${out}  -2
    Should Contain  ${success}  To manage this VM use
    # Now grab the IP address and return the name and ip for later use
    ${gotIP}=  Get Line  ${out}  -5
    @{gotIP}=  Split String  ${gotIP}  ${SPACE}
    ${ip}=  Remove String  @{gotIP}[5]  ,
    # Let's set a password so govc doesn't complain
    Remove Environment Variable  GOVC_PASSWORD
    Remove Environment Variable  GOVC_USERNAME
    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_URL  root:@${ip}
    ${out}=  Run  govc host.account.update -id root -password e2eFunctionalTest
    Should Be Empty  ${out}
    Log To Console  Successfully deployed new ESXi server - ${user}-${name}
    Close connection
    [Return]  ${user}-${name}  ${ip}

Deploy Nimbus vCenter Server
    [Arguments]  ${user}  ${password}
    ${name}=  Evaluate  'VC-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nDeploying Nimbus vCenter server: ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}

    ${out}=  Execute Command  nimbus-vcvadeploy --vcvaBuild 3634791 ${name}
    # Make sure the deploy actually worked
    ${success}=  Get Line  ${out}  -5
    Should Contain  ${success}  Overall Status: Succeeded
    # Now grab the IP address and return the name and ip for later use
    ${gotIP}=  Get Line  ${out}  -22
    ${ip}=  Fetch From Right  ${gotIP}  ${SPACE}

    Set Environment Variable  GOVC_INSECURE  1
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin!23
    Set Environment Variable  GOVC_URL  ${ip}
    Log To Console  Successfully deployed new vCenter server - ${user}-${name}
    Close connection
    [Return]  ${user}-${name}  ${ip}

Kill Nimbus Server
    [Arguments]  ${user}  ${password}  ${name}
    Open Connection  %{NIMBUS_GW}
    Login  ${user}  ${password}
    ${out}=  Execute Command  nimbus-ctl kill '${name}'
    Close connection
