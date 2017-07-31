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
Documentation  Test 6-02 - Verify default parameters
Resource  ../../resources/Util.robot

*** Test Cases ***
Delete with defaults
    Set Test Environment Variables

    ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD}
    Should Contain  ${ret}  vic-machine-linux delete failed:  resource pool
    Should Contain  ${ret}  /Resources/virtual-container-host' not found

Wrong Password No Panic
  Set Test Environment Variables
  ${ret}=  Run  bin/vic-machine-linux create --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux create failed

  ${ret}=  Run  bin/vic-machine-linux delete --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux delete failed

  ${ret}=  Run  bin/vic-machine-linux inspect --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux inspect  failed

  ${ret}=  Run  bin/vic-machine-linux ls --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux ls failed

  ${ret}=  Run  bin/vic-machine-linux upgrade --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux upgrade failed

  ${ret}=  Run  bin/vic-machine-linux configure --target %{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=INCORRECT
  Should Contain  ${ret}  vic-machine-linux configure failed
