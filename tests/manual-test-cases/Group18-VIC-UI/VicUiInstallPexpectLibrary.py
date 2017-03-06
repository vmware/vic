# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

import os.path
import pexpect
import time

class VicUiInstallPexpectLibrary(object):
    TIMEOUT_LIMIT = 180
    NGC_TESTS_TIMEOUT_LIMIT = 1800
    with open('testbed-information', 'r') as f:
        testbed_information = f.read().splitlines()

    IS_TESTING_VSPHERE65 = 'TEST_VSPHERE_VER=65' in testbed_information[0]
    INSTALLER_PATH = os.path.join(os.path.dirname(__file__), '../../..', 'ui', 'installer', 'HTML5Client') if IS_TESTING_VSPHERE65 else os.path.join(os.path.dirname(__file__), '../../..', 'ui', 'installer', 'VCSA')
    NGC_TESTS_PATH = os.path.join(os.path.dirname(__file__), '../../..', 'ui', 'vic-ui-h5c/uia/h5-plugin-tests/ui-automation/vic-uia') if IS_TESTING_VSPHERE65 else os.path.join(os.path.dirname(__file__), '../../..', 'ui', 'vic-uia/uia/vic-uia')

    def _prepare_and_spawn(self, operation, callback, force=False):
        try:
            executable = os.path.join(VicUiInstallPexpectLibrary.INSTALLER_PATH, operation + '.sh')
            if force:
                executable += ' --force'

            self._f = open(operation + '.log', 'wb')
            self._pchild = pexpect.spawn(executable, cwd = VicUiInstallPexpectLibrary.INSTALLER_PATH, timeout = VicUiInstallPexpectLibrary.TIMEOUT_LIMIT)
            self._pchild.logfile = self._f
            callback()
            self._f.close()

        except IOError as e:
            return 'Error: ' + e.value

    def _common_prompts(self, vcenter_user, vcenter_password, root_password):
        self._pchild.expect('Enter your vCenter Administrator Username: ')
        self._pchild.sendline(vcenter_user)
        self._pchild.expect('Enter your vCenter Administrator Password: ')
        self._pchild.sendline(vcenter_password)

    def install_vicui_without_webserver(self, vcenter_user, vcenter_password, root_password, force=False):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, root_password)
            match_index = self._pchild.expect(['root@.*', '.*continue connecting.*'])
            if match_index == 1:
                self._pchild.sendline('yes')
                self._pchild.expect('root@.*')

            self._pchild.sendline(root_password)
            self._pchild.expect('root@.*')
            self._pchild.sendline(root_password)
            self._pchild.expect('root@.*')
            self._pchild.sendline(root_password)
            #self._pchild.interact()
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('install', commands, force)

    def install_vicui_without_webserver_nor_bash(self, vcenter_user, vcenter_password, root_password):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, root_password)
            match_index = self._pchild.expect(['root@.*', '.*continue connecting.*'])
            if match_index == 1:
                self._pchild.sendline('yes')
                self._pchild.expect('root@.*')

            self._pchild.sendline(root_password)
            self._pchild.expect('.*When all done.*')
            #self._pchild.interact()
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('install', commands)

    def install_fails_for_wrong_vcenter_ip(self, vcenter_user, vcenter_password, root_password):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, root_password)
            self._pchild.expect('.*Error.*')
            #self._pchild.interact()
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('install', commands)

    def install_fails_at_extension_reg(self, vcenter_user, vcenter_password, root_password, is_nourl=True):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, root_password)
            if is_nourl == True:
                match_index = self._pchild.expect(['root@.*', '.*continue connecting.*'])
                if match_index == 1:
                    self._pchild.sendline('yes')
                    self._pchild.expect('root@.*')
                self._pchild.sendline(root_password)

            self._pchild.expect('.*Error.*')
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('install', commands)

    def uninstall_fails(self, vcenter_user, vcenter_password):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, None)
            self._pchild.expect('.*Error.*')
            #self._pchild.interact()
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('uninstall', commands)

    def uninstall_vicui(self, vcenter_user, vcenter_password):
        def commands():
            self._common_prompts(vcenter_user, vcenter_password, None)
            self._pchild.expect(['.*successful', 'Error! Could not unregister.*'])
            #self._pchild.interact()
            self._pchild.expect(pexpect.EOF)

        self._prepare_and_spawn('uninstall', commands)

    def run_ngc_tests(self, vcenter_user, vcenter_password):
        try:
            self._f = open('ngc_tests.log', 'wb')
            self._pchild = pexpect.spawn('mvn test -Denv.VC_ADMIN_USERNAME=' + vcenter_user + ' -Denv.VC_ADMIN_PASSWORD=' + vcenter_password, cwd = VicUiInstallPexpectLibrary.NGC_TESTS_PATH, timeout = VicUiInstallPexpectLibrary.NGC_TESTS_TIMEOUT_LIMIT)
            self._pchild.logfile = self._f
            self._pchild.expect(pexpect.EOF)
            self._f.close()

        except IOError as e:
            return 'Error: ' + e.value

    def run_hsuia_tests(self):
        try:
            self._f = open('ngc_tests.log', 'wb')
            self._pchild = pexpect.spawn('mvn clean compile exec:exec -e -Dhsuia.runlist=work/runlists/default.runlist', cwd = VicUiInstallPexpectLibrary.NGC_TESTS_PATH, timeout = VicUiInstallPexpectLibrary.NGC_TESTS_TIMEOUT_LIMIT)
            self._pchild.logfile = self._f
            self._pchild.expect(pexpect.EOF)
            self._f.close()

        except IOError as e:
            return 'Error: ' + e.value
