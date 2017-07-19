@ECHO OFF
REM Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
SETLOCAL

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
    ENDLOCAL
    EXIT /b 1
)

SET arg_name=%1
SET arg_value=%2

:read_vc_args
IF NOT "%1"=="" (
    IF "%1"=="-i" (
        SET target_vcenter_ip=%2
        SHIFT
    )
    IF "%1"=="-u" (
        SET vcenter_username=%2
        SHIFT
    )
    IF "%1"=="-p" (
        SET vcenter_password=%2
        SHIFT
    )
    SHIFT
    GOTO :read_vc_args
)

ECHO -------------------------------------------------------------
ECHO This script will upgrade vSphere Integrated Containers plugin
ECHO for vSphere Client (HTML) and vSphere Web Client (Flex).
ECHO.
ECHO Please provide connection information to the vCenter Server.
ECHO -------------------------------------------------------------
IF [%target_vcenter_ip%] == [] (
    SET /p target_vcenter_ip="Enter IP to target vCenter Server: "
)
IF [%vcenter_username%] == [] (
    SET /p vcenter_username="Enter your vCenter Administrator Username: "
)
IF [%vcenter_password%] == [] (
    GOTO :read_vc_password
) ELSE (
    GOTO :after_vc_info_read
)

:read_vc_password
SET "psCommand=powershell -Command "$pword = read-host 'Enter your vCenter Administrator Password' -AsSecureString ; ^ $BSTR=[System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($pword); ^ [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)""
FOR /f "usebackq delims=" %%p in (`%psCommand%`) do set vcenter_password=%%p

:after_vc_info_read
SET plugin_manager_bin=%parent%..\..\vic-ui-windows.exe
SET vcenter_reg_common_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

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
    ENDLOCAL
    EXIT /b 1
)

:validate_vc_thumbprint
SET /p accept_vc_thumbprint="Are you sure you trust the authenticity of this host [yes/no]? "
IF /I [%accept_vc_thumbprint%] == [yes] (
    SETLOCAL DISABLEDELAYEDEXPANSION
    GOTO check_existing_plugins
)
IF /I [%accept_vc_thumbprint%] == [no] (
    SET /p vc_thumbprint="Enter SHA-1 thumbprint of target VC: "
    SETLOCAL DISABLEDELAYEDEXPANSION
    GOTO check_existing_plugins
)
ECHO Please answer either "yes" or "no"
GOTO validate_vc_thumbprint

:check_existing_plugins
ECHO.
ECHO -------------------------------------------------------------
ECHO Checking existing plugins...
ECHO -------------------------------------------------------------
SET plugins_installed=0
REM check for h5c plugin
"%parent%..\..\vic-ui-windows.exe" info %vcenter_reg_common_flags% --key com.vmware.vic --thumbprint %vc_thumbprint% > scratch.tmp 2>&1
REM check for connection failure
TYPE scratch.tmp | findstr -c:"fail" > NUL
IF %ERRORLEVEL% EQU 0 (
    TYPE scratch.tmp
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    ENDLOCAL
    EXIT /b 1
)
REM check if plugin (h5c) is not registered
TYPE scratch.tmp | findstr -c:"is not registered" > NUL
IF %ERRORLEVEL% GTR 0 (
    SETLOCAL ENABLEDELAYEDEXPANSION
    TYPE scratch.tmp | findstr -r -c:"Version.*" > scratch2.tmp
    FOR /F "usebackq tokens=2 delims=INFO" %%C IN (scratch2.tmp) DO SET ver_string=%%C
    ECHO com.vmware.vic is already registered. Version: !ver_string:~11!
    SET /A "plugins_installed=%plugins_installed%+1"
    SET old_plugin_ver=!ver_string:~11!
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
    ENDLOCAL
    EXIT /b 1
)
REM check if plugin (flex) is not registered
TYPE scratch.tmp | findstr -c:"is not registered" > NUL
IF %ERRORLEVEL% GTR 0 (
    SETLOCAL ENABLEDELAYEDEXPANSION
    TYPE scratch.tmp | findstr -r -c:"Version.*" > scratch2.tmp
    FOR /F "usebackq tokens=2 delims=INFO" %%C IN (scratch2.tmp) DO SET ver_string=%%C
    ECHO com.vmware.vic.ui is already registered. Version: !ver_string:~11!
    SET /A "plugins_installed=%plugins_installed%+1"
    SET old_plugin_ver=!ver_string:~11!
    SETLOCAL DISABLEDELAYEDEXPANSION
)
REM no plugin is installed, prompt the user if the plugins should be installed fresh
IF %plugins_installed% EQU 0 (
    ECHO No VIC Engine UI plugin was found on the target VC
    GOTO confirm_fresh_install
) ELSE (
    ECHO The version you are about to install is '%version:"=%'.
    SETLOCAL ENABLEDELAYEDEXPANSION
    FOR /F "usebackq tokens=1,2,3,4 delims=." %%m IN ('%version:"=%') DO (
        SET new_plugin_build=%%p
        SET /A "new_plugin_ver=%%m*100+%%n*10+%%o"
    )
    FOR /F "usebackq tokens=1,2,3,4 delims=." %%m IN ('%old_plugin_ver%') DO (
        SET old_plugin_build=%%p
        SET /A "old_plugin_ver=%%m*100+%%n*10+%%o"
    )

    IF NOT "!new_plugin_build!"=="" IF NOT "!old_plugin_build!"=="" SET /A "version_comparison=!old_plugin_build!-!new_plugin_build!"
    IF [!version_comparison!]==[] (
        SET /A "version_comparison=!old_plugin_ver!-!new_plugin_ver!"
    )

    IF !version_comparison! GEQ 0 (
        ECHO.
        ECHO You are trying to install plugins of an older or same version. For changes to take effect,
        ECHO please restart the vSphere Web Client and vSphere Client services after upgrade is completed
        ECHO For instructions, please refer to https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/ts_ui_not_appearing.html
    )
    SETLOCAL DISABLEDELAYEDEXPANSION
    GOTO confirm_upgrade
)

