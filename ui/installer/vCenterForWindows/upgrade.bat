@ECHO OFF
REM Copyright 2016 VMware, Inc. All Rights Reserved.
REM
REM Licensed under the Apache License, Version 2.0 (the "License");
REM you may not use this file except in compliance with the License.
REM You may obtain a copy of the License at
REM
REM    http://www.apache.org/licenses/LICENSE-2.0
REM
REM Unless required by applicable law or agreed to in writing, software
REM distributed under the License is distributed on an "AS IS" BASIS,
REM WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
REM See the License for the specific language governing permissions and
REM limitations under the License.

SETLOCAL ENABLEEXTENSIONS
SET me=%~n0
SET parent=%~dp0

FOR /F "tokens=*" %%A IN (configs) DO (
    IF NOT %%A=="" (
        %%A
    )
)

IF [%target_vcenter_ip%] == [] (
    ECHO Error! vCenter IP cannot be empty. Please provide a valid IP in the configs file
    GOTO:EOF
)

SET "psCommand=powershell -Command "$pword = read-host 'Enter your vCenter Administrator Password' -AsSecureString ; ^
    $BSTR=[System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($pword); ^
        [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)""
FOR /f "usebackq delims=" %%p in (`%psCommand%`) do set vcenter_password=%%p

SET plugin_manager_bin=%parent%..\..\vic-machine-windows.exe something 
SET utils_path=%parent%utils\
SET vcenter_username=administrator@vsphere.local
SET vcenter_reg_common_flags=--url https://%target_vcenter_ip%/sdk/ --username %vcenter_username% --password %vcenter_password% --showInSolutionManager

IF EXIST _scratch_flags.txt (
    DEL _scratch_flags.txt
)

IF /I %vic_ui_host_url% NEQ NOURL (
    IF /I %vic_ui_host_url:~0,5%==https (
        SET vcenter_reg_common_flags=%vcenter_reg_common_flags% --serverThumbprint %vic_ui_host_thumbprint%
    )

    IF %vic_ui_host_url:~-1,1% NEQ / (
        SET vic_ui_host_url=%vic_ui_host_url%/
    )    
)

cd ..\vsphere-client-serenity
FOR /D %%i IN (*) DO (
    IF /I $vic_ui_host_url%==NOURL (
        "%utils_path%xml.exe" sel -t -o "--key " -v "/pluginPackage/@id" -o " --name \"" -v "/pluginPackage/@name" -o "\" --version " -v "/pluginPackage/@version" -o " --summary \"" -v "/pluginPackage/@description" -o "\" --company \"" -v "/pluginPackage/@vendor" -o "\" --pluginurl NOURL" -n %%i\plugin-package.xml >> ..\vCenterForWindows\_scratch_flags.txt
    ) ELSE (
        "%utils_path%xml.exe" sel -t -o "--key " -v "/pluginPackage/@id" -o " --name \"" -v "/pluginPackage/@name" -o "\" --version " -v "/pluginPackage/@version" -o " --summary \"" -v "/pluginPackage/@description" -o "\" --company \"" -v "/pluginPackage/@vendor" -o "\" --pluginurl %vic_ui_host_url%" -v "/pluginPackage/@id" -o "-" -v "/pluginPackage/@version" -o ".zip" -n %%i\plugin-package.xml >> ..\vCenterForWindows\_scratch_flags.txt
    )
)

ECHO Unregistering old VIC UI Plugins...
FOR /F "tokens=*" %%A IN (..\vCenterForWindows\_scratch_flags.txt) DO (
    IF NOT %%A=="" (
        REM %plugin_manager_bin% --unregister %vcenter_reg_common_flags% %%A 
        java -jar %parent%register-plugin.jar --unregister %vcenter_reg_common_flags% %%A
    )
)
"%utils_path%winscp.com" /command "open -hostkey=* sftp://%sftp_username%:%sftp_password%@%target_vcenter_ip%" "cd %target_vc_packages_path%" "rm com.vmware.vicui.*" "exit"

IF %sftp_supported% EQU 1 (
    ECHO Copying plugins...
    "%utils_path%winscp.com" /command "open -hostkey=* sftp://%sftp_username%:%sftp_password%@%target_vcenter_ip%" "put -filemask=|*.zip ..\vsphere-client-serenity\* %target_vc_packages_path%" "exit"
) ELSE (
    ECHO SFTP not enabled. You have to manually copy the com.vmware.vicui.* folder in \ui\vsphere-client-serenity to %VMWARE_CFG_DIR%\vsphere-client\vc-packages\vsphere-client-serenity
)

IF %ERRORLEVEL% GTR 0 (
    ECHO Error: Failed uploading plugin files! Check the configs file for correct connection credentials
    GOTO:EOF
)

ECHO Registering new VIC UI Plugins...
FOR /F "tokens=*" %%A IN (..\vCenterForWindows\_scratch_flags.txt) DO (
    IF NOT %%A=="" (
        REM %plugin_manager_bin% %vcenter_reg_common_flags% %%A
        java -jar %parent%register-plugin.jar %vcenter_reg_common_flags% %%A
    )
)

cd ..\vCenterForWindows
DEL _scratch_flags.txt

IF %ERRORLEVEL%==9009 (
    ECHO Error: java.exe was not found. Did you install Java?
)

ECHO Done
