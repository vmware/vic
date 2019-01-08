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
Documentation    This resource contains keywords which are helpful for testing storage quota.

*** Keywords ***
Get storage quota limit and usage
    [Arguments]  ${info}

    ${limitline}=  Get Lines Containing String  ${info}  VCH storage limit:
    ${usageline}=  Get Lines Containing String  ${info}  VCH storage usage:
    @{limitline}=  Split String  ${limitline}
    Length Should Be  ${limitline}  5
    @{usageline}=  Split String  ${usageline}
    Length Should Be  ${usageline}  5
    ${limitval}=  Convert To Number  @{limitline}[3]
    ${usageval}=  Convert To Number  @{usageline}[3]

    [Return]  ${limitval}  ${usageval}

Get storage usage
    [Arguments]  ${info}

    ${usageline}=  Get Lines Containing String  ${info}  VCH storage usage:
    @{usageline}=  Split String  ${usageline}
    Length Should Be  ${usageline}  5
    ${usageval}=  Convert To Number  @{usageline}[3]

    [Return]  ${usageval}
