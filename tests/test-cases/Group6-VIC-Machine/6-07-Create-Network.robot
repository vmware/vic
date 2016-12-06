*** Settings ***
Documentation  Test 6-07 - Verify vic-machine create network function
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Public network - default
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    ${vm}=  Get VM Name  %{VCH-NAME}
    ${info}=  Get VM Info  ${vm}
    Should Contain  ${info}  VM Network

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Public network - invalid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --public-network=AAAAAAAAAA ${vicmachinetls}
    Should Contain  ${output}  --public-network: network 'AAAAAAAAAA' not found
    Should Contain  ${output}  vic-machine-linux failed

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Public network - invalid vCenter
    Pass execution  Test not implemented

Public network - DHCP
    Pass execution  Test not implemented

Public network - valid
    Pass execution  asdf

Management network - none
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    ${status}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "management" is sharing NIC with "public"
    ${status2}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "public" is sharing NIC with "management"
    ${status3}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "public" is sharing NIC with "client"
    ${status4}=  Run Keyword And Return Status  Should Contain  ${output}  Network role "management" is sharing NIC with "client"
    Should Be True  ${status} | ${status2} | ${status3} | ${status4}
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Management network - invalid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --management-network=AAAAAAAAAA ${vicmachinetls}
    Should Contain  ${output}  --management-network: network 'AAAAAAAAAA' not found
    Should Contain  ${output}  vic-machine-linux failed

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Management network - invalid vCenter
    Pass execution  Test not implemented

Management network - unreachable
    Pass execution  Test not implemented

Management network - valid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --management-network=%{PUBLIC_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Bridge network - vCenter none
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} ${vicmachinetls}
    Should Contain  ${output}  FAILURE

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}


Bridge network - ESX none
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Pass Execution  Test skipped on VC

    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Bridge network - invalid
    Pass execution  asdf
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=AAAAAAAAAA ${vicmachinetls}
    Should Contain  ${output}  --bridge-network: network 'AAAAAAAAAA' not found
    Should Contain  ${output}  vic-machine-linux failed

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - invalid vCenter
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Pass execution  Test not implemented

Bridge network - non-DPG
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Pass Execution  Test skipped on ESXi

    Pass execution  Test not implemented

Bridge network - valid
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Bridge network - reused port group
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --management-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --client-network=%{BRIDGE_NETWORK} ${vicmachinetls}
    Should Contain  ${output}  the bridge network must not be shared with another network role

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - invalid IP settings
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server

    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --bridge-network-range 1.1.1.1 ${vicmachinetls}
    Should Contain  ${output}  Error parsing bridge network ip range

    # Delete the portgroup added by env vars keyword
    Cleanup VCH Bridge Network  %{VCH-NAME}

Bridge network - valid
    Pass execution  Test not implemented

Container network invalid 1
    Pass execution  Test not implemented

Container network invalid 2
    Pass execution  Test not implemented

Container network 1
    Pass execution  Test not implemented

Container network 2
    Pass execution  Test not implemented

Network mapping invalid
    Pass execution  Test not implemented

Network mapping gateway invalid
    Pass execution  Test not implemented

Network mapping IP invalid
    Pass execution  Test not implemented

DNS format invalid
    Pass execution  Test not implemented

Network mapping
    Pass execution  Test not implemented

VCH static IP - Static public
    Pass execution  Test not implemented

VCH static IP - Static client
    Pass execution  Test not implemented

VCH static IP - Static management
    Pass execution  Test not implemented

VCH static IP - different port groups 1
    Pass execution  Test not implemented

VCH static IP - different port groups 2
    Pass execution  Test not implemented

VCH static IP - same port group
    Pass execution  Test not implemented

VCH static IP - same subnet for multiple port groups
    Pass execution  Test not implemented
