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

OS=$(uname)
PLUGIN_BUNDLES=''
VCENTER_ADMIN_USERNAME="administrator@vsphere.local"
VCENTER_SDK_URL="https://${VCENTER_IP}/sdk/"
COMMONFLAGS="--target $VCENTER_SDK_URL --user $VCENTER_ADMIN_USERNAME --password $VCENTER_ADMIN_PASSWORD"
WEBCLIENT_PLUGINS_FOLDER="/etc/vmware/vsphere-client/vc-packages/vsphere-client-serenity/"
OLD_PLUGIN_FOLDERS=''
FORCE_INSTALL=''

case $1 in
    "-f")
        COMMONFLAGS="$COMMONFLAGS --force"
        ;;
    "--force")
        COMMONFLAGS="$COMMONFLAGS --force"
        ;;
esac

if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-linux"
fi

if [[ $VIC_UI_HOST_URL != 'NOURL' ]] ; then
    if [[ ${VIC_UI_HOST_URL:0:5} == 'https' ]] ; then
        COMMONFLAGS="${COMMONFLAGS} --serverThumbprint ${VIC_UI_HOST_THUMBPRINT}"
    elif [[ ${VIC_UI_HOST_URL:0:5} == 'HTTPS' ]] ; then
        COMMONFLAGS="${COMMONFLAGS} --serverThumbprint ${VIC_UI_HOST_THUMBPRINT}"
    fi

    if [[ ${VIC_UI_HOST_URL: -1: 1} != "/" ]] ; then
        VIC_UI_HOST_URL="$VIC_UI_HOST_URL/"
    fi
fi

check_prerequisite () {
    if [[ ! -d ../vsphere-client-serenity ]] ; then
        echo "Error! VIC UI plugin bundle was not found. Please try downloading the VIC UI installer again"
        exit 1
    fi
}

parse_and_register_plugins () {
    for d in ../vsphere-client-serenity/* ; do
        if [[ -d $d ]] ; then
            echo "Reading plugin-package.xml..."

            while IFS='' read -r p_line; do
                eval "local $p_line"
            done < $d/vc_extension_flags

            if [[ ! -d "../vsphere-client-serenity/${key}-${version}" ]] ; then
                rename_package_folder $d "../vsphere-client-serenity/$key-$version"
            fi

            local plugin_url="$VIC_UI_HOST_URL"
            if [[ $plugin_url != 'NOURL' ]] ; then
                if [[ ! -f "../vsphere-client-serenity/${key}-${version}.zip" ]] ; then
                    echo "File ${key}-${version}.zip does not exist!"
                    exit 1
                fi
                local plugin_url="$plugin_url$key-$version.zip"
            fi
            
            local plugin_flags="--key $key --name $name --version $version --summary $summary --company $company --url $plugin_url"
            echo "----------------------------------------"
            echo "Registering vCenter Server Extension..."
            echo "----------------------------------------"

            $PLUGIN_MANAGER_BIN install $COMMONFLAGS $plugin_flags
            if [[ $OLD_PLUGIN_FOLDERS -eq "" ]] ; then
                OLD_PLUGIN_FOLDERS="$key-*"
            else
                OLD_PLUGIN_FOLDERS="$OLD_PLUGIN_FOLDERS $key-*"
            fi

            if [[ $? > 0 ]] ; then
                echo "Error! Could not register plugin with vCenter Server. Please see the message above"
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

upload_packages () {
    for d in ../vsphere-client-serenity/* ; do
        if [[ -d $d ]] ; then
            PLUGIN_BUNDLES+="$d "
        fi
    done

    echo "-------------------------------------------------------------"
    echo "Copying plugin contents to the vSphere Web Client..."
    echo "Please enter the root password for your machine running VCSA"
    echo "-------------------------------------------------------------"
    scp -qr $PLUGIN_BUNDLES root@$VCENTER_IP:/tmp/
    if [[ $? > 0 ]] ; then
        echo "Error! Could not upload the VIC plugins to the target VCSA"
        exit 1
    fi
}

update_ownership () {
    echo "--------------------------------------------------------------"
    echo "Updating ownership of the plugin folders..."
    echo "Please enter the root password for your machine running VCSA"
    echo "--------------------------------------------------------------"
    local PLUGIN_BUNDLES_WITHOUT_PREFIX=$(echo $PLUGIN_BUNDLES | sed 's/\.\.\/vsphere\-client\-serenity\///g')
    ssh -t root@$VCENTER_IP "mkdir -p $WEBCLIENT_PLUGINS_FOLDER; cd $WEBCLIENT_PLUGINS_FOLDER; rm -rf $OLD_PLUGIN_FOLDERS; cp -rf /tmp/$PLUGIN_BUNDLES_WITHOUT_PREFIX $WEBCLIENT_PLUGINS_FOLDER; chown -R vsphere-client:users /etc/vmware/vsphere-client/vc-packages"
    if [[ $? > 0 ]] ; then
        echo "Error! Failed to update the ownership of folders. Please manually set them to vsphere-client:users"
        exit 1
    fi
}

# Check if plugin is located properly 
check_prerequisite

# Read from each plugin bundle the plugin-package.xml file and register a vCenter Server Extension based off of it
# Also, rename the folders such that they follow the convention of $PLUGIN_KEY-$PLUGIN_VERSION
parse_and_register_plugins

# if VIC_UI_HOST_URL is NOURL
if [[ $VIC_UI_HOST_URL == "NOURL" ]] ; then
    # Upload the folders to VCSA's vSphere Web Client plugins cache folder
    upload_packages
    # Chown the uploaded folders from root to vsphere-client
    update_ownership
fi

echo "--------------------------------------------------------------"
if [[ $? > 0 ]] ; then
    echo "Error! There was a problem in the VIC UI registration process"
    exit 1
else
    echo "VIC UI registration was successful"
    exit 0
fi
