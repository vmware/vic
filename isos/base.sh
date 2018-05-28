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
echo "Usage: $0 -p package-name(tgz) [-c package-cache]" 1>&2
exit 1
}

while getopts "c:p:r:" flag
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

        c)
            # Optional. Offline cache of package packages
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

REPO=${REPO:-photon-2.0}

if [[ $DRONE_BUILD_NUMBER && $DRONE_BUILD_NUMBER > 0 ]]; then
    # THIS SHOULD BE MOVED TO .drone.yml AS IT OVERRIDES THE -r OPTION
    REPO="ci"
fi

REPODIR="$DIR/base/repos/${REPO}/"
PACKAGE_MANAGER=$(cat $REPODIR/repo-spec.json | jq -r '.packagemanager')
PACKAGE_MANAGER=${PACKAGE_MANAGER:-tdnf}

# prep the build system
# ensure_apt_packages cpio rpm tar ca-certificates xz-utils

PKGDIR=$(mktemp -d)

# initialize the bundle
initialize_bundle $PKGDIR

ln -s /lib $(rootfs_dir $PKGDIR)/lib64

# preform a repo customization if needed
[ -f $REPODIR/base.sh ] && . $REPODIR/base.sh

setup_pm $REPODIR $PKGDIR $PACKAGE_MANAGER $REPO

# install the core packages
CORE_PKGS=$(cat $REPODIR/repo-spec.json | jq -r '.packages.base')
echo "Install core packages"
package_cached -c $cache -u -p $PKGDIR install $CORE_PKGS --nogpgcheck -y

# determine the kernel package
KERNEL=$(cat $REPODIR/repo-spec.json | jq -r '.kernel')
# if the kernel isn't a file then install it with the other packages
if [ -f "$(pwd)/$KERNEL" ]; then
    echo "Using kernel package from directory: $(pwd)/$KERNEL"
    KERNEL=$(pwd)/$KERNEL
    (
        cd $(rootfs_dir $PKGDIR)
        rpm2cpio $KERNEL | cpio -idm
    )
else
    echo "Using kernel file RPM package: $KERNEL"
    package_cached -c $cache -u -p $PKGDIR install $KERNEL --nogpgcheck -y
fi

# Issue 3858: find all kernel modules and unpack them and run depmod against that directory
find $(rootfs_dir $PKGDIR)/lib/modules -name "*.ko.xz" -exec xz -d {} \;
KERNEL_VERSION=$(basename $(rootfs_dir $PKGDIR)/lib/modules/*)
chroot $(rootfs_dir $PKGDIR) depmod $KERNEL_VERSION 

# strip the cache from the resulting image
package_cached -c $cache -p $PKGDIR clean all

# move kernel into bootfs /boot directory so that syslinux could load it
mv $(rootfs_dir $PKGDIR)/boot/vmlinuz-* $(bootfs_dir $PKGDIR)/boot/vmlinuz64
# try copying over the other boot files - rhel kernel seems to need a side car configuration file and System map
find $(rootfs_dir $PKGDIR)/boot -type f  -exec cp {} $(bootfs_dir $PKGDIR)/boot/ \;

# https://www.freedesktop.org/wiki/Software/systemd/InitrdInterface/
touch $(rootfs_dir $PKGDIR)/etc/initrd-release

# package up the result
pack $PKGDIR $PACKAGE
