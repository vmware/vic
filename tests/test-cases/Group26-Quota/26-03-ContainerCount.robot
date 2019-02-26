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
Documentation     Suite 26-03 - ContainerCount
Resource          ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  additional-args='--storage-quota=15'
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  30 minutes

*** Test Cases ***
Create a VCH with container count and inspect shows the container count
    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --containers=1

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

Create second container and get container count exceeding failure
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  125
    Should Contain  ${output}  Container count exceeds limit 1

Configure VCH with a larger container count of 2
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --containers 2
    Should Contain  ${output}  Completed successfully

    ${output}=  Run  bin/vic-machine-linux inspect config --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT}
    Should Contain  ${output}  --containers=2

Create second container and succeed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  125
    Should Contain  ${output}  Container count exceeds limit 2

Configure VCH with unlimited container count of 0
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --containers 0
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0
