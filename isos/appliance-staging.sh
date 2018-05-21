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
echo "Usage: $0 -c package-cache(tgz) -p base-package(tgz) -o output-package(tgz)" 1>&2
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
            # Optional. Offline cache of packages
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

# Install VCH packages
# load the repo to use from the package if not explicit in env
REPO=${REPO:-$(cat $PKGDIR/repo.cfg)}
REPODIR="$DIR/base/repos/${REPO}/"
PACKAGE_MANAGER=${PACKAGE_MANAGER:-$(cat $REPODIR/repo-spec.json | jq -r '.packagemanager')}
PACKAGE_MANAGER=${PACKAGE_MANAGER:-tdnf}
setup_pm $REPODIR $PKGDIR $PACKAGE_MANAGER $REPO

script=$(cat $REPODIR/repo-spec.json | jq -r '.packages_script_staging.appliance')
[ -n "$script" ] && package_cached -c $cache -u -p $PKGDIR $script --nogpgcheck -y

STAGING_PKGS=$(cat $REPODIR/repo-spec.json | jq -r '.packages.appliance')
package_cached -c $cache -u -p $PKGDIR install $STAGING_PKGS --nogpgcheck -y

# hack around toybox issues
# rm $(rootfs_dir $PKGDIR)/bin/{passwd,login,su}
# chroot $(rootfs_dir $PKGDIR) rpm -i --replacefiles /var/cache/tdnf/photon-2.0/packages/shadow-4.2.1-13.ph2.x86_64.rpm || true

# Give a permission to vicadmin to run iptables.
echo "vicadmin ALL=NOPASSWD: /sbin/iptables --list" >> $(rootfs_dir $PKGDIR)/etc/sudoers

# ensure we're not including a cache in the staging bundle
# but don't update the cache bundle we're using to install
package_cached  -p $PKGDIR clean all

# configure us for autologin of root
getty_dir=$(rootfs_dir $PKGDIR)/usr/lib/systemd/system/getty@tty1.service.d/
mkdir -p $getty_dir
cp ${DIR}/appliance/override.conf $getty_dir
# Use mingetty as our getty provider
ln -sf $(rootfs_dir $PKGDIR)/sbin/mingetty $(rootfs_dir $PKGDIR)/usr/bin/agetty 

# Disable SSH by default - this can be enabled via guest operations
rm $(rootfs_dir $PKGDIR)/usr/lib/systemd/system/sshd@.service
rm -f $(rootfs_dir $PKGDIR)/etc/systemd/system/multi-user.target.wants/sshd.service
# Allow root login via ssh
sed -i -e "s/\#*PermitRootLogin\s.*/PermitRootLogin yes/" $(rootfs_dir $PKGDIR)/etc/ssh/sshd_config

# Disable root login
# sed -i -e 's@:/bin/bash$@:/bin/false@' $(rootfs_dir $PKGDIR)/etc/passwd

# Allow chpasswd to change expired password when launched from vic-init
cp -f ${DIR}/appliance/chpasswd.pam $(rootfs_dir $PKGDIR)/etc/pam.d/chpasswd
# Allow chage to be used with expired password when launched from vic-init
cp -f ${DIR}/appliance/chage.pam $(rootfs_dir $PKGDIR)/etc/pam.d/chage
# # Set shadow defaults
# cp -f ${DIR}/appliance/login.defs $(rootfs_dir $PKGDIR)/etc/
# cp -f ${DIR}/appliance/boot.local $(rootfs_dir $PKGDIR)/etc/rc.d/init.d/boot.local
# package up the result
pack $PKGDIR $OUT
