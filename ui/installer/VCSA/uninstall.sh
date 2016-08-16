#!/bin/bash -e
# Copyright 2016 VMware, Inc. All Rights Reserved.
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

read -p "Are you running vCenter 5.5? (y/n): " IS_VCENTER_5_5
if [[ $(echo $IS_VCENTER_5_5 | grep -i "y") ]] ; then
    IS_VCENTER_5_5=1
    WEBCLIENT_PLUGINS_FOLDER="/var/lib/vmware/vsphere-client/vc-packages/vsphere-client-serenity/"
else
    WEBCLIENT_PLUGINS_FOLDER="/etc/vmware/vsphere-client/vc-packages/vsphere-client-serenity/"
fi

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

parse_and_unregister_plugins () {
    for d in ../vsphere-client-serenity/* ; do
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
                echo "Error! Could not unregister plugin with vCenter Server. Please see the message above."
                exit 1
            fi

            if [[ $PLUGIN_FOLDERS -eq "" ]] ; then
                PLUGIN_FOLDERS="$key-*"
            else
                PLUGIN_FOLDERS="$PLUGIN_FOLDERS $key-*"
            fi
        fi
    done

    echo "-------------------------------------------------------------"
    echo "Deleting plugin contents..."
    echo "Please enter the root password for your machine running VCSA"
    echo "-------------------------------------------------------------"
    ssh -t root@$VCENTER_IP "cd $WEBCLIENT_PLUGINS_FOLDER; rm -rf $PLUGIN_FOLDERS"
}

rename_package_folder () {
    mv $1 $2
    if [[ $? > 0 ]] ; then
        echo "Error! Could not rename folder"
        exit 1
    fi
}

# Read from each plugin bundle the plugin-package.xml file and register a vCenter Server Extension based off of it
# Also, rename the folders such that they follow the convention of $PLUGIN_KEY-$PLUGIN_VERSION
parse_and_unregister_plugins

if [[ $? > 0 ]] ; then
    echo "--------------------------------------------------------------"
    echo "There was a problem in the VIC UI unregistration process"
    exit 1
else
    echo "--------------------------------------------------------------"
    echo "VIC UI unregistration was successful"
fi
