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

nightly_list_var="5-1-Distributed-Switch \
5-2-Cluster \
5-3-Enhanced-Linked-Mode \
5-4-High-Availability \
5-5-Heterogeneous-ESXi \
5-6-1-VSAN-Simple \
5-6-2-VSAN-Complex \
5-7-NSX \
5-8-DRS \
5-10-Multiple-Datacenter \
5-11-Multiple-Cluster \
5-12-Multiple-VLAN \
5-13-Invalid-ESXi-Install \
5-14-Remove-Container-OOB \
5-15-NFS-Datastore \
5-16-iSCSI-Datastore \
5-17-FC-Datastore \
5-21-Datastore-Path \
5-22-NFS-Volume \
5-24-Non-vSphere-Local-Cluster \
5-25-OPS-User-Grant \
13-1-vMotion-VCH-Appliance \
13-2-vMotion-Container \
21-1-Whitelist \
21-2-Artifactory"

numberOfTests=($nightly_list_var)
numberOfTests=${#numberOfTests[@]}

input=$(gsutil ls -l gs://vic-engine-builds/vic_* | grep -v TOTAL | sort -k2 -r | head -n1 | xargs | cut -d ' ' -f 3 | cut -d '/' -f 4)
buildNumber=${input:4}

n=0
   until [ $n -ge 5 ]
   do
      echo "Retry.. $n"
      echo "Downloading gcp file $input"
      wget https://storage.googleapis.com/vic-engine-builds/$input
      if [ -f "$input" ]
      then
      echo "File found.."
      break
      else
      echo "File NOT found"
      fi
      n=$[$n+1]
      sleep 15
   done

n=0
   until [ $n -ge 5 ]
   do
      mkdir bin
      echo "Extracting .tar.gz"
      tar xvzf $input -C bin/ --strip 1
      if [ -f "bin/vic-machine-linux" ]
      then
      echo "tar extraction complete.."
      canContinue="Yes"
      break
      else
      echo "tar extraction failed"
      canContinue="No"
      rm -rf bin
      fi
      n=$[$n+1]
      sleep 15
   done

if [ $canContinue = "No" ]
then
echo "Tarball extraction failed..quitting the run"
break
else
echo "Tarball extraction passed, Running nightlies test.."

echo "Deleting .tar.gz vic file"
rm $input

DATE=`date +%m_%d_%H_%M_`

nightlystatus=()
count=0

# There should not be any VMs existing prior to running this test
sshpass -p $NIMBUS_PASSWORD ssh -o StrictHostKeyChecking=no $NIMBUS_USER@$NIMBUS_GW nimbus-ctl kill '\*'

for i in $nightly_list_var; do
    #Clean up any previous runs creds
    rm -rf VCH-0-*
    echo "Executing nightly test $i vSphere 6.5"
    pybot --removekeywords TAG:secret -d 65/$i --suite $i tests/manual-test-cases/

    if [ $? -eq 0 ]
    then
    echo "Passed"
    nightlystatus[$count]="Pass"
    else
    echo "Failed"
    nightlystatus[$count]="FAIL"
    fi

    mv *.log 65/$i
    mv *.zip 65/$i
    ((count++))
    echo $count
done

# See if any VMs leaked and clean them up if so
sshpass -p $NIMBUS_PASSWORD ssh -o StrictHostKeyChecking=no $NIMBUS_USER@$NIMBUS_GW nimbus-ctl list
sshpass -p $NIMBUS_PASSWORD ssh -o StrictHostKeyChecking=no $NIMBUS_USER@$NIMBUS_GW nimbus-ctl kill '\*'

for i in $nightly_list_var; do
    #Clean up any previous runs creds
    rm -rf VCH-0-*
    echo "Executing nightly test $i on vSphere 6.0"
    pybot --removekeywords TAG:secret --variable ESX_VERSION:ob-5251623 --variable VC_VERSION:ob-5112509 -d 60/$i --suite $i tests/manual-test-cases/

    if [ $? -eq 0 ]
    then
    echo "Passed"
    nightlystatus[$count]="Pass"
    else
    echo "Failed"
    nightlystatus[$count]="FAIL"
    fi

    mv *.log 60/$i
    mv *.zip 60/$i
    ((count++))
    echo $count
done

# See if any VMs leaked
sshpass -p $NIMBUS_PASSWORD ssh -o StrictHostKeyChecking\=no $NIMBUS_USER@$NIMBUS_GW nimbus-ctl list

# Setting the NSX test status to Not Implemented.
nightlystatus[7+$numberOfTests]="N/A"

for i in "${nightlystatus[@]}"
do
    echo $i
    if [ $i = "Pass" ]
    then
    buildStatus="Passed"
    echo "Test Passed!"
    else
    buildStatus="Failed!"
    echo "Test failed, setting global test status to Failed!"
    break
    fi
done

echo "Global Nightly Test Status $buildStatus"

sh tests/nightly/upload-logs.sh $DATE$input

rm nightly_mail.html

nightly_list_var=($nightly_list_var)
cat <<EOT >> nightly_mail.html
To: mwilliamson-staff-adl@vmware.com
To: rashok@vmware.com
Subject: VIC Engine Nightly Build $buildNumber
From: VIC Nightly
MIME-Version: 1.0
Content-Type: text/html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <style>
	    tr.d0 td {
	  background-color:#E0E0E0;
	  color: black;
	}
	tr.d1 td {
	  background-color:#FFFFFF;
	  color: black;
	}
        tr.d2 td {
	  background-color:#66c2ff;
	  color: black;
	}
     </style>
  </head>
  <body>
    <table class="body-wrap">
      <tr>
        <td></td>
        <td class="container" width="600">
          <div class="content">
            <table class="main" width="100%" cellpadding="0" cellspacing="0">
              <tr>
                <td class="content-wrap">
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr class="d2">
                      <td>
                      vSphere v6.5 - VIC Build $buildNumber
                      </td>
                      <td>
                      </td>
                    </tr>
                    `for ((i=0; i < ${#nightly_list_var[@]}; ++i)); do echo "<tr class=\"d$(($i%2))\"><td>${nightly_list_var[$i]}: </td><td>${nightlystatus[$i]}</td></tr>"; done`
                  </table>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        <a href='https://storage.cloud.google.com/vic-ci-logs/functional_logs_$DATE$input.zip?authuser=1'>https://storage.cloud.google.com/vic-ci-logs/functional_logs_$DATE$input.zip?authuser=1</a>
                      </td>
                    </tr>
                  </table>
                  <hr>
                </td>
              </tr>
            </table>
          </div>
        </td>
        <td></td>
      </tr>
    </table>
  <table class="body-wrap">
      <tr>
        <td></td>
        <td class="container" width="600">
          <div class="content">
            <table class="main" width="100%" cellpadding="0" cellspacing="0">
              <tr>
                <td class="content-wrap">
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr class="d2">
                      <td>
                      vSphere v6.0 - VIC Build $buildNumber
                      </td>
                      <td>
                      </td>
                    </tr>
                    `for ((i=0; i < ${#nightly_list_var[@]}; ++i)); do echo "<tr class=\"d$(($i%2))\"><td>${nightly_list_var[$i]}: </td><td>${nightlystatus[(($i+${#nightly_list_var[@]}))]}</td></tr>"; done`
                  </table>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        <a href='https://storage.cloud.google.com/vic-ci-logs/functional_logs_$DATE$input.zip?authuser=1'>https://storage.cloud.google.com/vic-ci-logs/functional_logs_$DATE$input.zip?authuser=1</a>
                      </td>
                    </tr>
                  </table>
                  <hr>
                </td>
              </tr>
            </table>
          </div>
        </td>
        <td></td>
      </tr>
    </table>
  </body>
</html>
EOT

# Emails an HTML report of the test run results using SendMail.
sendmail -t < nightly_mail.html
fi

# Saves test results to reporting server
testresultsdb="vic-nightly.db"
rm $testresultsdb
scp $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns-db/$testresultsdb .

for i in $nightly_list_var; do
python -m dbbot.run -k  -b $testresultsdb 65/$i/output.xml
ssh $REPORTING_USER@$REPORTING_SERVER_URL mkdir -p /export/drone-test-results/testruns/$buildNumber-nightly/65/$i
scp 65/$i/log.html $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns/$buildNumber-nightly/65/$i
scp 65/$i/report.html $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns/$buildNumber-nightly/65/$i

python -m dbbot.run -k  -b $testresultsdb 60/$i/output.xml
ssh $REPORTING_USER@$REPORTING_SERVER_URL mkdir -p /export/drone-test-results/testruns/$buildNumber-nightly/60/$i
scp 60/$i/log.html $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns/$buildNumber-nightly/60/$i
scp 60/$i/report.html $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns/$buildNumber-nightly/60/$i
done

scp $testresultsdb $REPORTING_USER@$REPORTING_SERVER_URL:/export/drone-test-results/testruns-db/$testresultsdb

