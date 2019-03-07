# Copyright 2018 VMware, Inc. All Rights Reserved.
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
Documentation     Suite 26-02 - Reconfigure
Resource          ../../resources/Util.robot
Resource          ../../resources/Group26-Storage-Quota-Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  10 minutes

*** Test Cases ***
Create a VCH without storage quota and docker info shows the storage usage
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  storage limit

    ${usageval}=  Get storage usage  ${output}
    Should Be Equal As Integers  ${usageval}  0

Create a container and docker info shows the storage usage has changed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${usageval}=  Get storage usage  ${output}

    Should Be True  ${usageval} > 9
    Set Suite Variable  ${pre_usage_val}  ${usageval}

Configure VCH with a larger storage quota of 15GB
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --storage-quota 15
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}
    Should Be Equal As Integers  ${limitval}  15
    Should Be Equal As Integers  ${usageval}  ${pre_usage_val}

Create second container and get storage quota exceeding failure
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  125
    Should Contain  ${output}  Storage quota exceeds
