# Copyright 2018 VMware, Inc. All Rights Reserved.
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
Documentation     Test 23-04 - VCH Inspect
Resource          ../../resources/Util.robot
Resource          ../../resources/Group23-VIC-Machine-Service-Util.robot
Suite Setup       Setup
Suite Teardown    Teardown

Default Tags

*** Keywords ***
# TODO [AngieCris]: Uncomment it after ops user is enabled on CI (after https://github.com/vmware/vic/pull/7892 merges)
#Install VIC Appliance With Ops Credentials
#    Install VIC Appliance To Test Server  certs=${true}  additional-args=--ops-user=%{VCH_OPS_USERNAME} --ops-password=%{VCH_OPS_PASSWORD} --ops-grant-perms --debug 1

Setup
    Start VIC Machine Server
#    Install VIC Appliance With Ops Credentials
    Install VIC Appliance To Test Server

    ${id}=  Get VCH ID  %{VCH-NAME}
    ${dc-id}=  Get Datacenter ID

    Set Suite Variable  ${VCH-ID}  ${id}
    Set Suite Variable  ${DC-ID}  ${dc-id}


Teardown
    Stop VIC Machine Server
    Cleanup VIC Appliance On Test Server

Inspect VCH
    Get Path Under Target  vch/${VCH-ID}

Inspect VCH Using Session
    Get Path Under Target Using Session  vch/${VCH-ID}

Inspect VCH Within Datacenter
    Get Path Under Target    datacenter/${DC-ID}/vch/${VCH-ID}

Inspect VCH Within Datacenter Using Session
    Get Path Under Target Using Session    datacenter/${DC-ID}/vch/${VCH-ID}

Verify VCH Inspect Output
    # basic
    Property Should Be Equal        .debug                                                1
    Property Should Be Equal        .name                                                 %{VCH-NAME}

    # networks
    Property Should Be Equal        .network.bridge.ip_range                              172.16.0.0/12
    Property Should Not Be Empty    .network.bridge.port_group
    Property Should Be Equal        .network.public.nameservers[0]                        null
    Property Should Not Be Empty    .network.public.port_group
    Property Should Be Equal        .network.container[0].alias                           public
    Property Should Not Be Empty    .network.container[0].port_grou

    # cert
    ${domain}=  Get Environment Variable  DOMAIN  ''
    Run Keyword If  $domain != ''   Property Should Contain    .auth.server.certificate.pem  -----BEGIN CERTIFICATE-----
    Property Should Be Equal        .auth.server.private_key.pem                           null

    # compute
    Property Should Not Be Empty    .compute.resource.id

    # storage
    Property Should Be Equal        .storage.base_image_size.value                         8000000
    Property Should Be Equal        .storage.base_image_size.units                         KB

    Property Length Should Be       .storage.image_stores                                  1
    Property Should Contain         .storage.image_stores[0]                               %{TEST_DATASTORE}
    Property Length Should Be       .storage.volume_stores                                 1
    Property Should Contain         .storage.volume_stores[0].datastore                    %{TEST_DATASTORE}/%{VCH-NAME}-VOL
    Property Should Be Equal        .storage.volume_stores[0].label                        default

    # TODO [AngieCris]: uncomment this after #7892 merges
#    # ops creds
#    Property Should Be Equal        .endpoint.operations_credentials.user                  %{VCH_OPS_USERNAME}
#    Property Should Be Equal        .endpoint.operations_credentials.password              null
#    Property Should Be Equal        .endpoint.operations_credentials.grant_permissions     true

    # connection
    Property Should Be Equal        .runtime.docker_host                                   %{DOCKER_HOST}
    Property Should Be Equal        .runtime.admin_portal                                  %{VIC-ADMIN}
    Property Should Be Equal        .runtime.power_state                                   poweredOn
    Property Should Contain         .runtime.upgrade_status                                Up to date

    # version
    ${version}=  Get Service Version String
    Property Should Be Equal        .version                                               ${version}


*** Test Cases ***
Get VCH
    Inspect VCH

    Verify Return Code
    Verify Status Ok
    Verify VCH Inspect Output


Get VCH Using Session
    Inspect VCH Using Session

    Verify Return Code
    Verify Status Ok
    Verify VCH Inspect Output


Get VCH Within Datacenter
    Inspect VCH Within Datacenter

    Verify Return Code
    Verify Status Ok
    Verify VCH Inspect Output


Get VCH Within Datacenter Using Session
    Inspect VCH Within Datacenter Using Session

    Verify Return Code
    Verify Status Ok
    Verify VCH Inspect Output


Get VCH Within Invalid Datacenter
    Get Path Under Target    datacenter/INVALID/vch/${VCH-ID}

    Verify Return Code
    Verify Status Not Found


Get Invalid VCH
    Get Path Under Target  /vch/INVALID

    Verify Return Code
    Verify Status Not Found


Get Invalid VCH Within Datacenter
    Get Path Under Target  /datacenter/${DC-ID}/vch/INVALID
