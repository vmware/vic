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
Documentation   Test 1-37 - Docker Run As USER
Resource        ../../resources/Util.robot
Suite Setup     Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Run as NewUser in NewGroup
     ${rc}    ${output}=    Run And Return Rc And Output    docker %{VCH-PARAMS} run gigawhitlocks/1-37-docker-user-newuser-newgroup:latest
     Should Be Equal As Integers    ${rc}       0
     Should Match Regexp            ${output}   uid=\\d+\\\(newuser\\\)\\s+gid=\\d+\\\(newuser\\\)\\s+groups=\\d+\\\(newuser\\\)

Run as UID 2000
     ${rc}    ${output}=    Run And Return Rc And Output    docker %{VCH-PARAMS} run gigawhitlocks/1-37-docker-user-uid-2000:latest
     Should Be Equal As Integers    ${rc}       0
     Should Contain                 ${output}   uid=2000 gid=0(root)

Run as UID:GID 2000:2000
     ${rc}    ${output}=    Run And Return Rc And Output    docker %{VCH-PARAMS} run gigawhitlocks/1-37-docker-user-uid-gid-2000-2000:latest
     Should Be Equal As Integers    ${rc}       0
     Should Contain                 ${output}   uid=2000 gid=2000

