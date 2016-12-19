*** Settings ***
Documentation  Test 5-15 - NFS Datastore
Resource  ../../resources/Util.robot
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup

*** Test Cases ***
Test
    Log To Console  \nStarting test...
    Create a Simple VC Cluster  datacenter  cls

    ${name}  ${ip}=  Deploy Nimbus NFS Datastore  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    ${out}=  Run  govc datastore.create -mode readWrite -type nfs -name nfsDatastore -remote-host ${ip} -remote-path /store cls
    Should Be Empty  ${out}

    Set Environment Variable  TEST_DATASTORE  nfsDatastore
    Install VIC Appliance To Test Server

    Run Regression Tests

    Cleanup VIC Appliance On Test Server
