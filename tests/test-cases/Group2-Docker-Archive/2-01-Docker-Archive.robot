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
Documentation  Test 2-01 - Docker Archive
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  additional-args=--debug=3
Suite Teardown  Cleanup VIC Appliance Together With tmp dir

*** Keywords ***
Cleanup VIC Appliance Together With tmp dir 
    Cleanup VIC Appliance On Test Server
    Run  rm -rf /tmp/compare
    Run  rm -rf /tmp/pull
    Run  rm -rf /tmp/save

Compare Tar File Content
    [Arguments]  ${tarA}  ${tarB}
    ${rc}  ${output}=  Run And Return Rc And Output  tar -tvf ${tarA} > /tmp/compare/a
    Should Be Equal As Integers  ${rc}  0
    Should Be Empty  ${output}
    ${out}=  Run  cat /tmp/compare/a
    Log  ${out}

    ${rc}  ${output}=  Run And Return Rc And Output  tar -tvf ${tarB} > /tmp/compare/b
    Should Be Equal As Integers  ${rc}  0
    Should Be Empty  ${output}
    ${out}=  Run  cat /tmp/compare/b
    Log  ${out}

    ${rc}  ${output}=  Run And Return Rc And Output  diff /tmp/compare/a /tmp/compare/b
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    Should Be Empty  ${output}

Compare Files Digest in Tar
    [Arguments]  ${tarA}  ${tarB}
    Run  mkdir /tmp/compare/fileA
    Run  mkdir /tmp/compare/fileB
    ${rc}  ${output}=  Run And Return Rc And Output  tar -xvf ${tarA} -C /tmp/compare/fileA
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  tar -xvf ${tarB} -C /tmp/compare/fileB
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    ${out}=  Run  find /tmp/compare/fileA -type f
    ${files}=  Split To Lines  ${out}
    :FOR  ${fileA}  IN  @{files}
    \   ${fileB}=  Replace String  ${fileA}  /tmp/compare/fileA  /tmp/compare/fileB
    \   Log  ${fileB}
    \   ${rc}  ${output}=  Run And Return Rc And Output  sha256sum ${fileA}
    \   Should Be Equal As Integers  ${rc}  0
    \   ${digestA}=  Split String  ${output}
    \   ${rc}  ${output}=  Run And Return Rc And Output  sha256sum ${fileB}
    \   Should Be Equal As Integers  ${rc}  0
    \   ${digestB}=  Split String  ${output}
    \   Should Be Equal As Strings  @{digestA}[0]  @{digestB}[0]

*** Test Cases *** 
Docker Archive Download
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull ${ubuntu}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    Log To Console  \n imagec pull ${ubuntu}
    ${rc}  ${output}=  Run And Return Rc And Output  bin/imagec -insecure-skip-verify -reference ${ubuntu} -destination /tmp/pull -standalone -debug -operation pull
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

    Log To Console  \n imagec save ${ubuntu}
    ${imagestore}=  Run  govc datastore.ls %{VCH-NAME}/VIC/
    ${rc}  ${output}=  Run And Return Rc And Output  bin/imagec -insecure-skip-verify -reference ${ubuntu} -destination /tmp/save -standalone -debug -operation save -host %{VCH-IP}:2380 -image-store ${imagestore}
    Log  ${output}
    Should Be Equal As Integers  ${rc}  0

Compare Tar Files
    ${out}=  Run  find /tmp/pull -name *.tar
    ${files}=  Split To Lines  ${out}
    Run  mkdir /tmp/compare
    :FOR  ${pullfile}  IN  @{files}
    \   ${pullbase}=  Run  basename ${pullfile}
    \   ${rc}  ${savefile}=  Run And Return Rc And Output  find /tmp/save -name ${pullbase}
    \   Should Be Equal As Integers  ${rc}  0
    \   Compare Tar File Content  ${pullfile}  ${savefile}
    \   Compare Files Digest in Tar  ${pullfile}  ${savefile}

# TODO: refactor dir to use image name, add more image test
