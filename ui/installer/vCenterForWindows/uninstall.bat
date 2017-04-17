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
SET vcenter_unreg_flags=--target https://%target_vcenter_ip%/sdk/ --user %vcenter_username% --password %vcenter_password%

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

ECHO Unregistering VIC UI Plugins...

"%plugin_manager_bin%" remove %vcenter_unreg_flags% --key %key:"=%
IF %ERRORLEVEL% NEQ 0 (
    ECHO Error! Could not unregister plugin from vCenter Server. Please see the message above 
    GOTO:EOF
)

ECHO.
ECHO VIC UI was successfully uninstalled. Make sure to log out of vSphere Web Client if are logged in, and log back in.
