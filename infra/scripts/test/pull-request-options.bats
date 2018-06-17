#!/usr/bin/env bats
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

load pull-request-options

# array of all known directives, used for iterating over the available set
# in tests
all_directives=("${directive_skip_unit}"
"${directive_focused_unit}"
"${directive_skip_functional}"
"${directive_all_functional}"
"${directive_specific_functional}"
"${directive_specific_integration}"
"${directive_all_integration}")

@test "get pr body from github" {
    body="$(get-pr-body 8064)"
    [ "${#body}" -ge 0 ]
}

@test "get directives - single enabled set" {
    local original_function=$(declare -f get-pr-body)

    for directive in ${all_directives[@]}; do
        echo "# testing directive ${directive}"
        get-pr-body () {
            echo "[X] <!-- directive:${directive} -->"
        }

        directives=$(get-enabled-pr-directives mock-pr)

        echo "# directives=${directives[@]}"
        [ "${directives[@]}" == "${directive}" ]
    done
}


@test "get directives - single disabled set" {
    local original_function=$(declare -f get-pr-body)

    for directive in ${all_directives[@]}; do
        echo "# testing directive ${directive}"
        get-pr-body () {
            echo "[ ] <!-- directive:${directive} -->"
        }

        directives=$(get-enabled-pr-directives mock-pr)

        echo "# directives=${directives[@]}"
        [ "${directives[@]}" == "" ]
    done
}

@test "get directives - multiple enabled, single disabled" {
    local original_function=$(declare -f get-pr-body)

    get-pr-body () {
        for directive in ${directive_skip_unit} ${directive_skip_functional}; do
            echo "[X] <!-- directive:${directive} -->"
        done
        echo "[ ] <!-- directive:${directive_all_integration} -->"
    }

    directives=$(get-enabled-pr-directives mock-pr)

    echo "# directives=${directives[@]}"
    [ "${directives[@]}" == "${directive_skip_unit} ${directive_skip_functional}" ]
}


@test "get directives - single with parameters" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_functional}

    get-pr-body () {
        echo "[X] <!-- directive:${label}${directive_parameter_begin} -->"
        echo '```'
        echo "testA"
        echo "groupB/testC"
        echo '```'
        echo "<!-- directive:${label}${directive_parameter_end} -->"
    }

    directives=$(get-enabled-pr-directives mock-pr)

    echo "# directives=${directives[@]}"
    [ "${directives[@]}" == "${label}" ]
}

@test "get parameters - get parameters for single directive" {
    local label=${directive_specific_functional}

    test_parameters=("testA" "groupB/testC")
    get-pr-body () {
        echo "[X] <!-- directive:${label}${directive_parameter_begin} -->"
        echo '```'
        for param in "${test_parameters[@]}"; do
            echo "${param}"
        done
        echo '```'
        echo "<!-- directive:${label}${directive_parameter_end} -->"
    }

    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"
    result="${parameters[@]}"
    reference="${test_parameters[@]}"

    [ "${result}" == "${reference}" ]
}

@test "get parameters - get parameters for multiple directives" {
    local labelA=${directive_specific_functional}
    local labelB=${directive_specific_integration}

    test_parametersA=("testA" "groupB/testC")
    test_parametersB=("testB" "groupD/testE")
    get-pr-body () {
        echo "[X] <!-- directive:${labelA}${directive_parameter_begin} -->"
        echo '```'
        for param in "${test_parametersA[@]}"; do
            echo "${param}"
        done
        echo '```'
        echo "<!-- directive:${labelA}${directive_parameter_end} -->"
        echo ""
        echo "garbage line in PR body"
        echo ""
        echo "[X] <!-- directive:${labelB}${directive_parameter_begin} -->"
        echo '```'
        for param in "${test_parametersB[@]}"; do
            echo "${param}"
        done
        echo '```'
        echo "<!-- directive:${labelB}${directive_parameter_end} -->"
    }

    parametersA=$(get-pr-directive-parameters mock-pr ${labelA})
    parametersB=$(get-pr-directive-parameters mock-pr ${labelB})

    echo "# parametersA=${parametersA[@]}"
    echo "# parametersB=${parametersB[@]}"

    resultA="${parametersA[@]}"
    referenceA="${test_parametersA[@]}"
    resultB="${parametersA[@]}"
    referenceB="${test_parametersA[@]}"

    [ "${resultA}" == "${referenceA}" ]
    [ "${resultB}" == "${referenceB}" ]
}


@test "multiple occurrances - last wins" {
    skip
}

@test "multiple occurrances - with params and without" {
    skip
}
