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
Documentation  Test 11-03 - Upgrade-InsecureRegistry
Resource  ../../resources/Util.robot

*** Variables ***
${test_vic_version}  7315
${vic_success}  Installer completed successfully
${docker_bridge_network}  bridge
${docker_daemon_default_port}  2375
${http_harbor_name}  integration-test-harbor-http
${https_harbor_name}  integration-test-harbor-https
${default_local_docker_endpoint}  unix:///var/run/docker-local.sock

*** Keywords ***
Setup Test Environment
    [arguments]  ${insecure_registry}
    ${handle}  ${docker_daemon_pid}=  Start Docker Daemon Locally  --insecure-registry ${insecure_registry}
    Setup VCH And Registry  ${insecure_registry}
    [Return]  ${handle}  ${docker_daemon_pid}

Add Project On Registry
    [Tags]  secret
    [Arguments]  ${registry_ip}  ${protocol}  ${user}=admin  ${password}=%{TEST_PASSWORD}
    # Harbor API: https://github.com/vmware/harbor/blob/master/docs/swagger.yaml
    Run Keyword If  '${protocol}' == 'https'  Set Test Variable  ${curl_option}  --insecure
    Run Keyword If  '${protocol}' == 'http'  Set Test Variable  ${curl_option}  ${EMPTY}
    :FOR  ${i}  IN RANGE  12
    \   ${rc}=  Run And Return Rc  curl ${curl_option} -u ${user}:${password} -H "Content-Type: application/json" -X POST -d '{"project_name": "test","public": 1}' ${protocol}://${registry_ip}/api/projects
    \   Return From Keyword If  '${rc}' == '0'
    \   Sleep  10s
    Fail  Failed to add project on registry!

Setup VCH And Registry
    [Tags]  secret
    [Arguments]  ${registry_ip}  ${registry_user}=admin  ${registry_password}=%{TEST_PASSWORD}  ${docker}=DOCKER_API_VERSION=1.23 docker
    ${rc}  ${output}=  Run And Return Rc And Output  echo "From busybox" | ${docker} -H ${default_local_docker_endpoint} build -t ${registry_ip}/test/busybox -
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Successfully built
    Log To Console  \nbusybox built successfully
    ${rc}  ${output}=  Run And Return Rc And Output  ${docker} -H ${default_local_docker_endpoint} login --username ${registry_user} --password ${registry_password} ${registry_ip}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Login Succeeded
    Log To Console  \nLogin successfully
    ${rc}=  Run And Return Rc  ${docker} -H ${default_local_docker_endpoint} push ${registry_ip}/test/busybox
    Should Be Equal As Integers  ${rc}  0
    Log To Console  \nbusybox pushed successfully

Test VCH And Registry
    [Arguments]  ${vch_endpoint}  ${registry_ip}  ${docker}=DOCKER_API_VERSION=1.23 docker
    ${rc}  ${output}=  Run And Return Rc And Output  ${docker} -H ${vch_endpoint} pull ${registry_ip}/test/busybox
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Digest:
    Should Contain  ${output}  Status:
    Should Not Contain  ${output}  Error response from daemon

Cleanup Test Environment
    [Arguments]  ${handle}  ${docker_daemon_pid}  ${harbor_name}  ${harbor_ip}  ${docker}=DOCKER_API_VERSION=1.23 docker
    ${rc}=  Run And Return Rc  ${docker} -H ${default_local_docker_endpoint} rmi ${harbor_ip}/test/busybox
    Should Be Equal As Integers  ${rc}  0
    Kill Local Docker Daemon  ${handle}  ${docker_daemon_pid}
    Cleanup Harbor  ${harbor_name}
    Clean up VIC Appliance And Local Binary

*** Test Cases ***
Upgrade VCH with Harbor On HTTP
    ${status}=  Get State Of Github Issue  4999
    Run Keyword If  '${status}' == 'closed'  Fail  Test 11-03-Upgrade-InsecureRegistry.robot needs to be updated now that Issue #4999 has been resolved
    #${harbor_ip}=  Install Harbor To Test Server  ${http_harbor_name}
    #Add Project On Registry  ${harbor_ip}  http
    #${handle}  ${docker_daemon_pid}=  Setup Test Environment  ${harbor_ip}

    #Install VIC with version to Test Server  ${test_vic_version}  --insecure-registry ${harbor_ip} --no-tls

    #Test VCH And Registry  %{VCH-IP}:%{VCH-PORT}  ${harbor_ip}

    #Upgrade
    #Check Upgraded Version
    #Test VCH And Registry  %{VCH-IP}:%{VCH-PORT}  ${harbor_ip}

    #[Teardown]  Cleanup Test Environment  ${handle}  ${docker_daemon_pid}  ${http_harbor_name}  ${harbor_ip}

Upgrade VCH with Harbor On HTTPS
    ${status}=  Get State Of Github Issue  4999
    Run Keyword If  '${status}' == 'closed'  Fail  Test 11-03-Upgrade-InsecureRegistry.robot needs to be updated now that Issue #4999 has been resolved
    #${harbor_ip}=  Install Harbor To Test Server  ${https_harbor_name}  https
    #Add Project On Registry  ${harbor_ip}  https
    #${handle}  ${docker_daemon_pid}=  Setup Test Environment  ${harbor_ip}

    #${harbor_cert}=  Fetch Harbor Self Signed Cert  ${harbor_ip}
    #Install VIC with version to Test Server  ${test_vic_version}  --insecure-registry ${harbor_ip} --no-tls --registry-ca ${harbor_cert}

    #Test VCH And Registry  %{VCH-IP}:%{VCH-PORT}  ${harbor_ip}

    #Upgrade
    #Check Upgraded Version
    #Test VCH And Registry  %{VCH-IP}:%{VCH-PORT}  ${harbor_ip}

    #[Teardown]  Cleanup Test Environment  ${handle}  ${docker_daemon_pid}  ${https_harbor_name}  ${harbor_ip}
