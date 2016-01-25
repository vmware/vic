#!/bin/bash

NAME=kernel

# set default SRCDIR & BINDIR for local builds
if [ "${SRCDIR}" == "" ]; then
  SRCDIR="$(cd `dirname "$0"` && pwd)"
fi

if [ "${BINDIR}" == "" ]; then
  BINDIR=${SRCDIR}/../binary
fi

export JOB=${JOB_NAME:-$NAME}_${BUILD_NUMBER:-local_build}
DATE=$(date -u +%Y/%m/%d_@_%H:%M:%S)

echo SRCDIR=${SRCDIR}
echo BINDIR=${BINDIR}
mkdir -p ${BINDIR}

git_args="--git-dir=$SRCDIR/.git --work-tree=$SRCDIR"
branch_name="$(git $git_args symbolic-ref HEAD 2>/dev/null)" ||
branch_name="detached_head"     # detached HEAD
BRANCH=${branch_name##refs/heads/}
SHA=$(git $git_args rev-parse --short HEAD)

docker build -t ${NAME}-build ${SRCDIR}
SUCCESS=$?
if [ $SUCCESS -eq 0 ]; then
  BUILD_ID="$DATE@$BRANCH:$SHA"
  docker run --name=$JOB -e NAME=$NAME -e BUILD_ID="$BUILD_ID" ${NAME}-build:latest && {
    docker cp ${JOB}:/usr/src/photon/RPMS/x86_64/linux-esx-4.2.0-10.x86_64.rpm ${BINDIR}
  }

  SUCCESS=$?
fi

# clean up now the build's complete
docker rm -v ${JOB}

if [ $SUCCESS -ne 0 ]; then
  echo "Build failed for $NAME: $SUCCESS"
fi

# make the return value for the script reflect the status
test $SUCCESS -eq 0
