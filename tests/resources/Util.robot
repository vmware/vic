*** Settings ***
Library  OperatingSystem
Library  String
Library  Collections
Library  requests
Library  Process

*** Variables ***
${bin-dir}  /drone/src/github.com/vmware/vic/bin

*** Keywords ***
Install VIC Appliance To Test Server
    [Tags]  secret
    [Arguments]  ${certs}=true
    ${name}=  Evaluate  'VCH-' + str(random.randint(1000,9999))  modules=random
    Log To Console  \nInstalling VCH to test server...
    Set Suite Variable  ${vch-name}  ${name}
    ${output}=  Run  bin/vic-machine-linux create --name ${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-datastore=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --generate-cert=${certs} --password=%{TEST_PASSWORD} --force=true --bridge-network=network --compute-resource=%{TEST_RESOURCE}
    ${status}  ${message} =  Run Keyword And Ignore Error  Should Contain  ${output}  Installer completed successfully...
    Run Keyword If  "${status}" == "FAIL"  Fail  Installing the VIC appliance failed, wrong credentials or network problems?
    ${line}=  Get Line  ${output}  -2
    # Parse output when we are not using TLS
    ${ret}=  Run Keyword If  "${certs}" == "false"  Fetch From Right  ${line}  ] DOCKER_HOST=
    Run Keyword If  "${certs}" == "false"  Set Suite Variable  ${params}  ${ret}
    # Parse output when we are using TLS
    ${ret}=  Run Keyword If  "${certs}" == "true"  Fetch From Right  ${line}  ] docker
    ${ret}=  Run Keyword If  "${certs}" == "true"  Remove String  ${ret}  info
    Run Keyword If  "${certs}" == "true"  Set Suite Variable  ${params}  ${ret}
    Log To Console  Installer completed successfully: ${vch-name}...

Cleanup VIC Appliance On Test Server
    ${uuid}=  Run  govc vm.info -json\=true ${vch-name} | jq -r '.VirtualMachines[0].Config.Uuid'

    ${output}=  Run  govc vm.destroy ${vch-name}
    ${output}=  Run  govc pool.destroy %{GOVC_RESOURCE_POOL}/${vch-name}
    ${output}=  Run  govc datastore.rm ${vch-name}
    ${output}=  Run  rm -f ${vch-name}-*.pem
    ${output}=  Run  govc datastore.rm VIC/${uuid}

Get State Of Github Issue
    [Arguments]  ${num}
    ${result} =  get  https://api.github.com/repos/vmware/vic/issues/${num}
    Should Be Equal  ${result.status_code}  ${200}
    ${status} =  Get From Dictionary  ${result.json()}  state
    [Return]  ${status}
    
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
