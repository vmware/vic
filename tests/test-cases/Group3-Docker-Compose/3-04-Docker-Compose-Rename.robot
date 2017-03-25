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
Documentation  Test 3-04 - Docker Compose Rename
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: busybox\n${SPACE}${SPACE}command: ["/bin/top"] 
${newyml}  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: ubuntu\n${SPACE}${SPACE}command: ["date"]

*** Test Cases ***
Compose up -d --force-recreate
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    Run  echo '${yml}' > compose-rename.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file compose-rename.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file compose-rename.yml up -d --force-recreate
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Compose up -d with a new image
    Run  echo '${newyml}' > compose-rename.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file compose-rename.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
