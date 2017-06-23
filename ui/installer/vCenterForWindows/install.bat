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
SETLOCAL DISABLEDELAYEDEXPANSION

SET me=%~n0
SET parent=%~dp0

FOR /F "tokens=*" %%A IN (configs) DO (
    IF NOT %%A=="" (
        %%A
    )
)

IF NOT EXIST configs (
    ECHO -------------------------------------------------------------
    ECHO Error! Configs file is missing. Please try downloading the VIC UI installer again
    EXIT /b 1
)

ECHO -------------------------------------------------------------
ECHO This script will install vSphere Integrated Containers plugin
ECHO for vSphere Client (HTML) and vSphere Web Client (Flex).
ECHO.
ECHO Please provide connection information to the vCenter Server.
ECHO -------------------------------------------------------------
SET /p target_vcenter_ip="Enter IP to target vCenter Server: "
SET /p vcenter_username="Enter your vCenter Administrator Username: "
SET "psCommand=powershell -Command "$pword = read-host 'Enter your vCenter Administrator Password' -AsSecureString ; ^
    $BSTR=[System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($pword); ^
        [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)""
FOR /f "usebackq delims=" %%p in (`%psCommand%`) do set vcenter_password=%%p

SET plugin_manager_bin=%parent%..\..\vic-ui-windows.exe
SET vcenter_reg_common_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

IF [%1] == [--force] (
    SET force_install=1
) ELSE (
    SET force_install=0
)

REM read plugin-manifest
FOR /F "tokens=1,2 delims==" %%A IN (..\plugin-manifest) DO (
    IF NOT %%A=="" (
        CALL SET %%A=%%B
    )
)

REM add a forward slash to vic_ui_host_url if its last character is not '/'
IF [%vic_ui_host_url:~-1%] NEQ [/] (
    SET vic_ui_host_url=%vic_ui_host_url%/
)

REM replace space delimiters with colon delimiters
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /F "tokens=*" %%D IN ('ECHO %vic_ui_host_thumbprint%^| powershell -Command "$input.replace(' ', ':')"') DO (
    SET vic_ui_host_thumbprint=%%D
)
SETLOCAL DISABLEDELAYEDEXPANSION

REM entry routine
GOTO retrieve_vc_thumbprint

:retrieve_vc_thumbprint
"%parent%..\..\vic-ui-windows.exe" info %vcenter_reg_common_flags% --key com.vmware.vic.noop > scratch.tmp 2> NUL
TYPE scratch.tmp | findstr -c:"Failed to verify certificate" > NUL
IF %ERRORLEVEL% EQU 0 (
    SETLOCAL ENABLEDELAYEDEXPANSION
    FOR /F "usebackq tokens=2 delims=(" %%B IN (scratch.tmp) DO SET vc_thumbprint=%%B
    SET vc_thumbprint=!vc_thumbprint:~11,-1!
    ECHO.
    ECHO SHA-1 key fingerprint of host '%target_vcenter_ip%' is '!vc_thumbprint!'
    GOTO validate_vc_thumbprint
) ELSE (
    TYPE scratch.tmp
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)

:validate_vc_thumbprint
SET /p accept_vc_thumbprint="Are you sure you trust the authenticity of this host [yes/no]? "
IF /I [%accept_vc_thumbprint%] == [yes] (
    SETLOCAL DISABLEDELAYEDEXPANSION
    IF %force_install% NEQ 1 (
        GOTO check_existing_plugins
    ) ELSE (
        GOTO parse_and_register_plugins
    )
)
IF /I [%accept_vc_thumbprint%] == [no] (
    SET /p vc_thumbprint="Enter SHA-1 thumbprint of target VC: "
    SETLOCAL DISABLEDELAYEDEXPANSION
    IF %force_install% NEQ 1 (
        GOTO check_existing_plugins
    ) ELSE (
        GOTO parse_and_register_plugins
    )
)
ECHO Please answer either "yes" or "no"
GOTO validate_vc_thumbprint

:check_existing_plugins
ECHO.
ECHO -------------------------------------------------------------
ECHO Checking existing plugins...
ECHO -------------------------------------------------------------
SET can_install_continue=1
REM check for h5c plugin
"%parent%..\..\vic-ui-windows.exe" info %vcenter_reg_common_flags% --key com.vmware.vic --thumbprint %vc_thumbprint% > scratch.tmp 2>&1
REM check for connection failure
TYPE scratch.tmp | findstr -c:"fail" > NUL
IF %ERRORLEVEL% EQU 0 (
    TYPE scratch.tmp
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)
REM check if plugin (h5c) is not registered
TYPE scratch.tmp | findstr -c:"is not registered" > NUL
IF %ERRORLEVEL% GTR 0 (
    SETLOCAL ENABLEDELAYEDEXPANSION
    TYPE scratch.tmp | findstr -r -c:"Version.*" > scratch2.tmp
    FOR /F "usebackq tokens=2 delims=INFO" %%C IN (scratch2.tmp) DO SET ver_string=%%C
    ECHO com.vmware.vic is already registered. Version: !ver_string:~11!
    REM force flag condition
    SET can_install_continue=0
    SETLOCAL DISABLEDELAYEDEXPANSION
)
REM check for flex plugin
"%parent%..\..\vic-ui-windows.exe" info %vcenter_reg_common_flags% --key com.vmware.vic.ui --thumbprint %vc_thumbprint% > scratch.tmp 2>&1
REM check for connection failure
TYPE scratch.tmp | findstr -c:"fail" > NUL
IF %ERRORLEVEL% EQU 0 (
    TYPE scratch.tmp
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)
REM check if plugin (flex) is not registered
TYPE scratch.tmp | findstr -c:"is not registered" > NUL
IF %ERRORLEVEL% GTR 0 (
    SETLOCAL ENABLEDELAYEDEXPANSION
    TYPE scratch.tmp | findstr -r -c:"Version.*" > scratch2.tmp
    FOR /F "usebackq tokens=2 delims=INFO" %%C IN (scratch2.tmp) DO SET ver_string=%%C
    ECHO com.vmware.vic.ui is already registered. Version: !ver_string:~11!
    REM force flag condition
    SET can_install_continue=0
    SETLOCAL DISABLEDELAYEDEXPANSION
)
REM if either plugin is installed kill the script
IF %can_install_continue% EQU 0 (
    ECHO -------------------------------------------------------------
    ECHO Error! At least one plugin is already registered with the target VC.
    ECHO Run upgrade.bat, or install.bat --force instead.
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)
ECHO No VIC Engine UI plugin was detected. Continuing to install the plugins.
GOTO parse_and_register_plugins

:parse_and_register_plugins
REM remove obsolete plugin key if it ever exists
"%plugin_manager_bin%" remove %vcenter_reg_common_flags% --key com.vmware.vicui.Vicui --thumbprint %vc_thumbprint% > NUL 2> NUL
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to register vCenter Extension %name:"=%-H5Client...
ECHO -------------------------------------------------------------
SET plugin_reg_flags=%vcenter_reg_common_flags% --name "%name:"=%-H5Client" --thumbprint %vc_thumbprint% --version %version:"=% --summary %summary% --company %company% --key %key_h5c:"=% --url %vic_ui_host_url%files/%key_h5c:"=%-v%version:"=%.zip --server-thumbprint %vic_ui_host_thumbprint%
IF %force_install% EQU 1 (
    SET plugin_reg_flags=%plugin_reg_flags% --force
)
"%plugin_manager_bin%" install %plugin_reg_flags%
IF %ERRORLEVEL% NEQ 0 (
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to register vCenter Extension %name:"=%-FlexClient...
ECHO -------------------------------------------------------------
SET plugin_reg_flags=%vcenter_reg_common_flags% --name "%name:"=%-FlexClient" --thumbprint %vc_thumbprint% --version %version:"=% --summary %summary% --company %company% --key %key_flex:"=% --url %vic_ui_host_url%files/%key_flex:"=%-v%version:"=%.zip --server-thumbprint %vic_ui_host_thumbprint%
IF %force_install% EQU 1 (
    SET plugin_reg_flags=%plugin_reg_flags% --force
)
"%plugin_manager_bin%" install %plugin_reg_flags%
IF %ERRORLEVEL% NEQ 0 (
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    EXIT /b 1
)
GOTO end

:end
DEL scratch*.tmp 2>NUL
ECHO --------------------------------------------------------------
ECHO VIC Engine UI installer exited successfully
