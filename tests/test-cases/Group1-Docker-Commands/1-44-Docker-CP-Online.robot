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
    Create Directory  ${CURDIR}/bar
    Create File  ${CURDIR}/content   fake file content for testing only
    ${rc}  ${output}=  Run And Return Rc And Output  dd if=/dev/zero of=${CURDIR}/largefile.txt count=1024 bs=1024
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol1
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
    Cleanup VIC Appliance On Test Server

*** Test Cases ***
Copy a directory from online container to host, dst path doesn't exist
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name online ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online sh -c 'mkdir newdir && echo "testing" > /newdir/test.txt'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir ${CURDIR}/newdir
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Directory Should Exist  ${CURDIR}/newdir
    File Should Exist  ${CURDIR}/newdir/test.txt
    Remove Directory  ${CURDIR}/newdir  recursive=True

Copy the content of a directory from online container to host
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/. ${CURDIR}/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    File Should Exist  ${CURDIR}/bar/test.txt
    Remove File  ${CURDIR}/bar/test.txt

Copy a file from online container to host, overwrite dst file
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/test.txt ${CURDIR}/foo.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  Get File  ${CURDIR}/foo.txt
    Should Contain  ${content}   testing

Copy a file from host to online container, dst directory doesn't exist
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt online:/doesnotexist/
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  no such directory

Copy a file and directory from host to online container
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
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
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name online -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start online
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
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
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
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp offline:/vol1 ${CURDIR}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Directory Should Exist  ${CURDIR}/vol1
    Directory Should Exist  ${CURDIR}/vol1/bar
    File Should Exist  ${CURDIR}/vol1/content
    Remove Directory  ${CURDIR}/vol1  recursive=True
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a large file to an online container, dst is a volume
    #${status}=  Get State Of Github Issue  5606
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5606 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt online:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls -l /vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  1048576
    Should Contain  ${output}  largefile.txt

Copy a non-existent file out of an online container
    #${status}=  Get State Of Github Issue  5717
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5717 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/dne ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error

Copy a non-existent directory out of an online container
    #${status}=  Get State Of Github Issue  5717
    #Run Keyword If  '${status}' == 'closed'  Fail  Test 1-44-Docker-CP-Online.robot needs to be updated now that Issue #5717 has been resolved
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/dne/. ${CURDIR}
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
