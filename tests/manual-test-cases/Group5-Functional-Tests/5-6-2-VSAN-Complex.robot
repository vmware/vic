*** Settings ***
Documentation  Test 5-6-2 - VSAN-Complex
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Complex VSAN
    ${name}=  Evaluate  'vic-vsan-complex-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  --noSupportBundles --vcvaBuild ${VC_VERSION} --esxPxeDir ${ESX_VERSION} --esxBuild ${ESX_VERSION} --testbedName vcqa-vsan-complex-pxeBoot-vcva --runName ${name}
    ${out}=  Split To Lines  ${out}
    :FOR  ${line}  IN  @{out}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${line}  .vcva-${VC_VERSION}' is up. IP:
    \   ${ip}=  Run Keyword If  ${status}  Fetch From Right  ${line}  ${SPACE}
    \   Run Keyword If  ${status}  Set Test Variable  ${vc-ip}  ${ip}
    \   Exit For Loop If  ${status}

    Set Global Variable  @{list}  ${user}-${name}.vcva-${VC_VERSION}  ${user}-${name}.esx.0  ${user}-${name}.esx.1  ${user}-${name}.esx.2  ${user}-${name}.esx.3  ${user}-${name}.esx.4  ${user}-${name}.esx.5  ${user}-${name}.esx.6  ${user}-${name}.esx.7  ${user}-${name}.nfs.0  ${user}-${name}.iscsi.0

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Create A Distributed Switch  vcqaDC

    Create Three Distributed Port Groups  vcqaDC

    Add Host To Distributed Switch  /vcqaDC/host/cluster-vsan-1

    Log To Console  Enable DRS and VSAN on the cluster
    ${out}=  Run  govc cluster.change -drs-enabled /vcqaDC/host/cluster-vsan-1
    Should Be Empty  ${out}

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    ${datastore}=  Run  govc ls -t Datastore host/cluster-vsan-1/* | grep -v local | xargs -n1 basename | sort | uniq | grep vsan
    Set Environment Variable  TEST_DATASTORE  "${datastore}"
    Set Environment Variable  TEST_RESOURCE  cluster-vsan-1
    Set Environment Variable  TEST_TIMEOUT  30m

    Install VIC Appliance To Test Server

    Run Regression Tests
