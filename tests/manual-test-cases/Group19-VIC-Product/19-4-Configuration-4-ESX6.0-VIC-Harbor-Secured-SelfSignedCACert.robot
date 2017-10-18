# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 19-4 - Configuration 4 ESX 6.0 VIC Harbor Secured Self Signed
Resource  ../../resources/Util.robot
Suite Setup  19-4-Setup
Suite Teardown  19-4-Teardown

*** Variables ***
${developer}  test-user-developer
${guest}  test-user-guest
${password}  Test-password-1
${newPassword}  Test-new-password-1
${project}  vic-harbor
${image}  busybox
${container_name}  busy
${tag}  new
${developerRole}  Developer
${guestRole}  Guest
${developerEmail}  developer@test.com
${developerEmail2}  developer2@test.com
${guestEmail}  guest@test.com
${developerFullName}  Vic Developer
${guestFullName}  Vic Guest
${comments}  comments
${developer2}  test-user-developer2

*** Keywords ***
19-4-Setup
    Set Environment Variable  ESX_VERSION  3620759
    ${esx1}  ${esx1-ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}
    Set Environment Variable  TEST_URL  ${esx1-ip}
    Set Environment Variable  TEST_USERNAME  root
    Set Environment Variable  TEST_PASSWORD  e2eFunctionalTest
    Set Environment Variable  BRIDGE_NETWORK  VM Network
    Set Environment Variable  PUBLIC_NETWORK  "VM Network"
    Set Environment Variable  TEST_DATACENTER  ${EMPTY}
    Set Environment Variable  TEST_RESOURCE  ${EMPTY}
    Set Environment Variable  TEST_TIMEOUT  30m
    
    Set Suite Variable  ${ESX1}  ${esx1}
    Set Suite Variable  ${ESX1-IP}  ${esx1-ip}
    Set Global Variable  @{list}  ${esx1}

    Install Harbor To Test Server  name=19-4-harbor  protocol=https
    Log To Console  Harbor installer completed successfully...
    
    Install Harbor Self Signed Cert
    Install VIC Appliance To Test Server  vol=default --registry-ca=/etc/docker/certs.d/%{HARBOR_IP}/ca.crt  certs=${false}
    Create Project And Three Users For It  developer=${developer}  developer2=${developer2}  developerEmail=${developerEmail}  developerEmail2=${developerEmail2}  developerFullName=${developerFullName}  password=${password}  comments=${comments}  guest=${guest}  developerRole=${developerRole}  guestRole=${guestRole}  project=${project}  public=${False}
    Remove Environment Variable  DOCKER_HOST

19-4-Teardown
    Run Keyword And Continue On Failure  Cleanup VIC Appliance On Test Server
    ${out}=  Run Keyword And Continue On Failure  Run  govc vm.destroy 19-4-harbor
    Run Keyword And Continue On Failure  Nimbus Cleanup  ${list}

*** Test Cases ***
Test Pos001 Admin Operations
    Basic Docker Command With Harbor  user=admin  password=%{TEST_PASSWORD}  project=${project}  image=${image}  container_name=${container_name}

Test Pos002 Developer Operations
    Basic Docker Command With Harbor  user=${developer}  password=${password}  project=${project}  image=${image}  container_name=${container_name}
    
Test Neg001 Developer Operations
    # Docker login
    Log To Console  \nRunning docker login admin...
    ${rc}  ${output}=  Run And Return Rc And Output  docker login -u ${guest} -p ${password} %{HARBOR_IP}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Login Succeeded
    Should Not Contain  ${output}  Error response from daemon

    # Docker pull from harbor through VCH, ensure guest could pull
    Log To Console  docker pull from harbor using VCH...
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker pull from dockerhub
    Log To Console  docker pull from dockerhub...
    ${rc}  ${output}=  Run And Return Rc And Output  docker pull %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker tag image
    Log To Console  docker tag...
    ${rc}  ${output}=  Run And Return Rc And Output  docker tag %{HARBOR_IP}/${project}/${image} %{HARBOR_IP}/${project}/${image}:${tag}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker push image should fail
    Log To Console  push image...
    ${rc}  ${output}=  Run And Return Rc And Output  docker push %{HARBOR_IP}/${project}/${image}:${tag}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  unauthorized: authentication required
    Should Not Contain  ${output}  digest:

