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
Documentation  Test 1-43 - Docker CP Offline
Resource  ../../resources/Util.robot
Suite Setup  Set up test files and install VIC appliance to test server
Suite Teardown  Clean up test files and VIC appliance to test server

*** Keywords ***
Set up test files and install VIC appliance to test server
    Install VIC Appliance To Test Server
    Create File  ${CURDIR}/foo.txt   hello world
    Create Directory  ${CURDIR}/bar
    Create File  ${CURDIR}/content   fake file content for testing only
    ${rc}  ${output}=  Run And Return Rc And Output  dd if=/dev/zero of=${CURDIR}/largefile.txt count=4096 bs=4096
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name smallVol --opt Capacity=1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Clean up test files and VIC appliance to test server
    Remove File  ${CURDIR}/foo.txt
    Remove File  ${CURDIR}/content
    Remove File  ${CURDIR}/largefile.txt
    Remove Directory  ${CURDIR}/bar  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm smallVol
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Cleanup VIC Appliance On Test Server

Start container and inspect directory
    [Arguments]  ${containerName}  ${directory}
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start ${containerName}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls ${directory}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    [Return]  ${output}

*** Test Cases ***
Copy a file from host to offline container root dir
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  offline  /
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline sh -c 'rm /foo.txt'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a directory from offline container to host cwd
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline sh -c 'mkdir testdir && echo "file content" > /testdir/fakefile'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/testdir ${CURDIR}/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Directory Should Exist  ${CURDIR}/testdir
    File Should Exist  ${CURDIR}/testdir/fakefile
    Remove Directory  ${CURDIR}/testdir  recursive=True

Copy a directory from host to offline container, dst path doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar offline:/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  offline  /
    Should Contain  ${output}   bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a non-existent file out of an offline container
    ${status}=  Get State Of Github Issue  5717
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-43-Docker-CP-Offline.robot needs to be updated now that Issue #5717 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/dne ${CURDIR}
    #Should Not Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  Error

Copy a non-existent directory out of an offline container
    ${status}=  Get State Of Github Issue  5717
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-43-Docker-CP-Offline.robot needs to be updated now that Issue #5717 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/dne/. ${CURDIR}
    #Should Not Be Equal As Integers  ${rc}  0
    #Should Contain  ${output}  Error

Copy a non-existent directory into an offline container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/dne/ offline:/
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  no such file or directory
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a large file that exceeds the container volume into an offline container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i -v smallVol:/small --name offline ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt offline:/small
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a volume
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  offline  /vol1
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a nested volume with 2 levels
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 -v vol2:/vol1/vol2 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/vol1/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  offline  /vol1/vol2
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a nested volume with 3 levels
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 -v vol2:/vol1/vol2 -v vol3:/vol1/vol2/vol3 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/vol1/vol2/vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${output}=  Start container and inspect directory  offline  /vol1/vol2/vol3
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Concurrent copy: repeat copy a small file from host to offline container several times
    ${status}=  Get State Of Github Issue  5742
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-43-Docker-CP-Offline.robot needs to be updated now that Issue #5742 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/foo1
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/foo2
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/foo3
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${output}=  Start container and inspect directory  offline  /
    #Should Contain  ${output}  foo1
    #Should Contain  ${output}  foo2
    #Should Contain  ${output}  foo3
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop offline
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error

Concurrent copy: repeat copy a large file from host to offline container several times
    ${status}=  Get State Of Github Issue  5742
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-43-Docker-CP-Offline.robot needs to be updated now that Issue #5742 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt offline:/vol1/lg1
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt offline:/vol1/lg2
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt offline:/vol1/lg3
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${output}=  Start container and inspect directory  offline  /vol1
    #Should Contain  ${output}  lg1
    #Should Contain  ${output}  lg2
    #Should Contain  ${output}  lg3
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop offline
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error

Concurrent copy: repeat copy a large file from offline container to host several times
    ${status}=  Get State Of Github Issue  5742
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-43-Docker-CP-Offline.robot needs to be updated now that Issue #5742 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1/lg1 ${CURDIR}/large1
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1/lg1 ${CURDIR}/large2
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1/lg1 ${CURDIR}/large3
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #File Should Exist  ${CURDIR}/large1
    #File Should Exist  ${CURDIR}/large2
    #File Should Exist  ${CURDIR}/large3
    #Remove File  ${CURDIR}/large1
    #Remove File  ${CURDIR}/large2
    #Remove File  ${CURDIR}/large3
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    