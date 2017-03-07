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
Documentation  Test 3-01 - Docker Compose LEMP
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Compose LEMP Server
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0

    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300
    # must set CURL_CA_BUNDLE to work around Compose bug https://github.com/docker/compose/issues/3365
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}

    ${vch_ip}=  Get Environment Variable  VCH_IP  %{VCH-IP}
    Log To Console  \nThe VCH IP is %{VCH-IP}

    Run  cat %{GOPATH}/src/github.com/vmware/vic/demos/compose/webserving-app/docker-compose.yml | sed -e "s/192.168.60.130/${vch_ip}/g" > lemp-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker-compose %{VCH-PARAMS} --file lemp-compose.yml up -d
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
