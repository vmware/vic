*** Settings ***
Documentation  Common keywords used by VIC UI installation & uninstallation test suites
Resource  ../../resources/Util.robot
Library  VicUiInstallPexpectLibrary.py

*** Variables ***
${TEST_VC_VERSION}          6.0
${TEST_VC_ROOT_PASSWORD}    vmware
${TIMEOUT}                  10 minutes

${SELENIUM_SERVER_PORT}     4444
${SELENIUM_BROWSER}         *firefox
${DATACENTER_NAME}          Datacenter
${CLUSTER_NAME}             Cluster
${DATASTORE_TYPE}           NFS
${DATASTORE_NAME}           fake
${DATASTORE_IP}             1.1.1.1
${CONTAINER_VM_NAME}        sharp_feynman-d39db0a231f2f639a073814c2affc03e4737d9ad361649069eb424e6c4e09b52

*** Keywords ***
Load Nimbus Testbed Env
    Should Exist  testbed-information
    ${envs}=  OperatingSystem.Get File  testbed-information
    @{envs}=  Split To Lines  ${envs}
    :FOR  ${item}  IN  @{envs}
    \  @{kv}=  Split String  ${item}  =
    \  Set Environment Variable  @{kv}[0]  @{kv}[1]
    \  Set Suite Variable  \$@{kv}[0]  @{kv}[1]
    Set Suite Variable  ${TEST_VC_USERNAME}  %{TEST_USERNAME}
    Set Suite Variable  ${TEST_VC_PASSWORD}  %{TEST_PASSWORD}

Destroy Testbed
    Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  %{NIMBUS_USER}-ESX-UITEST-*
    Run Keyword And Ignore Error  Kill Nimbus Server  %{NIMBUS_USER}  %{NIMBUS_PASSWORD}  ${NIMBUS_USER}-VC-UITEST-*

Check Config And Install VCH
    Run Keyword  Set Absolute Script Paths
    Load Nimbus Testbed Env
    Install VIC Appliance To Test Server  ../../../bin/%{BUILD_NUMBER}/vic-machine-linux  ../../../bin/%{BUILD_NUMBER}/appliance.iso  ../../../bin/%{BUILD_NUMBER}/bootstrap.iso
    Set Environment Variable  VCH_VM_NAME  %{VCH-NAME}

Set Absolute Script Paths
    # TODO: Since Docker environment is always Linux, it would be impossible to directly test the Windows script in the Drone CI system. Rather, the test could be done manually on Windows
    ${UI_INSTALLERS_ROOT}=  Run  pwd
    ${UI_INSTALLERS_ROOT}=  Join Path  ${UI_INSTALLERS_ROOT}  ../../../ui/installer
    Run Keyword If  os.sep == '/'  Set Suite Variable  ${UI_INSTALLER_PATH}  ${UI_INSTALLERS_ROOT}/VCSA  ELSE  Set Suite Variable  ${UI_INSTALLER_PATH}  ${UI_INSTALLERS_ROOT}/vCenterForWindows
    Should Exist  ${UI_INSTALLER_PATH}
    ${configs_content}=  OperatingSystem.GetFile  ${UI_INSTALLER_PATH}/configs
    Set Suite Variable  ${configs}  ${configs_content}

    # set exact paths for installer and uninstaller scripts
    Set Script Filename  INSTALLER_SCRIPT_PATH  ./install
    Set Script Filename  UNINSTALLER_SCRIPT_PATH  ./uninstall

Set Script Filename
    [Arguments]    ${suite_varname}  ${script_name}
    ${SCRIPT_FILENAME}=  Run Keyword If  os.sep == '/'  Set Variable  ${script_name}.sh  ELSE  Set Variable  ${script_name}.bat
    ${SCRIPT_FILENAME}=  Join Path  ${UI_INSTALLER_PATH}  ${SCRIPT_FILENAME}
    Set Suite Variable  \$${suite_varname}  ${SCRIPT_FILENAME}

Set Vcenter Ip
    # Populate VCENTER_IP with ${TEST_VC_IP}
    Remove File  ${UI_INSTALLER_PATH}/configs
    ${results}=  Replace String Using Regexp  ${configs}  VCENTER_IP=.*  VCENTER_IP=\"${TEST_VC_IP}\"
    ${results2}=  Run Keyword If  ${TEST_VC_VERSION} == '5.5'  Replace String Using Regexp  ${results}  IS_VCENTER_5_5=.*  IS_VCENTER_5_5=1  ELSE  Set Variable  ${results}

    Create File  ${UI_INSTALLER_PATH}/configs  ${results2}
    ${check}=  OperatingSystem.Get File  ${UI_INSTALLER_PATH}/configs
    Should Contain  ${check}  ${TEST_VC_IP}

Unset Vcenter Ip
    # Revert the configs file back to what it was
    #Remove File  ${UI_INSTALLER_PATH}/configs
    ${results}=  Replace String Using Regexp  ${configs}  VCENTER_IP=.*  VCENTER_IP=\"\"
    ${results}=  Replace String Using Regexp  ${results}  IS_VCENTER_5_5=.*  IS_VCENTER_5_5=0
    #Generate Config  ${UI_INSTALLER_PATH}/configs  '${results}'
    Run  echo '${results}' > ${UI_INSTALLER_PATH}/configs
    Should Exist  ${UI_INSTALLER_PATH}/configs

Force Remove Vicui Plugin
    Uninstall Vicui  ${TEST_VC_USERNAME}  ${TEST_VC_PASSWORD}
    ${output}=  OperatingSystem.GetFile  uninstall.log
    Should Match Regexp  ${output}  (unregistration was successful|failed to find target plugin)
    Remove File  uninstall.log

Rename Folder
    [Arguments]  ${old}  ${new}
    Move Directory  ${old}  ${new}
    Should Exist  ${new}

Cleanup Installer Environment
    # Reverts the configs file and make sure the folder containing the UI binaries has its original name that might've been left modified due to a test failure
    Unset Vcenter Ip
    @{folders}=  OperatingSystem.List Directory  ${UI_INSTALLER_PATH}/..  vsphere-client-serenity*
    Run Keyword If  ('@{folders}[0]' != 'vsphere-client-serenity')  Rename Folder  ${UI_INSTALLER_PATH}/../@{folders}[0]  ${UI_INSTALLER_PATH}/../vsphere-client-serenity

Uninstall VCH
    Log To Console  Gathering logs from the test server...
    Gather Logs From Test Server
    Log To Console  Deleting the VCH appliance...
    ${rc}  ${output}=  Run Secret VIC Machine Delete Command  %{VCH-NAME}  ../../../bin/%{BUILD_NUMBER}/vic-machine-linux
    Check Delete Success  %{VCH-NAME}
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Completed successfully
    ${output}=  Run  rm -f %{VCH-NAME}-*.pem
    ${out}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run  govc host.portgroup.remove %{VCH-NAME}-bridge
