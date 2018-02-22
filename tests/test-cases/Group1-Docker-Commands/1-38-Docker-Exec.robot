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
Documentation  Test 1-38 - Docker Exec
Resource  ../../resources/Util.robot
Suite Setup  Conditional Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  20 minutes

*** Test Cases ***
Exec -d
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -d ${id} /bin/touch tmp/force
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec ${id} /bin/ls -al /tmp/force
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  force

Exec Echo
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  5
    \   ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec ${id} /bin/echo "Help me, Obi-Wan Kenobi. You're my only hope."
    \   Should Be Equal As Integers  ${rc}  0
    \   Should Be Equal As Strings  ${output}  Help me, Obi-Wan Kenobi. You're my only hope.

Exec Echo -i
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  5
    \   ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -i ${id} /bin/echo "Your eyes can deceive you. Don't trust them."
    \   Should Be Equal As Integers  ${rc}  0
    \   Should Be Equal As Strings  ${output}  Your eyes can deceive you. Don't trust them.

Exec Echo -t
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d ${busybox} /bin/top -d 600
    Should Be Equal As Integers  ${rc}  0
    :FOR  ${idx}  IN RANGE  0  5
    \   ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec -t ${id} /bin/echo "Do. Or do not. There is no try."
    \   Should Be Equal As Integers  ${rc}  0
    \   Should Be Equal As Strings  ${output}  Do. Or do not. There is no try.

Exec Sort
    ${status}=  Get State Of Github Issue  5479
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-38-Docker-Exec.robot needs to be updated now that Issue #5479 has been resolved
    #${rc}  ${tmp}=  Run And Return Rc And Output  mktemp -d -p /tmp
    #Should Be Equal As Integers  ${rc}  0
    #${fifo}=  Catenate  SEPARATOR=/  ${tmp}  fifo
    #${rc}  ${output}=  Run And Return Rc And Output  mkfifo ${fifo}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top -d 600
    #Should Be Equal As Integers  ${rc}  0
    #:FOR  ${idx}  IN RANGE  0  5
    #\     Start Process  docker %{VCH-PARAMS} exec ${output} /bin/sort < ${fifo}  shell=True  alias=custom
    #\     Run  echo one > ${fifo}
    #\     ${ret}=  Wait For Process  custom
    #\     Log  ${ret.stderr}
    #\     Should Be Empty  ${ret.stdout}
    #\     Should Be Equal As Integers  ${ret.rc}  0
    #\     Should Be Empty  ${ret.stderr}
    #Run  rm -rf ${tmp}

Exec Sort -i
    ${status}=  Get State Of Github Issue  5479
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-38-Docker-Exec.robot needs to be updated now that Issue #5479 has been resolved
    #${rc}  ${tmp}=  Run And Return Rc And Output  mktemp -d -p /tmp
    #Should Be Equal As Integers  ${rc}  0
    #${fifo}=  Catenate  SEPARATOR=/  ${tmp}  fifo
    #${rc}  ${output}=  Run And Return Rc And Output  mkfifo ${fifo}
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top -d 600
    #Should Be Equal As Integers  ${rc}  0
    #:FOR  ${idx}  IN RANGE  0  5
    #\     Start Process  docker %{VCH-PARAMS} exec -i ${output} /bin/sort < ${fifo}  shell=True  alias=custom
    #\     Run  echo one > ${fifo}
    #\     ${ret}=  Wait For Process  custom
    #\     Log  ${ret.stderr}
    #\     Should Be Equal  ${ret.stdout}  one
    #\     Should Be Equal As Integers  ${ret.rc}  0
    #\     Should Be Empty  ${ret.stderr}
    #Run  rm -rf ${tmp}

Exec NonExisting
    ${status}=  Get State Of Github Issue  5479
    Run Keyword If  '${status}' == 'closed'  Fail  Test 1-38-Docker-Exec.robot needs to be updated now that Issue #5479 has been resolved
    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #Should Not Contain  ${output}  Error
    #${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -d busybox /bin/top -d 600
    #Should Be Equal As Integers  ${rc}  0
    #:FOR  ${idx}  IN RANGE  0  5
    #\   ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} exec ${id} /NonExisting
    #\   Should Be Equal As Integers  ${rc}  0
    #\   Should Contain  ${output}  no such file or directory

