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
5-5-Heterogenous-ESXi \
5-6-1-VSAN-Simple \
5-6-2-VSAN-Complex \
5-7-NSX \
5-8-DRS \
5-10-Multiple-Datacenter \
5-11-Multiple-Cluster \
5-12-Multiple-VLAN \
5-13-Invalid-ESXi-Install \
5-14-Remove-Container-OOB \
13-1-vMotion-VCH-Appliance \
13-2-vMotion-Container"

echo "Removing VIC directory if present"
echo "Cleanup logs from previous run"

rm -rf *.zip *.log
rm -rf bin

for i in $nightly_list_var; do
    echo "Removing folder $i"
    rm -rf $i
done

input=$(wget -O - https://vmware.bintray.com/vic-repo |tail -n5 |head -n1 |cut -d':' -f 2 |cut -d'.' -f 3| cut -d'>' -f 2)
buildNumber=${input:4}

echo "Downloading bintray file $input"
wget https://vmware.bintray.com/vic-repo/$input.tar.gz

mkdir bin

echo "Extracting .tar.gz"
tar xvzf $input.tar.gz -C bin/ --strip 1

echo "Deleting .tar.gz vic file"
rm $input.tar.gz

nightlystatus=()
count=0

for i in $nightly_list_var; do
    echo "Executing nightly test $i"
    drone exec --trusted -e test="pybot -d $i --suite $i tests/manual-test-cases/" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

    if [ $? -eq 0 ]
    then
    echo "Passed"
    nightlystatus[$count]="PASS"
    else
    echo "Failed"
    nightlystatus[$count]="FAIL!"
    fi

    mv *.log $i
    mv *.zip $i
    ((count++))
    echo $count
done

for i in "${nightlystatus[@]}" 
do
    echo $i
    if [ $i = "PASS" ]
    then
    buildStatus=0
    echo "Test Passed!"
    else
    buildStatus=1
    echo "Test failed, setting global test status to Failed!"
    break
    fi	
done

echo "Global Nightly Test Status $buildStatus"

drone exec --trusted -e test="sh tests/nightly/upload-logs.sh $input" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

rm nightly_mail.html

if [ $buildStatus -eq 0 ]
then
echo "Success"
cat <<EOT >> nightly_mail.html
To: mwilliamson-staff-adl@vmware.com
To: rashok@vmware.com
Subject: VIC Nightly Run #$buildNumber
From: VIC Nightly
MIME-Version: 1.0
Content-Type: text/html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />

    <style>
      * {
        margin: 0;
        padding: 0;
        font-family: "Helvetica Neue", "Helvetica", Helvetica, Arial, sans-serif;
        box-sizing: border-box;
        font-size: 14px;
      }

      body {
        -webkit-font-smoothing: antialiased;
        -webkit-text-size-adjust: none;
        width: 100% !important;
        height: 100%;
        line-height: 1.6;
        background-color: #f6f6f6;
      }

      table td {
        vertical-align: top;
      }

      .body-wrap {
        background-color: #f6f6f6;
        width: 100%;
      }

      .container {
        display: block !important;
        max-width: 600px !important;
        margin: 0 auto !important;
        /* makes it centered */
        clear: both !important;
      }

      .content {
        max-width: 600px;
        margin: 0 auto;
        display: block;
        padding: 20px;
      }

      .main {
        background: #fff;
        border: 1px solid #e9e9e9;
        border-radius: 3px;
      }

      .content-wrap {
        padding: 20px;
      }

      .content-block {
        padding: 0 0 20px;
      }

      .header {
        width: 100%;
        margin-bottom: 20px;
      }

      h1, h2, h3 {
        font-family: "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
        color: #000;
        margin: 40px 0 0;
        line-height: 1.2;
        font-weight: 400;
      }

      h1 {
        font-size: 32px;
        font-weight: 500;
      }

      h2 {
        font-size: 24px;
      }

      h3 {
        font-size: 18px;
      }

      hr {
        border: 1px solid #e9e9e9;
        margin: 20px 0;
        height: 1px;
        padding: 0;
      }

      p,
      ul,
      ol {
        margin-bottom: 10px;
        font-weight: normal;
      }

      p li,
      ul li,
      ol li {
        margin-left: 5px;
        list-style-position: inside;
      }

      a {
        color: #348eda;
        text-decoration: underline;
      }

      .last {
        margin-bottom: 0;
      }

      .first {
        margin-top: 0;
      }

      .padding {
        padding: 10px 0;
      }

      .aligncenter {
        text-align: center;
      }

      .alignright {
        text-align: right;
      }

      .alignleft {
        text-align: left;
      }

      .clear {
        clear: both;
      }

      .alert {
        font-size: 16px;
        color: #fff;
        font-weight: 500;
        padding: 20px;
        text-align: center;
        border-radius: 3px 3px 0 0;
      }

      .alert a {
        color: #fff;
        text-decoration: none;
        font-weight: 500;
        font-size: 16px;
      }

      .alert.alert-warning {
        background: #ff9f00;
      }

      .alert.alert-bad {
        background: #d0021b;
      }

      .alert.alert-good {
        background: #68b90f;
      }

      @media only screen and (max-width: 640px) {
        h1,
        h2,
        h3 {
          font-weight: 600 !important;
          margin: 20px 0 5px !important;
        }

        h1 {
          font-size: 22px !important;
        }

        h2 {
          font-size: 18px !important;
        }

        h3 {
          font-size: 16px !important;
        }

        .container {
          width: 100% !important;
        }

        .content,
        .content-wrapper {
          padding: 10px !important;
        }
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
                  <td class="alert alert-good">
                    <a href="{{ system.link_url }}/{{ repo.owner }}/{{ repo.name }}/{{ build.number }}">
                      Successful build #$buildNumber
                    </a>
                  </td>
              </tr>
              <tr>
                <td class="content-wrap">
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        Distributed Switch:
                      </td>
                      <td>
                        ${nightlystatus[0]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Cluster:
                      </td>
                      <td>
                        ${nightlystatus[1]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Enhanced Linked Mode:
                      </td>
                      <td>
                        ${nightlystatus[2]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Heterogenous:
                      </td>
                      <td>
                        ${nightlystatus[3]}
                      </td>
                    </tr>
		    <tr>
                      <td>
                        VSAN Simple:
                      </td>
                      <td>
                        ${nightlystatus[4]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        VSAN Complex:
                      </td>
                      <td>
                        ${nightlystatus[5]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        NSX:
                      </td>
                      <td>
                        ${nightlystatus[6]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        DRS:
                      </td>
                      <td>
                        ${nightlystatus[7]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple Datacenter:
                      </td>
                      <td>
                        ${nightlystatus[8]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple Cluster:
                      </td>
                      <td>
                        ${nightlystatus[9]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple VLAN:
                      </td>
                      <td>
                        ${nightlystatus[10]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Invalid ESXi Install:
                      </td>
                      <td>
                        ${nightlystatus[11]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Remove Container OOB:
                      </td>
                      <td>
                        ${nightlystatus[12]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        vMotion VCH Appliance:
                      </td>
                      <td>
                        ${nightlystatus[13]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        vMotion Containers:
                      </td>
                      <td>
                        ${nightlystatus[14]}
                      </td>
                    </tr>
                  </table>
                  <hr>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        <a href='https://storage.cloud.google.com/vic-ci-logs/functional_logs_$input.zip?authuser=1'>https://storage.cloud.google.com/vic-ci-logs/functional_logs_$input.zip?authuser=1</a>
                      </td>
                    </tr>
                  </table>
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
else
echo "Failure"
cat <<EOT >> nightly_mail.html
To: mwilliamson-staff-adl@vmware.com
To: rashok@vmware.com
Subject: VIC Nightly Run #$buildNumber
From: VIC Nightly
MIME-Version: 1.0
Content-Type: text/html
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />

    <style>
      * {
        margin: 0;
        padding: 0;
        font-family: "Helvetica Neue", "Helvetica", Helvetica, Arial, sans-serif;
        box-sizing: border-box;
        font-size: 14px;
      }

      body {
        -webkit-font-smoothing: antialiased;
        -webkit-text-size-adjust: none;
        width: 100% !important;
        height: 100%;
        line-height: 1.6;
        background-color: #f6f6f6;
      }

      table td {
        vertical-align: top;
      }

      .body-wrap {
        background-color: #f6f6f6;
        width: 100%;
      }

      .container {
        display: block !important;
        max-width: 600px !important;
        margin: 0 auto !important;
        /* makes it centered */
        clear: both !important;
      }

      .content {
        max-width: 600px;
        margin: 0 auto;
        display: block;
        padding: 20px;
      }

      .main {
        background: #fff;
        border: 1px solid #e9e9e9;
        border-radius: 3px;
      }

      .content-wrap {
        padding: 20px;
      }

      .content-block {
        padding: 0 0 20px;
      }

      .header {
        width: 100%;
        margin-bottom: 20px;
      }

      h1, h2, h3 {
        font-family: "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
        color: #000;
        margin: 40px 0 0;
        line-height: 1.2;
        font-weight: 400;
      }

      h1 {
        font-size: 32px;
        font-weight: 500;
      }

      h2 {
        font-size: 24px;
      }

      h3 {
        font-size: 18px;
      }

      hr {
        border: 1px solid #e9e9e9;
        margin: 20px 0;
        height: 1px;
        padding: 0;
      }

      p,
      ul,
      ol {
        margin-bottom: 10px;
        font-weight: normal;
      }

      p li,
      ul li,
      ol li {
        margin-left: 5px;
        list-style-position: inside;
      }

      a {
        color: #348eda;
        text-decoration: underline;
      }

      .last {
        margin-bottom: 0;
      }

      .first {
        margin-top: 0;
      }

      .padding {
        padding: 10px 0;
      }

      .aligncenter {
        text-align: center;
      }

      .alignright {
        text-align: right;
      }

      .alignleft {
        text-align: left;
      }

      .clear {
        clear: both;
      }

      .alert {
        font-size: 16px;
        color: #fff;
        font-weight: 500;
        padding: 20px;
        text-align: center;
        border-radius: 3px 3px 0 0;
      }

      .alert a {
        color: #fff;
        text-decoration: none;
        font-weight: 500;
        font-size: 16px;
      }

      .alert.alert-warning {
        background: #ff9f00;
      }

      .alert.alert-bad {
        background: #d0021b;
      }

      .alert.alert-good {
        background: #68b90f;
      }

      @media only screen and (max-width: 640px) {
        h1,
        h2,
        h3 {
          font-weight: 600 !important;
          margin: 20px 0 5px !important;
        }

        h1 {
          font-size: 22px !important;
        }

        h2 {
          font-size: 18px !important;
        }

        h3 {
          font-size: 16px !important;
        }

        .container {
          width: 100% !important;
        }

        .content,
        .content-wrapper {
          padding: 10px !important;
        }
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
                  <td class="alert alert-bad">
                    <a href="{{ system.link_url }}/{{ repo.owner }}/{{ repo.name }}/{{ build.number }}">
                      Failed build #$buildNumber
                    </a>
                  </td>
              </tr>
              <tr>
                <td class="content-wrap">
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        Distributed Switch:
                      </td>
                      <td>
                        ${nightlystatus[0]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Cluster:
                      </td>
                      <td>
                        ${nightlystatus[1]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Enhanced Linked Mode:
                      </td>
                      <td>
                        ${nightlystatus[2]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Heterogenous:
                      </td>
                      <td>
                        ${nightlystatus[3]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        VSAN Simple:
                      </td>
                      <td>
                        ${nightlystatus[4]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        VSAN Complex:
                      </td>
                      <td>
                        ${nightlystatus[5]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        NSX:
                      </td>
                      <td>
                        ${nightlystatus[6]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        DRS:
                      </td>
                      <td>
                        ${nightlystatus[7]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple Datacenter:
                      </td>
                      <td>
                        ${nightlystatus[8]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple Cluster:
                      </td>
                      <td>
                        ${nightlystatus[9]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Multiple VLAN:
                      </td>
                      <td>
                        ${nightlystatus[10]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Invalid ESXi Install:
                      </td>
                      <td>
                        ${nightlystatus[11]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Remove Container OOB:
                      </td>
                      <td>
                        ${nightlystatus[12]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        vMotion VCH Appliance:
                      </td>
                      <td>
                        ${nightlystatus[13]}
                      </td>
                    </tr>
                    <tr>
                      <td>
                        vMotion Containers:
                      </td>
                      <td>
                        ${nightlystatus[14]}
                      </td>
                    </tr>
                  </table>
                  <hr>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
                        <a href='https://storage.cloud.google.com/vic-ci-logs/functional_logs_$input.zip?authuser=1'>https://storage.cloud.google.com/vic-ci-logs/functional_logs_$input.zip?authuser=1</a>
                      </td>
                    </tr>
                  </table>
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
fi

# Emails an HTML report of the test run results using SendMail.
sendmail -t < nightly_mail.html
