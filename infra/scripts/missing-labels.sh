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

HEADERS=("${HEADERS[@]}" "Authorization: token ${GITHUB_TOKEN?"GitHub API token must be supplied"}")

HEADER_ARGS=("${HEADERS[@]/#/"-H '"}")
HEADER_ARGS=("${HEADER_ARGS/%/"'"}")

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
    : ${1?"Usage: ${FUNCNAME[0]} LABEL"}

    args=("-w %{http_code}\n" "${HEADER_ARGS[@]}" "${CURL_ARGS[@]}")
    code=$(curl "${args[@]}" "${API_ENDPOINT%/}/${REPO}/labels/$1" | tail -n1)

    [ $code -eq 200 ]
}

# Updates the description and color associated with an existing label
#
# Arguments:
# 1: the label name
# 2: the label description
# 3: the label color
#
# Returns:
# N/A
#
# Exits:
# 0: the operation succeeded
# 1: the operation failed
label-update () {
    : ${2?"Usage: $0 LABEL DESCRIPTION [COLOR]"}

    if [ -z $3 ]
    then
        data="{\"description\": \"$2\"}"
    else
        data="{\"description\": \"$2\", \"color\": \"$3\"}"
    fi
    args=("--data" "${data}" "-X PATCH" "-w %{http_code}\n" "${HEADER_ARGS[@]}" "${CURL_ARGS[@]}")
    code=$(curl "${args[@]}" "${API_ENDPOINT%/}/${REPO}/labels/$1" | tail -n1)
    
    [ $code -eq 200 ]
}


