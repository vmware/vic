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

# Build the base of a bootable ISO

# exit on failure and configure debug, include util functions
set -e && [ -n "$DEBUG" ] && set -x
DIR=$(dirname $(readlink -f "$0"))
. $DIR/base/utils.sh


function usage() {
echo "Usage: $0 -p package-name(tgz) [-c yum-cache]" 1>&2
exit 1
}

while getopts "c:p:r:k:" flag
do
    case $flag in

        p)
            # Required. Package name
            PACKAGE="$OPTARG"
            ;;

        r)
            # Optional. Name of repo set in base/repos to use
            REPO="$OPTARG"
            ;;

        k)
            # Optional. Allows provision of custom kernel rpm
            # assumes it contains suitable /boot/vmlinuz-* and /lib/modules/... files
            CUSTOM_KERNEL_RPM="${OPTARG}"
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
if [ ! -z "$*" -o -z "$PACKAGE" ]; then
    usage
fi

REPO=${REPO:-photon-1.0}

# prep the build system
ensure_apt_packages cpio rpm tar ca-certificates xz-utils

PKGDIR=$(mktemp -d)

# initialize the bundle
initialize_bundle $PKGDIR

# base filesystem setup
mkdir -p $(rootfs_dir $PKGDIR)/{etc/yum,etc/yum.repos.d}
ln -s /lib $(rootfs_dir $PKGDIR)/lib64

# TODO: look at moving these prep pieces into the repo as functions
# sourced and run at appropriate stages
# work arounds for incorrect filesystem-1.0-13.ph2 package
if [ "$REPO" == "photon-2.0" ]; then
    mkdir -p $(rootfs_dir $PKGDIR)/{run,var}
    ln -s /run $(rootfs_dir $PKGDIR)/var/run
fi

# work arounds for libgcc.x86_64 0:4.4.7-18.el6 needing /dev/null
# it looks like udev package will create this but has it's own issues
if [ "$REPO" == "centos-6.9" ]; then
    # cpio will return an error on open syscall if lib64 is a symlink
    rm -f $(rootfs_dir $PKGDIR)/lib64
    mkdir -p $(rootfs_dir $PKGDIR)/{dev,lib64}
    mknod $(rootfs_dir $PKGDIR)/dev/null c 1 3
    chmod 666 $(rootfs_dir $PKGDIR)/dev/null
fi

if [[ $DRONE_BUILD_NUMBER && $DRONE_BUILD_NUMBER > 0 ]]; then
    # THIS SHOULD BE MOVED TO .drone.yml AS IT OVERRIDES THE -r OPTION
    REPOS="ci"
fi

# select the repo directory and populate the basic yum config
REPODIR="$DIR/base/repos/${REPO}/"
cp -a $REPODIR/*.repo $(rootfs_dir $PKGDIR)/etc/yum.repos.d/
cp $DIR/base/yum.conf $(rootfs_dir $PKGDIR)/etc/yum/
# allow future stages to know which repo this is using
echo "$REPO" > $PKGDIR/repo.cfg

# determine the kernel package
KERNEL=$(cat $REPODIR/kernel.pkg | awk '/^[^#]/{print}')
# if the kernel isn't a file then install it with the other packages
if [ -n "$CUSTOM_KERNEL_RPM" ]; then
    echo "Using kernel package from environment: $CUSTOM_KERNEL_RPM"
    KERNEL=$CUSTOM_KERNEL_RPM
elif [ ! -f $KERNEL ]; then
    echo "Using kernel repo package: $KERNEL"
    KERNEL_PKG=$KERNEL
else
    echo "Using kernel file package: $KERNEL"
    KERNEL="$(pwd)/$KERNEL"
fi

# install the core packages - strip lines starting with #
CORE_PKGS=$(cat $REPODIR/base.pkgs | awk '/^[^#]/{print}')
yum_cached -c $cache -u -p $PKGDIR install $CORE_PKGS $KERNEL_PKG --nogpgcheck -y

# check for raw kernel override
if [ -z "$KERNEL_PKG" ]; then
    (
        cd $(rootfs_dir $PKGDIR)
        rpm2cpio $KERNEL | cpio -idm --extract-over-symlinks
    )
fi

# Issue 3858: find all kernel modules and unpack them and run depmod against that directory
find $(rootfs_dir $PKGDIR)/lib/modules -name "*.ko.xz" -exec xz -d {} \;
KERNEL_VERSION=$(basename $(rootfs_dir $PKGDIR)/lib/modules/*)
chroot $(rootfs_dir $PKGDIR) depmod $KERNEL_VERSION 

# strip the cache from the resulting image
yum_cached -c $cache -p $PKGDIR clean all

# move kernel into bootfs /boot directory so that syslinux could load it
mv $(rootfs_dir $PKGDIR)/boot/vmlinuz-* $(bootfs_dir $PKGDIR)/boot/vmlinuz64
# try copying over the other boot files - rhel kernel seems to need a side car configuration file
find $(rootfs_dir $PKGDIR)/boot -type f | xargs cp -t $(bootfs_dir $PKGDIR)/boot/

# https://www.freedesktop.org/wiki/Software/systemd/InitrdInterface/
touch $(rootfs_dir $PKGDIR)/etc/initrd-release

# package up the result
pack $PKGDIR $PACKAGE
