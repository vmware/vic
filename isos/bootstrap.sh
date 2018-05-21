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

# File list to estimate the size of the target tempfs in bootstrap
tempfs_target_list=('/lib/modules/*' \
    '/bin/tether' \
    '/bin/unpack' \
    '/sbin/*tables*' \
    '/lib/libm.*'\
    '/lib/libm-*' \
    '/lib/libgcc_s*' \
    '/lib/libip*tc*' \
    '/lib/libxtables*' \
    '/lib/libdl*' \
    '/lib/libc.so*'\
    '/lib/libc-*' \
    '/lib64/ld-*' \
    '/usr/lib/iptables' \
    '/lib/libhavege.so.1' \
    '/usr/sbin/haveged')

# Build the bootstrap filesystem ontop of the base

# exit on failure
set -e && [ -n "$DEBUG" ] && set -x
DIR=$(dirname $(readlink -f "$0"))
. $DIR/base/utils.sh

function usage() {
echo "Usage: $0 -p staged-package(tgz) -b binary-dir -d <activates debug when set>" 1>&2
exit 1
}

while getopts "p:b:d:o:" flag
do
    case $flag in

        p)
            # Required. Package name
            package="$OPTARG"
            ;;

        b)
            # Required. Target for iso and source for components
            BIN="$OPTARG"
            ;;
        d)
            # Optional. directs script to make a debug iso instead of a production iso.
            debug="$OPTARG"
            ;;
        o)
            # Optional. Name of the generated ISO
            ISONAME="$OPTARG"
            ;;
        *)

            usage
            ;;
    esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "$package" -o -z "${BIN}" ]; then
    usage
fi

#################################################################
# Above: arg parsing and setup
# Below: the image authoring
#################################################################

PKGDIR=$(mktemp -d)

unpack $package $PKGDIR

# load the repo to use from the package if not explicit in env
REPO=${REPO:-$(cat $PKGDIR/repo.cfg)}
REPODIR="$DIR/base/repos/${REPO}/"
PACKAGE_MANAGER=${PACKAGE_MANAGER:-$(cat $REPODIR/repo-spec.json | jq -r '.packagemanager')}
PACKAGE_MANAGER=${PACKAGE_MANAGER:-tdnf}
setup_pm $REPODIR $PKGDIR $PACKAGE_MANAGER $REPO


#selecting the init script as our entry point.
if [ -v debug ]; then
    export ISONAME=${ISONAME:-bootstrap-debug.iso}
    cp ${DIR}/bootstrap/bootstrap.debug $(rootfs_dir $PKGDIR)/bin/bootstrap
    cp ${BIN}/rpctool $(rootfs_dir $PKGDIR)/sbin/
else
    export ISONAME=${ISONAME:-bootstrap.iso}
    cp ${DIR}/bootstrap/bootstrap $(rootfs_dir $PKGDIR)/bin/bootstrap
fi

# copy in our components
cp ${BIN}/tether-linux $(rootfs_dir $PKGDIR)/bin/tether
cp ${BIN}/unpack $(rootfs_dir $PKGDIR)/bin/unpack

if [ -d $(rootfs_dir $PKGDIR)/etc/systemd ]; then
    echo "Preparing systemd for bootstrap"

    # copy in systemd entry script
    cp ${DIR}/bootstrap/systemd-init $(rootfs_dir $PKGDIR)/bin/init

    # kick off our components at boot time
    mkdir -p $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants
    cp ${DIR}/bootstrap/tether.service $(rootfs_dir $PKGDIR)/etc/systemd/system/
    cp ${DIR}/appliance/vic.target $(rootfs_dir $PKGDIR)/etc/systemd/system/
    ln -s /etc/systemd/system/tether.service $(rootfs_dir $PKGDIR)/etc/systemd/system/vic.target.wants/
    ln -sf /etc/systemd/system/vic.target $(rootfs_dir $PKGDIR)/etc/systemd/system/default.target

    # disable networkd given we manage the link state directly
    rm -f $(rootfs_dir $PKGDIR)/etc/systemd/system/multi-user.target.wants/systemd-networkd.service
    rm -f $(rootfs_dir $PKGDIR)/etc/systemd/system/sockets.target.wants/systemd-networkd.socket

    # do not use the systemd dhcp client
    rm -f $(rootfs_dir $PKGDIR)/etc/systemd/network/*

    # some systemd distros (centos-7) do not use systemd-networkd
    [ -e $(rootfs_dir $PKGDIR)/etc/systemd/network/ ] && cp ${DIR}/base/no-dhcp.network $(rootfs_dir $PKGDIR)/etc/systemd/network/
else
    echo "Preparing systemV init for bootstrap"
    
    # copy in sysv-init entry script
    cp ${DIR}/bootstrap/sysv-init $(rootfs_dir $PKGDIR)/bin/init

    # kick off our components at boot time
    cp ${DIR}/bootstrap/tether $(rootfs_dir $PKGDIR)/etc/rc.d/init.d/
    chmod +x $(rootfs_dir $PKGDIR)/etc/rc.d/init.d/tether
    ln -sf /etc/rc.d/init.d/tether $(rootfs_dir $PKGDIR)/etc/rc.d/rc1.d/S90tether

    # set the default run level for sysvinit
    # networking disabled on rc1
    cp ${DIR}/bootstrap/inittab $(rootfs_dir $PKGDIR)/etc/inittab
fi
cp ${REPODIR}/init.sh $(rootfs_dir $PKGDIR)/bin/repoinit

# compute the size of the target tempfs,
# the list of directories/files in ${tempfs_target_list} should
# match the directories/files that are actually copied into tempfs
# by the script isos/bootstrap/bootstrap
target_list=$(rootfs_prepend $PKGDIR "${tempfs_target_list[@]}")
size=$(du -m --total ${target_list} | tail -1 | cut -f 1)
# 20% overhead should give a little more than 80M for stripped binaries
overhead=$(( size / 5 ))
size=$(( size + overhead ))
echo Total tempfs size: ${size}
# save the list of directories/files, for future usage
echo "${tempfs_target_list[@]}" > $(rootfs_dir $PKGDIR)/.tempfs_list
echo ${size} > $(rootfs_dir $PKGDIR)/.tempfs_size

INIT=$(cat $REPODIR/repo-spec.json | jq -r '.init')
generate_iso $PKGDIR $BIN/$ISONAME $INIT
