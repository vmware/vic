*** Settings ***
Documentation  Test 5-10 - Multiple Datacenters
Resource  ../../resources/Util.robot
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    ${esx1}  ${esx4-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    ${esx2}  ${esx5-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    Set Global Variable  @{list}  ${esx1}  ${esx2}

    Create a Simple VC Cluster  datacenter1  cls1

    Log To Console  Create datacenter2 on the VC
    ${out}=  Run  govc datacenter.create datacenter2
    Should Be Empty  ${out}
    ${out}=  Run  govc host.add -hostname=${esx4-ip} -username=root -dc=datacenter2 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Log To Console  Create datacenter3 on the VC
    ${out}=  Run  govc datacenter.create datacenter3
    Should Be Empty  ${out}
    ${out}=  Run  govc host.add -hostname=${esx5-ip} -username=root -dc=datacenter3 -password=e2eFunctionalTest -noverify=true
    Should Contain  ${out}  OK

    Set Environment Variable  TEST_DATACENTER  /datacenter1
    Set Environment Variable  GOVC_DATACENTER  /datacenter1
    Install VIC Appliance To Test Server  certs=${false}  vol=default

    Run Regression Tests

    Cleanup VIC Appliance On Test Server
