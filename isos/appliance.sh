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
echo "Usage: $0 -p staged-package(tgz) -b binary-dir" 1>&2
exit 1
}

while getopts "p:b:" flag
do
    case $flag in

        p)
            # Required. Package name
            PACKAGE="$OPTARG"
            ;;

        b)
            # Required. Target for iso and source for components
            BIN="$OPTARG"
            ;;

        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "$PACKAGE" -o -z "${BIN}" ]; then
    usage
fi

PKGDIR=$(mktemp -d)

# unpackage base package
unpack $PACKAGE $PKGDIR

#################################################################
# Above: arg parsing and setup
# Below: the image authoring
#################################################################

## systemd configuration
# create systemd vic target
cp ${DIR}/appliance/vic.target $(rootfs_dir $PKGDIR)/etc/systemd/system/
cp ${DIR}/appliance/vic-init.service $(rootfs_dir $PKGDIR)/etc/systemd/system/
cp ${DIR}/appliance/nat.service $(rootfs_dir $PKGDIR)/etc/systemd/system/
cp ${DIR}/appliance/nat-setup $(rootfs_dir $PKGDIR)/etc/systemd/scripts

mkdir -p $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants
ln -s /etc/systemd/system/vic-init.service $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants/
ln -s /etc/systemd/system/nat.service $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants/
ln -s /etc/systemd/system/multi-user.target $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants/

# disable networkd given we manage the link state directly
rm -f $(rootfs_dir $PKGDIR)/etc/systemd/system/multi-user.target.wants/systemd-networkd.service
rm -f $(rootfs_dir $PKGDIR)/etc/systemd/system/sockets.target.wants/systemd-networkd.socket

# change the default systemd target to launch VIC
ln -sf /etc/systemd/system/vic.target $(rootfs_dir $PKGDIR)/etc/systemd/system/default.target
# update the multi-user target to launch VIC - this launches sshd as well
#ln -s /etc/systemd/system/vic.target $(rootfs_dir $PKGDIR)/etc/systemd/system/multi-user.target.wants/vic.target

# do not use the systemd dhcp client
rm -f $(rootfs_dir $PKGDIR)/etc/systemd/network/*
cp ${DIR}/base/no-dhcp.network $(rootfs_dir $PKGDIR)/etc/systemd/network/

# do not use the default iptables rules - nat-setup supplants this
rm -f $(rootfs_dir $PKGDIR)/etc/systemd/network/*

#
# Set up vicadmin user
#

chroot $(rootfs_dir $PKGDIR) groupadd -g 1000 vicadmin
chroot $(rootfs_dir $PKGDIR) useradd -u 1000 -g 1000 -m -d /home/vicadmin -s /bin/false vicadmin
cp -R ${DIR}/vicadmin/* $(rootfs_dir $PKGDIR)/home/vicadmin
chown -R 1000:1000 $(rootfs_dir $PKGDIR)/home/vicadmin

## main VIC components
# TEMP: imagec wrapper
cp ${DIR}/appliance/imagec.sh $(rootfs_dir $PKGDIR)/sbin/imagec
# tether based init
cp ${BIN}/vic-init $(rootfs_dir $PKGDIR)/sbin/vic-init

cp ${BIN}/imagec $(rootfs_dir $PKGDIR)/sbin/imagec.bin
cp ${BIN}/{docker-engine-server,port-layer-server,vicadmin} $(rootfs_dir $PKGDIR)/sbin/

echo "net.ipv4.ip_forward = 1" > $(rootfs_dir $PKGDIR)/usr/lib/sysctl.d/50-vic.conf

## Generate the ISO
# Select systemd for our init process
generate_iso $PKGDIR $BIN/appliance.iso /lib/systemd/systemd
