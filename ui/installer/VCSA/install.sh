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

cleanup () {
    unset VCENTER_ADMIN_USERNAME
    unset VCENTER_ADMIN_PASSWORD
    unset BASH_ENABLED_ON_VCSA
    if [[ -f tmp.txt ]] ; then
        rm tmp.txt
    fi
}

# check for the configs file
if [[ ! -f "configs" ]] ; then
    echo "Error! Configs file is missing. Please try downloading the VIC UI installer again"
    echo ""
    cleanup
    exit 1
fi

# load configs variable into env
while IFS='' read -r line; do
    eval $line
done < ./configs

# check for the VC IP
if [[ $VCENTER_IP == "" ]] ; then
    echo "Error! vCenter IP cannot be empty. Please provide a valid IP in the configs file"
    cleanup
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

BASH_ENABLED_ON_VCSA=1

if [[ $VIC_UI_HOST_URL = "NOURL" ]] ; then
    echo "----------------------------------------"
    echo "Checking if VCSA has Bash shell enabled..."
    echo "Please enter the root password"
    echo "----------------------------------------"

    ssh root@$VCENTER_IP -t "scp" > tmp.txt 2>&1
    if [[ $(cat tmp.txt | grep -i "unknown command") ]] ; then
        BASH_ENABLED_ON_VCSA=0
    elif [[ $(cat tmp.txt | grep -i "Permission denied (publickey,password).") ]] ; then
        echo "----------------------------------------"
        echo "Error! Root password is incorrect"
        cleanup
        exit 1
    elif [[ $SIMULATE_NO_BASH_SUPPORT -eq 1 ]] ; then
        BASH_ENABLED_ON_VCSA=0
    fi
fi

OS=$(uname)
VCENTER_SDK_URL="https://${VCENTER_IP}/sdk/"
COMMONFLAGS="--target $VCENTER_SDK_URL --user $VCENTER_ADMIN_USERNAME --password $VCENTER_ADMIN_PASSWORD"
OLD_PLUGIN_FOLDERS=""
FORCE_INSTALL=""
WEBCLIENT_PLUGINS_FOLDER=""

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
    if [[ ${VIC_UI_HOST_URL: -1: 1} != "/" ]] ; then
        VIC_UI_HOST_URL="$VIC_UI_HOST_URL/"
    fi
fi

