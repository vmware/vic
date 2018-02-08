#!/bin/bash
sed -i 's|remotes/origin/*|remotes/origin/*\n\tfetch = +refs/pull/*/head:refs/remotes/origin/pr/|g' .git/config
git fetch origin && git checkout pr/$DRONE_PULL_REQUEST
