#!/bin/bash -x
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
#

# set to true to enable additional test logging for bats runs
my-debug () {    
#    echo "$@" >&3
    return
}

# explicit directives we know how to process
directive_skip_unit="skip-unit"
directive_focused_unit="focused-unit"
directive_skip_integration="skip-integration"
directive_all_integration="all-integration"
directive_specific_integration="specific-integration"
directive_specific_scenario="specific-scenario"
directive_all_scenario="all-scenario"
directive_fast_fail="fast-fail"
directive_shared_datastore="shared-datastore"
directive_ops_user="ops-user"
directive_parallel_jobs="parallel-jobs"

# the suffixes combined with directives for parameter blocks
directive_parameter_begin="-begin"
directive_parameter_end="-end"
directive_parameter_inline="-inline"


# Returns the comment body of a github PR.
# Arguments
# 1: the PR number
#
# Returns:
# body of the PR as a string
get-pr-body () {
    ${GITHUB_AUTOMATION_API_KEY:?Automation key must be provided in environment}

    declare -A cachedPrs

    body=cachedPrs[$1]
    body="${body:-$(curl -q https://api.github.com/repos/vmware/vic/pulls/${1}?access_token=${GITHUB_AUTOMATION_API_KEY} 2>/dev/null | jq -r '.body')}"
    cachedPrs[$1]=body
    
    echo "${body}"
}

# Returns the enabled (checked) directives for a PR
# Arguments
# 1: the PR number
#
# Returns:
# string array of enabled directives
get-enabled-pr-directives () {
    local body="$(get-pr-body $1)"
    declare -a matches
    while read line; do
        processed=$(echo ${line} | sed -n 's/\[X\] <!-- directive:\(.*\) -->/\1/p')
        if [ -n "${processed}" ]; then
            match=$(echo ${processed} | sed -re "s/${directive_parameter_begin}|${directive_parameter_end}//")

            # this is "returned" to the calling function
            my-debug "# adding directive: ${match}"
            matches+=("${match}")
        fi
    done <<<"${body}"

    echo "${matches[@]}"
    return
}


# Returns any parameters provided with the directive. These are present in the body as a newline separated list.
# Arguments
# 1: the PR number
# 2: the directive label
#
# Returns:
# string array of parameters
get-pr-directive-parameters () {
    local body="$(get-pr-body $1)"
    local label="$2"
    local accumulating=false
    declare -a matches
    declare -a parameters

    while read line; do
        # get the begin/end directive pairs

        # TODO: check for the directive occurring more than once
        # TODO: check for the directive occurring without parameters
        if ! ${accumulating}; then
            local processed=$(echo ${line} | sed -n 's/\[X\] <!-- directive:\(.*\) -->\(.*\)/\1/p')

            if [ "${label}${directive_parameter_begin}" == "${processed}" ]; then
                accumulating=true
            fi

            if [ "${label}${directive_parameter_inline}" == "${processed}" ]; then
                local trailing=$(echo ${line} | sed -n 's/\[X\] <!-- directive:\(.*\) -->\(.*\)/\2/p')
                # note that this is not appending to parameters, but replacing them
                parameters=$(get-parameter-from-inline-content "$trailing")

                my-debug "# assigning parameter: ${parameters[@]}"
            fi

            # TODO: check for the inline single parameter case
            # e.g. [X] <!-- directive:parallel-jobs-inline -->parallel-jobs=`4`
            continue
        fi

        local processed=$(echo ${line} | sed -n 's/<!-- directive:\(.*\) -->/\1/p')
        if [ "${label}${directive_parameter_end}" == "${processed}" ]; then
            # no longer in the parameter block, but keep processing as we may have multiple instances
            # of the directive
            accumulating=false
            continue
        fi

        if [ "${line}" == "\`\`\`" ]; then
            continue
        fi

        my-debug "# accumulating parameter: ${line}"
        parameters+=(${line})
    done <<<"${body}"

    echo "${parameters[@]}"
    return
}

# Extracts the parameter value for an inline parameter. This is needed as there is likely decoration around the actual
# value. To do so we will have to define some rules about how the value is separated from both prefix and suffix
# decoration.
# For now we assume anything in ` is value and if there are multiple of those we use the first only
# Arguments:
# 1: all line content after the directive
get-parameter-from-inline-content () {
    local inline="'${*}'"
    inline="${inline#*\`}"
    inline="${inline%%\`*}"

    echo ${inline}
}