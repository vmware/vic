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
ECHO This script will uninstall vSphere Integrated Containers plugin
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
SET vcenter_unreg_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

REM read plugin-manifest
FOR /F "tokens=1,2 delims==" %%A IN (..\plugin-manifest) DO (
    IF NOT %%A=="" (
        CALL SET %%A=%%B
    )
)

REM entry routine
GOTO retrieve_vc_thumbprint

:retrieve_vc_thumbprint
"%parent%..\..\vic-ui-windows.exe" info %vcenter_unreg_flags% --key com.vmware.vic.noop > scratch.tmp 2> NUL
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
    GOTO parse_and_unregister_plugins
)
IF /I [%accept_vc_thumbprint%] == [no] (
    SET /p vc_thumbprint="Enter SHA-1 thumbprint of target VC: "
    SETLOCAL DISABLEDELAYEDEXPANSION
    GOTO parse_and_unregister_plugins
)
ECHO Please answer either "yes" or "no"
GOTO validate_vc_thumbprint

:parse_and_unregister_plugins
SET uninstall_successful=1
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to unregister vCenter Extension %name:"=%-H5Client...
ECHO -------------------------------------------------------------
"%plugin_manager_bin%" remove %vcenter_unreg_flags% --key com.vmware.vic --thumbprint %vc_thumbprint%
IF %ERRORLEVEL% NEQ 0 (
    SET uninstall_successful=0
)
ECHO.
ECHO -------------------------------------------------------------
ECHO Preparing to unregister vCenter Extension %name:"=%-FlexClient...
ECHO -------------------------------------------------------------
"%plugin_manager_bin%" remove %vcenter_unreg_flags% --key com.vmware.vic.ui --thumbprint %vc_thumbprint%
IF %ERRORLEVEL% NEQ 0 (
    SET uninstall_successful=0
)
GOTO end

:end
DEL scratch*.tmp 2>NUL
IF %uninstall_successful% NEQ 1 (
    ECHO -------------------------------------------------------------
    ECHO Error! Could not register plugin with vCenter Server. Please see the message above
    ENDLOCAL
    EXIT /b 1
)
ECHO --------------------------------------------------------------
ECHO VIC Engine UI uninstaller exited successfully
ENDLOCAL
