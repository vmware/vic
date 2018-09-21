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
set -e -o pipefail +h && [ -n "$DEBUG" ] && set -x

IFS=$'\n'

ENFORCE=('cmd/' 'lib/' 'pkg/')
RULES_FILE=.godeps_rules

# Returns a list of all packages under the supplied list of directories
# Arguments
# *: directories to search under
#
# Returns:
# go packages under those directories
find-packages () {
    find "${@}" -type f -name '*.go' -exec dirname {} \; | sort --unique
}

# Returns a list of all rules files
# Arguments
# *: directories to search under
#
# Returns:
# rule files within searched directories
find-rules () {
    find "${@}" -type f -name "$RULES_FILE" | sort --unique
}


# Returns the path to the "nearest" rule file for a given package
# Arguments
# 1: package path
#
# Returns:
# path to rule file (defaulting to /dev/null)
find-rule () {
    local path="${1?Package path must be provided}"
    shift

    while [[ "$path" != "." ]]; do
        if [ -e "$path/$RULES_FILE" ]; then
            echo "$path/$RULES_FILE"
            return
        fi

        path="$(dirname "$path")"
    done

    echo /dev/null
}

# Returns the rules from the specified rules file, omitting comments
# Arguments
# 1: path to rules file
#
# Returns:
# contents of the file, excluding blank lines and lines beginning with "#"
get-rules () {
    local rules="${1?Rules file must be provided}"
    shift

    grep -v -e '^$' -e '^#' "$rules"
}

# Returns all dependencies for a given package
# Arguments
# 1: package path
#
# Returns:
# a list of direct and transitive dependencies
get-deps () {
    local package="${1?Package must be provided}"
    shift

    infra/scripts/go-deps.sh "$package"
}

# Returns any invalid dependencies for a given package by filtering the full set
# of dependencies based on the supplied rules
# Arguments
# 1: package path
# 2: path to rules file
#
# Returns:
# a list of invalid dependencies
filter-deps () {
    local package="${1?Package must be provided}"
    local rules="${2?Rules file must be provided}"
    shift 2

    get-deps "$package" | grep -v -e "^$package/*" -f <(get-rules "$rules") || true
}

# Returns any unused rules for a given set of packages (does not identify duplicate rules!)
# Arguments
# 1: the dependencies
# 2: path to rules file
#
# Returns:
# a list of rules which match none of the dependencies
filter-rules () {
    local deps="${1?Dependencies must be provided}"
    local rules="${2?Rules file must be provided}"
    shift 2

    local unmatched=()
    for rule in $(get-rules "$rules"); do
        matches=$(echo "$deps" | grep -c -e "$rule" || true)

        if [ "$matches" = "0" ]; then
            unmatched+=("$rule")
        fi
    done

    echo "${unmatched[*]}"
}


rc=0

for package in $(find-packages "${ENFORCE[@]}"); do
    rules="$(find-rule "$package")"

    invalid=$(filter-deps "$package" "$rules")

    if [ ! -z "$invalid" ]; then
        echo "Unexpected dependencies in $package:"
        echo "$invalid" | sed -e 's/^/  /'
        echo "See $rules for details."
        echo ""

        rc=$(( rc | 1 ))
    fi
done

export VIC_CACHE_DEPS=true # We've just calculated all of these; no need to do so again

for rules in $(find-rules "${ENFORCE[@]}"); do
    directory=$(dirname "$rules")

    deps=""
    for package in $(find-packages "$directory"); do
        if [ "$rules" != "$(find-rule "$package")" ]; then
            continue # This package is using a more specific rules file
        fi
        deps+="$(get-deps "$package")"
    done

    unused="$(filter-rules "$deps" "$rules")"

    if [ ! -z "$unused" ]; then
        echo "Unused rules in $rules:"
        echo "$unused" | sed -e 's/^/  /'
        echo ""

        rc=$(( rc | 2 ))
    fi
done

exit $rc
