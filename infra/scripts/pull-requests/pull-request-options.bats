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
"${directive_skip_integration}"
"${directive_all_integration}"
"${directive_specific_integration}"
"${directive_specific_scenario}"
"${directive_all_scenario}")

# convert a boolean to the string used in PR body for enabled/disabled
# 1: string:[enabled|disabled] - is enabled
bool_to_directive_state () {
    if [ "$1" == "enabled" ]; then
        echo "X"
    else
        echo " "
    fi
}

# 1: string:[enabled|disabled] - is directive enabled
# 2: string - directive name
build_directive () {
    echo "[$(bool_to_directive_state $1)] <!-- directive:${2} -->"
}

# 1: [enabled|disabled] - is directive enabled
# 2: string - directive name
# *: parameters - each arg will be on a new line
build_directive_with_parameters () {
    directive=${2}
    echo "[$(bool_to_directive_state $1)] <!-- directive:${directive}${directive_parameter_begin} -->"
    shift 2
    echo '```'
    for param in ${@}; do
        echo ${param}
    done
    echo '```'
    echo "<!-- directive:${directive}${directive_parameter_end} -->"
}

# 1: [enabled|disabled] - is directive enabled
# 2: string - directive name
# 3: parameter - parameter to embed in line
build_directive_with_inline_parameter () {
    directive=${2}
    echo -n "[$(bool_to_directive_state $1)] <!-- directive:${directive}${directive_parameter_inline} -->"
    echo -n "prefix padding...: "
    echo -n '`'${3}'`'
    echo -n "... suffix padding...!"
}

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

# Repeats the same tests as above to assert the correctness of the helper
# functions for generating PR bodies. This test is for directives without
# parameters
@test "get directives - (testing helper function)" {
    local original_function=$(declare -f get-pr-body)

    for directive in ${all_directives[@]}; do
        echo "# testing directive ${directive}"
        get-pr-body () {
            build_directive enabled ${directive}
        }

        directives=$(get-enabled-pr-directives mock-pr)

        echo "# directives=${directives[@]}"
        [ "${directives[@]}" == "${directive}" ]
    done

    for directive in ${all_directives[@]}; do
        echo "# testing directive ${directive}"
        get-pr-body () {
            build_directive disabled ${directive}
        }

        directives=$(get-enabled-pr-directives mock-pr)

        echo "# directives=${directives[@]}"
        [ "${directives[@]}" == "" ]
    done
}

@test "get directives - multiple enabled, single disabled" {
    local original_function=$(declare -f get-pr-body)

    get-pr-body () {
        for directive in ${directive_skip_unit} ${directive_skip_integration}; do
            build_directive enabled ${directive}
        done
        build_directive disabled ${directive_all_scenario}
    }

    directives=$(get-enabled-pr-directives mock-pr)

    echo "# directives=${directives[@]}"
    [ "${directives[@]}" == "${directive_skip_unit} ${directive_skip_integration}" ]
}

# Most basic test for directive with parameters. Does not use the helper functions
# to explicitly control input.
@test "get directives - single with parameters" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_integration}

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

# Most basic test for parameters from directive. Does not use the helper functions
# to explicitly control input.
@test "get parameters - get parameters for single directive" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_integration}

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

# Asserts correctness of the helper function for directives with parameters, repeating
# the control tests for the parameter cases.
@test "get directives - (test parameter helper function)" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_integration}

    test_parameters=("testA" "groupB/testC")
    get-pr-body () {
        build_directive_with_parameters enabled ${label} "${test_parameters[@]}"
    }

    # check plain directive
    directives=$(get-enabled-pr-directives mock-pr)

    echo "# directives=${directives[@]}"
    [ "${directives[@]}" == "${label}" ]

    # check paramters
    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"
    result="${parameters[@]}"
    reference="${test_parameters[@]}"

    [ "${result}" == "${reference}" ]
}

@test "get parameter - inline parameter from line" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_parallel_jobs}

    test_parameter="testA"
    get-pr-body () {
        build_directive_with_inline_parameter enabled ${label} "${test_parameter}"
    }

    # check paramters
    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"
    result="${parameters[@]}"
    reference="${test_parameter}"

    [ "${result}" == "${reference}" ]
}

@test "get parameter - disabled inline parameter" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_parallel_jobs}

    test_parameter="testA"
    get-pr-body () {
        build_directive_with_inline_parameter disabled ${label} "${test_parameter}"
    }

    # check paramters
    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"
    result="${parameters[@]}"

    [ "${result}" == "" ]
}

@test "get parameters - get parameters for multiple directives" {
    local original_function=$(declare -f get-pr-body)
    local labelA=${directive_specific_integration}
    local labelB=${directive_specific_scenario}

    test_parametersA=("testA" "groupB/testC")
    test_parametersB=("testB" "groupD/testE")
    get-pr-body () {
        build_directive_with_parameters enabled ${labelA} "${test_parametersA[@]}"
        echo ""
        echo "garbage line in PR body"
        echo ""
        build_directive_with_parameters enabled ${labelB} "${test_parametersB[@]}"
        echo "trailing comment line in PR body"
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

# If multiple blocks of paramters for the same directive are present and enabled,
# the parameters returned should be the combined set.
@test "multiple occurrances - combined parameters, all selected" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_integration}

    test_parametersA=("testA" "groupB/testC")
    test_parametersB=("testB" "groupD/testE")

    get-pr-body () {
        build_directive_with_parameters enabled ${label} "${test_parametersA[@]}"
        echo ""
        echo "garbage line in PR body"
        echo ""
        build_directive_with_parameters enabled ${label} "${test_parametersB[@]}"
        echo "trailing comment line in PR body"
    }

    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"

    result="${parameters[@]}"
    reference="${test_parametersA[@]} ${test_parametersB[@]}"

    [ "${result}" == "${reference}" ]
}

@test "multiple occurrances - combined parameters, some not selected" {
    local original_function=$(declare -f get-pr-body)
    local label=${directive_specific_integration}

    test_parametersA=("testA" "groupB/testC")
    test_parametersB=("testB" "groupD/testE")
    test_parametersC=("testC" "groupF/testG")

    get-pr-body () {
        build_directive_with_parameters enabled ${label} "${test_parametersA[@]}"
        echo ""
        echo "garbage line in PR body"
        build_directive_with_parameters disabled ${label} "${test_parametersC[@]}"
        echo ""
        build_directive_with_parameters enabled ${label} "${test_parametersB[@]}"
        echo "trailing comment line in PR body"
    }

    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"

    result="${parameters[@]}"
    # should not include set C which is disabled
    reference="${test_parametersA[@]} ${test_parametersB[@]}"

    [ "${result}" == "${reference}" ]
}

@test "multiple occurrances - inline, last enabled wins" {
    local label=${directive_specific_integration}

    test_parameterA="testA"
    test_parameterB="testB"
    test_parameterC="testC"

    get-pr-body () {
        build_directive_with_inline_parameter enabled ${label} "${test_parameterA}"
        echo ""
        echo "garbage line in PR body"
        build_directive_with_inline_parameter enabled ${label} "${test_parameterB}"
        echo ""
        build_directive_with_inline_parameter disabled ${label} "${test_parameterC}"
        echo "trailing comment line in PR body"
    }

    parameters=$(get-pr-directive-parameters mock-pr ${label})

    echo "# parameters=${parameters[@]}"

    result="${parameters[@]}"
    # should not include set C which is disabled
    reference="${test_parameterB}"

    [ "${result}" == "${reference}" ]
}

@test "get directives - multiple instances, single enabled result" {
    skip
}

@test "get directives - is specific directive enabled" {
    skip
}
