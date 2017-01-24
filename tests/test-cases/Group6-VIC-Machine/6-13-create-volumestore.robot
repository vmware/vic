*** Settings ***
Documentation  Test 6-13 - Verify proper volume store option behavior
Resource  ../../resources/Util.robot

*** Test Cases ***
Create with default volume store
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vsd %{TEST_DATASTORE}/test
       Should Contain  ${output}  Installer completed successfully
       Get Docker Params  ${output}  ${true}
       Log To Console  Installer completed successfully: %{VCH-NAME}\n

       ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
       Should Be Equal As Integers  ${rc}  0
       Should Contain  ${output}  VolumeStores: default
       Cleanup VIC Appliance On Test Server

Create default store with volume store flag
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vs %{TEST_DATASTORE}/test:default
       Should Contain  ${output}  Installer completed successfully
       Get Docker Params  ${output}  ${true}
       Log To Console  Installer completed successfully: %{VCH-NAME}\n

       ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
       Should Be Equal As Integers  ${rc}  0
       Should Contain  ${output}  VolumeStores: default
       Cleanup VIC Appliance On Test Server

Create default store using both flags with the same path
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vsd %{TEST_DATASTORE}/test -vs %{TEST_DATASTORE}/test:default
       Should Contain  ${output}  Installer completed successfully
       Get Docker Params  ${output}  ${true}
       Log To Console  Installer completed successfully: %{VCH-NAME}\n

       ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
       Should Be Equal As Integers  ${rc}  0
       Should Contain  ${output}  VolumeStores: default
       Cleanup VIC Appliance On Test Server

Create default store using both flags without the same path
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vsd %{TEST_DATASTORE}/testStore1 -vs %{TEST_DATASTORE}/testStore2:default
       Should Contain  ${output}  vic-machine-linux create failed: Error occurred while processing volume stores: Multiple paths were tagged for the same label(default), volumestore labels can only have one distinct path. The preferred method of setting the default store is the --default-volume-store flag.

Create with normal volume store
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vs %{TEST_DATASTORE}/test:TestStore
       Should Contain  ${output}  Installer completed successfully
       Get Docker Params  ${output}  ${true}
       Log To Console  Installer completed successfully: %{VCH-NAME}\n

       ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
       Should Be Equal As Integers  ${rc}  0
       Should Contain  ${output}  VolumeStores: TestStore
       Cleanup VIC Appliance On Test Server

Create with both volume store flags
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vsd %{TEST_DATASTORE}/test -vs %{TEST_DATASTORE}/test2:TestStore
       Should Contain  ${output}  Installer completed successfully
       Get Docker Params  ${output}  ${true}
       Log To Console  Installer completed successfully: %{VCH-NAME}\n

       ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
       Should Be Equal As Integers  ${rc}  0
       Should Contain  ${output}  VolumeStores: TestStore default
       Cleanup VIC Appliance On Test Server

Create with overlapping paths to single label
       Set Test Environment Variables
       Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
       Run Keyword And Ignore Error  Cleanup Datastore On Test Server

       ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls} -vs %{TEST_DATASTORE}/test:TestStore -vs %{TEST_DATASTORE}/test2:TestStore
       Should Contain  ${output}   vic-machine-linux create failed: Error occurred while processing volume stores: Multiple paths were tagged for the same label(TestStore), volumestore labels can only have one distinct path.
