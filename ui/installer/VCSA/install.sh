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

echo "----------------------------------------"
echo "Checking if VCSA has Bash shell enabled..."
echo "Please enter the root password"
echo "----------------------------------------"

if [[ $(ssh root@$VCENTER_IP -t "scp" 2> /dev/null 2>&1 | grep -i "unknown command") ]] ; then
    BASH_ENABLED_ON_VCSA=0
else
    BASH_ENABLED_ON_VCSA=1
fi

OS=$(uname)
PLUGIN_BUNDLES=''
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
    if [[ $BASH_ENABLED_ON_VCSA -eq 0 ]] ; then
        if [[ -f extension_already_registered ]] ; then
            return 0
        fi
    fi

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

            if [[ $OLD_PLUGIN_FOLDERS -eq "" ]] ; then
                OLD_PLUGIN_FOLDERS="$key-*"
            else
                OLD_PLUGIN_FOLDERS="$OLD_PLUGIN_FOLDERS $key-*"
            fi

            $PLUGIN_MANAGER_BIN install $COMMONFLAGS $plugin_flags

            if [[ $? > 0 ]] ; then
                echo "Error! Could not register plugin with vCenter Server. Please see the message above"
                exit 1
            fi
        fi
    done

    touch extension_already_registered
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

    if [[ $BASH_ENABLED_ON_VCSA -eq 1 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Copying plugin contents to the vSphere Web Client..."
        echo "Please enter the root password"
        echo "-------------------------------------------------------------"
        scp -qr $PLUGIN_BUNDLES root@$VCENTER_IP:/tmp/
        if [[ $? > 0 ]] ; then
            echo "Error! Could not upload the VIC plugins to the target VCSA"
            exit 1
        fi
    else
        echo "-------------------------------------------------------------"
        echo "WARNING!! Bash shell is required to complete the installation"
        echo "automatically, however it is not enabled on this machine."
        echo "Installation cannot be completed automatically. Please follow the"
        echo "steps below, or refer to instructions in the UI installation guide."
        echo ""
        echo ""
        echo "1) Make the {vic_extracted_folder}/ui/vsphere-client-serenity folder"
        echo "   accessible by SSH on the system where you are running this script."
        echo ""
        echo "2) Connect to the target VCSA using the root account"
        echo ""
        echo "3) When greeted with the \"Command>\" prompt, type the following in order:"
        echo "   shell.set --enabled True"
        echo "   shell"
        echo ""
        echo "4) \"scp\" the folder mentioned in 1) to /etc/vmware/vsphere-client/vc-packages/"
        echo "   on the target VCSA's location. Create any missing folder missing if any"
        echo "   (e.g. vc-packages, vsphere-client-serenity)"
        echo ""
        echo "5) Change owner of the vc-packages folder by typing the following command:"
        echo "   chown -R vsphere-client:users /etc/vmware/vsphere-client/vc-packages"
        echo ""
        echo "6) When all done, log out of the shell and log into the Web Client"

        cleanup
        exit 1
    fi
}

update_ownership () {
    echo "--------------------------------------------------------------"
    echo "Updating ownership of the plugin folders..."
    echo "Please enter the root password"
    echo "--------------------------------------------------------------"
    local PLUGIN_BUNDLES_WITHOUT_PREFIX=$(echo $PLUGIN_BUNDLES | sed 's/\.\.\/vsphere\-client\-serenity\///g')
    ssh -t root@$VCENTER_IP "mkdir -p $WEBCLIENT_PLUGINS_FOLDER; cd $WEBCLIENT_PLUGINS_FOLDER; rm -rf $OLD_PLUGIN_FOLDERS; cp -rf /tmp/$PLUGIN_BUNDLES_WITHOUT_PREFIX $WEBCLIENT_PLUGINS_FOLDER; chown -R vsphere-client:users /etc/vmware/vsphere-client/vc-packages"
    if [[ $? > 0 ]] ; then
        echo "Error! Failed to update the ownership of folders. Please manually set them to vsphere-client:users"
        exit 1
    fi
}

cleanup () {
    unset VCENTER_ADMIN_USERNAME
    unset VCENTER_ADMIN_PASSWORD
    unset BASH_ENABLED_ON_VCSA
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

cleanup
rm extension_already_registered

if [[ $? > 0 ]] ; then
    echo "--------------------------------------------------------------"
    echo "Error! There was a problem removing a temporary file. Please manually remove extension_already_registered if it exists"
else
    echo "--------------------------------------------------------------"
    echo "VIC UI registration was successful"
fi
