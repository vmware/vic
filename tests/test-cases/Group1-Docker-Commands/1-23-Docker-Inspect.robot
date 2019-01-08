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
Documentation  Test 1-23 - Docker Inspect
Resource  ../../resources/Util.robot
Suite Setup  Conditional Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  20 minutes

*** Keywords ***
Get container inspect status
    [Arguments]  ${container}
    ${rc}  ${status}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container} -f '{{.State.Status}}'
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${status}

*** Test Cases ***
Simple docker inspect of image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying type
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=image ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying incorrect type
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=container ${busybox}
    Should Be Equal As Integers  ${rc}  1
    ${out}=  Run Keyword If  '${busybox}' == 'busybox'  Should Contain  ${output}  Error: No such container: busybox
    ${out}=  Run Keyword Unless  '${busybox}' == 'busybox'  Should Contain  ${output}  Error: No such container: wdc-harbor-ci.eng.vmware.com/default-project/busybox

Simple docker inspect of container
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container specifying type
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=container ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container check cmd and image name
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox} /bin/bash
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${config}=  Get From Dictionary  ${output[0]}  Config
    ${image}=  Get From Dictionary  ${config}  Image
    Should Contain  ${image}  busybox
    ${cmd}=  Get From Dictionary  ${config}  Cmd
    Should Be Equal As Strings  ${cmd}  [u'/bin/bash']

Docker inspect container specifying incorrect type
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=image ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such image: ${container}

Docker inspect container with multiple networks
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create net-one
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create net-two
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=two-net-test --net=net-one busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network connect net-two two-net-test
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start two-net-test
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' two-net-test
    Should Contain  ${out}  net-two
    Should Contain  ${out}  net-one
    Should Be Equal As Integers  ${rc}  0

Docker inspect container with correct gateway
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name=test-gateway -d ${busybox} sleep 600
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range .NetworkSettings.Networks}}{{.Gateway}}{{end}}' test-gateway
    Should Contain  ${out}  172.16.0.1
    Should Not Contain  ${out}  /

Docker inspect invalid object
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect fake
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such object: fake

Docker inspect non-nil volume
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=test-with-volume -v /var/lib/test busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Config.Volumes}}' test-with-volume
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  /var/lib/test
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect test-with-volume | jq '.[]|.["Config"]|.["Volumes"]|keys[0]'
    Should Be Equal As Integers  ${rc}  0
    ${mount}=  Split String  ${out}  :
	${volID}=  Get Substring  @{mount}[0]  1
    Log To Console  Find volume ${volID} in container inspect
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume ls
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  ${volID}

Inspect RepoDigest is valid
    ${rc}  Run And Return Rc  docker %{VCH-PARAMS} rmi ${busybox}
    ${rc}  ${busybox_digest}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox} | grep Digest | awk '{print $2}'
    Should Be Equal As Integers  ${rc}  0
    Should Not Be Empty  ${busybox_digest}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.RepoDigests}}' ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  ${busybox_digest}

Docker inspect mount data
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name=named-volume
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=mount-data-test -v /mnt/test -v named-volume:/mnt/named busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Mounts}}' mount-data-test
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${out}  /mnt/test
    Should Contain  ${out}  /mnt/named

Docker inspect container status
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox} /bin/sh -c 'a=0; while [ $a -lt 90 ]; do echo "line $a"; a=`expr $a + 1`; sleep 2; done;'
    Should Be Equal As Integers  ${rc}  0
    # keyword at top of file
    ${created}=  Get container inspect status  ${container}
    Should Contain  ${created}  created
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0
    # keyword at top of file
    ${running}=  Get container inspect status  ${container}
    Should Contain  ${running}  running
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${container}
    Should Be Equal As Integers  ${rc}  0
    # keyword at top of file
    ${stopped}=  Get container inspect status  ${container}
    Should Contain  ${stopped}  exited

Docker inspect container with specified hostname
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -h testhostname --name=test-hostname -d ${busybox} sleep 600
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{.Config.Hostname}}' test-hostname
    Should Be Equal  ${out}  testhostname

Docker inspect for macaddress
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name=test-macaddress -d busybox sleep 600
    Should Be Equal As Integers  ${rc}  0
    Sleep  60s
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range .NetworkSettings.Networks}}{{.MacAddress}}{{end}}' ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Match Regexp  ${output}  ^([0-9a-f]{2}[:-]){5}([0-9a-f]{2})$

Docker inspect container with specified DNS
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --dns=8.8.8.8 --name=test-with-specified-dns busybox sleep 600
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  8.8.8.8
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0

Docker inspect container with multiple specified DNS
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --dns=8.8.8.8 --dns=8.8.8.9 --name=test-with-multiple-specified-dns busybox sleep 600
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  8.8.8.8
    Should Contain  ${output}  8.8.8.9
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${container}
    Should Be Equal As Integers  ${rc}  0

Docker inspect for cpu and memory settings
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${hostconfig}=  Get From Dictionary  ${output[0]}  HostConfig
    ${mem}=  Get From Dictionary  ${hostconfig}  Memory
    Should Be Equal As Integers  ${mem}  2048
    ${cpu}=  Get From Dictionary  ${hostconfig}  CpusetCpus
    Should Be Equal As Integers  ${cpu}  2
