@echo off
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

@IF not defined ANT_HOME (
   @echo BUILD FAILED: You must set the env variable ANT_HOME to your Apache Ant folder
   goto end
)
@IF not defined FLEX_HOME (
   @echo BUILD FAILED: You must set the env variable FLEX_HOME to your Flex SDK folder,
   @echo for instance: set FLEX_HOME="C:\Program Files (x86)\Adobe\Adobe Flash Builder 4.6\sdks\4.6.0"
   goto end
)
@IF not defined VSPHERE_SDK_HOME (
   @echo BUILD FAILED: You must set the env variable VSPHERE_SDK_HOME to your vSphere Web Client SDK folder
   goto end
)
@IF not exist "%VSPHERE_SDK_HOME%\libs\vsphere-client-lib.swc" (
   @echo BUILD FAILED: VSPHERE_SDK_HOME is not set to a valid vSphere Web Client SDK folder
   @echo %VSPHERE_SDK_HOME%\libs\vsphere-client-lib.swc is missing
   goto end
)

@call "%ANT_HOME%\bin\ant" -f build-flex.xml

:end
