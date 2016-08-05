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

SET /p vcenter_username="Enter your vCenter Administrator Username: "
SET "psCommand=powershell -Command "$pword = read-host 'Enter your vCenter Administrator Password' -AsSecureString ; ^
    $BSTR=[System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($pword); ^
        [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)""
FOR /f "usebackq delims=" %%p in (`%psCommand%`) do set vcenter_password=%%p

SET plugin_manager_bin=%parent%..\..\vic-ui-windows.exe
SET utils_path=%parent%utils\
SET vcenter_reg_common_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

IF [%1] == [--force] (
   SET vcenter_reg_common_flags=%vcenter_reg_common_flags% --force 
)

IF /I %vic_ui_host_url% NEQ NOURL (
    IF /I %vic_ui_host_url:~0,5%==https (
        SET vcenter_reg_common_flags=%vcenter_reg_common_flags% --serverThumbprint %vic_ui_host_thumbprint%
    )

    IF %vic_ui_host_url:~-1,1% NEQ / (
        SET vic_ui_host_url=%vic_ui_host_url%/
    )    
) ELSE (
    IF %sftp_supported% EQU 1 (
        ECHO Copying plugins...
        "%utils_path%winscp.com" /command "open -hostkey=* sftp://%sftp_username%:%sftp_password%@%target_vcenter_ip%" "put -filemask=|*.zip ..\vsphere-client-serenity\* %target_vc_packages_path%" "exit"
    ) ELSE (
        ECHO SFTP not enabled. You have to manually copy the com.vmware.vicui.* folder in \ui\vsphere-client-serenity to %VMWARE_CFG_DIR%\vsphere-client\vc-packages\vsphere-client-serenity
    )
)

IF %ERRORLEVEL% GTR 0 (
    ECHO Error: Failed uploading plugin files! Check the configs file for correct connection credentials
    GOTO:EOF
)

IF EXIST _scratch_flags.txt (
    DEL _scratch_flags.txt
)

cd ..\vsphere-client-serenity
FOR /D %%i IN (*) DO (
    IF /I %vic_ui_host_url%==NOURL (
        "%utils_path%xml.exe" sel -t -o "--key " -v "/pluginPackage/@id" -o " --name \"" -v "/pluginPackage/@name" -o "\" --version " -v "/pluginPackage/@version" -o " --summary \"" -v "/pluginPackage/@description" -o "\" --company \"" -v "/pluginPackage/@vendor" -o "\" --url NOURL" -n %%i\plugin-package.xml >> ..\vCenterForWindows\_scratch_flags.txt
    ) ELSE (
        "%utils_path%xml.exe" sel -t -o "--key " -v "/pluginPackage/@id" -o " --name \"" -v "/pluginPackage/@name" -o "\" --version " -v "/pluginPackage/@version" -o " --summary \"" -v "/pluginPackage/@description" -o "\" --company \"" -v "/pluginPackage/@vendor" -o "\" --url %vic_ui_host_url%" -v "/pluginPackage/@id" -o "-" -v "/pluginPackage/@version" -o ".zip" -n %%i\plugin-package.xml >> ..\vCenterForWindows\_scratch_flags.txt
    )
)

ECHO Registering VIC UI Plugins...
FOR /F "tokens=*" %%A IN (..\vCenterForWindows\_scratch_flags.txt) DO (
    IF NOT %%A=="" (
        %plugin_manager_bin% install %vcenter_reg_common_flags% %%A
        IF %ERRORLEVEL% NEQ 0 (
            ECHO Error! Could not register plugin with vCenter Server. Please see the message above 
            GOTO:EOF
        )
    )
)

cd ..\vCenterForWindows
DEL _scratch_flags.txt

ECHO VIC UI was successfully installed. Make sure to log out of vSphere Web Client if you are logged in, and log back in.
