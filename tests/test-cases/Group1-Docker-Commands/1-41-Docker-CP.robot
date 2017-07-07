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
Documentation  Test 1-41 - Docker CP
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple case: copy a file and directory from host to online container root dir
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name test1 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Create File  ${CURDIR}/foo.txt   hello world
    Create Directory  ${CURDIR}/bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar test1:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt
    Should Contain  ${output}  bar
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 sh -c 'rm /foo.txt && rmdir /bar'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

# missing src file case 1. missing src dir case 2. src host dst container
Simple case: copy a directory from host to online container, dst path doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar test1:/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}   bar

Simple case: copy a directory from online container to host, dst path doesn't exit
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 sh -c 'mkdir newdir && echo "testing" > /newdir/test.txt'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp test1:/newdir ${CURDIR}/newdir
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Directory Should Exist  ${CURDIR}/newdir
    File Should Exist  ${CURDIR}/newdir/test.txt

Simple case: copy the content of a directory from online container to host
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp test1:/newdir/. ${CURDIR}/bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    File Should Exist  ${CURDIR}/bar/test.txt

Simple case: copy a file from online container to host, overwrite dst file
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp test1:/newdir/test.txt ${CURDIR}/foo.txt
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${content}=  Get File  ${CURDIR}/foo.txt
    Should Contain  ${content}   testing

Simple case: copy a file from online container to host, dst directory doesn't exist
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/doesnotexist/
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  no such directory
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 sh -c 'rm -rf /newdir && rmdir /bar'
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Simple case: copy a file from host to offline container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop -t 0 test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a volume
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name test1 -v vol1:/vol1 ${busybox}
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start test1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     Should Contain  ${output}  foo.txt

Copy a directory from host to online container, dst is a volume
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/bar test1:/vol1/
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /vol1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error
     Should Contain  ${output}  bar
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f test1
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a nested volume with 2 levels
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name test1 -v vol1:/vol1 -v vol2:/vol1/vol2 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/vol1/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /vol1/vol2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Copy a file from host to offline container, dst is a nested volume with 3 levels
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} volume create --name vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create -i --name test1 -v vol1:/vol1 -v vol2:/vol1/vol2 -v vol3:/vol1/vol2/vol3 ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} cp ${CURDIR}/foo.txt test1:/vol1/vol2/vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start test1
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec test1 ls /vol1/vol2/vol3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Should Contain  ${output}  foo.txt

#====== haven't been added to md ========
#Copy a file from host to online container, dst is a volume shared with another online container


#Copy a file from host to offline container, dst is a volume shared with an online container


Clean up current directory
    Remove File  ${CURDIR}/foo.txt
    Remove Directory  ${CURDIR}/bar recursive=True
    Remove Directory  ${CURDIR}/newdir recursive=True