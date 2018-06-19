#!/bin/bash
# Copyright 2018 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License. 

DEFAULT_API_ENDPOINT="https://api.github.com/repos/"
DEFAULT_HEADERS=("Accept: application/vnd.github.symmetra-preview+json")
DEFAULT_CURL_ARGS=("-s")
DEFAULT_REPO="vmware/vic"

API_ENDPOINT=${API_ENDPOINT:-${DEFAULT_API_ENDPOINT}}
HEADERS=${HEADERS:-${DEFAULT_HEADERS}}
CURL_ARGS=${CURL_ARGS:-${DEFAULT_CURL_ARGS}}
REPO=${REPO:-${DEFAULT_REPO}}

# Determines whether a label already exists
#
# Arguments:
# 1: the label name
#
# Returns:
# N/A
#
# Exits:
# 0: the label exists
# 1: the label does not exist
label-exists () {
    header_args=("${HEADERS[@]/#/"-H '"}")
    header_args=("${header_args/%/"'"}")
    args=("-w %{http_code}\n" "${header_args[@]}" "${CURL_ARGS[@]}")
    code=$(curl "${args[@]}" ${API_ENDPOINT%/}/${REPO}/labels/$1 | tail -n1)
    
    [ $code -eq 200 ]
}


