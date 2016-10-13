*** Settings ***
Library  OperatingSystem
Library  String
Library  Collections
Library  requests
Library  Process
Library  SSHLibrary  1 minute  prompt=bash-4.1$
Library  DateTime

*** Variables ***
${bin-dir}  ${CURDIR}/../../bin

*** Keywords ***
Set Test Environment Variables
    # Finish setting up environment variables
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  DRONE_BUILD_NUMBER
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  DRONE_BUILD_NUMBER  0
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  BRIDGE_NETWORK
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  BRIDGE_NETWORK  network
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  EXTERNAL_NETWORK
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  EXTERNAL_NETWORK  'VM Network'

    @{URLs}=  Split String  %{TEST_URL_ARRAY}
    ${len}=  Get Length  ${URLs}
    ${IDX}=  Evaluate  %{DRONE_BUILD_NUMBER} \% ${len}

    Set Environment Variable  TEST_URL  @{URLs}[${IDX}]
    Log To Console  %{TEST_URL}
    Set Environment Variable  GOVC_URL  %{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}

    ${host}=  Run  govc ls host
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  TEST_RESOURCE
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  TEST_RESOURCE  ${host}/Resources
    Set Environment Variable  GOVC_RESOURCE_POOL  %{TEST_RESOURCE}
    Set Environment Variable  GOVC_DATASTORE  %{TEST_DATASTORE}

    ${about}=  Run  govc about
    ${status}=  Run Keyword And Return Status  Should Contain  ${about}  VMware ESXi
    Run Keyword If  ${status}  Set Environment Variable  HOST_TYPE  ESXi
    Run Keyword Unless  ${status}  Set Environment Variable  HOST_TYPE  VC

    # set the TLS config options suitable for vic-machine in this env
    ${domain}=  Get Environment Variable  DOMAIN  ''
    Run Keyword If  '${domain}' == ''  Set Suite Variable  ${vicmachinetls}  '--no-tlsverify'
    Run Keyword If  '${domain}' != ''  Set Suite Variable  ${vicmachinetls}  '--tls-cname=*.${domain}'

Set Test VCH Name
    ${name}=  Evaluate  'VCH-%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
    Set Suite Variable  ${vch-name}  ${name}

Get Docker Params
    # Get VCH docker params e.g. "-H 192.168.218.181:2376 --tls"
    [Arguments]  ${output}  ${certs}
    @{output}=  Split To Lines  ${output}
    :FOR  ${item}  IN  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  DOCKER_HOST=
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}

    # Ensure we start from a clean slate with docker env vars
    Remove Environment Variable  DOCKER_HOST  DOCKER_TLS_VERIFY  DOCKER_CERT_PATH

    # Split the log log into pieces, discarding the initial log decoration, and assign to env vars
    ${logdeco}  ${vars}=  Split String  ${line}  ${SPACE}  1
    ${vars}=  Split String  ${vars}
    :FOR  ${var}  IN  @{vars}
    \   ${varname}  ${varval}=  Split String  ${var}  =
    \   Set Environment Variable  ${varname}  ${varval}

    ${dockerHost}=  Get Environment Variable  DOCKER_HOST

    @{hostParts}=  Split String  ${dockerHost}  :
    ${ip}=  Strip String  @{hostParts}[0]
    ${port}=  Strip String  @{hostParts}[1]
    Set Suite Variable  ${vch-ip}  ${ip}
    Set Suite Variable  ${vch-port}  ${port}

    ${proto}=  Set Variable If  ${port} == 2376  "https"  "http"
    Set Suite Variable  ${proto}

    Run Keyword If  ${port} == 2376  Set Suite Variable  ${params}  -H ${dockerHost} --tls
    Run Keyword If  ${port} == 2375  Set Suite Variable  ${params}  -H ${dockerHost}


    :FOR  ${item}  IN  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  http
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}

    ${rest}  ${vic-admin}=  Split String From Right  ${line}
    Set Suite Variable  ${vic-admin}

