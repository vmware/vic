#!/bin/bash
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
#

# check for the configs file
if [[ ! -f "configs" ]] ; then
    echo "Error! Configs file is missing. Please try downloading the VIC UI installer again"
    echo ""
    exit 1
fi

# load configs variable into env
while IFS='' read -r line; do
    eval $line
done < ./configs

# check for the VC IP
if [[ $VCENTER_IP == "" ]] ; then
    echo "Error! vCenter IP cannot be empty. Please provide a valid IP in the configs file"
    exit 1
fi

# check for the pllugin manifest file
if [[ ! -f ../plugin-manifest ]] ; then
    echo "Error! Plugin manifest was not found!"
    cleanup
    exit 1
fi

# load plugin manifest into env
while IFS='' read -r p_line; do
    eval "$p_line"
done < ../plugin-manifest

read -p "Enter your vCenter Administrator Username: " VCENTER_ADMIN_USERNAME
echo -n "Enter your vCenter Administrator Password: "
read -s VCENTER_ADMIN_PASSWORD
echo ""

OS=$(uname)
PLUGIN_BUNDLES=''
VCENTER_SDK_URL="https://${VCENTER_IP}/sdk/"
COMMONFLAGS="--target $VCENTER_SDK_URL --user $VCENTER_ADMIN_USERNAME --password $VCENTER_ADMIN_PASSWORD"
PLUGIN_FOLDERS=''

if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-linux"
fi

check_prerequisite () {
    # if PLUGIN_TYPE is not specified default to html5 plugin
    if [[ $PLUGIN_TYPE = 'flex' ]] ; then
        PLUGIN_TYPE=flex
        key=$key_flex
    else
        PLUGIN_TYPE=html5
        key=$key_h5c
    fi

    if [[ $(curl -v --head https://$VCENTER_IP -k 2>&1 | grep -i "could not resolve host") ]] ; then
        echo "Error! Could not resolve the hostname. Please make sure you set VCENTER_IP correctly in the configuration file"
        exit 1
    fi
}

parse_and_unregister_plugins () {
    local plugin_flags="--key $key"
    echo "-------------------------------------------------------------"
    echo "Unregistering vCenter Server Extension..."
    echo "-------------------------------------------------------------"
    $PLUGIN_MANAGER_BIN remove $COMMONFLAGS $plugin_flags

    if [[ $? > 0 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Could not unregister plugin from vCenter Server. Please see the message above."
        exit 1
    fi
}

check_prerequisite
parse_and_unregister_plugins

echo "--------------------------------------------------------------"
echo "VIC UI unregistration was successful"
echo ""
