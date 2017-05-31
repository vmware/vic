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
Configure VCH
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --http-proxy http://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  http://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128

    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --https-proxy https://proxy.vmware.com:3128
    Should Contain  ${output}  Completed successfully
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTPS_PROXY
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  https://proxy.vmware.com:3128
    ${rc}  ${output}=  Run And Return Rc And Output  govc vm.info -e %{VCH-NAME} | grep HTTP_PROXY
    Should Be Equal As Integers  ${rc}  1
    Should Not Contain  ${output}  proxy.vmware.com:3128
