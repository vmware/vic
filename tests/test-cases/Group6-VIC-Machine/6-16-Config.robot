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
# limitations under the License

*** Settings ***
Documentation  Test 6-16 - Verify vic-machine configure
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Keywords ***
Verify VCH Debug State
    [Arguments]  ${expected}
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep guestinfo.vice./init/diagnostics/debug
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${expected}

Verify Container Debug State
    [Arguments]  ${vm}  ${expected}
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e ${vm} | grep guestinfo.vice./diagnostics/debug
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${expected}

*** Test Cases ***
Configure VCH
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --http-proxy http://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  http://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128
    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --conf
    Should Contain  ${output}  --http-proxy=http://proxy.vmware.com:3128
	Should Not Contain  ${output}  --https-proxy

    Run Regression Tests

Configure debug state
    Verify VCH Debug state  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${vm1}=  Get VM display name  ${id1}
    Verify Container Debug State  ${vm1}  1
    ${output}=  Run  bin/vic-machine-linux configure --debug 0 --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    Verify VCH Debug state  0
    Verify Container Debug State  ${vm1}  1
    ${rc}  ${id2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${vm2}=  Get VM display name  ${id2}
    Verify Container Debug State  ${vm2}  0
    ${output}=  Run  bin/vic-machine-linux configure --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=%{VCH-NAME} --conf --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Be Equal As Integers  0  ${rc}
    Should Contain  ${output}  --debug=1
