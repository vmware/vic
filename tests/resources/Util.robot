*** Settings ***
Library  OperatingSystem
Library  String
Library  Collections
Library  requests
Library  Process
Library  SSHLibrary  1 minute  prompt=bash-4.1$
Library  DateTime
Resource  Nimbus-Util.robot

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
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  TEST_DATACENTER
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  TEST_DATACENTER  ${SPACE}

    @{URLs}=  Split String  %{TEST_URL_ARRAY}
    ${len}=  Get Length  ${URLs}
    ${IDX}=  Evaluate  %{DRONE_BUILD_NUMBER} \% ${len}

    Set Environment Variable  TEST_URL  @{URLs}[${IDX}]
    Set Environment Variable  GOVC_URL  %{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}
    # TODO: need an integration/vic-test image update to include the about.cert command
    #${rc}  ${thumbprint}=  Run And Return Rc And Output  govc about.cert -k | jq -r .ThumbprintSHA1
    ${rc}  ${thumbprint}=  Run And Return Rc And Output  openssl s_client -connect $(govc env -x GOVC_URL_HOST):443 </dev/null 2>/dev/null | openssl x509 -fingerprint -noout | cut -d= -f2
    Should Be Equal As Integers  ${rc}  0
    Set Environment Variable  TEST_THUMBPRINT  ${thumbprint}
    Log To Console  \nTEST_URL=%{TEST_URL}

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
    # Set a unique bridge network for each VCH that has a random VLAN ID
    ${vlan}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Evaluate  str(random.randint(1, 4093))  modules=random
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vlan=${vlan} -vswitch vSwitch0 ${vch-name}-bridge
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Environment Variable  BRIDGE_NETWORK  ${vch-name}-bridge

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


    :FOR  ${index}  ${item}  IN ENUMERATE  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  http
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  Published ports can be reached at
    \   ${idx} =  Evaluate  ${index} + 1
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${ext-ip}  @{output}[${idx}]

    ${rest}  ${ext-ip} =  Split String From Right  ${ext-ip}
    ${ext-ip} =  Strip String  ${ext-ip}
    Set Suite Variable  ${ext-ip}  ${ext-ip}

    ${rest}  ${vic-admin}=  Split String From Right  ${line}
    Set Suite Variable  ${vic-admin}

Install VIC Appliance To Test Server
    [Arguments]  ${vic-machine}=bin/vic-machine-linux  ${appliance-iso}=bin/appliance.iso  ${bootstrap-iso}=bin/bootstrap.iso  ${certs}=${true}  ${vol}=default
    Set Test Environment Variables
    # disable firewall
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.esxcli network firewall set -e false
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Run Keyword And Ignore Error  Cleanup Dangling Networks On Test Server
    Run Keyword And Ignore Error  Cleanup Dangling vSwitches On Test Server

    # Install the VCH now
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run VIC Machine Command  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    Log  ${output}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${certs}
    Log To Console  Installer completed successfully: ${vch-name}

Run VIC Machine Command
    [Tags]  secret
    [Arguments]  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    ${output}=  Run Keyword If  ${certs}  Run  ${vic-machine} create --debug 1 --name=${vch-name} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} ${vicmachinetls}
    Return From Keyword If  ${certs}  ${output}

    ${output}=  Run Keyword Unless  ${certs}  Run  ${vic-machine} create --debug 1 --name=${vch-name} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} --no-tlsverify
    [Return]  ${output}

Cleanup VIC Appliance On Test Server
    Log To Console  Gathering logs from the test server ${vch-name}
    Gather Logs From Test Server
    Log To Console  Deleting the VCH appliance ${vch-name}
    ${output}=  Run VIC Machine Delete Command
    Run Keyword And Ignore Error  Cleanup VCH Bridge Network  ${vch-name}
    [Return]  ${output}

Cleanup VCH Bridge Network
    [Arguments]  ${vch-name}
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove ${vch-name}-bridge
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.info
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Not Contain  ${out}  ${vch-name}-bridge

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
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux delete --name=${vch-name} --target=%{TEST_URL}%{TEST_DATACENTER} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    [Return]  ${rc}  ${output}

Run VIC Machine Delete Command
    ${rc}  ${output}=  Run Secret VIC Machine Delete Command  ${vch-name}
    Wait Until Keyword Succeeds  6x  5s  Check Delete Success  ${vch-name}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully
    [Return]  ${output}

Cleanup Datastore On Test Server
    ${out}=  Run  govc datastore.ls
    ${items}=  Split To Lines  ${out}
    :FOR  ${item}  IN  @{items}
    \   ${build}=  Split String  ${item}  -
    \   # Skip any item that is not associated with integration tests
    \   Continue For Loop If  '@{build}[0]' != 'VCH'
    \   # Skip any item that is still running
    \   ${state}=  Get State Of Drone Build  @{build}[1]
    \   Continue For Loop If  '${state}' == 'running'
    \   Log To Console  Removing the following item from datastore: ${item}
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
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.vswitch.info | grep VCH
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
    [Tags]  secret
    ${out}=  Run  curl -k -D vic-admin-cookies -Fusername=%{TEST_USERNAME} -Fpassword=%{TEST_PASSWORD} ${vic-admin}/authentication
    Log  ${out}
    ${out}=  Run  curl -k -b vic-admin-cookies ${vic-admin}/container-logs.zip -o ${SUITE NAME}-${vch-name}-container-logs.zip
    Log  ${out}
    Remove File  vic-admin-cookies

Gather Logs From ESX Server
    Environment Variable Should Be Set  TEST_URL
    ${out}=  Run  govc logs.download

Change Log Level On Server
    [Arguments]  ${level}
    ${out}=  Run  govc host.option.set Config.HostAgent.log.level ${level}
    Should Be Empty  ${out}

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

Wait Until Container Stops
    [Arguments]  ${container}
    :FOR  ${idx}  IN RANGE  0  30
    \   ${out}=  Run  docker ${params} inspect ${container} | grep Status
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  exited
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

Get VM Info
    [Arguments]  ${vm}
    ${rc}  ${out}=  Run And Return Rc And Output  govc vm.info -r ${vm}
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${out}

Run Regression Tests
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} pull busybox
    Should Be Equal As Integers  ${rc}  0
    # Pull an image that has been pulled already
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
    # get docker_cert_path or empty string if it's unset
    ${docker_cert_path}=  Get Environment Variable  DOCKER_CERT_PATH  ${EMPTY}
    # Ensure container logs are correctly being gathered for debugging purposes
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk ${vic-admin}/authentication -XPOST -F username=%{TEST_USERNAME} -F password=%{TEST_PASSWORD} -D /tmp/cookies-${vch-name}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk ${vic-admin}/container-logs.tar.gz -b /tmp/cookies-${vch-name} | tar tvzf -
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

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rmi busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} images
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  busybox

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

Power Off VM OOB
    [Arguments]  ${vm}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -off ${vch-name}/"${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -off "${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Log To Console  Waiting for VM to power off ...
    Wait Until VM Powers Off  "${vm}"

Power On VM OOB
    [Arguments]  ${vm}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -on ${vch-name}/"${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -on "${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Log To Console  Waiting for VM to power on ...
    Wait Until VM Powers On  ${vm}

Destroy VM OOB
    [Arguments]  ${vm}
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.destroy ${vch-name}/"*-${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.destroy "*-${vm}"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0

Get Datacenter Name
    ${out}=  Run  govc datacenter.info
    ${out}=  Split To Lines  ${out}
    ${name}=  Fetch From Right  @{out}[0]  ${SPACE}
    [Return]  ${name}