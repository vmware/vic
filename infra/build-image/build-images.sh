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
set -ex;

REV=$(git rev-parse --verify --short=8 HEAD)
REPO="gcr.io/eminent-nation-87317"
IMAGE="vic-build-image"

for DEP in yum tdnf; do
    docker build -t "$IMAGE-$DEP:$REV" -f ./infra/build-image/Dockerfile.$DEP .
    docker tag "$IMAGE-$DEP:$REV" "$REPO/$IMAGE:$DEP"
    docker tag "$IMAGE-$DEP:$REV" "$REPO/$IMAGE:$DEP-$REV"
    # gcloud docker -- push "$REPO/$IMAGE:$DEP"
    # gcloud docker -- push "$REPO/$IMAGE:$DEP-$REV"
done
