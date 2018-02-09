#!/bin/bash
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

# Need to determine what options we support and how they're specified for:
# * reviewing PRs
# * feature branches
# * master
# * custom branch from contributor repo
#   * not yet a PR
#   * used for review or for active development

# standard decoration: exit on failure and configure debug
set -e && [ -n "$DEBUG" ] && set -x
DIR=$(dirname $(readlink -f "$0"))

function usage() {
echo "Usage: $0 -p pull-request OR -b github_user/branch" 1>&2
exit 1
}

while getopts "p:b:" flag
do
    case $flag in

        p)
            # Optional. Pull request number
            pull_req="$OPTARG"
            github_user=vmware
            ;;

        b)
            # Optional. Branch specifier
            # bit of an abuse of dirname/basename but hey, / separators
            github_user=$(dirname $OPTARG)
            branch=$(basename $OPTARG)
            ;;

        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "${github_user}" -o -n "${pull_req}" -a -n "${branch}" ]; then
    usage
fi

mkdir -p ${SRCDIR:?Expected script to be run with SRCDIR set} && cd ${SRCDIR} && git init .
git remote add origin https://github.com/${github_user}/vic

if [ ! -z ${pull_req} ]; then
    git fetch origin --depth=5 -v refs/pull/${pull_req}/head:refs/remotes/origin/pr/${pull_req}
    git checkout pr/${pull_req}
elif [ ! -b ${branch} ]; then
    git fetch origin --depth=5 -v refs/heads/${branch}:refs/remotes/origin/${branch}
    git checkout origin/${branch}
fi

# drop to interactive shell
exec bash