Concurrent Simple Exec
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error

     ${suffix}=  Evaluate  '%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
     Set Test Variable  ${ExecSimpleContainer}  Exec-simple-${suffix}
     ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name ${ExecSimpleContainer} ${busybox} sleep 30
     Should Be Equal As Integers  ${rc}  0

     :FOR  ${idx}  IN RANGE  1  5
     \   Start Process  docker %{VCH-PARAMS} exec ${id} /bin/ls  alias=exec-simple-%{VCH-NAME}-${idx}  shell=true

     :FOR  ${idx}  IN RANGE  1  5
     \   ${result}=  Wait For Process  exec-simple-%{VCH-NAME}-${idx}  timeout=40s
     \   Should Be Equal As Integers  ${result.rc}  0
     \   # if any of these are missing check to see if the busy box fs changed first.
     \   Should Contain  ${result.stdout}  bin
     \   Should Contain  ${result.stdout}  dev
     \   Should Contain  ${result.stdout}  etc
     \   Should Contain  ${result.stdout}  home
     \   Should Contain  ${result.stdout}  lib
     \   Should Contain  ${result.stdout}  lost+found
     \   Should Contain  ${result.stdout}  mnt
     \   Should Contain  ${result.stdout}  proc
     \   Should Contain  ${result.stdout}  root
     \   Should Contain  ${result.stdout}  run
     \   Should Contain  ${result.stdout}  sbin
     \   Should Contain  ${result.stdout}  sys
     \   Should Contain  ${result.stdout}  tmp
     \   Should Contain  ${result.stdout}  usr
     \   Should Contain  ${result.stdout}  var

     ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} wait ${id}
     Should Be Equal As Integers  ${rc}  0


Exec During Poweroff Of A Container Performing A Long Running Task
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error

     ${suffix}=  Evaluate  '%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
     Set Test Variable  ${ExecPowerOffContainerLong}  Exec-Poweroff-${suffix}
     ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name ${ExecPoweroffContainerLong} ${busybox} /bin/top
     Should Be Equal As Integers  ${rc}  0

     :FOR  ${idx}  IN RANGE  1  10
     \   Start Process  docker %{VCH-PARAMS} exec ${id} /bin/ls  alias=exec-%{VCH-NAME}-${idx}  shell=true


     Sleep  10s
     ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} stop ${id}
     Should Be Equal As Integers  ${rc}  0

     ${combinedErr}=  Set Variable
     ${combinedOut}=  Set Variable

     :FOR  ${idx}  IN RANGE  1  10
     \   ${result}=  Wait For Process  exec-%{VCH-NAME}-${idx}  timeout=2 mins
     \   ${combinedErr}=  Catenate  ${combinedErr}  ${result.stderr}${\n}
     \   ${combinedOut}=  Catenate  ${combinedOut}  ${result.stdout}${\n}

     Should Contain  ${combinedErr}  Container (${id}) is not running

     # We should get atleast one successful exec...
     Should Contain  ${combinedOut}  bin
     Should Contain  ${combinedOut}  dev
     Should Contain  ${combinedOut}  etc
     Should Contain  ${combinedOut}  home
     Should Contain  ${combinedOut}  lib
     Should Contain  ${combinedOut}  lost+found
     Should Contain  ${combinedOut}  mnt
     Should Contain  ${combinedOut}  proc
     Should Contain  ${combinedOut}  root
     Should Contain  ${combinedOut}  run
     Should Contain  ${combinedOut}  sbin
     Should Contain  ${combinedOut}  sys
     Should Contain  ${combinedOut}  tmp
     Should Contain  ${combinedOut}  usr
     Should Contain  ${combinedOut}  var


Exec During Poweroff Of A Container Performing A Short Running Task
     ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${busybox}
     Should Be Equal As Integers  ${rc}  0
     Should Not Contain  ${output}  Error

     ${suffix}=  Evaluate  '%{DRONE_BUILD_NUMBER}-' + str(random.randint(1000,9999))  modules=random
     Set Test Variable  ${ExecPoweroffContainerShort}  Exec-Poweroff-${suffix}
     ${rc}  ${id}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name ${ExecPoweroffContainerShort} ${busybox} sleep 20
     Should Be Equal As Integers  ${rc}  0

     :FOR  ${idx}  IN RANGE  1  5
     \   Start Process  docker %{VCH-PARAMS} exec ${id} /bin/ls  alias=exec-%{VCH-NAME}-${idx}  shell=true

     ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} wait ${id}
     Should Be Equal As Integers  ${rc}  0

     ${combinedErr}=  Set Variable
     ${combinedOut}=  Set Variable

     :FOR  ${idx}  IN RANGE  1  5
     \   ${result}=  Wait For Process  exec-%{VCH-NAME}-${idx}  timeout=2 mins
     \   ${combinedErr}=  Catenate  ${combinedErr}  ${result.stderr}${\n}
     \   ${combinedOut}=  Catenate  ${combinedOut}  ${result.stdout}${\n}

     Should Contain  ${combinedErr}  Container (${id}) is not running

     # We should get atleast one successful exec...
     Should Contain  ${combinedOut}  bin
     Should Contain  ${combinedOut}  dev
     Should Contain  ${combinedOut}  etc
     Should Contain  ${combinedOut}  home
     Should Contain  ${combinedOut}  lib
     Should Contain  ${combinedOut}  lost+found
     Should Contain  ${combinedOut}  mnt
     Should Contain  ${combinedOut}  proc
     Should Contain  ${combinedOut}  root
     Should Contain  ${combinedOut}  run
     Should Contain  ${combinedOut}  sbin
     Should Contain  ${combinedOut}  sys
     Should Contain  ${combinedOut}  tmp
     Should Contain  ${combinedOut}  usr
     Should Contain  ${combinedOut}  var
