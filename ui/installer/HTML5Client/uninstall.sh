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

if [[ ! -f "configs" ]] ; then
    echo "Error! Configs file is missing. Please try downloading the VIC UI installer again"
    echo ""
    exit 1
fi

CONFIGS_FILE="configs"
while IFS='' read -r line; do
    eval $line
done < $CONFIGS_FILE

if [[ $VCENTER_IP == "" ]] ; then
    echo "Error! vCenter IP cannot be empty. Please provide a valid IP in the configs file"
    exit 1
fi

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
    if [[ ! -d ../plugin-packages ]] ; then
        echo "Error! VIC UI plugin bundle was not found. Please try downloading the VIC UI installer again"
        exit 1
    fi

    if [[ $(curl -v --head https://$VCENTER_IP -k 2>&1 | grep -i "could not resolve host") ]] ; then
        echo "Error! Could not resolve the hostname. Please make sure you set VCENTER_IP correctly in the configuration file"
        exit 1
    fi
}

parse_and_unregister_plugins () {
    for d in ../plugin-packages/* ; do
        if [[ -d $d ]] ; then
            echo "Reading plugin-package.xml..."

            while IFS='' read -r p_line; do
                eval "local $p_line"
            done < $d/vc_extension_flags

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
        fi
    done
}

rename_package_folder () {
    mv $1 $2
    if [[ $? > 0 ]] ; then
        echo "Error! Could not rename folder"
        exit 1
    fi
}

check_prerequisite
parse_and_unregister_plugins

echo "--------------------------------------------------------------"
echo "VIC UI unregistration was successful"
echo ""
