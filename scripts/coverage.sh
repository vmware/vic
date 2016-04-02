#!/bin/sh
# Copyright 2016 VMware, Inc. All Rights Reserved.
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

# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` does not work
# with multiple packages, see https://code.google.com/p/go/issues/detail?id=6909
#
# Usage: script/coverage [--html]
#
#     --html        Create HTML report and open it in browser
#

set -e

workdir=../.cover
profile="$workdir/cover.out"
mode=count

generate_cover_data() {
    rm -rf "$workdir"
    mkdir "$workdir"
    echo "Generating coverage report for: $@"
    for dir in $@; do
      pkgs=$(go list $dir/... | grep -v /vendor/)
      for pkg in $pkgs; do
          f="$workdir/$(echo $pkg | tr / -).cover"
          go test -v -covermode="$mode" -coverprofile="$f" "$pkg"
      done
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" $workdir/*.cover >>"$profile"
}

show_cover_report() {
    go tool cover -${1}="$profile"
}

TEST_DIRS=$@
generate_cover_data $TEST_DIRS
show_cover_report func
