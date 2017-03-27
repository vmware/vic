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

*** Test Cases ***
Compose basic
    ${yml} =  Set Variable  version: "2"\nservices:\n${SPACE}web:\n${SPACE}${SPACE}image: python:2.7\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5000:5000"\n${SPACE}${SPACE}depends_on:\n${SPACE}${SPACE}- redis\n${SPACE}redis:\n${SPACE}${SPACE}image: redis\n${SPACE}${SPACE}ports:\n${SPACE}${SPACE}- "5001:5001"

    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    Run  echo '${yml}' > basic-compose.yml
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{COMPOSE-PARAMS} network create vic_default
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
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Compose up test
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/UpperCaseDir/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0  
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Abort on container exit 0
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/abort-on-container-exit-0/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Abort on container exit 1 
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/abort-on-container-exit-1/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  1
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Build path override dir
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/build-path-override-dir/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Build path
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/build-path/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Bundle with digests
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/bundle-with-digests/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Commands composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/commands-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Default env file
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/default-env-file/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Echo services
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/echo-services/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Entrypoint composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/entrypoint-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Entrypoint dockerfile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/entrypoint-dockerfile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Env file
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/env-file/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Environment composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/environment-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Environment interpolation
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/environment-interpolation/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Exit code from
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/exit-code-from/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Expose composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/expose-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Extends
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/extends/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Healthcheck
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/healthcheck/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Invalid Compose file
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/invalid-composefile/invalid.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Links composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/links-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Logging composefile legacy
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/logging-composefile-legacy/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Logging composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/logging-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Logs composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/logs-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Logs tail composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/logs-tail-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Longer filename composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/longer-filename-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Multiple composefiles
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/multiple-composefiles/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Net container
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/net-container/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Networks
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/networks/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

No links composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/no-links-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

No services
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/no-services/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Override files
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/override-files/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Ports composefile scale
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/ports-composefile-scale/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Ports composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/ports-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Restart
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/restart/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Run workdir
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/run-workdir/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Simple composefile volume ready
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/simple-composefile-volume-ready/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Simple composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/simple-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Simple dockerfile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/simple-dockerfile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Simple failing dockerfile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/simple-failing-dockerfile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Sleeps composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/sleeps-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Stop signal composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/stop-signal-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Top
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/top/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Unicode environment
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/unicode-environment/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

User composefile
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/user-composefile/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

V1 config
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/v1-config/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

V2 Dependencies
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/v2-dependencies/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

V2 full
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/v2-full/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

V2 simple
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/v2-simple/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

V3 full
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/v3-full/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Volume path interpolation
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/volume-path-interpolation/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Volume path
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/volume-path/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Volume
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/volume/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Volumes from container
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/volumes-from-container/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down

Volumes
    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} -f %{GOPATH}/src/github.com/vmware/vic/tests/compose/fixtures/volumes/docker-compose.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} down
