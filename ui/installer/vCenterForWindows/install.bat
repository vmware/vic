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
SET vcenter_reg_common_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

IF [%1] == [--force] (
   SET vcenter_reg_common_flags=%vcenter_reg_common_flags% --force
)

FOR /F "tokens=1,2 delims==" %%A IN (..\plugin-manifest) DO (
    IF NOT %%A=="" (
        CALL SET %%A=%%B
    )
)

IF %target_vc_version% EQU 6.5 (
    SET key=%key_h5c%
) ELSE (
    SET key=%key_flex%
)

SET PLUGIN_FLAGS=--key %key:"=% --name %name% --version %version:"=% --summary %summary% --company %company%
SET PLUGIN_URL=%vic_ui_host_url%%key%-v%version%.zip
IF /I %vic_ui_host_url% NEQ NOURL (
    IF %vic_ui_host_url:~-1,1% NEQ / (
        SET vic_ui_host_url=%vic_ui_host_url%/
    )

    IF /I %vic_ui_host_url:~0,5%==https (
        SET vcenter_reg_common_flags=%vcenter_reg_common_flags% --server-thumbprint %vic_ui_host_thumbprint%
    )
    SET PLUGIN_FLAGS=%PLUGIN_FLAGS% --url %PLUGIN_URL:"=%
) ELSE (
    SET PLUGIN_FLAGS=%PLUGIN_FLAGS% --url NOURL
)

ECHO Registering VIC UI Plugins...

"%plugin_manager_bin%" install %PLUGIN_FLAGS% %vcenter_reg_common_flags%
IF %ERRORLEVEL% NEQ 0 (
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    GOTO:EOF
)
ECHO.
ECHO Installation of UI plugin succeeded
ECHO.

IF /I %vic_ui_host_url% EQU NOURL (
    ECHO =============================
    IF %target_vc_version% EQU 6.5 (
        ECHO If you are installing the plugin for the vSphere Web Client (Flex Client^^^)
        ECHO To finish installation, copy the com.vmware.vic* folder from \ui\vsphere-client-serenity to %VMWARE_CFG_DIR%\vsphere-client\vc-packages\vsphere-client-serenity. Create any missing folders in between if necessary.
        ECHO Once all done, log out of vSphere Web Client and then log back in.
        ECHO.
        ECHO WARNING: If you are installing the plugin for the vSphere Client (H5 Client^^^)
        ECHO The plugin won't install correctly without a web server for the vSphere Client at this moment.
        ECHO Please see issue #4279 (https://github.com/vmware/vic/issues/4279^^^)
    ) ELSE (
        IF %target_vc_version% EQU 6.0 (
            ECHO NEXT STEP for vCenter 6.0 users
            ECHO With the current version of VIC running on vCenter for Windows, the com.vmware.vic.* folder needs to be manually copied from \ui\vsphere-client-serenity to %VMWARE_CFG_DIR%\vsphere-client\vc-packages\vsphere-client-serenity. If you have not done so, please copy it now.
        ) ELSE (
            ECHO NEXT STEP for vCenter 5.5 users
            ECHO VIC UI may run on a vCenter 5.5 setup, but is NOT officially supported. Use it at your own risk. To proceed, copy the com.vmware.vic.* folder to %PROGRAMDATA%\VMware\vSphere Web Client\vc-packages\vsphere-client-serenity instead.
        )
        ECHO Once all done, log out of vSphere Web Client and then log back in.
    )
    ECHO =============================
)
