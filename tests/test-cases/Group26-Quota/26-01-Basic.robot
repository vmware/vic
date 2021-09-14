# Copyright 2018 VMware, Inc. All Rights Reserved.
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
Documentation     Suite 26-01 - Basic
Resource          ../../resources/Util.robot
Resource          ../../resources/Group26-Storage-Quota-Util.robot
Suite Setup  Install VIC Appliance To Test Server  additional-args='--storage-quota=15'
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  30 minutes

*** Test Cases ***
Create a VCH with storage quota and docker info shows the storage quota and storage usage
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}
    Should Be Equal As Integers  ${limitval}  15
    Should Be Equal As Integers  ${usageval}  0

Create a container and docker info shows the storage usage has changed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  15
    Should Be True  ${usageval} > 9
    Set Suite Variable  ${pre_usage_val}  ${usageval}

Pull a debian image and docker info shows the storage usage has changed
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${debian}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}
    Should Be Equal As Integers  ${limitval}  15
    Should Be True  ${usageval} > ${pre_usage_val}
    Set Suite Variable  ${pre_usage_val}  ${usageval}

Create second container and get storage quota exceeding failure
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  125
    Should Contain  ${output}  Storage quota exceeds

Configure VCH with a larger storage quota of 35GB
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --storage-quota 35
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}
    Should Be Equal As Integers  ${limitval}  35

Create second container successfully
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id --name secondc ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be True  ${usageval} > 18

Remove a container successfully with storage usage changes
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rm -f secondc
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be Equal As Integers  ${usageval}  ${pre_usage_val}

Create a busybox container with memory of 4GB successfully
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -m=4g -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be True  ${usageval} > 20

Create a debian container and commit to an image successfully
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d --name commit1 ${debian} tail -f /dev/null
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be True  ${usageval} > 29
    Set Suite Variable  ${pre_usage_val}  ${usageval}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec commit1 apt-get update
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec commit1 apt-get install nano -y
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop -t1 commit1
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} commit commit1 debian-nano
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be True  ${usageval} > ${pre_usage_val}

Delete an image successfully with storage usage decreased
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} rmi debian-nano
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0

    ${limitval}  ${usageval}=  Get storage quota limit and usage  ${output}

    Should Be Equal As Integers  ${limitval}  35
    Should Be Equal As Integers  ${usageval}  ${pre_usage_val}

Create a busybox continainer afer unsetting storage quota
    ${output}=  Run  bin/vic-machine-linux configure --name=%{VCH-NAME} --target=%{TEST_URL}%{TEST_DATACENTER} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --timeout %{TEST_TIMEOUT} --storage-quota 0
    Should Contain  ${output}  Completed successfully

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -id ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} info
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  storage limit

    ${usageval}=  Get storage usage  ${output}

    Should Be True  ${usageval} > 38
