# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  This resource contains all keywords related to creating, deleting, maintaining an instance of Admiral

*** Keywords ***
Install Admiral
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -p 8282:8282 --name admiral vmware/admiral
    Should Be Equal As Integers  0  ${rc}
    Set Environment Variable  ADMIRAL-IP  %{VCH-IP}:8282
    :FOR  ${idx}  IN RANGE  0  10
    \   ${out}=  Run  curl %{ADMIRAL-IP}
    \   ${status}=  Run Keyword And Return Status  Should Contain  ${out}  <body class="admiral-default">
    \   Return From Keyword If  ${status}
    \   Sleep  5
    Fail  Install Admiral failed: Admiral endpoint failed to respond to curl

Cleanup Admiral
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f admiral
    Should Be Equal As Integers  0  ${rc}