:confirm_upgrade
ECHO.
SET /p accept_install="Are you sure you want to continue (yes/no)? "
IF /I [%accept_install%] == [yes] (
    GOTO parse_and_force_register_plugins
)
IF /I [%accept_install%] == [no] (
    ECHO -------------------------------------------------------------
    ECHO Error! Upgrade was cancelled by user
    DEL scratch*.tmp 2>NUL
    ENDLOCAL
    EXIT /b 1
)
ECHO Please answer either "yes" or "no"
GOTO confirm_upgrade

:confirm_fresh_install
SET /p accept_install="Do you want to install the plugins [yes/no]? "
IF /I [%accept_install%] == [yes] (
    GOTO parse_and_force_register_plugins
)
IF /I [%accept_install%] == [no] (
    ECHO -------------------------------------------------------------
    ECHO Error! Upgrade was cancelled by user
    DEL scratch*.tmp 2>NUL
    ENDLOCAL
    EXIT /b 1
)
ECHO Please answer either "yes" or "no"
GOTO confirm_fresh_install

:parse_and_force_register_plugins
REM remove obsolete plugin key if it ever exists
"%plugin_manager_bin%" remove %vcenter_reg_common_flags% --key com.vmware.vicui.Vicui --thumbprint %vc_thumbprint% > NUL 2> NUL
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to register vCenter Extension %name:"=%-H5Client...
ECHO -------------------------------------------------------------
SET plugin_reg_flags=%vcenter_reg_common_flags% --force --name "%name:"=%-H5Client" --thumbprint %vc_thumbprint% --version %version:"=% --summary "Plugin for %name:"=%-H5Client" --company %company% --key %key_h5c:"=% --url %vic_ui_host_url%files/%key_h5c:"=%-v%version:"=%.zip --server-thumbprint %vic_ui_host_thumbprint%
"%plugin_manager_bin%" install %plugin_reg_flags%
IF %ERRORLEVEL% NEQ 0 (
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    ENDLOCAL
    EXIT /b 1
)
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to register vCenter Extension %name:"=%-FlexClient...
ECHO -------------------------------------------------------------
SET plugin_reg_flags=%vcenter_reg_common_flags% --force --name "%name:"=%-FlexClient" --thumbprint %vc_thumbprint% --version %version:"=% --summary "Plugin for %name:"=%-FlexClient" --company %company% --key %key_flex:"=% --url %vic_ui_host_url%files/%key_flex:"=%-v%version:"=%.zip --server-thumbprint %vic_ui_host_thumbprint%
"%plugin_manager_bin%" install %plugin_reg_flags%
IF %ERRORLEVEL% NEQ 0 (
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    DEL scratch*.tmp 2>NUL
    ENDLOCAL
    EXIT /b 1
)
GOTO end

:end
DEL scratch*.tmp 2>NUL
ECHO --------------------------------------------------------------
ECHO VIC Engine UI upgrader exited successfully
ENDLOCAL
