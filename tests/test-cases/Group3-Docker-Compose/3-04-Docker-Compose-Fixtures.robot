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
Documentation  Test 3-03 - Docker Compose Keywords.
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${True}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Command
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/command/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/command/docker-compose.yml down

Container Name
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/container_name/docker-compose.yml up
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  my-web-container exited with code 0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/container_name/docker-compose.yml down

Depends On
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/depends_on/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/depends_on/docker-compose.yml down

Env File
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/env-file/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/env-file/docker-compose.yml down

Environment
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/environment-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/environment-composefile/docker-compose.yml down

Extends
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/extends/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/extends/docker-compose.yml down

Group Add
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/group_add/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/group_add/docker-compose.yml down

Labels
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/labels/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/labels/docker-compose.yml down

Links
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/links-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/links-composefile/docker-compose.yml down

Networks
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/docker-compose.yml down

Networks- Default
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/default-network-config.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/default-network-config.yml down

Networks- External Default
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{COMPOSE-PARAMS} network create composetest_external_network
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/external-default.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/external-default.yml down

Networks-Bridge
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/bridge.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/bridge.yml down

Networks-External
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{COMPOSE-PARAMS} network create networks_foo
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{COMPOSE-PARAMS} network create networks_bar
    Should Be Equal As Integers  ${rc}  0    

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/external-networks.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/external-networks.yml down

Networks-Label
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-label.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-label.yml down

Networks-mode
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-mode.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-mode.yml down

Networks-static-address
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-static-addresses.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/networks/network-static-addresses.yml down

Ports
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/ports-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/ports-composefile/docker-compose.yml down

Stop Signal
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/stop-signal-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/stop-signal-composefile/docker-compose.yml down

Volumes
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/volumes/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/configs/volumes/docker-compose.yml down


