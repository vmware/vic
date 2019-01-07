# Copyright 2016-2019 VMware, Inc. All Rights Reserved.
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
Documentation  Test 11-05 - Upgrade from 1.4.3 to latest for testing photon version
Resource  ../../resources/Util.robot
Suite Setup  Disable Ops User And Install VIC To Test Server
Suite Teardown  Re-Enable Ops User And Clean Up VIC Appliance
Default Tags

*** Variables ***
${volume1}=  volume1
${volume2}=  volume2
${volume3}=  volume3
${nginx_test1}=  nginx-test1
${nginx_test2}=  nginx-test2
${nginx_test3}=  nginx-test3
${mydata}=  /mydata
${port}=  80
${run-as-ops-user}=  ${EMPTY}
${old_version}=  v1.4.3

*** Keywords ***
Disable Ops User And Install VIC To Test Server
    ${run-as-ops-user}=  Get Environment Variable  RUN_AS_OPS_USER  0
    Set Environment Variable  RUN_AS_OPS_USER  0
    Install VIC with version to Test Server  ${old_version}  additional-args=--cpu-reservation 1 --cpu-shares normal --memory-reservation 1 --memory-shares normal --endpoint-cpu 1 --endpoint-memory 2048 --base-image-size 8GB --volume-store %{TEST_DATASTORE}/VCH-${old_version}-VOL:volumes_${old_version} --bridge-network-range 172.16.0.0/12 --container-network-firewall vm-network:published --certificate-key-size 2048

Re-Enable Ops User And Clean Up VIC Appliance
    Set Environment Variable  RUN_AS_OPS_USER  ${run-as-ops-user}
    Clean up VIC Appliance And Local Binary

Create Volumes Using Different Volume Store
    ${rc}  ${content1}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt VolumeStore=volumes_${old_version} --opt Capacity=2G --name=${volume1}
    Log  ${content1}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${content1}  ${volume1}

    ${rc}  ${content2}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --opt Capacity=4G --name=${volume2}
    Log  ${content2}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${content2}  ${volume2}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${volume1}
    Should Contain  ${output}  ${volume2}

Create Containers With Volumes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull nginx
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${volume1_container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v ${volume1}:${mydata} ${busybox} sh -c "echo '<p>HelloWorld1</p>' > /${mydata}/index.html"
    Log  ${volume1_container}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${volume2_container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v ${volume2}:${mydata} ${busybox} sh -c "echo '<p>HelloWorld2</p>' > /${mydata}/index.html"
    Log  ${volume2_container}
    Should Be Equal As Integers  ${rc}  0

Create Containers Using volume And Bridge And Public Network
    [Arguments]  ${containerName}  ${volume}  ${curl_msg}  ${photon_version}
    ${rc}  ${public_container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i -p ${port} --net public --name ${containerName} -v ${volume}:/usr/share/nginx/html:ro nginx
    Log  ${public_container}
    Should Be Equal As Integers  ${rc}  0    

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network connect bridge ${public_container}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Be Empty  ${output}
   
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerName}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    
    ${public_container_shortID}=  Get container shortID  ${public_container}
    ${container_ip}=  Get Container IP  %{VCH-PARAMS}  ${containerName}  public
    Set Suite Variable  ${nginx_ip}  ${container_ip}
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk http://${container_ip}:${port}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${curl_msg}
 
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format='{{(index .NetworkSettings.Networks "bridge").IPAddress}}' ${public_container_shortID}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --format='{{(index .NetworkSettings.Networks "public").IPAddress}}' ${public_container_shortID}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${output}
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec ${public_container_shortID} sh -c "uname -a"
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${photon_version}

Actions After Upgrade
    # Check nginx still working
    ${rc}  ${output}=  Run And Return Rc And Output  curl -sk http://${nginx_ip}:${port}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  HelloWorld1
    
    # Create a new container using volume2
    Create Containers Using volume And Bridge And Public Network  ${nginx_test2}  ${volume2}  HelloWorld2  ph2-esx
   
    # Create a new volume named volume3
    ${rc}  ${content3}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=${volume3}
    Log  ${content3}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${content3}  ${volume3}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${volume1}
    Should Contain  ${output}  ${volume2}  
    Should Contain  ${output}  ${volume3}

    ${rc}  ${volume3_container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -v ${volume3}:${mydata} ${busybox} sh -c "echo '<p>HelloWorld3</p>' > /${mydata}/index.html"
    Log  ${volume3_container}
    Should Be Equal As Integers  ${rc}  0
    
    # Check running container and total container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -q
    ${running_cids}=  Split String  ${output}
    ${running_count}=  Get Length  ${running_cids}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Integers  ${running_count}  2
    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -aq
    ${total_cids}=  Split String  ${output}
    ${total_count}=  Get Length  ${total_cids}
    Should Be Equal As Integers  ${rc}  0
    Should Be Equal As Integers  ${total_count}  5

Actions Before Second Upgrade
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${nginx_test2}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Wait Until Container Stops  ${nginx_test2}

*** Test Cases ***
Test Photon Version
    Create Volumes Using Different Volume Store
    Create Containers With Volumes
    Create Containers Using volume And Bridge And Public Network  ${nginx_test1}  ${volume1}  HelloWorld1  ph1-esx
    Upgrade
    Check Upgraded Version
    Actions After Upgrade
    Rollback
    Check Original Version
    Create Containers Using volume And Bridge And Public Network  ${nginx_test3}  ${volume3}  HelloWorld3  ph1-esx
    Actions Before Second Upgrade
    Upgrade with ID
    Check Upgraded Version

    Remove All Containers
    Log To Console  Regression Tests...
    Run Regression Tests


