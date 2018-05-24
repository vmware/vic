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
Test Timeout  20 minutes

*** Keywords ***
Set up test files and install VIC appliance to test server
    Conditional Install VIC Appliance To Test Server
    Remove All Volumes
    Create File  ${CURDIR}/foo-online.txt   hello world
    Create File  ${CURDIR}/content-online   fake file content for testing only
    Create Directory  ${CURDIR}/bar-online
    Create Directory  ${CURDIR}/mnt-online
    Create Directory  ${CURDIR}/mnt-online/vol1-online
    Create Directory  ${CURDIR}/mnt-online/vol2-online
    Create File  ${CURDIR}/mnt-online/root-online.txt   rw layer file
    Create File  ${CURDIR}/mnt-online/vol1-online/v1-online.txt   vol1 file
    Create File  ${CURDIR}/mnt-online/vol2-online/v2-online.txt   vol2 file
    ${rc}  ${output}=  Run And Return Rc And Output  dd if=/dev/urandom of=${CURDIR}/largefile-online.txt count=1024 bs=1024
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create v1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create v2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Clean up test files and VIC appliance to test server
    Run Keyword and Continue on Failure  Remove File  ${CURDIR}/foo-online.txt
    Run Keyword and Continue on Failure  Remove File  ${CURDIR}/content-online
    Run Keyword and Continue on Failure  Remove File  ${CURDIR}/largefile-online.txt
    Run Keyword and Continue on Failure  Remove Directory  ${CURDIR}/bar-online  recursive=True
    Run Keyword and Continue on Failure  Remove Directory  ${CURDIR}/mnt-online  recursive=True
    Cleanup VIC Appliance On Test Server

*** Test Cases ***
Copy a directory from online container to host, destination path doesn't exist
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/. ${CURDIR}/bar-online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.File Should Exist  ${CURDIR}/bar-online/test.txt
    Remove File  ${CURDIR}/bar-online/test.txt

Copy a file from online container to host, overwrite destination file
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/test.txt ${CURDIR}/foo-online.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  OperatingSystem.Get File  ${CURDIR}/foo-online.txt
    Should Contain  ${content}   testing

Copy a file from host to online container, destination directory doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo-online.txt online:/doesnotexist/
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  no such directory

Copy a file and directory from host to online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo-online.txt online:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar-online online:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo-online.txt
    Should Contain  ${output}  bar-online
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a directory from host to online container, destination is a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it --name online_vol -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar-online online_vol:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online_vol ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  bar-online

Copy a file from host to offline container, destination is a volume shared with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/content-online offline:/vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online_vol ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  content-online

Copy a directory from offline container to host, destination is a volume shared with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1 ${CURDIR}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/vol1
    OperatingSystem.Directory Should Exist  ${CURDIR}/vol1/bar-online
    OperatingSystem.File Should Exist  ${CURDIR}/vol1/content-online
    Remove Directory  ${CURDIR}/vol1  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a large file to an online container, destination is a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile-online.txt online_vol:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online_vol ls -l /vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  1048576
    Should Contain  ${output}  largefile-online.txt

Copy a non-existent file out of an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online_vol:/dne/dne ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error

Copy a non-existent directory out of an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online_vol:/dne/. ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online_vol
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Concurrent copy: create processes to copy a small file from host to online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --name concurrent -v v1:/vol1 -d -it ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for small file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp ${CURDIR}/foo-online.txt concurrent:/foo-online-${idx}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Log  ${res.stderr}
    \   Log  ${res.stdout}
    \   Should Be Equal As Integers  ${res.rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec concurrent ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   Should Contain  ${output}  foo-online-${idx}

Concurrent copy: repeat copy a large file from host to online container several times
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for large file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp ${CURDIR}/largefile-online.txt concurrent:/vol1/lg-online-${idx}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Log  ${res.stderr}
    \   Log  ${res.stdout}
    \   Should Be Equal As Integers  ${res.rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec concurrent ls /vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   Should Contain  ${output}  lg-online-${idx}

Concurrent copy: repeat copy a large file from online container to host several times
    ${pids}=  Create List
    Log To Console  \nIssue 10 docker cp commands for large file
    :FOR  ${idx}  IN RANGE  0  10
    \   ${pid}=  Start Process  docker %{VCH-PARAMS} cp concurrent:/vol1/lg-online-${idx} ${CURDIR}  shell=True
    \   Append To List  ${pids}  ${pid}
    Log To Console  \nWait for them to finish and check their RC
    :FOR  ${pid}  IN  @{pids}
    \   Log To Console  \nWaiting for ${pid}
    \   ${res}=  Wait For Process  ${pid}
    \   Log  ${res.stderr}
    \   Log  ${res.stdout}
    \   Should Be Equal As Integers  ${res.rc}  0
    Log To Console  \nCheck if the copy operations succeeded
    :FOR  ${idx}  IN RANGE  0  10
    \   OperatingSystem.File Should Exist  ${CURDIR}/lg-online-${idx}
    \   Remove File  ${CURDIR}/lg-online-${idx}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f concurrent
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from host to an online container, destination includes several volumes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it -v A:/mnt/vol1 -v B:/mnt/vol2 --name subVol ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/mnt-online subVol:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec subVol find /mnt-online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  /mnt-online/root-online.txt
    Should Contain  ${output}  /mnt-online/vol1-online/v1-online.txt
    Should Contain  ${output}  /mnt-online/vol2-online/v2-online.txt

Sub volumes: copy from online container to host, source includes several volumes
    Remove Directory  ${CURDIR}/result1  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp subVol:/mnt-online ${CURDIR}/result1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/result1/vol1-online
    OperatingSystem.Directory Should Exist  ${CURDIR}/result1/vol2-online
    OperatingSystem.File Should Exist  ${CURDIR}/result1/root-online.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result1/vol1-online/v1-online.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result1/vol2-online/v2-online.txt
    Remove Directory  ${CURDIR}/result1  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from host to an offline container, destination includes a shared vol with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d -it -v vol1:/vol1 --name subVol_on ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i -v vol1:/mnt/vol1 --name subVol_off ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/mnt-online subVol_off:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start Container and Exec Command  subVol_off  find /mnt-online
    Should Contain  ${output}  /mnt-online/root-online.txt
    Should Contain  ${output}  /mnt-online/vol1-online/v1-online.txt
    Should Contain  ${output}  /mnt-online/vol2-online/v2-online.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop subVol_off
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Sub volumes: copy from an offline container to host, source includes a shared vol with an online container
    Remove Directory  ${CURDIR}/result2  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp subVol_off:/mnt-online ${CURDIR}/result2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    OperatingSystem.Directory Should Exist  ${CURDIR}/result2/vol1-online
    OperatingSystem.Directory Should Exist  ${CURDIR}/result2/vol2-online
    OperatingSystem.File Should Exist  ${CURDIR}/result2/root-online.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result2/vol1-online/v1-online.txt
    OperatingSystem.File Should Exist  ${CURDIR}/result2/vol2-online/v2-online.txt
    Remove Directory  ${CURDIR}/result2  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol_off
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f subVol_on
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Malformed tarball doesn't extract to wrong container from shared volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name one -v vol1:/vol1 -v vol2:/vol2 ${busybox} /bin/sh
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name two -it -v vol2:/vol2 -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  cat ${CURDIR}/../../resources/archive.tar.gz | docker %{VCH-PARAMS} cp - two:/vol1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop one
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start two
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec one ls /tmp
    Should Not Contain  ${output}  pingme
