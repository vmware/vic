#!/bin/bash
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

# Snippets taken from tests/integration-test.sh
rc="$?"
outfile="vic_ui_build_"$DRONE_BUILD_NUMBER"_"$DRONE_COMMIT".log"

#GC credentials
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

if [ -f "vic_ui_build.log" ] ; then
  mv vic_ui_build.log $outfile
  gsutil cp $outfile gs://vic-ci-logs
else
  echo "No log output file to upload"
fi

if [ -f "$keyfile" ] ; then
  rm -f $keyfile
fi

exit $rc
