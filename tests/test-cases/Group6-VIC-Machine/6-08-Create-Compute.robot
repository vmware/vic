*** Settings ***
Documentation  Test 6-08 - Verify vic-machine compute resources
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***

Compute resources - Default resource pool
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Compute resources - Non-Default resource pool
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name
    Log To Console  Creating resource pool %{TEST_RESOURCE}/test1
    Run  govc pool.create %{TEST_RESOURCE}/test1

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --compute-resource=%{TEST_RESOURCE}/test1 --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...
    ${out}=  Run  govc vm.info -json ${vch-name}
    Log To Console  Checking the VM Info for resource pool
    Should Contain  ${out}  POOL_PATH=%{TEST_RESOURCE}/test1
    
    Run Regression Tests
    Run  bin/vic-machine-linux delete --name=${vch-name} --target=%{TEST_URL}%{TEST_DATACENTER} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE}/test1 --timeout %{TEST_TIMEOUT}

Compute resources - Correct Absolute path ESXi
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --compute-resource=%{TEST_RESOURCE} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Installer completed successfully

    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: ${vch-name}...

    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Compute resources - InCorrect Absolute path ESXi
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --compute-resource=%{TEST_RESOURCE}/incorrect --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Suggested values for --compute-resource:
    Should Contain  ${output}  resource pool '%{TEST_RESOURCE}/incorrect' not found
    Should Contain  ${output}  vic-machine-linux failed: validation of configuration failed

    ${out}=  Run  govc ls vm
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}
    ${out}=  Run  govc datastore.ls
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}
    ${out}=  Run  govc ls host/*/Resources/*
    Log  ${out}
    Should Not Contain  ${out}  ${vch-name}

Compute resources - Wrong relative path multiple VC clusters and multiple available resource pools
    Pass execution  TODO in nightly

Compute resources - Correct relative path with single VC cluster
    Pass execution  TODO in the nightly functional run - Create with compute resource set to <cluster name> (real cluster name here)

Compute resources - Correct relative path with single VC cluster
    Pass execution  TODO in the nightly functional run - Create with compute resource set to RP1 (RP1 exists in cluster)

Compute resources - Correct relative path with single VC cluster
    Pass execution  TODO in the nightly functional run - Create with compute resource not set 

Compute resources - Correct relative path with multiple VC clusters and multiple available resource pools 
    Pass execution  TODO in the nightly functional run - Create with compute resource set to Cluster (Cluster exists)

Compute resources - Correct relative path with multiple VC clusters and multiple available resource pools 
    Pass execution  TODO in the nightly functional run - Create with compute resource set to Cluster/RP1 (Cluster/RP1 exists in cluster)
