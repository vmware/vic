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
Documentation  Test 6-09 - Verify vic-machine inspect functions
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Inspect VCH configuration
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
    Should Not Contain  ${output}  --insecure-registry
    Should Not Contain  ${output}  --cpu
    Should Be Equal As Integers  0  ${rc}

