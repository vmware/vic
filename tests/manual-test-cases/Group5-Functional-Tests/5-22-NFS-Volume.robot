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
Documentation  Test 5-22 - NFS Volume
Resource  ../../resources/Util.robot
Suite Setup  Setup ESX And NFS Suite
#Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}


*** Variables ***
${nfsVolumeStore}=  nfsVolumeStore
${nfsFakeVolumeStore}=  nfsFakeVolumeStore
${unnamedNFSVolContainer}  unnamedNFSvolContainer
${namedNFSVolContainer}  namednfsVolContainer
${createFileContainer}=  createFileContainer
#${addToFileContainer}=  addToFileContainer
${nfs_bogon_ip}=  198.51.100.1


*** Keywords ***
Setup ESX And NFS Suite
    Log To Console  \nStarting test...

    ${esx1}  ${esx1_ip}=  Deploy Nimbus ESXi Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    ${nfs}  ${nfs_ip}=  Deploy Nimbus NFS Datastore  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    Set Global Variable  @{list}  ${esx1}  ${nfs}
    Set Global Variable  ${ESX1}  ${esx1}
    Set Global Variable  ${ESX1_IP}  ${esx1_ip}
    Set Global Variable  ${NFS_IP}  ${nfs_ip}

Setup ENV Variables for VIC Appliance Install
    Log To Console  \nSetup Environment Variables for VIC Appliance To ESX\n

    Set Environment Variable  TEST_URL_ARRAY  ${ESX1_IP}
    Set Environment Variable  TEST_URL  ${ESX1_IP}
    Set Environment Variable  TEST_USERNAME  root
    Set Environment Variable  TEST_PASSWORD  ${NIMBUS_ESX_PASSWORD}
    Set Environment Variable  TEST_DATASTORE  datastore1
    Set Environment Variable  TEST_TIMEOUT  30m
    Set Environment Variable  HOST_TYPE  ESXi
    Remove Environment Variable  TEST_DATACENTER
    Remove Environment Variable  TEST_RESOURCE
    Remove Environment Variable  BRIDGE_NETWORK
    Remove Environment Variable  PUBLIC_NETWORK

Verify NFS Volume Basic Setup
    [Arguments]  ${volumeName}  ${containerName}  ${nfsIP}  ${rwORro}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${containerName} -v ${volumeName}:/mydata ${busybox} mount
    Should Be Equal As Integers  ${rc}  0
    #Log to Console  ${outputTemp}

    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -a ${containerName}
    #Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${nfsIP}://store/volumes/${volumeName}
    Should Contain  ${output}  /mydata type nfs (${rwORro}

    ${ContainerRC}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} wait ${containerName}
    Should Be Equal As Integers  ${ContainerRC}  0
    Should Not Contain  ${output}  Error response from daemon

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${containerName}
    Should Be Equal As Integers  ${rc}  0

Verify NFS Volume Already Created
    [Arguments]  ${containerVolName}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=${containerVolName} --opt VolumeStore=${nfsVolumeStore}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: A volume named ${containerVolName} already exists. Choose a different volume name.


*** Test Cases ***
VIC Appliance Install With Fake NFS Server
    Setup ENV Variables for VIC Appliance Install

    # Will only produce a warning in VCH creation output
    ${output}=  Install VIC Appliance To Test Server  certs=${false}  additional-args=--volume-store="nfs://${nfs_bogon_ip}/store?uid=0&gid=0:${nfsFakeVolumeStore}"
    Should Contain  ${output}  VolumeStore (${nfsFakeVolumeStore}) specified was not able to be established in the portlayer. Please check network and nfs server configurations.

VIC Appliance Install With Correct NFS Server
    Setup ENV Variables for VIC Appliance Install
    Log To Console  \nDeploy VIC Appliance To ESX

    # Should succeed
    ${output}=  Install VIC Appliance To Test Server  certs=${false}  additional-args=--volume-store="nfs://${NFS_IP}/store?uid=0&gid=0:${nfsVolumeStore}"
    Should Contain  ${output}  Installer completed successfully

Simple docker volume create
    #Pull image  ${busybox}

    ${rc}  ${volumeOutput}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt VolumeStore=${nfsVolumeStore}
    Should Be Equal As Integers  ${rc}  0

    Set Suite Variable  ${nfsUnNamedVolume}  ${volumeOutput}

    Verify NFS Volume Basic Setup  ${nfsUnNamedVolume}  ${unnamedNFSVolContainer}  ${NFS_IP}  rw

Docker volume create named volume
    ${rc}  ${volumeOutput}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name nfs-volume_%{VCH-NAME} --opt VolumeStore=${nfsVolumeStore}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${volumeOutput}  nfs-volume_%{VCH-NAME}

    Set Suite Variable  ${nfsNamedVolume}  ${volumeOutput}

    Verify NFS Volume Basic Setup  nfs-volume_%{VCH-NAME}  ${namedNFSVolContainer}  ${NFS_IP}  rw

Docker volume create already named volume
    Run Keyword And Ignore Error  Verify NFS Volume Already Created  ${nfsUnNamedVolume}

    Run Keyword And Ignore Error  Verify NFS Volume Already Created  ${nfsNamedVolume}

Docker volume create with possibly invalid name
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name="test!@\#$%^&*()" --opt VolumeStore=${nfsVolumeStore}
    Should Be Equal As Integers  ${rc}  1
    Should Be Equal As Strings  ${output}  Error response from daemon: volume name "test!@\#$%^&*()" includes invalid characters, only "[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed

Docker Single Write and Read to/from File from one Container using NFS Volume
    # Done with the same contianer for this test.
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${createFileContainer} -d -v ${nfsNamedVolume}:/mydata ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${createFileContainer} sh -c "echo 'The Texas and Chile flag look similar.\n' > /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${createFileContainer} sh -c "ls mydata/"
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  test_nfs_file.txt

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${createFileContainer} sh -c "cat mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  The Texas and Chile flag look similar.

Docker multiple writes from multiple containers and read from one
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -v ${nfsNamedVolume}:/mydata ${busybox} sh -c "echo 'The Chad and Romania flag look the same.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -v ${nfsNamedVolume}:/mydata ${busybox} sh -c "echo 'The Luxembourg and the Netherlands flag look exactly the same.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -v ${nfsNamedVolume}:/mydata ${busybox} sh -c "echo 'Norway and Iceland have flags that are basically inverses of each other.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -v ${nfsNamedVolume}:/mydata ${busybox} sh -c "cat mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  The Texas and Chile flag look similar.
    Should Contain  ${output}  The Chad and Romania flag look the same.
    Should Contain  ${output}  The Luxembourg and the Netherlands flag look exactly the same.
    Should Contain  ${output}  Norway and Iceland have flags that are basically inverses of each other.

Simple docker volume inspect
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume inspect ${nfsNamedVolume}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Name
    Should Be Equal As Strings  ${id}  ${nfsNamedVolume}

Simple Volume ls test
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  vsphere
    Should Contain  ${output}  ${nfsNamedVolume}
    Should Contain  ${output}  ${nfsUnNamedVolume}
    Should Contain  ${output}  DRIVER
    Should Contain  ${output}  VOLUME NAME

Volume rm tests
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm ${nfsUnNamedVolume}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${nfsUnNamedVolume}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm ${nfsNamedVolume}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: volume ${nfsNamedVolume} in use by