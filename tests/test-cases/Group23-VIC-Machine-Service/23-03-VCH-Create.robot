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
Documentation     Test 23-03 - VCH Create
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Setup
Suite Teardown    Teardown
Default Tags


*** Keywords ***
Setup
    Start VIC Machine Server
    Set Test Environment Variables


Teardown
    Terminate All Processes    kill=True


Create VCH
    [Arguments]    ${data}
    Post Path Under Target    vch    ${data}


Create VCH Within Datacenter
    [Arguments]    ${data}
    ${dcID}=    Get Datacenter ID
    Post Path Under Target    datacenter/${dcID}/vch    ${data}


Inspect VCH ${name}
    ${RC}    ${OUTPUT}=    Run And Return Rc And Output    bin/vic-machine-linux inspect config --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user %{TEST_USERNAME} --password=%{TEST_PASSWORD} --name=${name} --format raw
    Should Be Equal As Integers    ${RC}    0
    Set Test Variable    ${OUTPUT}


Arguments Should Contain
    [Arguments]    ${expected}
    Should Contain    ${OUTPUT}    ${expected}


Get VCH
    [Arguments]    ${name}

    Get Path Under Target    vch
    ${id}=    Run    echo '${OUTPUT}' | jq -r '.vchs[] | select(.name=="${name}").id'

    Get Path Under Target    vch/${id}


Property Should Be Equal
    [Arguments]    ${jq}    ${expected}

    ${actual}=  Run    echo '${OUTPUT}' | jq -r '${jq}'
    Should Be Equal    ${actual}    ${expected}


Property Should Not Be Equal
    [Arguments]    ${jq}    ${expected}

    ${actual}=  Run    echo '${OUTPUT}' | jq -r '${jq}'
    Should Not Be Equal    ${actual}    ${expected}


Property Should Contain
    [Arguments]    ${jq}    ${expected}

    ${actual}=  Run    echo '${OUTPUT}' | jq -r '${jq}'
    Should Contain    ${actual}    ${expected}


*** Test Cases ***
Create minimal VCH
    Create VCH    '{"name":"%{VCH-NAME}-api-test-minimal","compute":{"resource":{"name":"%{TEST_RESOURCE}"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"]},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    201


    Inspect VCH %{VCH-NAME}-api-test-minimal

    Arguments Should Contain    --image-store=ds://%{TEST_DATASTORE}
    Arguments Should Contain    --bridge-network=%{BRIDGE_NETWORK}


    Get VCH    %{VCH-NAME}-api-test-minimal

    Property Should Be Equal        .name                                %{VCH-NAME}-api-test-minimal

    Property Should Not Be Equal    .compute.resource.id                 null

    Property Should Contain         .storage.image_stores[0]             %{TEST_DATASTORE}
    Property Should Be Equal        .storage.base_image_size.value       8000000
    Property Should Be Equal        .storage.base_image_size.units       KB

    Property Should Contain         .auth.server.certificate.pem         -----BEGIN CERTIFICATE-----
    Property Should Be Equal        .auth.server.private_key.pem         null

    Property Should Contain         .network.bridge.ip_range             172.16.0.0/12

    Property Should Contain         .runtime.power_state                 poweredOn
    Property Should Contain         .runtime.upgrade_status              Up to date

    [Teardown]    Run Secret VIC Machine Delete Command    %{VCH-NAME}-api-test-minimal


Create minimal VCH within datacenter
    Create VCH Within Datacenter    '{"name":"%{VCH-NAME}-api-test-dc","compute":{"resource":{"name":"%{TEST_RESOURCE}"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"]},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    201


    Inspect VCH %{VCH-NAME}-api-test-dc

    Arguments Should Contain    --image-store=ds://%{TEST_DATASTORE}
    Arguments Should Contain    --bridge-network=%{BRIDGE_NETWORK}


    Get VCH    %{VCH-NAME}-api-test-dc

    Property Should Be Equal        .name                                %{VCH-NAME}-api-test-dc

    Property Should Not Be Equal    .compute.resource.id                 null

    Property Should Contain         .storage.image_stores[0]             %{TEST_DATASTORE}
    Property Should Be Equal        .storage.base_image_size.value       8000000
    Property Should Be Equal        .storage.base_image_size.units       KB

    Property Should Contain         .auth.server.certificate.pem         -----BEGIN CERTIFICATE-----
    Property Should Be Equal        .auth.server.private_key.pem         null

    Property Should Contain         .network.bridge.ip_range             172.16.0.0/12

    Property Should Contain         .runtime.power_state                 poweredOn
    Property Should Contain         .runtime.upgrade_status              Up to date

    [Teardown]    Run Secret VIC Machine Delete Command    %{VCH-NAME}-api-test-dc


