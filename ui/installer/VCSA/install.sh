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
    while getopts ":fi:u:p:" o; do
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
            f)
                FORCE_INSTALL=1
                ;;
            *)
                echo Usage: $0 [-i vc_ip] [-u vc admin_username] [-p vc_admin_password] >&2
                exit 1
                ;;
        esac
    done
    shift $((OPTIND-1))

    echo "-------------------------------------------------------------"
    echo "This script will install vSphere Integrated Containers plugin"
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

# check for the plugin manifest file
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

# set binary to call based on os
if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-linux"
fi

# add a forward slash to VIC_UI_HOST_URL in case it misses it
if [[ ${VIC_UI_HOST_URL: -1: 1} != "/" ]] ; then
    VIC_UI_HOST_URL="$VIC_UI_HOST_URL/"
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

    # retrieve VC thumbprint
    VC_THUMBPRINT=$($PLUGIN_MANAGER_BIN info $COMMONFLAGS --key com.vmware.vic.noop 2>&1 | grep -o "(thumbprint.*)" | cut -c13-71)

    # verify the thumbprint of VC
    echo ""
    echo "SHA-1 key fingerprint of host '$VCENTER_IP' is '$VC_THUMBPRINT'"
    prompt_thumbprint_verification

    # replace space delimiters with colon delimiters
    VC_THUMBPRINT=$(echo $VC_THUMBPRINT | sed -e 's/[[:space:]]/\:/g')
}

# purpose of this function is to remove an outdated version of vic ui in case it's installed
remove_old_key_installation () {
    $PLUGIN_MANAGER_BIN remove $COMMONFLAGS --force --key com.vmware.vicui.Vicui > /dev/null 2> /dev/null
}

register_plugin() {
    local plugin_name=$1
    local plugin_key=$2
    local plugin_url="${VIC_UI_HOST_URL}files/"
    local plugin_flags="--version $version --company $company --url $plugin_url$plugin_key-v$version.zip"
    if [[ $FORCE_INSTALL -eq 1 ]] ; then
        plugin_flags="$plugin_flags --force"
    fi

    echo "-------------------------------------------------------------"
    echo "Preparing to register vCenter Extension $1..."
    echo "-------------------------------------------------------------"

    $PLUGIN_MANAGER_BIN install --key $plugin_key \
                                $COMMONFLAGS $plugin_flags \
                                --thumbprint $VC_THUMBPRINT \
                                --server-thumbprint $VIC_UI_HOST_THUMBPRINT \
                                --name "$plugin_name" \
                                --summary "Plugin for $plugin_name"
    if [[ $? > 0 ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Could not register plugin with vCenter Server. Please see the message above"
        exit 1
    fi
    echo ""
}

parse_and_register_plugins () {
    echo ""
    if [[ $FORCE_INSTALL -eq 1 ]] ; then
        register_plugin "$name-FlexClient" $key_flex
        register_plugin "$name-H5Client" $key_h5c
        return
    fi

    echo "-------------------------------------------------------------"
    echo "Checking existing plugins..."
    echo "-------------------------------------------------------------"
    local check_h5c=$($PLUGIN_MANAGER_BIN info $COMMONFLAGS --key $key_h5c --thumbprint $VC_THUMBPRINT 2>&1)
    local check_flex=$($PLUGIN_MANAGER_BIN info $COMMONFLAGS --key $key_flex --thumbprint $VC_THUMBPRINT 2>&1)

    if [[ $(echo $check_h5c | grep -oi "fail\|error") ]] ; then
        echo $check_h5c
        echo "Error! Failed to check the status of plugin. Please see the message above"
        exit 1
    fi

    if [[ $(echo $check_flex | grep -oi "fail\|error") ]] ; then
        echo $check_flex
        echo "Error! Failed to check the status of plugin. Please see the message above"
        exit 1
    fi

    local pattern="\([[:digit:]]\.\)\{2\}[[:digit:]]\(\.[[:digit:]]\{1,\}\)\{0,1\}"
    local h5c_plugin_version=$(echo $check_h5c | grep -oi "$pattern")
    local flex_plugin_version=$(echo $check_flex | grep -oi "$pattern")
    local existing_version=""

    if [[ -z $(echo $h5c_plugin_version$flex_plugin_version) ]] ; then
        echo "No VIC Engine UI plugin was detected. Continuing to install the plugins."
        echo ""
        register_plugin "$name-FlexClient" $key_flex
        register_plugin "$name-H5Client" $key_h5c
    else
        # assuming user always keeps the both plugins at the same version
        if [[ $(echo $check_h5c | grep -oi "is registered") ]] ; then
            echo "Plugin with key '$key_h5c' is already registered with VC. (Version: $h5c_plugin_version)"
        fi

        if [[ $(echo $check_flex | grep -oi "is registered") ]] ; then
            echo "Plugin with key '$key_flex' is already registered with VC (Version: $flex_plugin_version)"
        fi

        echo "-------------------------------------------------------------"
        echo "Error! At least one plugin is already registered with the target VC."
        echo "Please run upgrade.sh instead."
        echo ""
        exit 1
    fi
}

verify_plugin_url() {
    local plugin_key=$1
    local PLUGIN_BASENAME=$plugin_key-v$version.zip

    if [[ $BYPASS_PLUGIN_VERIFICATION ]] ; then
        return
    fi

    if [[ ! $(echo ${VIC_UI_HOST_URL:0:5} | grep -i "https") ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! VIC_UI_HOST_URL should always start with 'https' in the configs file"
        exit 1
    fi

    if [[ -z $VIC_UI_HOST_THUMBPRINT ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Please provide VIC_UI_HOST_THUMBPRINT in the configs file"
        exit 1
    fi

    local CURL_RESPONSE=$(curl --head $VIC_UI_HOST_URL$PLUGIN_BASENAME -k 2>&1)

    if [[ $(echo $CURL_RESPONSE | grep -i "could not resolve\|fail") ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Could not resolve the host provided. Please make sure the URL is correct"
        exit 1
    fi

    local RESPONSE_STATUS=$(echo $CURL_RESPONSE | grep -i "HTTP" | grep "4[[:digit:]][[:digit:]]\|500")
    if [[ $(echo $RESPONSE_STATUS | grep -oi "404") ]] ; then
        echo "-------------------------------------------------------------"
        echo "Error! Plugin bundle was not found. Please make sure \"$PLUGIN_BASENAME\" is available at \"$VIC_UI_HOST_URL\", and retry installing the plugin"
        exit 1
    fi
}

check_prerequisite
remove_old_key_installation
verify_plugin_url $key_flex
verify_plugin_url $key_h5c
parse_and_register_plugins

echo "--------------------------------------------------------------"
echo "VIC Engine UI installer exited successfully"
echo ""