Test Pos003 Two VCH With One Harbor
    # Add another VC
    ${VCH1-URL}=  Set Variable  %{VCH-PARAMS}

    Install VIC Appliance To Test Server  vol=default --registry-ca=/etc/docker/certs.d/%{HARBOR_IP}/ca.crt  certs=${false}  cleanup=${false}
    Remove Environment Variable  DOCKER_HOST
    ${VCH2-URL}=  Set Variable  %{VCH-PARAMS}

    # Docker login VCH1
    Log To Console  \nRunning docker login developer VCH1...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH1-URL} login -u ${developer} -p ${password} %{HARBOR_IP}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Login Succeeded
    Should Not Contain  ${output}  Error response from daemon

    # Docker login VCH2
    Log To Console  \nRunning docker login developer2 VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} login -u ${developer2} -p ${password} %{HARBOR_IP}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Login Succeeded
    Should Not Contain  ${output}  Error response from daemon

    # Docker pull from dockerhub
    Log To Console  docker pull from dockerhub...
    ${rc}  ${output}=  Run And Return Rc And Output  docker pull ${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker tag image
    Log To Console  docker tag...
    ${rc}  ${output}=  Run And Return Rc And Output  docker tag ${image} %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker push image
    Log To Console  push image...
    ${rc}  ${output}=  Run And Return Rc And Output  docker push %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  digest:
    Should Contain  ${output}  latest:
    Should Not Contain  ${output}  No such image:

    # Docker delete image in local registry
    Log To Console  docker rmi local registry...
    ${rc}  ${output}=  Run And Return Rc And Output  docker rmi -f %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Untagged

    # Make sure the image pushed by VCH1 could be used by VCH2
    # Docker pull from harbor using VCH2
    Log To Console  docker pull from dockerhub VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} pull %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker run image
    Log To Console  docker run VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} run --name ${container_name} %{HARBOR_IP}/${project}/${image} /bin/ash -c "dmesg;echo END_OF_THE_TEST" 
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  END_OF_THE_TEST

    # Docker rm container
    Log To Console  docker rm VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} rm -f ${container_name}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker create
    Log To Console  docker create VCH2...
    ${rc}  ${containerID}=  Run And Return Rc And Output  docker ${VCH2-URL} create --name ${container_name} -i %{HARBOR_IP}/${project}/${image} /bin/top
    Log  ${containerID}
    Should Be Equal As Integers  ${rc}  0

    # Docker start
    Log To Console  docker start VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} start ${container_name}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker attach 
    Log To Console  Starting process Docker attach VCH2...
    Start Process  docker ${VCH2-URL} attach ${container_name} < /tmp/fifo  shell=True  alias=custom
    Sleep  3
    Run  echo q > /tmp/fifo
    ${ret}=  Wait For Process  custom
    Log  ${ret}
    Should Be Equal As Integers  ${ret.rc}  0
    Should Be Empty  ${ret.stdout}
    Should Be Empty  ${ret.stderr}

    # Docker start  
    Log To Console  docker start VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} start ${container_name}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker stop
    Log To Console  docker stop VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} stop ${container_name}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    # Docker remove
    Log To Console  docker rm VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} rm -f ${container_name}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Keyword Succeeds  10x  6s  Check That VM Is Removed  ${container_name}
    Wait Until Keyword Succeeds  10x  6s  Check That Datastore Is Cleaned  ${container_name}

    # Docker delete image
    Log To Console  docker rmi VCH2...
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${VCH2-URL} rmi -f %{HARBOR_IP}/${project}/${image}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Untagged

# Test Pos004 Three Client Machines With One Harbor
