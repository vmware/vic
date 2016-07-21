#!/bin/bash
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

# Build the appliance filesystem ontop of the base

# exit on failure and configure debug, include util functions
set -e && [ -n "$DEBUG" ] && set -x
DIR=$(dirname $(readlink -f "$0"))
. $DIR/base/utils.sh


function usage() {
echo "Usage: $0 -c yum-cache(tgz) -p base-package(tgz) -o output-package(tgz)" 1>&2
exit 1
}

while getopts "c:p:o:" flag
do
    case $flag in

        p)
            # Required. Package name
            PACKAGE="$OPTARG"
            ;;

        o)
            # Required. Target for iso and source for components
            OUT="$OPTARG"
            ;;

        c)
            # Optional. Offline cache of yum packages
            cache="$OPTARG"
            ;;

        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "$PACKAGE" -o -z "${OUT}" ]; then
    usage
fi

PKGDIR=$(mktemp -d)

unpack $PACKAGE $PKGDIR

#################################################################
# Above: arg parsing and setup
# Below: the image authoring
#################################################################

# Install VCH base packages
#
# List stable packages here
#   e2fsprogs # for mkfs.ext4
#   procps-ng # for ps
#   iputils   # for ping
#   iproute2  # for ip
#   iptables  # for iptables
#   net-tools # for netstat
#   openssh   # for ssh server
#   sudo      # for sudo
#
# Temporary packages list here
#   systemd   # for convenience only at this time
#   tndf      # so we can deploy other packages into the appliance live - MUST BE REMOVED FOR SHIPPING
#   vim       # basic editing function
yum_cached -c $cache -u -p $PKGDIR install \
    haveged \
    systemd \
    openssh \
    iptables \
    e2fsprogs \
    procps-ng \
    iputils \
    iproute2 \
    iptables \
    net-tools \
    sudo \
    tdnf \
    vim \
	photon-release \
    -y --nogpgcheck

# https://www.freedesktop.org/wiki/Software/systemd/InitrdInterface/
touch $(rootfs_dir $PKGDIR)/etc/initrd-release

# ensure we're not including a cache in the staging bundle
# but don't update the cache bundle we're using to install
yum_cached -p $PKGDIR clean all

# configure us for autologin of root
#COPY override.conf $ROOTFS/etc/systemd/system/getty@.service.d/
# HACK until the issues with override.conf above are dealt with
pwhash=$(openssl passwd -1 -salt vic password)
sed -i -e "s/^root:[^:]*:/root:${pwhash}:/" $(rootfs_dir $PKGDIR)/etc/shadow

# 1218: Temporarily disable SSH for TP3
rm $(rootfs_dir $PKGDIR)/usr/lib/systemd/system/sshd@.service

# Allow root login via ssh
sed -i -e "s/\#*PermitRootLogin\s.*/PermitRootLogin yes/" $(rootfs_dir $PKGDIR)/etc/ssh/sshd_config

# package up the result
pack $PKGDIR $OUT
