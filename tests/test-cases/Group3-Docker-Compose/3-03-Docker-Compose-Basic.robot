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
Documentation  Test 3-03 - Docker Compose Basic
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${True}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Variables ***
${yml}  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: python:2.7\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5000:5000"\n${SPACE}${SPACE}depends_on:\n${SPACE}${SPACE}- redis\n${SPACE}redis:\n${SPACE}${SPACE}image: redis\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5001:5001"
${link-yml}  version: "2"\nservices:\n${SPACE}redis1:\n${SPACE}${SPACE}image: redis:alpine\n${SPACE}${SPACE}container_name: redis1\n${SPACE}${SPACE}ports: ["6379"]\n${SPACE}web1:\n${SPACE}${SPACE}image: busybox\n${SPACE}${SPACE}container_name: a.b.c\n${SPACE}${SPACE}links:\n${SPACE}${SPACE}- redis1:aaa\n${SPACE}${SPACE}command: ["ping", "aaa"]

*** Test Cases ***
Compose basic
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    Run  echo '${yml}' > basic-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create vic_default
    Log  ${output}
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml create
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml start
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml logs
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml stop
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Compose kill
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f basic-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f basic-compose.yml kill redis
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f basic-compose.yml down
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0

Compose Up while another container is running (ps filtering related)
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f basic-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0

Compose Up with link
    Run  echo '${link-yml}' > link-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file link-compose.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file link-compose.yml logs
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  PING aaa
    Should Not Contain  ${output}  bad address 'aaa'

Compose bundle creation
    ${rc}  Run And Return Rc  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml pull
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --file basic-compose.yml bundle
    Log  ${output}
    Should Contain  ${output}  Wrote bundle
    Should Be Equal As Integers  ${rc}  0
