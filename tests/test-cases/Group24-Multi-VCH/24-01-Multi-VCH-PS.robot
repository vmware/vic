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
Documentation  Test 6-04 - Verify vic-machine create basic use cases
Resource  ../../resources/Util.robot
Test Teardown  Extra Cleanup

*** Keywords ***
Extra Cleanup
      Manually Cleanup VCH  %{GROUP-14-MULTI-VCH-PS}
      Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create Multi VCH - Docker Ps Only Contains The Correct Containers
    ${container1}=  Evaluate  'cvm-vch1-' + str(random.randint(1000,9999))  modules=random
    ${container2}=  Evaluate  'cvm-vch2-' + str(random.randint(1000,9999))  modules=random

    Install VIC Appliance To Test Server
    Set Suite Variable  ${old-vm}  %{VCH-NAME}
    Set Suite Variable  ${old-vch-params}  %{VCH-PARAMS}

    # Avoid deleting certs
    Remove Environment Variable  VCH-NAME

    Set Environment Variable  GROUP-14-MULTI-VCH-PS  ${old-vm}

    Add VCH to Removal Exception List  ${old-vm}
    Install VIC Appliance To Test Server
    Remove VCH from Removal Exception List  ${old-vm}

    ${rc}=  Run And Return Rc  docker ${old-vch-params} create --name ${container1} ${busybox}
    Should Be Equal As Integers  ${rc}  0
    ${rc}=  Run And Return Rc  docker %{VCH-PARAMS} create --name ${container2} ${busybox}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${old-vch-params} ps -a
    Should Contain  ${output}  ${container1}
    Should Not Contain  ${output}  ${container2}

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps -a
    Should Contain  ${output}  ${container2}
    Should Not Contain  ${output}  ${container1}