Install VIC Appliance To Test Server
    [Arguments]  ${vic-machine}=bin/vic-machine-linux  ${appliance-iso}=bin/appliance.iso  ${bootstrap-iso}=bin/bootstrap.iso  ${certs}=${true}  ${vol}=default
    Set Test Environment Variables
    # disable firewall
    Run  govc host.esxcli network firewall set -e false
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Run Keyword And Ignore Error  Cleanup Dangling Networks On Test Server
    Run Keyword And Ignore Error  Cleanup Dangling vSwitches On Test Server
    Set Test VCH Name
    # Set a unique bridge network for each VCH that has a random VLAN ID
    ${vlan}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Evaluate  str(random.randint(1, 4093))  modules=random
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vlan=${vlan} -vswitch vSwitch0 ${vch-name}-bridge
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Environment Variable  BRIDGE_NETWORK  ${vch-name}-bridge

    # Install the VCH now
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run VIC Machine Command  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    Log  ${output}
    Get Docker Params  ${output}  ${certs}
    Log To Console  Installer completed successfully: ${vch-name}...

Run VIC Machine Command
    [Tags]  secret
    [Arguments]  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    ${output}=  Run Keyword If  ${certs}  Run  ${vic-machine} create --debug 1 --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} ${vicmachinetls} 
    Run Keyword If  ${certs}  Should Contain  ${output}  Installer completed successfully
    Return From Keyword If  ${certs}  ${output}

    ${output}=  Run Keyword Unless  ${certs}  Run  ${vic-machine} create --debug 1 --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} --no-tls
    Run Keyword Unless  ${certs}  Should Contain  ${output}  Installer completed successfully
    [Return]  ${output}

Cleanup VIC Appliance On Test Server
    Log To Console  Gathering logs from the test server...
    Gather Logs From Test Server
    Log To Console  Deleting the VCH appliance...
    ${output}=  Run VIC Machine Delete Command
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove ${vch-name}-bridge
    [Return]  ${output}

