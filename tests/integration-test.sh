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
#

set -x
gsutil version -l
set +x

dpkg -l > package.list


buildinfo=$(drone build info vmware/vic $DRONE_BUILD_NUMBER)

if [[ $DRONE_BRANCH == "master" || $DRONE_BRANCH == *"refs/tags"* || $DRONE_BRANCH == "releases/"* ]] && [[ $DRONE_REPO == "vmware/vic" ]] && [[ $DRONE_BUILD_EVENT == "push" || $DRONE_BUILD_EVENT == "tag" ]]; then
    echo "Running full CI for $DRONE_BUILD_EVENT on $DRONE_BRANCH"
	pybot --removekeywords TAG:secret --exclude skip tests/test-cases
elif grep -q "\[full ci\]" <(drone build info vmware/vic $DRONE_BUILD_NUMBER); then
    echo "Running full CI as per commit message"
	pybot --removekeywords TAG:secret --exclude skip tests/test-cases
elif (echo $buildinfo | grep -q "\[specific ci="); then
    echo "Running specific CI as per commit message"
	buildtype=$(echo $buildinfo | grep "\[specific ci=")
    testsuite=$(echo $buildtype | awk -v FS="(=|])" '{print $2}')
    pybot --removekeywords TAG:secret --suite $testsuite --suite 7-01-Regression tests/test-cases
else
        echo "Running regressions"
    pybot --removekeywords TAG:secret --exclude skip --include regression tests/test-cases
fi

rc="$?"

timestamp=$(date +%s)
outfile="integration_logs_"$DRONE_BUILD_NUMBER"_"$DRONE_COMMIT".zip"

zip -9 $outfile output.xml log.html report.html package.list *container-logs.zip *.log

# GC credentials
keyfile="/root/vic-ci-logs.key"
botofile="/root/.boto"
echo -en $GS_PRIVATE_KEY > $keyfile
chmod 400 $keyfile
echo "[Credentials]" >> $botofile
echo "gs_service_key_file = $keyfile" >> $botofile
echo "gs_service_client_id = $GS_CLIENT_EMAIL" >> $botofile
echo "[GSUtil]" >> $botofile
echo "content_language = en" >> $botofile
echo "default_project_id = $GS_PROJECT_ID" >> $botofile


if [ -f "$outfile" ]; then
  gsutil cp $outfile gs://vic-ci-logs
  source "$(dirname "${BASH_SOURCE[0]}")/save-test-results.sh"

  echo "----------------------------------------------"
  echo "View test logs:"
  echo "https://vic-logs.vcna.io/$DRONE_BUILD_NUMBER/"
  echo "Download test logs:"
  echo "https://console.cloud.google.com/m/cloudstorage/b/vic-ci-logs/o/$outfile?authuser=1"
  echo "----------------------------------------------"
else
  echo "No log output file to upload"
fi

if [ -f "$keyfile" ]; then
  rm -f $keyfile
fi

exit $rc
