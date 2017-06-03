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

*** Test Cases ***
Configure VCH debug state
    ${output}=  Check VM Guestinfo  %{VCH-NAME}  guestinfo.vice./init/diagnostics/debug
    Should Contain  ${output}  1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${id1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${vm1}=  Get VM display name  ${id1}
    ${output}=  Check VM Guestinfo  ${vm1}  guestinfo.vice./diagnostics/debug
    Should Contain  ${output}  1
    ${output}=  Run  bin/vic-machine-linux configure --debug 0 --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    ${output}=  Check VM Guestinfo  %{VCH-NAME}  guestinfo.vice./init/diagnostics/debug
    Should Contain  ${output}  0
    ${output}=  Check VM Guestinfo  ${vm1}  guestinfo.vice./diagnostics/debug
    Should Contain  ${output}  1
    ${rc}  ${id2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${vm2}=  Get VM display name  ${id2}
    ${output}=  Check VM Guestinfo  ${vm2}  guestinfo.vice./diagnostics/debug
    Should Contain  ${output}  0
    ${output}=  Run  bin/vic-machine-linux configure --debug 1 --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc snapshot.tree -vm %{VCH-NAME} | grep reconfigure
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Split To Lines  ${output}
    Length Should Be  ${output}  1
    ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=%{VCH-NAME} --conf --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT}
    Should Be Equal As Integers  0  ${rc}
    Should Contain  ${output}  --debug=1

Configure VCH https-proxy
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

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --https-proxy https://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128
    ${output}=  Run  bin/vic-machine-linux inspect --name=%{VCH-NAME} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" --thumbprint=%{TEST_THUMBPRINT} --conf
    Should Contain  ${output}  --https-proxy=https://proxy.vmware.com:3128
	Should Not Contain  ${output}  --http-proxy
