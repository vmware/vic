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
SET vcenter_unreg_flags=--url https://%target_vcenter_ip%/sdk/ --username %vcenter_username% --password %vcenter_password% --unregister

IF EXIST _scratch_flags.txt (
    DEL _scratch_flags.txt
)

cd ..\vsphere-client-serenity
FOR /D %%i IN (*) DO (
    "%utils_path%xml.exe" sel -t -o "--key " -v "/pluginPackage/@id" %%i\plugin-package.xml >> ..\vCenterForWindows\_scratch_flags.txt
)

ECHO Unregistering VIC UI Plugins...
FOR /F "tokens=*" %%A IN (..\vCenterForWindows\_scratch_flags.txt) DO (
    IF NOT %%A=="" (
        REM %plugin_manager_bin% %vcenter_unreg_flags% %%A
        java -jar "%parent%register-plugin.jar" %vcenter_unreg_flags% %%A
    )
    IF %sftp_supported% EQU 1 (
        "%utils_path%winscp.com" /command "open -hostkey=* sftp://%sftp_username%:%sftp_password%@%target_vcenter_ip%" "cd %target_vc_packages_path%" "rm com.vmware.vicui.*" "exit"
    ) ELSE (
        ECHO SFTP not enabled. You have to manually delete %VMWARE_CFG_DIR%\vsphere-client\vc-packages\vsphere-client-serenity\com.vmware.vicui.*
    )
)

cd ..\vCenterForWindows
DEL _scratch_flags.txt

IF %ERRORLEVEL%==9009 (
    ECHO Error: java.exe was not found. Did you install Java?
)

ECHO VIC UI was successfully uninstalled. Make sure to log out of vSphere Web Client if are logged in, and log back in.
