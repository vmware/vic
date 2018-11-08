#!/bin/sh
# Copyright 2017 VMware, Inc. All Rights Reserved.
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
# This file builds the docker images vic-build-image.tdnf vic-build-image.yum
# when there are changes to the dockerfiles that warrant an image update in gcr.
set -ex

REV=$(git rev-parse --verify --short=8 HEAD)
REGISTRY="gcr.io/eminent-nation-87317"
IMAGE="vic-local-repo"

for REPO in photon-2.0; do
    docker build -t "$IMAGE-$REPO:$REV" -f ./infra/integration-image/local-repo/Dockerfile.$REPO ./infra/integration-image/local-repo
    docker tag "$IMAGE-$REPO:$REV" "$REGISTRY/$IMAGE:$REPO"
    docker tag "$IMAGE-$REPO:$REV" "$REGISTRY/$IMAGE:$REPO-$REV"
    # gcloud docker -- push "$REGISTRY/$IMAGE:$REPO"
    # gcloud docker -- push "$REGISTRY/$IMAGE:$REPO-$REV"
done
