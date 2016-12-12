*** Settings ***
Documentation  Test 18-2 - VIC UI Uninstallation
Resource  ../../resources/Util.robot
Resource  ./vicui-common.robot
Test Teardown  Cleanup Installer Environment
Suite Setup  Check Config And Install VCH
Suite Teardown  Uninstall VCH

*** Test Cases ***
Ensure Vicui Is Installed Before Testing
    Set Vcenter Ip
    Install Vicui Without Webserver  ${TEST_VC_USERNAME}  ${TEST_VC_PASSWORD}  ${TEST_VC_ROOT_PASSWORD}  ${TRUE}
    ${output}=  OperatingSystem.GetFile  install.log
    ${passed}=  Run Keyword And Return Status  Should Contain  ${output}  was successful
    Run Keyword Unless  ${passed}  Copy File  install.log  install-fail-ensure-vicui-is-installed-before-testing.log
    Remove File  install.log
    Should Be True  ${passed}

Attempt To Uninstall With Configs File Missing
    # Rename the configs file and run the uninstaller script to see if it fails in an expected way
    Move File  ${UI_INSTALLER_PATH}/configs  ${UI_INSTALLER_PATH}/configs_renamed
    ${rc}  ${output}=  Run And Return Rc And Output  ${UNINSTALLER_SCRIPT_PATH}
    Run Keyword And Continue On Failure  Should Contain  ${output}  Configs file is missing
    Move File  ${UI_INSTALLER_PATH}/configs_renamed  ${UI_INSTALLER_PATH}/configs

Attempt To Uninstall With Plugin Missing
    # Rename the folder containing the VIC UI binaries and run the uninstaller script to see if it fails in an expected way
    Set Vcenter Ip
    Move Directory  ${UI_INSTALLER_PATH}/../vsphere-client-serenity  ${UI_INSTALLER_PATH}/../vsphere-client-serenity-a
    Uninstall Fails  ${TEST_VC_USERNAME}  ${TEST_VC_PASSWORD}
    ${output}=  OperatingSystem.GetFile  uninstall.log
    ${passed}=  Run Keyword And Return Status  Should Contain  ${output}  VIC UI plugin bundle was not found
    Run Keyword Unless  ${passed}  Copy File  uninstall.log  uninstall-fail-attempt-to-uninstall-with-plugin-missing.log
    Move Directory  ${UI_INSTALLER_PATH}/../vsphere-client-serenity-a  ${UI_INSTALLER_PATH}/../vsphere-client-serenity
    Remove File  uninstall.log
    Should Be True  ${passed}

Attempt To Uninstall With vCenter IP Missing
    # Leave VCENTER_IP empty and run the uninstaller script to see if it fails in an expected way
    Remove File  ${UI_INSTALLER_PATH}/configs
    ${results}=  Replace String Using Regexp  ${configs}  VCENTER_IP=.*  VCENTER_IP=\"\"
    Create File  ${UI_INSTALLER_PATH}/configs  ${results}
    ${rc}  ${output}=  Run And Return Rc And Output  cd ${UI_INSTALLER_PATH} && ${UNINSTALLER_SCRIPT_PATH}
    Run Keyword And Continue On Failure  Should Contain  ${output}  Please provide a valid IP

Attempt To Uninstall With Wrong Vcenter Credentials
    # Try uninstalling the plugin with wrong vCenter credentials and see if it fails in an expected way
    Set Vcenter Ip
    Uninstall Fails  ${TEST_VC_USERNAME}_nope  ${TEST_VC_PASSWORD}_nope
    ${output}=  OperatingSystem.GetFile  uninstall.log
    ${passed}=  Run Keyword And Return Status  Should Contain  ${output}  Cannot complete login due to an incorrect user name or password
    Run Keyword Unless  ${passed}  Copy File  uninstall.log  uninstall-fail-attempt-to-uninstall-with-wrong-vcenter-credentials.log
    Remove File  uninstall.log
    Should Be True  ${passed}

Uninstall Successfully
    Set Vcenter Ip
    Uninstall Vicui  ${TEST_VC_USERNAME}  ${TEST_VC_PASSWORD}
    ${output}=  OperatingSystem.GetFile  uninstall.log
    ${passed}=  Run Keyword And Return Status  Should Match Regexp  ${output}  unregistration was successful
    Run Keyword Unless  ${passed}  Copy File  uninstall.log  uninstall-fail-uninstall-successfully.log
    Remove File  uninstall.log
    Should Be True  ${passed}

Attempt To Uninstall Plugin That Is Already Gone
    Set Vcenter Ip
    Uninstall Fails  ${TEST_VC_USERNAME}  ${TEST_VC_PASSWORD}
    ${output}=  OperatingSystem.GetFile  uninstall.log
    ${passed}=  Run Keyword And Return Status  Should Contain  ${output}  failed to find target
    Run Keyword Unless  ${passed}  Copy File  uninstall.log  uninstall-fail-attempt-to-uninstall-plugin-that-is-already-gone.log
    Remove File  uninstall.log
    Should Be True  ${passed}
