*** Settings ***
Documentation  This resource contains all keywords related to creating, deleting, maintaining VCHs

*** Keywords ***
Set Test Environment Variables
    # Finish setting up environment variables
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  DRONE_BUILD_NUMBER
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  DRONE_BUILD_NUMBER  0
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  BRIDGE_NETWORK
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  BRIDGE_NETWORK  network
    ${status}  ${message}=  Run Keyword And Ignore Error  Environment Variable Should Be Set  PUBLIC_NETWORK
    Run Keyword If  '${status}' == 'FAIL'  Set Environment Variable  PUBLIC_NETWORK  'VM Network'
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
    Run Keyword If  '${domain}' == ''  Set Suite Variable  ${vicmachinetls}  --no-tlsverify
    Run Keyword If  '${domain}' != ''  Set Suite Variable  ${vicmachinetls}  --tls-cname=*.${domain}

    Set Test VCH Name
    # Set a unique bridge network for each VCH that has a random VLAN ID
    ${vlan}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Evaluate  str(random.randint(1, 4093))  modules=random
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.add -vlan=${vlan} -vswitch vSwitch0 %{VCH-NAME}-bridge
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Set Environment Variable  BRIDGE_NETWORK  %{VCH-NAME}-bridge

Set Test VCH Name
    ${name}=  Evaluate  'VCH-%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
    Set Environment Variable  VCH-NAME  ${name}

Set List Of Env Variables
    [Arguments]  ${vars}
    @{vars}=  Split String  ${vars}
    :FOR  ${var}  IN  @{vars}
    \   ${varname}  ${varval}=  Split String  ${var}  =
    \   Set Environment Variable  ${varname}  ${varval}

Parse Environment Variables
    [Arguments]  ${line}
    #  If using the old logging format
    ${status}=  Run Keyword And Return Status  Should Contain  ${line}  mINFO
    ${logdeco}  ${vars}=  Run Keyword If  ${status}  Split String  ${line}  ${SPACE}  1
    Run Keyword If  ${status}  Set List Of Env Variables  ${vars}
    Return From Keyword If  ${status}
    
    # Split the log log into pieces, discarding the initial log decoration, and assign to env vars
    ${logmon}  ${logday}  ${logyear}  ${logtime}  ${loglevel}  ${vars}=  Split String  ${line}  max_split=5
    Set List Of Env Variables  ${vars}

Get Docker Params
    # Get VCH docker params e.g. "-H 192.168.218.181:2376 --tls"
    [Arguments]  ${output}  ${certs}
    @{output}=  Split To Lines  ${output}
    :FOR  ${item}  IN  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  DOCKER_HOST=
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}

    # Ensure we start from a clean slate with docker env vars
    Remove Environment Variable  DOCKER_HOST  DOCKER_TLS_VERIFY  DOCKER_CERT_PATH

    Parse Environment Variables  ${line}

    ${dockerHost}=  Get Environment Variable  DOCKER_HOST

    @{hostParts}=  Split String  ${dockerHost}  :
    ${ip}=  Strip String  @{hostParts}[0]
    ${port}=  Strip String  @{hostParts}[1]
    Set Environment Variable  VCH-IP  ${ip}
    Set Environment Variable  VCH-PORT  ${port}

    :FOR  ${index}  ${item}  IN ENUMERATE  @{output}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  http
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${line}  ${item}
    \   ${status}  ${message}=  Run Keyword And Ignore Error  Should Contain  ${item}  Published ports can be reached at
    \   ${idx} =  Evaluate  ${index} + 1
    \   Run Keyword If  '${status}' == 'PASS'  Set Suite Variable  ${ext-ip}  @{output}[${idx}]

    ${rest}  ${ext-ip} =  Split String From Right  ${ext-ip}  ${SPACE}  1
    ${ext-ip} =  Strip String  ${ext-ip}
    Set Environment Variable  EXT-IP  ${ext-ip}

    ${rest}  ${vic-admin}=  Split String From Right  ${line}  ${SPACE}  1
    Set Environment Variable  VIC-ADMIN  ${vic-admin}

    Run Keyword If  ${port} == 2376  Set Environment Variable  VCH-PARAMS  -H ${dockerHost} --tls
    Run Keyword If  ${port} == 2375  Set Environment Variable  VCH-PARAMS  -H ${dockerHost}

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
    Run Keyword And Ignore Error  Cleanup Dangling Containers On Test Server

    # Install the VCH now
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run VIC Machine Command  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    Log  ${output}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${certs}
    Log To Console  Installer completed successfully: %{VCH-NAME}...

