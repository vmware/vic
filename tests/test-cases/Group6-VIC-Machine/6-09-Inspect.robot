# Copyright 2017 VMware, Inc. All Rights Reserved.
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
# limitations under the License.

*** Settings ***
Documentation  Test 6-09 - Verify vic-machine inspect functions
Resource  ../../resources/Util.robot
Test Teardown  Run Keyword If Test Failed  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Inspect VCH configuration
    Install VIC Appliance To Test Server

    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --configuration --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --name=%{VCH-NAME}
    Should Contain  ${output}  --debug=1
    Should Contain  ${output}  --name=%{VCH-NAME}
    Should Contain  ${output}  --target=https://%{TEST_URL}
    Should Contain  ${output}  --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --image-store=%{TEST_DATASTORE}
    Should Contain  ${output}  --compute-resource=%{TEST_RESOURCE}
    Should Contain  ${output}  --timeout=
    Should Contain  ${output}  --volume-store=%{TEST_DATASTORE}/test
    Should Contain  ${output}  --appliance-iso=bin/appliance.iso
    Should Contain  ${output}  --bootstrap-iso=bin/bootstrap.iso
    Should Contain  ${output}  --force=true
    Should Contain  ${output}  --bridge-network=%{BRIDGE_NETWORK}
    Should Contain  ${output}  --public-network=VM Network
    Should Not Contain  ${output}  --cpu
    Should Be Equal As Integers  0  ${rc}

    Cleanup VIC Appliance On Test Server

Verify inspect output for a full tls VCH
    Install VIC Appliance To Test Server

    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  DOCKER_CERT_PATH=${EXECDIR}/%{VCH-NAME}

    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --cert-path=%{VCH-NAME}
    Should Contain  ${output}  DOCKER_CERT_PATH=${EXECDIR}/%{VCH-NAME}

    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --cert-path=fakeDir
    Should Not Contain  ${output}  DOCKER_CERT_PATH=${EXECDIR}/%{VCH-NAME}
    Should Contain  ${output}  Unable to find valid client certs
    Should Contain  ${output}  DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments

    Cleanup VIC Appliance On Test Server

Verify inspect output for a --no-tls VCH
    Set Test Environment Variables
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --no-tls
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}

    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Not Contain  ${output}  DOCKER_CERT_PATH=${EXECDIR}/%{VCH-NAME}
    Should Not Contain  ${output}  Unable to find valid client certs
    Should Not Contain  ${output}  DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments

    Cleanup VIC Appliance On Test Server

Verify inspect output for a --no-tlsverify VCH
    Set Test Environment Variables
    ${output}=  Run  bin/vic-machine-linux create --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} --public-network=%{PUBLIC_NETWORK} --no-tlsverify
    Should Contain  ${output}  Installer completed successfully
    Get Docker Params  ${output}  ${true}

    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Not Contain  ${output}  DOCKER_CERT_PATH=${EXECDIR}/%{VCH-NAME}
    Should Not Contain  ${output}  Unable to find valid client certs
    Should Not Contain  ${output}  DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments

    Cleanup VIC Appliance On Test Server

