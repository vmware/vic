*** Settings ***
Documentation  Test 6-04 - Verify vic-machine create basic use cases
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create VCH - custom base disk
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    ${output}=  Run  docker %{VCH-PARAMS} logs $(docker %{VCH-PARAMS} start $(docker %{VCH-PARAMS} create --name customDiskContainer busybox /bin/df -h) && sleep 10) | grep /dev/sda | awk '{print $2}'
    # df shows GiB and vic-machine takes in GB so 6GB on cmd line == 5.5GB in df
    Should Be Equal As Strings  ${output}  5.5G
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f customDiskContainer
    Should Be Equal As Integers  ${rc}  0

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - URL without user and password
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  vSphere user must be specified

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Create VCH - target URL
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - operations user
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --ops-user=%{TEST_USERNAME} --ops-password=%{TEST_PASSWORD}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - specified datacenter
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Requires vCenter environment

    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --compute-resource=%{TEST_DATACENTER}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server


Create VCH - defaults
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Contain  ${output}  Installer completed successfully
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Get Docker Params  ${output}  ${true}
    ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Contain  ${output}  Installer completed successfully
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - full params
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} --volume-store=%{TEST_DATASTORE}/test:default ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - custom image store directory
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store %{TEST_DATASTORE}/vic-machine-test-images --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} ${vicmachinetls}

    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}
    ${output}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls
    Should Contain  ${output}  vic-machine-test-images

    Run Regression Tests
    Cleanup VIC Appliance On Test Server
    ${output}=  Run  GOVC_DATASTORE=%{TEST_DATASTORE} govc datastore.ls
    Should Not Contain  ${output}  vic-machine-test-images

Create VCH - long VCH name
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME}-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  exceeds the permitted 31 characters limit

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Create VCH - Existing VCH name
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  exists, to install with same name, please delete it first

    Cleanup VIC Appliance On Test Server

Create VCH - Existing VM name
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Create dummy VM
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.create -net=%{PUBLIC_NETWORK} %{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Get Docker Params  ${output}  ${true}
    Log  ${output}
    Should Contain  ${output}  Installer completed successfully
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Keyword And Ignore Error  Cleanup VIC Appliance On Test Server
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.destroy %{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    Cleanup VCH Bridge Network  %{VCH-NAME}

Create VCH - Existing RP on ESX
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Pass Execution  Test skipped on VC

    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Create dummy RP
    ${rc}  ${output}=  Run And Return Rc And Output  govc pool.create %{TEST_RESOURCE}/%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  Installer completed successfully
    Log  Installer completed successfully: %{VCH-NAME}

    Cleanup VIC Appliance On Test Server

    ${rc}  ${output}=  Run And Return Rc And Output  govc pool.destroy %{TEST_RESOURCE}/%{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0

Create VCH - Existing vApp on vCenter
    Pass execution  Test not implemented

Create VCH - defaults with --no-tls
    ${status}=  Get State Of Github Issue  3063

    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --no-tls
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}


    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - force accept target thumbprint
    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    # Test that --force without --thumbprint accepts the --target thumbprint
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --force --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Create VCH - Specified keys
    Pass execution  Test not implemented until vic-machine can poll status correctly

Create VCH - Invalid keys
    ${domain}=  Get Environment Variable  DOMAIN  ''
    Run Keyword If  '${domain}' == ''  Pass Execution  Skipping test - domain not set, won't generate keys

    Set Test Environment Variables
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}

    # Invalid server key
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --tls-ca="./%{VCH-NAME}/ca.pem" --cert="./%{VCH-NAME}/server-cert.pem" --key="./%{VCH-NAME}/ca.pem"
    Should Contain  ${output}  found a certificate rather than a key in the PEM for the private key

    # Invalid server cert
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --tls-ca="./%{VCH-NAME}/ca.pem" --cert="./%{VCH-NAME}/server-key.pem" --key="./%{VCH-NAME}/server-key.pem"
    Should Contain  ${output}  did find a private key; PEM inputs may have been switched

    # Invalid CA
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls} --tls-ca="./%{VCH-NAME}/key.pem" --cert="./%{VCH-NAME}/server-cert.pem" --key="./%{VCH-NAME}/server-key.pem"
    Should Contain  ${output}  Unable to load certificate authority data

    Cleanup VIC Appliance On Test Server

Create VCH - Reuse keys
    ${domain}=  Get Environment Variable  DOMAIN  ''
    Run Keyword If  '${domain}' == ''  Pass Execution  Skipping test - domain not set, won't generate keys

    Set Test Environment Variables

    # use one install to generate certificates
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    # remove the initial deployment
    ${ret}=  Run  bin/vic-machine-linux delete --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --name=%{VCH-NAME} --force
    Should Contain  ${ret}  Completed successfully

    # deploy using the same name - should reuse certificates
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully

    Should Contain  ${output}  Loaded server certificate
    Should Contain  ${output}  Loaded CA with default name from certificate path
    Should Contain  ${output}  Loaded client certificate with default name from certificate path

    Cleanup VIC Appliance On Test Server

Basic timeout
    ${status}=  Get State Of Github Issue  3241
    Run Keyword If  '${status}' == 'closed'  Fail  Test 6-04-Create-Basic.robot needs to be updated now that Issue #3241 has been resolved
    Log  Issue \#3241 is blocking implementation  WARN

#    Set Test Environment Variables
#    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
#    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
#
#    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --timeout 30s ${vicmachinetls} --public-network-ip=172.16.5.5/24 --public-network-gateway=172.16.5.1/24
#    Should Contain  ${output}  Create timed out
#
#    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE} --name ${vch-name}
#    Should Contain  ${ret}  Completed successfully
#    ${out}=  Run  govc ls vm
#    Should Not Contain  ${out}  %{VCH-NAME}

Basic VCH resource config
    Pass execution  Test not implemented

Invalid VCH resource config
    Pass execution  Test not implemented

Use resource pool
    Pass execution  Test not implemented

CPU reservation shares invalid
    Pass execution  Test not implemented

CPU reservation invalid
    Pass execution  Test not implemented

CPU reservation valid
    Pass execution  Test not implemented

Memory reservation shares invalid
    Pass execution  Test not implemented

Memory reservation invalid 1
    Pass execution  Test not implemented

Memory reservation invalid 2
    Pass execution  Test not implemented

Memory reservation invalid 3
    Pass execution  Test not implemented

Memory reservation valid
    Pass execution  Test not implemented

Extension installation
    Pass execution  Test not implemented

Install existing extension
    Pass execution  Test not implemented
