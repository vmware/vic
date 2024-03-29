#!/bin/bash
# Copyright 2016-2018 VMware, Inc. All Rights Reserved.
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

set -x
buildinfo=$(drone build info vmware/vic $DRONE_BUILD_NUMBER)
prNumber=$(drone build info --format "{{ .Ref }}" vmware/vic $DRONE_BUILD_NUMBER | cut -f 3 -d'/')
set +x
prBody=$(curl -H "Authorization: token $GITHUB_AUTOMATION_API_KEY" https://api.github.com/repos/vmware/vic/pulls/$prNumber | jq -r ".body")

if (echo $prBody | grep -q "\[fast fail\]"); then
    export FAST_FAILURE=1
else
    export FAST_FAILURE=0
fi

if (echo $prBody | grep -q "\[ops user\]"); then
    export RUN_AS_OPS_USER=1
else
    export RUN_AS_OPS_USER=0
fi

if (echo $prBody | grep -q "\[shared datastore="); then
    command=$(echo $prBody | grep "\[shared datastore=")
    datastore=$(echo $command | awk -F"\[shared datastore=" '{sub(/\].*/,"",$2);print $2}')
    export TEST_DATASTORE=$datastore
fi

bootstrapiso='bin/bootstrap-centos-6.9.iso'
bootstrapdir='centos6.9'
function test_custom_iso(){
    #Currently we only test 3 cases for custom-bootstrap-iso 
    pabot --verbose --processes $jobs --removekeywords TAG:secret -v BOOTSTRAP-ISO:$bootstrapiso --outputdir $bootstrapdir --suite 1-45-Docker-Container-Network --suite 1-19-Docker-Volume-Create --suite 7-01-Regression tests/test-cases
}

jobs="2"
if (echo $prBody | grep -q "\[parallel jobs="); then
    parallel=$(echo $prBody | grep "\[parallel jobs=")
    jobs=$(echo $parallel | awk -F"\[parallel jobs=" '{sub(/\].*/,"",$2);print $2}')
fi

if [[ $DRONE_BRANCH == "master" || $DRONE_BRANCH == "releases/"* ]] && [[ $DRONE_REPO == "vmware/vic" ]] && [[ $DRONE_BUILD_EVENT == "push" ]]; then
    echo "Running full CI for $DRONE_BUILD_EVENT on $DRONE_BRANCH"
    pabot --verbose --processes $jobs --removekeywords TAG:secret --exclude skip tests/test-cases && test_custom_iso
elif [[ $DRONE_REPO == "vmware/vic" ]] && [[ $DRONE_BUILD_EVENT == "tag" ]]; then
    echo "Running only Group11-Upgrade and 7-01-Regression for $DRONE_BUILD_EVENT on $DRONE_BRANCH"
    pabot --verbose --processes $jobs --removekeywords TAG:secret --suite Group11-Upgrade --suite 7-01-Regression tests/test-cases
elif (echo $prBody | grep -q "\[full ci\]"); then
    echo "Running full CI as per commit message"
    pabot --verbose --processes $jobs --removekeywords TAG:secret --exclude skip tests/test-cases && test_custom_iso
elif (echo $prBody | grep -q "\[specific ci="); then
    echo "Running specific CI as per commit message"
    buildtype=$(echo $prBody | grep "\[specific ci=")
    testsuite=$(echo $buildtype | awk -F"\[specific ci=" '{sub(/\].*/,"",$2);print $2}')
    pabot --verbose --processes $jobs --removekeywords TAG:secret --suite $testsuite --suite 7-01-Regression tests/test-cases && \
    pabot --verbose --processes $jobs --removekeywords TAG:secret -v BOOTSTRAP-ISO:$bootstrapiso --outputdir $bootstrapdir --exclude centos69-invalid --suite $testsuite --suite 7-01-Regression tests/test-cases
else
    echo "Running regressions"
    pabot --verbose --processes $jobs --removekeywords TAG:secret --exclude skip --include regression tests/test-cases && \
    pabot --verbose --processes $jobs --removekeywords TAG:secret -v BOOTSTRAP-ISO:$bootstrapiso --outputdir $bootstrapdir --exclude skip --include regression tests/test-cases
fi

rc="$?"

timestamp=$(date +%s)
outfile="integration_logs_"$DRONE_BUILD_NUMBER"_"$DRONE_COMMIT".zip"

( cd $bootstrapdir && for f in *; do mv $f ../$bootstrapdir"_"$f; done)
zip -9 -j $outfile $bootstrapdir"_"* output.xml log.html report.html package.list *container-logs*.zip *.log /var/log/vic-machine-server/vic-machine-server.log *.debug

LOG_SERVER_IP='10.158.214.81'
LOG_SERVER_USERNAME='root'
LOG_SERVER_PWD=${TEST_PASSWORD}
save_dir=${DRONE_BUILD_NUMBER}
unzip_cmd="unzip -q /var/www/html/${save_dir}/$outfile -d /var/www/html/${save_dir}/"
sshpass -p ${LOG_SERVER_PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${LOG_SERVER_USERNAME}@${LOG_SERVER_IP} mkdir -p /var/www/html/${save_dir}
sshpass -p ${LOG_SERVER_PWD} scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null $outfile ${LOG_SERVER_USERNAME}@${LOG_SERVER_IP}:/var/www/html/${save_dir}
sshpass -p ${LOG_SERVER_PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${LOG_SERVER_USERNAME}@${LOG_SERVER_IP} ${unzip_cmd}
sshpass -p ${LOG_SERVER_PWD} ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${LOG_SERVER_USERNAME}@${LOG_SERVER_IP} rm -rf /var/www/html/${save_dir}/$outfile

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
  ls -al $outfile
  n=0
  until [ $n -ge 5 ]
  do
    gsutil cp $outfile gs://vic-ci-logs && break
    n=$[$n+1]
    sleep 15
  done
  source "$(dirname "${BASH_SOURCE[0]}")/save-test-results.sh"

  echo "----------------------------------------------"
  echo "View test logs:"
#  echo "https://vic-logs.vcna.io/$DRONE_BUILD_NUMBER/"
  echo "Please visit to http://${LOG_SERVER_IP}/${save_dir}/ for looking over ci log."
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