Create complex VCH
    Create VCH    '{"name":"%{VCH-NAME}-api-test-complex","debug":3,"compute":{"cpu":{"limit":{"units":"MHz","value":2345},"reservation":{"units":"GHz","value":2},"shares":{"level":"high"}},"memory":{"limit":{"units":"MiB","value":1200},"reservation":{"units":"MiB","value":501},"shares":{"number":81910}},"resource":{"name":"%{TEST_RESOURCE}"}},"endpoint":{"cpu":{"sockets":2},"memory":{"units":"MiB","value":3072}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"],"volume_stores":[{"datastore":"ds://%{TEST_DATASTORE}/test-volumes/foo","label":"foo"}],"base_image_size":{"units":"B","value":16000000}},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"registry":{"image_fetch_proxy":{"http":"http://example.com","https":"https://example.com"},"insecure":["https://insecure.example.com"],"whitelist":["10.0.0.0/8"]},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}},"syslog_addr":"tcp://syslog.example.com:4444"}'

    Verify Return Code
    Verify Status    201


    Inspect VCH %{VCH-NAME}-api-test-complex

    Arguments Should Contain    --debug=3

    Arguments Should Contain    --cpu=2345
    Arguments Should Contain    --cpu-reservation=2000
    Arguments Should Contain    --cpu-shares=high
    Arguments Should Contain    --memory=1200
    Arguments Should Contain    --memory-reservation=501
    Arguments Should Contain    --memory-shares=81910

    Arguments Should Contain    --endpoint-cpu=2
    Arguments Should Contain    --endpoint-memory=3072

    Arguments Should Contain    --image-store=ds://%{TEST_DATASTORE}
    Arguments Should Contain    --volume-store=ds://%{TEST_DATASTORE}/test-volumes/foo:foo
    Arguments Should Contain    --base-image-size=16MB

    Arguments Should Contain    --bridge-network=%{BRIDGE_NETWORK}

    Arguments Should Contain    --insecure-registry=https://insecure.example.com
    Arguments Should Contain    --whitelist-registry=10.0.0.0/8
    Arguments Should Contain    --whitelist-registry=https://insecure.example.com


    Get VCH    %{VCH-NAME}-api-test-complex

    Property Should Be Equal        .name                                %{VCH-NAME}-api-test-complex
    Property Should Be Equal        .debug                               3
    Property Should Be Equal        .syslog_addr                         tcp://syslog.example.com:4444

    Property Should Not Be Equal    .compute.resource.id                 null
    Property Should Be Equal        .compute.cpu.limit.value             2345
    Property Should Be Equal        .compute.cpu.limit.units             MHz
    Property Should Be Equal        .compute.cpu.reservation.value       2000
    Property Should Be Equal        .compute.cpu.reservation.units       MHz
    Property Should Be Equal        .compute.cpu.shares.level            high
    Property Should Be Equal        .compute.memory.limit.value          1200
    Property Should Be Equal        .compute.memory.limit.units          MiB
    Property Should Be Equal        .compute.memory.reservation.value    501
    Property Should Be Equal        .compute.memory.reservation.units    MiB
    Property Should Be Equal        .compute.memory.shares.number        81910

    Property Should Be Equal        .endpoint.cpu.sockets                2
    Property Should Be Equal        .endpoint.memory.value               3072
    Property Should Be Equal        .endpoint.memory.units               MiB

    Property Should Contain         .storage.image_stores[0]             %{TEST_DATASTORE}
    Property Should Contain         .storage.volume_stores[0].datastore  %{TEST_DATASTORE}/test-volumes/foo
    Property Should Contain         .storage.volume_stores[0].label      foo
    Property Should Be Equal        .storage.base_image_size.value       16000
    Property Should Be Equal        .storage.base_image_size.units       KB

    Property Should Be Equal        .registry.image_fetch_proxy.http     http://example.com
    Property Should Be Equal        .registry.image_fetch_proxy.https    https://example.com
    Property Should Contain         .registry.insecure | join(" ")       https://insecure.example.com
    Property Should Contain         .registry.whitelist | join(" ")      https://insecure.example.com
    Property Should Contain         .registry.whitelist | join(" ")      10.0.0.0/8

    Property Should Contain         .auth.server.certificate.pem         -----BEGIN CERTIFICATE-----
    Property Should Be Equal        .auth.server.private_key.pem         null

    Property Should Be Equal        .network.bridge.ip_range             172.16.0.0/12

    Property Should Be Equal        .runtime.power_state                 poweredOn
    Property Should Be Equal        .runtime.upgrade_status              Up to date

    [Teardown]    Run Secret VIC Machine Delete Command    %{VCH-NAME}-api-test-complex


Fail to create VCH with invalid operations credentials
    Create VCH    '{"name":"%{VCH-NAME}-api-bad-ops","compute":{"resource":{"name":"%{TEST_RESOURCE}"}},"endpoint":{"operations_credentials":{"user":"invalid","password":"invalid"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"]},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    400

    Should Contain    ${output}    operations credentials


Fail to create VCH with invalid datastore
    Create VCH    '{"name":"%{VCH-NAME}-api-bad-storage","compute":{"resource":{"name":"%{TEST_RESOURCE}"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}-invalid"]},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    400

    Should Contain    ${output}    datastore


Fail to create VCH with invalid compute
    Create VCH    '{"name":"%{VCH-NAME}-api-bad-compute","compute":{"resource":{"name":"%{TEST_RESOURCE}-invalid"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"]},"network":{"bridge":{"ip_range":"172.16.0.0/12","port_group":{"name":"%{BRIDGE_NETWORK}"}},"public":{"port_group":{"name":"%{PUBLIC_NETWORK}"}}},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    400

    Should Contain    ${output}    compute resource


Fail to create VCH without network
    Create VCH    '{"name":"%{VCH-NAME}-api-bad-network","compute":{"resource":{"name":"%{TEST_RESOURCE}"}},"storage":{"image_stores":["ds://%{TEST_DATASTORE}"]},"auth":{"server":{"generate":{"cname":"vch.example.com","organization":["VMware, Inc."],"size":{"value":2048,"units":"bits"}}},"client":{"no_tls_verify": true}}}'

    Verify Return Code
    Verify Status    400

    Should Contain    ${output}    network