Run VIC Machine Command
    [Tags]  secret
    [Arguments]  ${vic-machine}  ${appliance-iso}  ${bootstrap-iso}  ${certs}  ${vol}
    ${output}=  Run Keyword If  ${certs}  Run  ${vic-machine} create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} ${vicmachinetls}
    Run Keyword If  ${certs}  Should Contain  ${output}  Installer completed successfully
    Return From Keyword If  ${certs}  ${output}

    ${output}=  Run Keyword Unless  ${certs}  Run  ${vic-machine} create --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=${appliance-iso} --bootstrap-iso=${bootstrap-iso} --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:${vol} --no-tlsverify
    Run Keyword Unless  ${certs}  Should Contain  ${output}  Installer completed successfully
    [Return]  ${output}

Run Secret VIC Machine Delete Command
    [Tags]  secret
    [Arguments]  ${vch-name}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux delete --name=${vch-name} --target=%{TEST_URL}%{TEST_DATACENTER} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    [Return]  ${rc}  ${output}

Run Secret VIC Machine Inspect Command
    [Tags]  secret
    [Arguments]  ${name}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=${name} --target=%{TEST_URL}%{TEST_DATACENTER} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --thumbprint=%{TEST_THUMBPRINT}
    [Return]  ${rc}  ${output}

