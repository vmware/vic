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
Documentation  Test 5-14 - Remove Container OOB
Resource  ../../resources/Util.robot
#Suite Teardown  Run Keyword And Ignore Error  Nimbus Cleanup  ${list}

*** Test Cases ***
Docker run an image from a container that was removed OOB
    ${status}=  Get State Of Github Issue  3845
    Run Keyword If  '${status}' == 'closed'  Fail  Test 5-14-Remove-Container-OOB.robot needs to be updated now that Issue #3845 has been resolved
    #${esx1}  ${esx2}  ${esx3}  ${vc}  ${esx1-ip}  ${esx2-ip}  ${esx3-ip}  ${vc-ip}=  Create a Simple VC Cluster
    #Set Global Variable  @{list}  ${esx1}  ${esx2}  ${esx3}  ${vc}

    #Install VIC Appliance To Test Server

    #${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    #Should Be Equal As Integers  ${rc}  0
    #${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd --name removeOOB busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
    #Destroy VM OOB  removeOOB*
    #${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run -itd busybox /bin/top
    #Should Be Equal As Integers  ${rc}  0
