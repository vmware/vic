#!/bin/bash

[ -n "$DEBUG" ] && set -x

NAME=container

# set default SRCDIR & BINDIR for local builds
if [ "${SRCDIR}" == "" ]; then
  SRCDIR="$(cd `dirname "$0"` && pwd)"
fi

if [ "${BINDIR}" == "" ]; then
  BINDIR=${SRCDIR}/../binary
fi
BINBASE=$(basename $BINDIR)

# allow building of specific targets
if [ "$#" == 0 ]; then
  TARGETS="${SRCDIR}/targets/*"
elif [ "$#" == 1 ]; then
  TARGETS="${SRCDIR}/targets/$1"
else
  TARGETS="${SRCDIR}/targets/$1"
  shift
  for arg in "$@"; do
    TARGETS+=" ${SRCDIR}/targets/$arg"
  done
fi

export JOB=${JOB_NAME:-$NAME}_${BUILD_NUMBER:-local_build}
DATE=$(date -u +%Y/%m/%d_@_%H:%M:%S)

echo SRCDIR=${SRCDIR}
echo BINDIR=${BINDIR}
mkdir -p ${BINDIR}

for i in $TARGETS; do
  TNAME=${NAME}-$(basename $i)

  git_args="--git-dir=$SRCDIR/.git --work-tree=$SRCDIR"
  branch_name="$(git $git_args symbolic-ref HEAD 2>/dev/null)" ||
  branch_name="detached_head"     # detached HEAD
  BRANCH=${branch_name##refs/heads/}
  SHA=$(git $git_args rev-parse --short HEAD)

  cp -r ${SRCDIR}/tether ${i}/tether

  # if there is a base build, run that
  if [ -d ${i}/base ]; then
    docker build -t ${TNAME}-base ${i}/base || ( echo "Base build failed for $i" && break )
    if [ $? -ne 0 ]; then
      echo "Base build failed for $i: $?"
      break
    fi
  fi
  docker build --no-cache -t ${TNAME}-build ${i}
  SUCCESS=$?


  if [ $SUCCESS -eq 0 ]; then
    BUILD_ID="$DATE@$BRANCH:$SHA"
    docker run --name=$JOB-${TNAME} -e BUILD_ID=$BUILD_ID -e BINBASE=${BINBASE} -e ISOOUT=/binary/${TNAME}.iso ${TNAME}-build:latest && {
      docker cp ${JOB}-${TNAME}:/tmp/${BINBASE} ${BINDIR}/..
    }

    SUCCESS=$?
  fi

  # clean up now the build's complete
  rm -fr ${i}/tether
  docker rm -v ${JOB}-${TNAME}

  if [ $SUCCESS -ne 0 ]; then
    echo "Build failed for $i: $SUCCESS"
    break
  fi
done

# make the return value for the script reflect the status
test $SUCCESS -eq 0
