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
Documentation    This resource contains keywords which are helpful dealing with HTTP services


*** Variables ***
${STATUS}        The HTTP status of the last curl invocation


*** Keywords ***
Verify Status
    [Arguments]    ${expected}
    Should Be Equal As Integers    ${expected}    ${STATUS}

Verify Status Ok
    Verify Status    200

Verify Status Created
    Verify Status    201

Verify Status Accepted
    Verify Status    202

Verify Status Bad Request
    Verify Status    400

Verify Status Not Found
    Verify Status    404

Verify Status Unprocessable Entity
    Verify Status    422

Verify Status Internal Server Error
    Verify Status    500