Run VIC Machine Delete Command
    ${rc}  ${output}=  Run Secret VIC Machine Delete Command  %{VCH-NAME}
    Wait Until Keyword Succeeds  6x  5s  Check Delete Success  %{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  rm -f %{VCH-NAME}-*.pem
    [Return]  ${output}

Run VIC Machine Inspect Command
    ${rc}  ${output}=  Run Secret VIC Machine Inspect Command  %{VCH-NAME}
    Get Docker Params  ${output}  ${true}

Gather Logs From Test Server
    [Tags]  secret
    ${out}=  Run  curl -k -D vic-admin-cookies -Fusername=%{TEST_USERNAME} -Fpassword=%{TEST_PASSWORD} %{VIC-ADMIN}/authentication
    Log  ${out}
    ${out}=  Run  curl -k -b vic-admin-cookies %{VIC-ADMIN}/container-logs.zip -o ${SUITE NAME}-%{VCH-NAME}-container-logs.zip
    Log  ${out}
    Remove File  vic-admin-cookies
    ${out}=  Run  govc datastore.download %{VCH-NAME}/vmware.log %{VCH-NAME}-vmware.log
    Should Contain  ${out}  OK
    ${out}=  Run  sshpass -p %{TEST_PASSWORD} scp %{TEST_USERNAME}@%{TEST_URL}:/var/log/vmkernel.log ./vmkernel.log

Check For The Proper Log Files
    [Arguments]  ${container}
    # Ensure container logs are correctly being gathered for debugging purposes
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk %{VIC-ADMIN}/authentication -XPOST -F username=%{TEST_USERNAME} -F password=%{TEST_PASSWORD} -D /tmp/cookies-%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc and Output  curl -sk %{VIC-ADMIN}/container-logs.tar.gz -b /tmp/cookies-%{VCH-NAME} | tar tvzf -
    Should Be Equal As Integers  ${rc}  0
    Log  ${output}
    Should Contain  ${output}  ${container}/output.log
    Should Contain  ${output}  ${container}/vmware.log
    Should Contain  ${output}  ${container}/tether.debug

Scrape Logs For the Password
    [Tags]  secret
    ${rc}=  Run And Return Rc  curl -sk %{VIC-ADMIN}/authentication -XPOST -F username=%{TEST_USERNAME} -F password=%{TEST_PASSWORD} -D /tmp/cookies-%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

    ${rc}=  Run And Return Rc  curl -sk %{VIC-ADMIN}/logs/port-layer.log -b /tmp/cookies-%{VCH-NAME} | grep -q "%{TEST_PASSWORD}"
    Should Be Equal As Integers  ${rc}  1
    ${rc}=  Run And Return Rc  curl -sk %{VIC-ADMIN}/logs/init.log -b /tmp/cookies-%{VCH-NAME} | grep -q "%{TEST_PASSWORD}"
    Should Be Equal As Integers  ${rc}  1
    ${rc}=  Run And Return Rc  curl -sk %{VIC-ADMIN}/logs/docker-personality.log -b /tmp/cookies-%{VCH-NAME} | grep -q "%{TEST_PASSWORD}"
    Should Be Equal As Integers  ${rc}  1
    ${rc}=  Run And Return Rc  curl -sk %{VIC-ADMIN}/logs/vicadmin.log -b /tmp/cookies-%{VCH-NAME} | grep -q "%{TEST_PASSWORD}"
    Should Be Equal As Integers  ${rc}  1

    Remove File  /tmp/cookies-%{VCH-NAME}

Cleanup VIC Appliance On Test Server
    Log To Console  Gathering logs from the test server %{VCH-NAME}
    Gather Logs From Test Server
    Log To Console  Deleting the VCH appliance %{VCH-NAME}
    ${output}=  Run VIC Machine Delete Command
    Run Keyword And Ignore Error  Cleanup VCH Bridge Network  %{VCH-NAME}
    [Return]  ${output}

Cleanup VCH Bridge Network
    [Arguments]  ${name}
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove ${name}-bridge
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.info
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Not Contain  ${out}  ${name}-bridge

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
    \   Wait Until Keyword Succeeds  6x  5s  Check Delete Success  ${item}

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
    \   Wait Until Keyword Succeeds  6x  5s  Check Delete Success  ${vm}

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

Get Scratch Disk From VM Info
    [Arguments]  ${vm}
    ${disks}=  Run  govc vm.info -json ${vm} | jq -r '.VirtualMachines[].Layout.Disk[].DiskFile[]'
    ${disks}=  Split To Lines  ${disks}
    :FOR  ${disk}  IN  @{disks}
    \   ${disk}=  Fetch From Right  ${disk}  ${SPACE}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${disk}  scratch.vmdk
    \   Return From Keyword If  ${status}  ${disk}

Cleanup Dangling Containers On Test Server
    ${vms}=  Run  govc ls vm
    ${vms}=  Split To Lines  ${vms}
    :FOR  ${vm}  IN  @{vms}
    \   # Ignore VCH's, we only care about containers at this point
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${vm}  VCH
    \   Continue For Loop If  ${status}
    \   ${disk}=  Get Scratch Disk From VM Info  ${vm}
    \   ${vch}=  Fetch From Left  ${disk}  /
    \   ${vch}=  Split String  ${vch}  -
    \   # Skip any VM that is not associated with integration tests
    \   Continue For Loop If  '@{vch}[0]' != 'VCH'
    \   ${state}=  Get State Of Drone Build  @{vch}[1]
    \   # Skip any VM that is still running
    \   Continue For Loop If  '${state}' == 'running'
    \   # Destroy the VM and remove it from datastore because it is a dangling container
    \   Log To Console  Cleaning up dangling container: ${vm}
    \   ${out}=  Run  govc vm.destroy ${vm}
    \   ${name}=  Fetch From Right  ${vm}  /
    \   ${out}=  Run  govc datastore.rm ${name}
    \   Wait Until Keyword Succeeds  6x  5s  Check Delete Success  ${name}
