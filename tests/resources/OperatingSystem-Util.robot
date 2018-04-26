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
Documentation    This resource extends functionality from the OperatingSystem library to simplify common patterns.
Library          OperatingSystem


*** Keywords ***

Run And Verify Rc
    [Arguments]    ${command}    ${expectedRC}=0

    ${rc}    ${output}=    Run And Return Rc And Output    ${command}
    Log    ${output}
    Should Be Equal As Integers    ${rc}    ${expectedRC}

    [Return]    ${output}
