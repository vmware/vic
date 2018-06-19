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
DEFAULT_REPO="vmware/vic-tasks"

API_ENDPOINT=${API_ENDPOINT:-${DEFAULT_API_ENDPOINT}}
HEADERS=("${HEADERS:-${DEFAULT_HEADERS}}")
CURL_ARGS=${CURL_ARGS:-${DEFAULT_CURL_ARGS}}
REPO=${REPO:-${DEFAULT_REPO}}

HEADERS=("${HEADERS[@]}" "Authorization: token ${GITHUB_TOKEN?"GitHub API token must be supplied"}")
HEADER_ARGS=("${HEADERS[@]/#/"-H"}")

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
    : ${2?"Usage: ${FUNCNAME[0]} LABEL DESCRIPTION [COLOR]"}

    if [ -z $3 ]
    then
        data="{\"description\": \"$2\"}"
    else
        data="{\"description\": \"$2\", \"color\": \"$3\"}"
    fi
    args=("--data" "${data}" "-XPATCH" "-w %{http_code}\n" "${HEADER_ARGS[@]}" "${CURL_ARGS[@]}")
    code=$(curl "${args[@]}" "${API_ENDPOINT%/}/${REPO}/labels/$1" | tail -n1)
    
    [ $code -eq 200 ]
}

# Creates a label with the given description and color
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
label-create () {
    : ${2?"Usage: ${FUNCNAME[0]} LABEL DESCRIPTION [COLOR]"}

    if [ -z $3 ]
    then
        data="{\"name\":\"$1\", \"description\": \"$2\"}"
    else
        data="{\"name\":\"$1\", \"description\": \"$2\", \"color\": \"$3\"}"
    fi
    args=("--data" "${data}" "-w %{http_code}\n" "${HEADER_ARGS[@]}" "${CURL_ARGS[@]}")
    code=$(curl "${args[@]}" "${API_ENDPOINT%/}/${REPO}/labels" | tail -n1)

    [ $code -eq 201 ]
}

# Creates a label with the given description and color, or updates one that exists
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
label-merge () {
    : ${2?"Usage: ${FUNCNAME[0]} LABEL DESCRIPTION [COLOR]"}

    if label-exists "$1"
    then
        label-update "$1" "$2" "$3"
    else
        label-create "$1" "$2" "$3"
    fi
}

merge-kinds () {
    typeset -A kinds
    kinds=(
        [debt]="Problems that increase the cost of other work"
        [defect]="Behavior that is inconsistent with what's intended"
        [defect_performance]="Behavior that is functionally correct, but performs worse than intended"
        [defect_regression]="Changed behavior that is inconsistent with what's intended"
        [defect_security]="A flaw or weakness that could lead to a violation of security policy"
        [enhancement]="Behavior that was intended, but we want to make better"
        [feature]="New functionality you could include in marketing material"
        [task]="Work not related to changing the functionality of the product"
        [question]="A request for information"
        [investigation]="A scoped effort to learn the answers to a set of questions which may include prototyping"
    )

    for label in "${!kinds[@]}"; do
        name="kind/${label/_/\/}"
        description="${kinds[$label]}"

        label-merge "$name" "$description" "bfd4f2"
    done
}

merge-source () {
    typeset -A sources
    sources=(
        [ci]="Found via a continuous integration failure"
        [customer]="Reported by a customer, directly or via an intermediary"
        [dogfooding]="Found via a dogfooding activity"
        [longevity]="Found via a longevity failure"
        [nightly]="Found via a nightly failure"
        [system-test]="Reported by the system testing team"
        [performance]="Reported by the performance testing team"
    )

    for label in "${!sources[@]}"; do
        name="source/${label/_/\/}"
        description="${sources[$label]}"

        label-merge "$name" "$description" "f9d0c4"
    done
}
