# Copyright 2016-2019 VMware, Inc. All Rights Reserved.
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
Documentation  Test 5-30 - docker compose ELK
Resource  ../../resources/Util.robot
Suite Setup  Deploy Environment
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Variables ***
${NIMBUS_LOCATION}  sc
${NIMBUS_LOCATION_FULL}  NIMBUS_LOCATION=${NIMBUS_LOCATION}

*** Keywords ***
Deploy Environment
    [Timeout]  110 minutes
    Run Keyword And Ignore Error  Nimbus Cleanup  ${list}  ${false}
    ${name}=  Evaluate  'vic-iscsi-' + str(random.randint(1000,9999))  modules=random
    ${out}=  Deploy Nimbus Testbed  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  spec=vic-cluster-2esxi-iscsi.rb  args= --plugin testng --noSupportBundles --vcvaBuild ${VC_VERSION} --esxBuild ${ESX_VERSION} --testbedName vic-cluster-iscsi --runName ${name}
    Set Suite Variable  @{list}  %{NIMBUS_PERSONAL_USER}-${name}.vc.0  %{NIMBUS_PERSONAL_USER}-${name}.esx.0  %{NIMBUS_PERSONAL_USER}-${name}.esx.1  %{NIMBUS_PERSONAL_USER}-${name}.iscsi.0

    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.vc.0 | grep %{NIMBUS_PERSONAL_USER}-${name}.vc.0
    ${vc-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.esx.0 | grep %{NIMBUS_PERSONAL_USER}-${name}.esx.0
    ${esx0-ip}=  Fetch From Right  ${out}  ${SPACE}
    
    ${out}=  Execute Command  ${NIMBUS_LOCATION_FULL} USER=%{NIMBUS_PERSONAL_USER} nimbus-ctl ip %{NIMBUS_PERSONAL_USER}-${name}.esx.1 | grep %{NIMBUS_PERSONAL_USER}-${name}.esx.1
    ${esx1-ip}=  Fetch From Right  ${out}  ${SPACE}

    Set Environment Variable  GOVC_URL  ${esx0-ip}
    Set Environment Variable  GOVC_USERNAME  root
    Set Environment Variable  GOVC_PASSWORD  e2eFunctionalTest
    Run  govc host.esxcli network firewall set -e false
    Set Environment Variable  GOVC_URL  ${esx1-ip}
    Run  govc host.esxcli network firewall set -e false

    Log To Console  Set environment variables up for GOVC
    Set Environment Variable  GOVC_URL  ${vc-ip}
    Set Environment Variable  GOVC_USERNAME  Administrator@vsphere.local
    Set Environment Variable  GOVC_PASSWORD  Admin\!23

    Log To Console  Deploy VIC to the VC cluster
    Set Environment Variable  TEST_URL_ARRAY  ${vc-ip}
    Set Environment Variable  TEST_USERNAME  Administrator@vsphere.local
    Set Environment Variable  TEST_PASSWORD  Admin\!23
    Set Environment Variable  BRIDGE_NETWORK  bridge
    Set Environment Variable  PUBLIC_NETWORK  vm-network
    Remove Environment Variable  TEST_DATACENTER
    Set Environment Variable  TEST_DATASTORE  sharedVmfs-0
    Set Environment Variable  TEST_RESOURCE  cls
    Set Environment Variable  TEST_TIMEOUT  30m

Check Service State
    [Arguments]  ${containerName}
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} ${containerName}
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' ${containerName}
    Log  ${out}
    Should Not Be Empty  ${out}
    Should Be Equal As Integers  ${rc}  0

*** Test Cases ***
Compose ELK Test
    Install VIC Appliance To Test Server  certs=${true}  vol=default  additional-args=--cpu-reservation 1 --cpu-shares normal --memory-reservation 1 --memory-shares normal --endpoint-cpu 1 --endpoint-memory 2048 --base-image-size 8GB --bridge-network-range 172.17.0.0/12 --container-network-firewall vm-network:published --certificate-key-size 2048

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --skip-hostname-check -f ${CURDIR}/../../../demos/compose/elk-app/docker-compose-elk.yml up -d
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0
    
    Check Service State  elastic1
    Check Service State  elastic2
    Check Service State  elastic3
    Check Service State  logstash1
    Check Service State  logstash2
    Check Service State  kibana1
    Check Service State  kibana2


    Sleep  10m
    Check Service State  elastic1
    Check Service State  elastic2
    Check Service State  elastic3
    Check Service State  logstash1
    Check Service State  logstash2
    Check Service State  kibana1
    Check Service State  kibana2

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose %{COMPOSE-PARAMS} --skip-hostname-check -f ${CURDIR}/../../../demos/compose/elk-app/docker-compose-elk.yml down
    Log  ${out}
    Should Be Equal As Integers  ${rc}  0    

    Cleanup VIC Appliance On Test Server

