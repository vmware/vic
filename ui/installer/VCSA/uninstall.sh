#!/bin/bash
# Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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

read_vc_information () {
    while getopts ":i:u:p:" o; do
        case "${o}" in
            i)
                VCENTER_IP=$OPTARG
                ;;
            u)
                VCENTER_ADMIN_USERNAME=$OPTARG
                ;;
            p)
                VCENTER_ADMIN_PASSWORD=$OPTARG
                ;;
            *)
                echo Usage: $0 [-i vc_ip] [-u vc admin_username] [-p vc_admin_password] >&2
                exit 1
                ;;
        esac
    done
    shift $((OPTIND-1))

    echo "-------------------------------------------------------------"
    echo "This script will uninstall vSphere Integrated Containers plugin"
    echo "for vSphere Client (HTML) and vSphere Web Client (Flex)."
    echo ""
    echo "Please provide connection information to the vCenter Server."
    echo "-------------------------------------------------------------"

    if [ -z $VCENTER_IP ] ; then
        read -p "Enter IP to target vCenter Server: " VCENTER_IP
    fi

    if [ -z $VCENTER_ADMIN_USERNAME ] ; then
        read -p "Enter your vCenter Administrator Username: " VCENTER_ADMIN_USERNAME
    fi

    if [ -z $VCENTER_ADMIN_PASSWORD ] ; then
        echo -n "Enter your vCenter Administrator Password: "
        read -s VCENTER_ADMIN_PASSWORD
        echo ""
    fi
}

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

read_vc_information $*

# replace space delimiters with colon delimiters
VIC_UI_HOST_THUMBPRINT=$(echo $VIC_UI_HOST_THUMBPRINT | sed -e 's/[[:space:]]/\:/g')

# check for the pllugin manifest file
if [[ ! -f ../plugin-manifest ]] ; then
    echo "Error! Plugin manifest was not found!"
    exit 1
fi

# load plugin manifest into env
while IFS='' read -r p_line; do
    eval "$p_line"
done < ../plugin-manifest

OS=$(uname)
VCENTER_SDK_URL="https://${VCENTER_IP}/sdk/"
COMMONFLAGS="--target $VCENTER_SDK_URL --user $VCENTER_ADMIN_USERNAME --password $VCENTER_ADMIN_PASSWORD"

if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-linux"
fi

prompt_thumbprint_verification() {
    read -p "Are you sure you trust the authenticity of this host (yes/no)? " SHOULD_ACCEPT_VC_FINGERPRINT
    if [[ ! $(echo $SHOULD_ACCEPT_VC_FINGERPRINT | grep -woi "yes\|no") ]] ; then
        echo Please answer either \"yes\" or \"no\"
        prompt_thumbprint_verification
        return
    fi

    if [[ $SHOULD_ACCEPT_VC_FINGERPRINT = "no" ]] ; then
        read -p "Enter SHA-1 thumbprint of target VC: " VC_THUMBPRINT
    fi
}

check_prerequisite () {
    # check if the provided VCENTER_IP is a valid vCenter Server host
    local CURL_RESPONSE=$(curl -sLk https://$VCENTER_IP)
    if [[ ! $(echo $CURL_RESPONSE | grep -oi "vmware vsphere") ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! vCenter Server was not found at host $VCENTER_IP"
        exit 1
    fi

    #retrieve VC thumbprint
    VC_THUMBPRINT=$($PLUGIN_MANAGER_BIN info $COMMONFLAGS --key com.vmware.vic.noop 2>&1 | grep -o "(thumbprint.*)" | cut -c13-71)

    # verify the thumbprint of VC
    echo ""
    echo "SHA-1 key fingerprint of host '$VCENTER_IP' is '$VC_THUMBPRINT'"
    prompt_thumbprint_verification

    # replace space delimiters with colon delimiters
    VC_THUMBPRINT=$(echo $VC_THUMBPRINT | sed -e 's/[[:space:]]/\:/g')
}

unregister_plugin() {
    local plugin_name=$1
    local plugin_key=$2
    echo "-------------------------------------------------------------"
    echo "Preparing to unregister vCenter Extension $plugin_name..."
    echo "-------------------------------------------------------------"
    local plugin_check_results=$($PLUGIN_MANAGER_BIN info $COMMONFLAGS --key $plugin_key --thumbprint $VC_THUMBPRINT 2>&1)
    if [[ $(echo $plugin_check_results | grep -oi "fail\|error") ]] ; then
        echo $plugin_check_results
        echo "-------------------------------------------------------------"
        echo "Error! Failed to check the status of plugin. Please see the message above"
        exit 1
    fi

    if [[ $(echo $plugin_check_results | grep -oi "is not registered") ]] ; then
        echo "Warning! Plugin with key '$plugin_key' is not registered with VC!"
        echo "Uninstallation was skipped"
        echo ""
        return
    fi

    $PLUGIN_MANAGER_BIN remove --key $plugin_key $COMMONFLAGS --thumbprint $VC_THUMBPRINT

    if [[ $? > 0 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Could not unregister plugin with vCenter Server. Please see the message above"
        exit 1
    fi
    echo ""
}

parse_and_unregister_plugins () {
    echo ""
    unregister_plugin "$name-FlexClient" $key_flex
    unregister_plugin "$name-H5Client" $key_h5c
}

check_prerequisite
parse_and_unregister_plugins

echo "--------------------------------------------------------------"
echo "VIC Engine UI uninstaller exited successfully"
echo ""