check_prerequisite () {
    # if PLUGIN_TYPE is not specified default to html5 plugin
    if [[ $PLUGIN_TYPE = 'flex' ]] ; then
        PLUGIN_TYPE=flex
        if [[ $IS_VCENTER_5_5 -eq 1 ]] ; then
            WEBCLIENT_PLUGINS_FOLDER="/var/lib/vmware/vsphere-client/vc-packages/vsphere-client-serenity/"
        else
            WEBCLIENT_PLUGINS_FOLDER="/etc/vmware/vsphere-client/vc-packages/vsphere-client-serenity/"
        fi
        key=$key_flex
    else
        PLUGIN_TYPE=html5
        WEBCLIENT_PLUGINS_FOLDER="/usr/lib/vmware-vsphere-ui/plugin-packages/"
        key=$key_h5c
    fi

    if [[ $VIC_UI_HOST_URL = "NOURL" ]] ; then
        if [[ $PLUGIN_TYPE = "html5" ]] ; then
            plugin_folder="plugin-packages"
        else
            plugin_folder="vsphere-client-serenity"
        fi

        if [[ ! -d ../$plugin_folder ]] ; then
            echo "Error! VIC UI plugin bundle was not found. Please try downloading the VIC UI installer again"
            cleanup
            exit 1
        fi
    fi

    if [[ $(curl -v --head https://$VCENTER_IP -k 2>&1 | grep -i "could not resolve host") ]] ; then
        echo "Error! Could not resolve the hostname. Please make sure you set VCENTER_IP correctly in the configuration file"
        cleanup
        exit 1
    fi
}

# purpose of this function is to remove an outdated version of vic ui in case it's installed
remove_old_key_installation () {
    if [[ ! $(ls -la ../vsphere-client-serenity/ | grep -i "com.vmware.vicui") ]] ; then
        $PLUGIN_MANAGER_BIN remove $COMMONFLAGS --key com.vmware.vicui.Vicui > /dev/null
    fi
}

parse_and_register_plugins () {
    local plugin_url=$VIC_UI_HOST_URL
    local plugin_flags="--key $key --name $name --version $version --summary $summary --company $company --url $plugin_url"

    if [[ $VIC_UI_HOST_URL != "NOURL" ]] ; then
        #TODO: factor these in integration test 18-1
        if [[ ! $(echo ${VIC_UI_HOST_URL:0:5} | grep -i "https") ]] ; then
            echo "-------------------------------------------------------------"
            echo "Error! VIC_UI_HOST_URL should always start with 'https' in the configs file"
            cleanup
            exit 1
        fi

        if [[ -z $VIC_UI_HOST_THUMBPRINT ]] ; then
            echo "-------------------------------------------------------------"
            echo "Error! Please provide VIC_UI_HOST_THUMBPRINT in the configs file"
            cleanup
            exit 1
        fi
            
        local plugin_flags="$plugin_flags$key-v$version.zip --server-thumbprint $VIC_UI_HOST_THUMBPRINT"
    fi

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
        echo "-------------------------------------------------------------"
        echo "Error! Could not register plugin with vCenter Server. Please see the message above"
        cleanup
        exit 1
    fi
}

upload_packages () {
    if [[ $BASH_ENABLED_ON_VCSA -eq 1 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Copying plugin contents to the vSphere Web Client..."
        echo "Please enter the root password of vSphere Client Appliance"
        echo "-------------------------------------------------------------"

        scp -qr ../$plugin_folder/$key-v$version root@$VCENTER_IP:/tmp/
        if [[ $? > 0 ]] ; then
            echo "Error! Could not upload the VIC plugins to the target VCSA"
            cleanup
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
        echo "1) Make the {vic_extracted_folder}/ui/plugin-packages folder"
        echo "   accessible by SSH on the system where you are running this script."
        echo ""
        echo "2) Connect to the target VCSA using the root account"
        echo ""
        echo "3) When greeted with the \"Command>\" prompt, type the following in order:"
        echo "   shell.set --enabled True"
        echo "   shell"
        echo ""
        echo "4) \"scp\" the folder mentioned in 1) to /usr/lib/vmware-vsphere-ui/plugin-packages/"
        echo "   on the target VCSA's location. Create any missing folder missing if any"
        echo ""
        echo "5) Change owner of the plugin-packages folder by typing the following command:"
        echo "   chown -R vsphere-ui:root /usr/lib/vmware-vsphere-ui/plugin-packages"
        echo ""
        echo "6) When all done, log out of the shell and log into the vSphere Client"

        cleanup
        exit 1
    fi
}

update_ownership () {
    echo "--------------------------------------------------------------"
    echo "Updating ownership of the plugin folders..."
    echo "Please enter the root password"
    echo "--------------------------------------------------------------"
    local SSH_COMMANDS_STR="mkdir -p $WEBCLIENT_PLUGINS_FOLDER; cd $WEBCLIENT_PLUGINS_FOLDER; rm -rf $OLD_PLUGIN_FOLDERS; cp -rf /tmp/$key-v$version $WEBCLIENT_PLUGINS_FOLDER"
    if [[ $PLUGIN_TYPE == "flex" ]] ; then
        if [[ ! $IS_VCENTER_5_5 -eq 1 ]] ; then
            local SSH_COMMANDS_STR="$SSH_COMMANDS_STR; chown -R vsphere-client:users /etc/vmware/vsphere-client/vc-packages"
        fi
    else
        local SSH_COMMANDS_STR="$SSH_COMMANDS_STR; chown -R vsphere-ui:root $WEBCLIENT_PLUGINS_FOLDER"
    fi

    ssh -t root@$VCENTER_IP $SSH_COMMANDS_STR
    if [[ $? > 0 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Failed to finalize the plugin installation."
        cleanup
        exit 1
    fi
}

verify_plugin_url() {
    local PLUGIN_BASENAME=$key-v$version.zip
    local CURL_RESPONSE=$(curl -v --head $VIC_UI_HOST_URL$PLUGIN_BASENAME -k 2>&1)
    local RESPONSE_STATUS=$(echo $CURL_RESPONSE | grep -E "HTTP\/.*\s4\d{2}\s.*")

    if [[ $(echo $CURL_RESPONSE | grep -i "could not resolve host") ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Could not resolve the host provided. Please make sure the URL is correct"
        cleanup
        exit 1

    elif [[ ! $(echo $RESPONSE_STATUS | wc -w) -eq 0 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Plugin was not found in the web server. Please make sure you have uploaded \"$PLUGIN_BASENAME\" to \"$VIC_UI_HOST_URL\", and retry installing the plugin"
        cleanup
        exit 1
    fi
}

check_prerequisite
remove_old_key_installation

if [[ $VIC_UI_HOST_URL = "NOURL" ]] ; then
    parse_and_register_plugins
    upload_packages
    update_ownership
else
    verify_plugin_url
    parse_and_register_plugins
fi

cleanup

echo "--------------------------------------------------------------"
echo "VIC UI registration was successful"
echo ""
