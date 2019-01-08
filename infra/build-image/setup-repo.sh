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
echo "Usage: $0 -p pull-request OR -b github_user/branch OR -t tag OR -c commit" 1>&2
exit 1
}

args_num=0
while getopts "t:c:p:b:h" flag
do
    case $flag in

        t)
            # Optional. Tag name
            localbranch="$OPTARG"
            github_user=vmware
            (( args_num += 1 ))
            ;;
        c)
            # Optional. Commit ID
            localbranch="$OPTARG"
            github_user=vmware
            (( args_num += 1 ))
            ;;
        p)
            # Optional. Pull request number
            pull_req="$OPTARG"
            github_user=vmware

            refspec=refs/pull/${pull_req}/head:refs/remotes/origin/pr/${pull_req}
            localbranch=pr/${pull_req}
            (( args_num += 1 ))
            ;;

        b)
            # Optional. Branch specifier
            github_user=${OPTARG%%/*}
            branch=${OPTARG#*/}

            refspec=refs/heads/${branch}:refs/remotes/origin/${branch}
            localbranch=origin/${branch}
            (( args_num += 1 ))
            ;;
        h)
            usage
            ;;
        *)
            break
            ;;
    esac
done

if [ "$args_num" -gt 1 ]; then
    usage
fi

shift $((OPTIND-1))

mkdir -p ${SRCDIR:?Expected script to be run with SRCDIR set} && cd ${SRCDIR}
url=https://github.com/${github_user}/vic

if ! git remote show -n origin >/dev/null 2>&1 ; then
    git init . && git remote add origin ${url}
fi

if [ -n "${refspec}" ]; then
    # we don't limit the depth as that was resulting in failure to find the most recent tag for a given commit
    git fetch origin -v ${refspec}
    if [ "$(git rev-parse --abbrev-ref HEAD)" != "$branch" ]; then
        git checkout -b ${branch} ${localbranch}
    fi
fi

if [ -n "${localbranch}" ]; then
    # try to fast-forward but just checkout if that fails
    git pull --ff-only || git checkout ${localbranch}
fi

if [ $# -ne 0 ]; then
    exec bash -c "$*"
else
    # drop to interactive shell after running any 
    exec bash
fi