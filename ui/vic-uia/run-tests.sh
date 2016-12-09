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

set -o pipefail

echo "Removing VIC directory if present"
echo "Cleanup logs from previous run"
rm -rf 18-1-VIC-UI-Installer 18-2-VIC-UI-Uninstaller 18-3-VIC-UI-NGC-tests ui/installer/vsphere-client-serenity 2>/dev/null
rm -rf *.zip *.log vic_*.tar.gz tests/manual-test-cases/Group18-VIC-UI/*.log
for f in $(find ui/vic-uia/ -name "\$*") ; do
    rm $f
done

input=$(wget -O - https://vmware.bintray.com/vic-repo |tail -n5 |head -n1 |cut -d':' -f 2 |cut -d'.' -f 3| cut -d'>' -f 2)
buildNumber=${input:4}

echo "Downloading bintray file $input"
wget https://vmware.bintray.com/vic-repo/$input.tar.gz

mkdir -p bin/$buildNumber

echo "Extracting .tar.gz"
tar xvzf $input.tar.gz -C bin/$buildNumber --strip 1

echo "Deleting .tar.gz vic file"
rm $input.tar.gz

cp bin/$buildNumber/vic-ui-linux ui/
cp -rf bin/$buildNumber/ui/vsphere-client-serenity ui/installer/

drone exec --trusted -e buildNumber=$buildNumber -e test="cd tests/manual-test-cases/Group18-VIC-UI && robot -C ansi setup-testbed.robot && robot -C ansi -d ../../../18-1-VIC-UI-Installer 18-1-VIC-UI-Installer.robot" -E ui/vic-uia/nightly_ui_tests_secrets.yml --yaml ui/vic-uia/ui-tests.yml
if [ $? -eq 0 ] ; then
    echo "Passed"
    InstallerStatus="Passed"
else
    echo "Failed"
    InstallerStatus="FAILED!"
fi
cp tests/manual-test-cases/Group18-VIC-UI/*.log 18-1-VIC-UI-Installer/ >/dev/null && rm tests/manual-test-cases/Group18-VIC-UI/*.log

drone exec --trusted -e buildNumber=$buildNumber -e test="cd tests/manual-test-cases/Group18-VIC-UI && robot -C ansi setup-testbed.robot && robot -C ansi -d ../../../18-2-VIC-UI-Uninstaller 18-2-VIC-UI-Uninstaller.robot" -E ui/vic-uia/nightly_ui_tests_secrets.yml --yaml ui/vic-uia/ui-tests.yml
if [ $? -eq 0 ] ; then
    echo "Passed"
    UninstallerStatus="Passed"
else
    echo "Failed"
    UninstallerStatus="FAILED!"
fi
cp tests/manual-test-cases/Group18-VIC-UI/*.log 18-2-VIC-UI-Uninstaller/ >/dev/null && rm tests/manual-test-cases/Group18-VIC-UI/*.log

drone exec --trusted -e buildNumber=$buildNumber -e test="apt-get update && apt-get install -yq maven && cd tests/manual-test-cases/Group18-VIC-UI && robot -C ansi setup-testbed.robot && mvn install -f ../../../ui/vic-uia/pom.xml ; robot -C ansi -d ../../../18-3-VIC-UI-NGC-tests 18-3-VIC-UI-NGC-tests.robot" -E ui/vic-uia/nightly_ui_tests_secrets.yml --yaml ui/vic-uia/ui-tests.yml
if [ $? -eq 0 ] ; then
    echo "Passed"
    NGCTestStatus="Passed"
else
    echo "Failed"
    NGCTestStatus="FAILED!"
fi
cp tests/manual-test-cases/Group18-VIC-UI/*.log 18-3-VIC-UI-NGC-tests/ >/dev/null

rm -rf bin/$buildNumber

if [[ $InstallerStatus = "Passed" && $UninstallerStatus = "Passed" && $NGCTestStatus = "Passed" ]] ; then
    buildStatus=0
else
    buildStatus=1
fi

if [ $buildStatus -eq 0 ] ; then
    cat <<EOF > nightly_ui_mail.html
To: kjosh@vmware.com
Subject: VIC Nightly UI Run #$buildNumber
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
                        Installer:
                      </td>
                      <td>
                        $InstallerStatus
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Uninstaller:
                      </td>
                      <td>
                        $UninstallerStatus
                      </td>
                    </tr>
                    <tr>
                      <td>
                        NGC Test:
                      </td>
                      <td>
                        $NGCTestStatus
                      </td>
                    </tr>
                  </table>
                  <hr>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
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
EOF
else
    cat <<EOF > nightly_ui_mail.html
To: kjosh@vmware.com
Subject: VIC Nightly UI Run #$buildNumber
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
                        Installer:
                      </td>
                      <td>
                        $InstallerStatus
                      </td>
                    </tr>
                    <tr>
                      <td>
                        Uninstaller:
                      </td>
                      <td>
                        $UninstallerStatus
                      </td>
                    </tr>
                    <tr>
                      <td>
                        NGC Test:
                      </td>
                      <td>
                        $NGCTestStatus
                      </td>
                    </tr>
                  </table>
                  <hr>
                  <table width="100%" cellpadding="0" cellspacing="0">
                    <tr>
                      <td>
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
EOF
fi

sendmail -t < nightly_ui_mail.html
