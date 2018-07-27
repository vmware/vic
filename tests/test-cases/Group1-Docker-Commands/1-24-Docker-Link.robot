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
Documentation  Test 1-24 - Docker Link
Resource  ../../resources/Util.robot
Suite Setup  Conditional Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  20 minutes

*** Test Cases ***
Link and alias
    # NOTE: we are using debian for the ping operations because busybox seems to have intermittent bugs in the net utilities.

    # In the case of nslookup the busybox implementation is querying both the A and AAAA records. If the AAAA response comes back
    # first then it doesn't report the A response, however if the other order then you get both responses output.

    # link support for container on bridge network only
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name foo ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --link foo:bar ${debian} ping -c1 bar
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create jedi
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${debian}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it -d --net jedi --name first ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net jedi ${debian} ping -c1 first
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # cannot reach first from another network
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run ${debian} ping -c1 first
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain Any  ${output}  unknown host  Name or service not known

    # the link
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net jedi --link first:1st ${debian} ping -c1 1st
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # cannot reach first using c1 from another container
    # first run a container that has the alias "c1" for the "first" container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --net jedi --link first:1st ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    # check if we can use alias "c1" from another container
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net jedi ${debian} ping -c1 1st
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain Any  ${output}  unknown host  Name or service not known

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it -d --net jedi --net-alias 2nd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # the alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net jedi ${debian} ping -c1 2nd
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # another container with same network alias
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -it -d --net jedi --net-alias 2nd ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # pinned the version of busybox to one with a working nslookup. See https://github.com/docker-library/busybox/issues/48
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run --net jedi --name lookup busybox:1.28.4 nslookup 2nd
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} logs lookup
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Address 1
    Should Contain  ${output}  Address 2