Check Delete Success
    [Arguments]  ${vch-name}
    ${out}=  Run  govc ls vm
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}
    ${out}=  Run  govc datastore.ls
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}
    ${out}=  Run  govc ls host/*/Resources/*
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}

Run Secret VIC Machine Delete Command
    [Tags]  secret
    [Arguments]  ${vch-name}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux delete --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    [Return]  ${rc}  ${output}

Run VIC Machine Delete Command
    ${rc}  ${output}=  Run Secret VIC Machine Delete Command  ${vch-name}
    Check Delete Success  ${vch-name}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  rm -f ${vch-name}-*.pem
    [Return]  ${output}

Cleanup Datastore On Test Server
    ${out}=  Run  govc datastore.ls
    ${lines}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{lines}
    \   Continue For Loop If  '${item}' == 'VIC'
    \   ${contents}=  Run  govc datastore.ls ${item}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${contents}  vmx
    \   Continue For Loop If  ${status}
    \   ${out}=  Run  govc datastore.rm ${item}

Cleanup Dangling VMs On Test Server
    ${out}=  Run  govc ls vm
    ${vms}=  Split To Lines  ${out}
    :FOR  ${vm}  IN  @{vms}
    \   ${vm}=  Fetch From Right  ${vm}  /
    \   ${build}=  Split String  ${vm}  -
    \   # Skip any VM that is not associated with integration tests
    \   Continue For Loop If  '@{build}[0]' != 'VCH'
    \   # Skip any VM that is still running
    \   ${state}=  Get State Of Drone Build  @{build}[1]
    \   Continue For Loop If  '${state}' == 'running'
    \   ${uuid}=  Run  govc vm.info -json\=true ${vm} | jq -r '.VirtualMachines[0].Config.Uuid'
    \   Log To Console  Destroying dangling VCH: ${vm}
    \   ${rc}  ${output}=  Run Secret VIC Machine Delete Command  ${vm}

Cleanup Dangling Networks On Test Server
    ${out}=  Run  govc ls network
    ${nets}=  Split To Lines  ${out}
    :FOR  ${net}  IN  @{nets}
    \   ${net}=  Fetch From Right  ${net}  /
    \   ${build}=  Split String  ${net}  -
    \   # Skip any Network that is not associated with integration tests
    \   Continue For Loop If  '@{build}[0]' != 'VCH'
    \   # Skip any Network that is still running
    \   ${state}=  Get State Of Drone Build  @{build}[1]
    \   Continue For Loop If  '${state}' == 'running'
    \   ${uuid}=  Run  govc host.portgroup.remove ${net}
    
Cleanup Dangling vSwitches On Test Server
    ${out}=  Run  govc host.vswitch.info | grep VCH
    ${nets}=  Split To Lines  ${out}
    :FOR  ${net}  IN  @{nets}
    \   ${net}=  Fetch From Right  ${net}  ${SPACE}
    \   ${build}=  Split String  ${net}  -
    \   # Skip any vSwitch that is not associated with integration tests
    \   Continue For Loop If  '@{build}[0]' != 'VCH'
    \   # Skip any vSwitch that is still running
    \   ${state}=  Get State Of Drone Build  @{build}[1]
    \   Continue For Loop If  '${state}' == 'running'
    \   ${uuid}=  Run  govc host.vswitch.remove ${net}

Gather Logs From Test Server
    Variable Should Exist  ${params}
    ${params}=  Strip String  ${params}
    ${status}  ${message}=  Run Keyword And Ignore Error  Should Not Contain  ${params}  --tls
    # Non-certificate case
    ${ip}=  Run Keyword If  '${status}'=='PASS'  Split String  ${params}  :
    Run Keyword If  '${status}'=='PASS'  Run  wget ${vic-admin}/container-logs.zip -O ${SUITE NAME}-${vch-name}-container-logs.zip
    # Certificate case
    ${ip}=  Run Keyword If  '${status}'=='FAIL'  Split String  ${params}  ${SPACE}
    ${ip}=  Run Keyword If  '${status}'=='FAIL'  Split String  @{ip}[1]  :
    Run Keyword If  '${status}'=='FAIL'  Run  wget --no-check-certificate ${vic-admin}/container-logs.zip -O ${vch-name}-container-logs.zip

Gather Logs From ESX Server
    Environment Variable Should Be Set  TEST_URL
    ${out}=  Run  govc logs.download

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

Get State Of Drone Build
    [Arguments]  ${num}
    ${out}=  Run  drone build info vmware/vic ${num}
    ${lines}=  Split To Lines  ${out}
    [Return]  @{lines}[2]

Get Title of Drone Build
    [Arguments]  ${num}
    ${out}=  Run  drone build info vmware/vic ${num}
    ${lines}=  Split To Lines  ${out}
    [Return]  @{lines}[-1]

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
    [Arguments]  ${user}  ${password}  ${version}=3620759
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
    Log To Console  Successfully deployed new ESXi server - ${user}-${name}
    Close connection
    [Return]  ${user}-${name}  ${ip}

Deploy Nimbus vCenter Server
    [Arguments]  ${user}  ${password}  ${version}=3634791
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
    Gather Logs From Test Server
    Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  *

Wait Until Container Stops
    [Arguments]  ${container}
    :FOR  ${idx}  IN RANGE  0  30
    \   ${out}=  Run  docker ${params} ps --filter status=running --no-trunc
    \   ${status}=  Run Keyword And Return Status  Should Not Contain  ${out}  ${container}
    \   Return From Keyword If  ${status}
    \   Sleep  1
    Fail  Container did not stop within 30 seconds

Wait Until VM Powers Off
    [Arguments]  ${vm}
    :FOR  ${idx}  IN RANGE  0  30
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc vm.info ${vch-name}/${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'VC'  Set Test Variable  ${out}  ${ret}
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc vm.info ${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Test Variable  ${out}  ${ret}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  poweredOff
    \   Return From Keyword If  ${status}
    \   Sleep  1
    Fail  VM did not power off within 30 seconds

Wait Until VM Is Destroyed
    [Arguments]  ${vm}
    :FOR  ${idx}  IN RANGE  0  30
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc ls vm/${vch-name}/${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'VC'  Set Test Variable  ${out}  ${ret}
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc ls vm/${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Test Variable  ${out}  ${ret}
    \   ${status}=  Run Keyword And Return Status  Should Be Empty  ${out}
    \   Return From Keyword If  ${status}
    \   Sleep  1
    Fail  VM was not destroyed within 30 seconds

Wait Until VM Powers On
    [Arguments]  ${vm}
    :FOR  ${idx}  IN RANGE  0  30
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  govc vm.info ${vch-name}/${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'VC'  Set Test Variable  ${out}  ${ret}
    \   ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc vm.info ${vm}
    \   Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Test Variable  ${out}  ${ret}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  poweredOn
    \   Return From Keyword If  ${status}
    \   Sleep  1
    Fail  VM did not power on within 30 seconds

Get VM IP
    [Arguments]  ${vm}
    ${rc}  ${out}=  Run And Return Rc And Output  govc vm.ip ${vm}
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${out}

Get VM Name
    [Arguments]  ${vm}
    ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Set Variable  ${vm}/${vm}
    ${ret}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Variable  ${vm}
    [Return]  ${ret}

    
Get VM Host Name
    [Arguments]  ${vm}
    ${vm}=  Get VM Name  ${vm}
    ${ret}=  Run  govc vm.info ${vm}/${vm}
    Set Test Variable  ${out}  ${ret}
    ${out}=  Split To Lines  ${out}
    ${host}=  Fetch From Right  @{out}[-1]  ${SPACE} 
    [Return]  ${host}

Run Unit Tests
    [Tags]  secret
    Set Environment Variable  VIC_ESX_TEST_URL  %{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}
    Log To Console  \nls vendor/github.com/vmware/govmomi/vim25/methods:
    ${output}=  Run  ls vendor/github.com/vmware/govmomi/vim25/methods
    Log To Console  ${output}
    Log To Console  Execute the unit tests...
    ${output}=  Run  make -j3 test
    Log To Console  ${output}
    Should Not Contain  ${output}  FAIL
    Should Not Contain  ${output}  [build failed]

Run Regression Tests
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  busybox
    ${rc}  ${container}=  Run And Return Rc And Output  docker ${params} create busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  /bin/top
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Container Stops  ${container}
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Exited
    # Ensure container logs are correctly being gathered for debugging purposes
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk ${vic-admin}/container-logs.tar.gz | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/output.log
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm ${container}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  /bin/top

    # Check for regression for #1265
    ${rc}  ${container1}=  Run And Return Rc And Output  docker ${params} create -it busybox /bin/top
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container2}=  Run And Return Rc And Output  docker ${params} create -it busybox
    Should Be Equal As Integers  ${rc}  0
    ${shortname}=  Get Substring  ${container2}  1  12
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    ${lines}=  Get Lines Containing String  ${output}  ${shortname}
    Should Not Contain  ${lines}  /bin/top
    ${rc}=  Run And Return Rc  docker ${params} rm ${container1}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker ${params} rm ${container2}
    Should Be Equal As Integers  ${rc}  0

    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  busybox

Put Host Into Maintenance Mode
    ${rc}  ${output}=  Run And Return Rc And Output  govc host.maintenance.enter -host.ip=%{TEST_URL}
    Should Contain  ${output}  entering maintenance mode... OK

Remove Host From Maintenance Mode
    ${rc}  ${output}=  Run And Return Rc And Output  govc host.maintenance.exit -host.ip=%{TEST_URL}
    Should Contain  ${output}  exiting maintenance mode... OK

Hit Nginx Endpoint
    [Arguments]  ${vch-ip}  ${port}
    ${rc}  ${output}=  Run And Return Rc And Output  wget ${vch-ip}:${port}
    Should Be Equal As Integers  ${rc}  0

Run Docker Info
    [Arguments]  ${docker-params}
    ${rc}=  Run And Return Rc  docker ${docker-params} info
    Should Be Equal As Integers  ${rc}  0
