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
Documentation  Test 1-44 - Docker CP Online
Resource  ../../resources/Util.robot
Suite Setup  Set up test files and install VIC appliance to test server
Suite Teardown  Clean up test files and VIC appliance to test server

*** Keywords ***
Set up test files and install VIC appliance to test server
    Install VIC Appliance To Test Server
    Create File  ${CURDIR}/foo.txt   hello world
    Create File  ${CURDIR}/content   fake file content for testing only
    Create Directory  ${CURDIR}/bar
    Create Directory  ${CURDIR}/mnt
    Create Directory  ${CURDIR}/mnt/vol1
    Create Directory  ${CURDIR}/mnt/vol2
    Create File  ${CURDIR}/mnt/root.txt   rw layer file
    Create File  ${CURDIR}/mnt/vol1/v1.txt   vol1 file
    Create File  ${CURDIR}/mnt/vol2/v2.txt   vol2 file
    ${rc}  ${output}=  Run And Return Rc And Output  dd if=/dev/zero of=${CURDIR}/largefile.txt count=1024 bs=1024
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name v1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name v2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Clean up test files and VIC appliance to test server
    Remove File  ${CURDIR}/foo.txt
    Remove File  ${CURDIR}/content
    Remove File  ${CURDIR}/largefile.txt
    Remove Directory  ${CURDIR}/bar  recursive=True
    Remove Directory  ${CURDIR}/mnt  recursive=True
    Cleanup VIC Appliance On Test Server

Start container and inspect directory
    [Arguments]  ${containerName}  ${directory}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerName}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec ${containerName} ls ${directory}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    [Return]  ${output}

*** Test Cases ***
Copy a directory from online container to host, dst path doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it --name online ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online sh -c 'mkdir newdir && echo "testing" > /newdir/test.txt'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir ${CURDIR}/newdir
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/newdir
    OperatingSystem.File Should Exist  ${CURDIR}/newdir/test.txt
    Remove Directory  ${CURDIR}/newdir  recursive=True

Copy the content of a directory from online container to host
    ${status}=  Get State Of Github Issue  5796
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5796 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/. ${CURDIR}/bar
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #OperatingSystem.File Should Exist  ${CURDIR}/bar/test.txt
    #Remove File  ${CURDIR}/bar/test.txt

Copy a file from online container to host, overwrite dst file
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/test.txt ${CURDIR}/foo.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  OperatingSystem.Get File  ${CURDIR}/foo.txt
    Should Contain  ${content}   testing

Copy a file from host to online container, dst directory doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt online:/doesnotexist/
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  no such directory

Copy a file and directory from host to online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt online:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar online:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt
    Should Contain  ${output}  bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a directory from host to online container, dst is a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it --name online -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar online:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  bar

Copy a file from host to offline container, dst is a volume shared with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/content offline:/vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  content

Copy a directory from offline container to host, dst is a volume shared with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1 ${CURDIR}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/vol1
    OperatingSystem.Directory Should Exist  ${CURDIR}/vol1/bar
    OperatingSystem.File Should Exist  ${CURDIR}/vol1/content
    Remove Directory  ${CURDIR}/vol1  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a large file to an online container, dst is a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt online:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls -l /vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  1048576
    Should Contain  ${output}  largefile.txt

Copy a non-existent file out of an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/dne ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error

Copy a non-existent directory out of an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/dne/. ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error

Concurrent copy: create processes to copy a small file from host to online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name concurrent -v v1:/vol1 -d -it ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for small file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt concurrent:/foo-${idx}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec concurrent ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   Should Contain  ${output}  foo-${idx}

Concurrent copy: repeat copy a large file from host to offline container several times
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for large file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt concurrent:/vol1/lg-${idx}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec concurrent ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   Should Contain  ${output}  lg-${idx}

Concurrent copy: repeat copy a large file from offline container to host several times
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for large file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp concurrent:/vol1/lg-${idx} ${CURDIR}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Should Be Equal As Integers  ${res.rc}  0
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   OperatingSystem.File Should Exist  ${CURDIR}/lg-${idx}
    \   Remove File  ${CURDIR}/lg-${idx}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f concurrent
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from host to an online container, dst includes several volumes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it -v v1:/mnt/vol1 -v v2:/mnt/vol2 --name subVol ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/mnt subVol:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol ls /mnt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  root.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol ls /mnt/vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  v1.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol ls /mnt/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  v2.txt

Sub volumes: copy from online container to host, src includes several volumes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp subVol:/mnt ${CURDIR}/result
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/result/vol1
    OperatingSystem.Directory Should Exist  ${CURDIR}/result/vol2
    OperatingSystem.File Should Exist  ${CURDIR}/result/root.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result/vol1/v1.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result/vol2/v2.txt
    Remove Directory  ${CURDIR}/result  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from host to an offline container, dst includes a shared vol with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it -v vol1:/vol1 --name subVol_on ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i -v vol1:/mnt/vol1 --name subVol_off ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/mnt subVol_off:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  subVol_off  /mnt
    Should Contain  ${output}  root.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol_off ls /mnt/vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  v1.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol_off ls /mnt/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  v2.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop subVol_off
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from an offline container to host, src includes a shared vol with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp subVol_off:/mnt ${CURDIR}/result
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/result/vol1
    OperatingSystem.Directory Should Exist  ${CURDIR}/result/vol2
    OperatingSystem.File Should Exist  ${CURDIR}/result/root.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result/vol1/v1.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result/vol2/v2.txt
    Remove Directory  ${CURDIR}/result  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol_off
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
