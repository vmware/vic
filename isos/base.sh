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
echo "Usage: $0 -p package-name(tgz) [-c yum-cache] -k kernel-rpm" 1>&2
exit 1
}

while getopts "c:p:k:" flag
do
    case $flag in

        k)
            # Required.  Path to kernel rpm
            KERNEL="$OPTARG"
            ;;

        p)
            # Required. Package name
            PACKAGE="$OPTARG"
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
if [ ! -z "$*" -o -z "$PACKAGE" -o -z "$KERNEL" ]; then
    usage
fi

# prep the build system
ensure_apt_packges cpio rpm tar ca-certificates

KERNEL_TREE=$(mktemp -d)
RPM=$(mktemp)
PKGDIR=$(mktemp -d)

# unpack kernel and modules
rpm2cpio $KERNEL | (cd $KERNEL_TREE && cpio -id) && rm -f $RPM

# initialize the bundle
initialize_bundle $PKGDIR $KERNEL_TREE/boot/vmlinuz-esx-4.2.0

# base filesystem setup
mkdir -p $(rootfs_dir $PKGDIR)/{etc/yum,etc/yum.repos.d} $KERNEL_TREE
ln -s /lib $(rootfs_dir $PKGDIR)/lib64
cp $DIR/base/*.repo $(rootfs_dir $PKGDIR)/etc/yum.repos.d/
cp $DIR/base/yum.conf $(rootfs_dir $PKGDIR)/etc/yum/

# install the core packages
yum_cached -c $cache -u -p $PKGDIR install filesystem coreutils bash --nogpgcheck -y
# strip the cache from the resulting image
yum_cached -c $cache -p $PKGDIR clean all

# add the kernel modules now the filesystem is installed
rm -rf $(rootfs_dir $PKGDIR)/{lib/modules/,lib/firmware/}
mv $KERNEL_TREE/lib/modules/ $(rootfs_dir $PKGDIR)/lib/
rm -fr $KERNEL_TREE

# package up the result
pack $PKGDIR $PACKAGE
