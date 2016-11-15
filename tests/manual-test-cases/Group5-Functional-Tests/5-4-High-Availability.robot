*** Settings ***
Documentation  Test 5-4 - High Availability
Resource  ../../resources/Util.robot
Suite Teardown  Nimbus Cleanup

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX1}  ${esx1}
    ${esx2}  ${esx2-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX2}  ${esx2}
    ${esx3}  ${esx3-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${ESX3}  ${esx3}

    ${vc}  ${vc-ip}=  Deploy Nimbus vCenter Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Suite Variable  ${VC}  ${vc}

    Log To Console  Create a datacenter on the VC
    ${out}=  Run  govc datacenter.create ha-datacenter
    Should Be Empty  ${out}

    Log To Console  Create a cluster on the VC
    ${out}=  Run  govc cluster.create cls
    Should Be Empty  ${out}

    Log To Console  Add ESX host to the VC
    ${out}=  Run  govc cluster.add -hostname=${esx1-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx2-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK
    ${out}=  Run  govc cluster.add -hostname=${esx3-ip} -username=root -dc=ha-datacenter -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create a distributed switch
    ${out}=  Run  govc dvs.create -dc=ha-datacenter test-ds
    Should Contain  ${out}  OK

    Log To Console  Create three new distributed switch port groups for management and vm network traffic
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds management
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds vm-network
    Should Contain  ${out}  OK
    ${out}=  Run  govc dvs.portgroup.add -nports 12 -dc=ha-datacenter -dvs=test-ds bridge
    Should Contain  ${out}  OK

    Log To Console  Add all the hosts to the distributed switch
    ${out}=  Run  govc dvs.add -dvs=test-ds -pnic=vmnic1 /ha-datacenter/host/cls
    Should Contain  ${out}  OK

    Log To Console  Enable HA on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled -ha-enabled /ha-datacenter/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Keyword And Ignore Error  Run Regression Tests

    ${output}=  Run  govc vm.info ${vch-name}/${vch-name}
    @{output}=  Split To Lines  ${output}
    ${curHost}=  Fetch From Right  @{output}[-1]  ${SPACE}

    # Abruptly power off the host
    Open Connection  ${curHost}  prompt=:~]
    Login  root  e2eFunctionalTest
    ${out}=  Execute Command  poweroff -d 0 -f
    Close connection

    Run Regression Tests
