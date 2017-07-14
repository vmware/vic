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
Documentation  Test 1-43 - Docker CP
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Set up test files, directories and volumes
    Create File  ${CURDIR}/foo.txt   hello world
    Create Directory  ${CURDIR}/bar
    Create File  ${CURDIR}/content   fake file content for testing only
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file and directory from host to offline container root dir
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar offline:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt
    Should Contain  ${output}  bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline sh -c 'rm /foo.txt && rmdir /bar'
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
    Remove Directory  ${CURDIR}/testdir recursive=True
    Remove Directory  ${CURDIR}/testdir

Copy a directory from host to offline container, dst path doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar offline:/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}   bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a directory from online container to host, dst path doesn't exit
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
    Remove Directory  ${CURDIR}/newdir recursive=True
    Remove Directory  ${CURDIR}/newdir

Copy the content of a directory from online container to host
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/. ${CURDIR}/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    File Should Exist  ${CURDIR}/bar/test.txt
    Remove File  ${CURDIR}/bar/test.txt

Copy a file from online container to host, overwrite dst file
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp online:/newdir/test.txt ${CURDIR}/foo.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  Get File  ${CURDIR}/foo.txt
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
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a volume
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt offline:/vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start offline
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls /vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls /vol1/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
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
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec offline ls /vol1/vol2/vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

# depends on disk in use error and online working
Copy a file from host to offline container, dst is a volume shared with an online container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name online -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name offline -v vol1:/vol1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start online
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
    Directory Should Exist  ${CURDIR}/vol1
    Directory Should Exist  ${CURDIR}/vol1/bar
    Directory Should Exist  ${CURDIR}/vol1/vol2
    File Should Exist  ${CURDIR}/vol1/foo.txt
    File Should Exist  ${CURDIR}/vol1/content
    Remove Directory  ${CURDIR}/vol1 recursive=True
    Remove Directory  ${CURDIR}/vol1
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f offline
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a large file to an online container, dst is a volume
    ${output}=  Run  dd if=/dev/zero of=${CURDIR}/largefile.txt count=4096 bs=1024
    Should Contain  ${output}  4194304 bytes
    File Should Exist  ${CURDIR}/largefile.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/largefile.txt online:/vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec online ls -l /vol1/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  4194304
    Should Contain  ${output}  largefile.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f online
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Remove File  ${CURDIR}/largefile.txt

Clean up current directory
    Remove File  ${CURDIR}/foo.txt
    Remove File  ${CURDIR}/content
    Remove Directory  ${CURDIR}/bar recursive=True
    Remove Directory  ${CURDIR}/bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume rm vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error