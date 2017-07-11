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
Suite Setup  Setup NFS Server
Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}


*** Variables ***
${createFileContainer}=  createFileContainer
${addToFileContainer}=  addToFileContainer


*** Keywords ***
Setup NFS Server
    Log To Console  \nStarting test...
    ${esx3}  ${esx4}  ${esx5}  ${vc}  ${esx3-ip}  ${esx4-ip}  ${esx5-ip}  ${vc-ip}=  Create a Simple VC Cluster  datacenter1  cls1

    ${nfs_default_name}  ${nfs_default_ip}=  Deploy Nimbus NFS Datastore  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}

    Set Global Variable  @{list}  ${esx3}  ${esx4}  ${esx5}  ${vc}  ${nfs_default_name}

    Set Suite Variable  ${NFS_DEFAULT_IP}  ${nfs_default_ip}
    Set Suite Variable  ${NFS_DEFAULT_NAME}  ${nfs_default_name}

    #Need to set nfs volume at vic creation time
    Install VIC Appliance To Test Server  additional-args=--volume-store="nfs://${NFS_DEFAULT_IP}/store?uid=0&gid=0:nfsVolumeStoreDefault"

Verify NFS Volume Basic Setup
    [Arguments]  ${prevOutput}  ${ContainerName}  ${nfsIP}  ${rwORro}

    ${rc}  ${outputTemp}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${ContainerName} -d -v ${prevOutput}:/mydata ${busybox} mount
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start -a ${ContainerName}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${nfsIP}://store/volumes/${prevOutput}
    Should Contain  ${output}  /mydata type nfs (${rwORro}

    ${ContainerRC}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} wait ${ContainerName}
    Should Be Equal As Integers  ${ContainerRC}  0
    Should Not Contain  ${output}  Error response from daemon

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm ${outputTemp}
    Should Be Equal As Integers  ${rc}  0

Verify NFS Volume Already Created
    [Arguments]  ${containerVolName}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=${containerVolName} --opt VolumeStore=nfsVolumeStoreDefault
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: A volume named ${containerVolName} already exists. Choose a different volume name.


*** Test Cases ***
Simple docker volume create
    Pull image  ${busybox}

    ${rc}  ${outputDefault}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt VolumeStore=nfsVolumeStoreDefault
    Should Be Equal As Integers  ${rc}  0

    Set Suite Variable  ${UnnamedNFSVolContainer}  unnamednfsVolContainer
    Set Suite Variable  ${nfsDefaultVolume}  ${outputDefault}

    Verify NFS Volume Basic Setup  ${nfsDefaultVolume}  ${UnnamedNFSVolContainer}  ${NFS_DEFAULT_IP}  rw

Docker volume create named volume
    ${rc}  ${outputDefault}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name nfs_default_%{VCH-NAME} --opt VolumeStore=nfsVolumeStoreDefault
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Strings  ${outputDefault}  nfs_default_%{VCH-NAME}

    Set Suite Variable  ${namedNFSVolContainer}  namednfsVolContainer
    Set Suite Variable  ${nfsNamedVolume}  ${outputDefault}

    Verify NFS Volume Basic Setup  ${nfsNamedVolume}  ${namedNFSVolContainer}  ${NFS_DEFAULT_IP}  rw

Docker volume create already named volume
    Verify NFS Volume Already Created  ${nfsDefaultVolume}

    Verify NFS Volume Already Created  ${nfsNamedVolume}

Docker volume create with possibly invalid name
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=test??? --opt VolumeStore=nfsVolumeStoreDefault
    Should Be Equal As Integers  ${rc}  1
    Should Be Equal As Strings  ${output}  Error response from daemon: volume name "test???" includes invalid characters, only "[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed

Docker single write to file from one container in named volume
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

Docker multiple write from multiple containers and read from one
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name ${addToFileContainer} -d -v ${nfsNamedVolume}:/mydata ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${addToFileContainer} sh -c "echo 'The Chad and Romania flag look the same.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${createFileContainer} sh -c "echo 'The Luxembourg and the Netherlands flag look exactly the same.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${addToFileContainer} sh -c "cat mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  The Chad and Romania flag look the same.
    Should Contain  ${output}  The Luxembourg and the Netherlands flag look exactly the same.

Docker write from one container and read from another in named volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${addToFileContainer} sh -c "echo 'Norway and Iceland have flags that are basically inverses of each other.\n' >> /mydata/test_nfs_file.txt"
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${createFileContainer} sh -c "cat mydata/test_nfs_file.txt"
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
    Should Contain  ${output}  ${nfsDefaultVolume}
    Should Contain  ${output}  DRIVER
    Should Contain  ${output}  VOLUME NAME

Volume rm tests
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm ${nfsDefaultVolume}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  ${nfsDefaultVolume}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm ${nfsNamedVolume}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error response from daemon: volume ${nfsNamedVolume} in use by