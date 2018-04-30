# Copyright 2018 VMware, Inc. All Rights Reserved.
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

*** Settings ***
Documentation     Suite 25-02 - Reconfigure
Resource          ../../resources/Util.robot
Resource          ../../resources/Group25-Host-Affinity-Util.robot
Test Setup        Set Test Environment Variables
Test Teardown     Cleanup
Default Tags


*** Test Cases ***
Configuring a VCH does not affect affinity
    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Verify Group Contains VMs    %{VCH-NAME}    1

    Configure VCH without modifying affinity

    Verify Group Contains VMs    %{VCH-NAME}    1

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4


Configuring a VCH without a VM group does not affect affinity
    [Teardown]    Cleanup VIC Appliance On Test Server

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Verify Group Not Found       %{VCH-NAME}

    Configure VCH without modifying affinity

    Verify Group Not Found       %{VCH-NAME}

    Create Three Containers

    Verify Group Not Found       %{VCH-NAME}


Enabling affinity affects existing container VMs
    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Verify Group Not Found       %{VCH-NAME}

    Create Three Containers

    Verify Group Not Found       %{VCH-NAME}

    Configure VCH to enable affinity

    Verify Group Contains VMs    %{VCH-NAME}    4


Enabling affinity affects subsequent container VMs
    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables

    Verify Group Not Found       %{VCH-NAME}

    Configure VCH to enable affinity

    Verify Group Contains VMs    %{VCH-NAME}    1

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4


Disabling affinity affects existing container VMs
    [Teardown]    Cleanup VIC Appliance On Test Server

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Verify Group Contains VMs    %{VCH-NAME}    1

    Create Three Containers

    Verify Group Contains VMs    %{VCH-NAME}    4

    Configure VCH to disable affinity

    Verify Group Not Found       %{VCH-NAME}


Disabling affinity affects subsequent container VMs
    [Teardown]    Cleanup VIC Appliance On Test Server

    Verify Group Not Found       %{VCH-NAME}

    Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--affinity-vm-group

    Verify Group Contains VMs    %{VCH-NAME}    1

    Configure VCH to disable affinity

    Verify Group Not Found       %{VCH-NAME}

    Create Three Containers

    Verify Group Not Found       %{VCH-NAME}


With operations user, enabling affinity affects existing container VMs
    ${status}=  Get State Of Github Issue  7722
    Run Keyword If  '${status}' == 'closed'  Fail  Test 25-02-Upgrade.robot needs to be updated now that Issue #7722 has been resolved
    Log  Issue \#7722 is blocking implementation  WARN
    [Teardown]    No Operation

    #Verify Group Not Found       %{VCH-NAME}
    #
    #Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--ops-user=${ops_user_name} --ops-password=${ops_user_password} --ops-grant-perms
    #
    #Verify Group Not Found       %{VCH-NAME}
    #
    #Create Three Containers
    #
    #Verify Group Not Found       %{VCH-NAME}
    #
    #Configure VCH to enable affinity
    #
    #Verify Group Contains VMs    %{VCH-NAME}    4


With operations user, enabling affinity affects subsequent container VMs
    ${status}=  Get State Of Github Issue  7722
    Run Keyword If  '${status}' == 'closed'  Fail  Test 25-02-Upgrade.robot needs to be updated now that Issue #7722 has been resolved
    Log  Issue \#7722 is blocking implementation  WARN
    [Teardown]    No Operation

    #Verify Group Not Found       %{VCH-NAME}
    #
    #Install VIC Appliance To Test Server With Current Environment Variables    additional-args=--ops-user=${ops_user_name} --ops-password=${ops_user_password} --ops-grant-perms
    #
    #Verify Group Not Found       %{VCH-NAME}
    #
    #Configure VCH to enable affinity
    #
    #Verify Group Contains VMs    %{VCH-NAME}    1
    #
    #Create Three Containers
    #
    #Verify Group Contains VMs    %{VCH-NAME}    4
