*** Settings ***
Documentation  Test 5-16 - iSCSI Datastore
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
iSCSI Datastore
    ${name}=  Evaluate  'vic-iscsi-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --customizeTestbed '/esx desiredPassword=e2eFunctionalTest' --noSupportBundles --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vcqa-sdrs-iscsi-fullInstall-vcva --runName vic-iscsi
    Set Global Variable  @{list}  %{NIMBUS_USER}-vic-iscsi.vcva-${VC_VERSION}  %{NIMBUS_USER}-vic-iscsi.esx.0  %{NIMBUS_USER}-vic-iscsi.esx.1  %{NIMBUS_USER}-vic-iscsi.iscsi.0
    Should Contain  ${out}  "deployment_result"=>"PASS"

    ${out}=  Execute Command  nimbus-ctl ip %{NIMBUS_USER}-vic-iscsi.vcva-${VC_VERSION} | grep %{NIMBUS_USER}-vic-iscsi.vcva-${VC_VERSION}
    ${vc-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  nimbus-ctl ip %{NIMBUS_USER}-vic-iscsi.esx.0 | grep %{NIMBUS_USER}-vic-iscsi.esx.0
    ${esx0-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  nimbus-ctl ip %{NIMBUS_USER}-vic-iscsi.esx.1 | grep %{NIMBUS_USER}-vic-iscsi.esx.1
    ${esx1-ip}=  Fetch From Right  ${out}  ${SPACE}

    Set Environment Variable  GOVC_URL  ${esx0-ip}
    Set Environment Variable  GOVC_USERNAME  root
    Set Environment Variable  GOVC_PASSWORD  e2eFunctionalTest
    Run  govc host.esxcli network firewall set -e false
    Set Environment Variable  GOVC_URL  ${esx1-ip}
    Run  govc host.esxcli network firewall set -e false

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Create A Distributed Switch  vcqaDC

    Create Three Distributed Port Groups  vcqaDC

    Add Host To Distributed Switch  /vcqaDC/host/cls

    Log To Console  Enable DRS on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cls
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server

    Run Regression Tests