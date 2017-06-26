# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 6-05 - Verify vic-machine create validation function
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Suggest resources - Invalid datacenter
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/WOW --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Suggested datacenters:
    Should Contain  ${output}  vic-machine-linux create failed:

Suggest resources - Invalid target path
    Log To Console  \nRunning vic-machine create
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/MUCH/DATACENTER --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} ${vicmachinetls}
    Should Contain  ${output}  Suggested datacenters:
    Should Contain  ${output}  vic-machine-linux create failed:

Create VCH - target thumbprint verification
    Log To Console  \nRunning vic-machine create - thumbprint verification
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    ${output}=  Run  bin/vic-machine-linux create --thumbprint=NOPE --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --image-store=ENOENT ${vicmachinetls}
    Should Contain  ${output}  thumbprint does not match

Default image datastore
    # This test case is dependent on the ESX environment having only one datastore
    Log To Console  \nRunning vic-machine create - default image datastore
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Log  ${output}
    Should Contain  ${output}  Using default datastore
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}...
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Custom image datastore
    # This test case is dependent on the ESX environment having only one datastore
    Log To Console  \nRunning vic-machine create - custom image datastore
    Set Test Environment Variables
    # Attempt to cleanup old/canceled tests
    Run Keyword And Ignore Error  Cleanup Dangling VMs On Test Server
    Run Keyword And Ignore Error  Cleanup Datastore On Test Server
    Set Test VCH Name

    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --image-store=%{TEST_DATASTORE}/long/weird/path ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Log  ${output}
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}...
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Trailing slash works as expected
    Log To Console  \nInstalling VCH to test server...
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --insecure-registry harbor.ci.drone.local
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    Log To Console  Installer completed successfully: %{VCH-NAME}...
    Run Regression Tests
    Cleanup VIC Appliance On Test Server

Whitelist registries - blocked registry wildcard domain
    Set Test Environment Variables
    Log To Console  \nInstalling VCH to test server...
    # *.docker.io
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --whitelist-registry *.docker.io
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    # try a docker pull from docker.io; this should fail
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Not Be Equal As Integers  ${rc}  0
    Cleanup VIC Appliance On Test Server

Whitelist registries - blocked registry ip address of valid registry fqdn
    Set Test Environment Variables
    # ip address of docker.io
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --whitelist-registry 52.200.132.201
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    # try a docker pull from docker.io; this should fail
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Not Be Equal As Integers  ${rc}  0
    Cleanup VIC Appliance On Test Server

Whitelist registries - allowed registry fqdn
    Set Test Environment Variables
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --whitelist-registry docker.io
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    # try a docker pull from docker.io; this should succeed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Cleanup VIC Appliance On Test Server

Whitelist registries - allowed registry wildcard domain
    Set Test Environment Variables
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target=%{TEST_URL}/ --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} ${vicmachinetls} --whitelist-registry *docker.io
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}
    # try a docker pull from docker.io; this should succeed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Cleanup VIC Appliance On Test Server
