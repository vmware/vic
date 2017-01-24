#!/bin/sh
# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if [ -z "$ANT_HOME" ] || [ ! -f "${ANT_HOME}"/bin/ant ]
then
   echo BUILD FAILED: You must set the environment variable ANT_HOME to your Apache Ant folder
   exit 1
fi

if [ -z "$VSPHERE_SDK_HOME" ] || [ ! -f "${VSPHERE_SDK_HOME}"/libs/vsphere-client-lib.jar ]
then
   echo BUILD FAILED: You must set the environment variable VSPHERE_SDK_HOME to your vSphere Client SDK folder
   exit 1
fi

if [ -z "$FLEX_HOME" ] || [ ! -f "$FLEX_HOME"/bin/mxmlc ]
 then
   echo Using the Adobe Flex SDK files bundled with the vSphere Client SDK
   export FLEX_HOME="${VSPHERE_SDK_HOME}"/resources/flex_sdk_4.6.0.23201_vmw
fi

"${ANT_HOME}"/bin/ant -f build-war.xml

exit 0
