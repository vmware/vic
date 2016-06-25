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

dpkg -l > package.list

if [ $1 = "skip-docker" ]; then
    pybot --removekeywords TAG:secret --exclude docker tests/test-cases
else
    pybot --removekeywords TAG:secret tests/test-cases
fi

rc="$?"

timestamp=$(date +%s)
outfile="integration_test_logs_$timestamp.tar"

tar cf $outfile log.html package.list *container-logs.tar.gz

if [ -f imagec.log ]; then
  tar --append imagec.log -f $outfile
fi

# GC credentials
set +x
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
set -x

if [ -f "$outfile" ]; then
  gsutil cp $outfile gs://vic-ci-logs
else
  echo "No log output file to upload"
fi

if [ -f "$keyfile" ]; then
  rm -f $keyfile
fi

exit $rc
