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

# TEMP: imagec wrapper
cp ${DIR}/appliance/imagec.sh $(rootfs_dir $PKGDIR)/sbin/imagec

# TEMP: tether based init - configured for debug
cp ${BIN}/vch-init $(rootfs_dir $PKGDIR)/sbin/vch-debug

# kick off our components at boot time
cp ${DIR}/appliance/launcher.sh $(rootfs_dir $PKGDIR)/bin/
cp ${DIR}/appliance/launcher.service $(rootfs_dir $PKGDIR)/etc/systemd/system/
ln -s /etc/systemd/system/launcher.service $(rootfs_dir $PKGDIR)/etc/systemd/system/multi-user.target.wants/launcher.service

cp ${BIN}/imagec $(rootfs_dir $PKGDIR)/sbin/imagec.bin
cp ${BIN}/{docker-engine-server,port-layer-server,rpctool,vicadmin} $(rootfs_dir $PKGDIR)/sbin/

# Select systemd for our init process
generate_iso $PKGDIR $BIN/appliance.iso /lib/systemd/systemd
