#!/bin/bash
# Build the appliance filesystem ontop of the base

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
      set -x
fi

DIR=$(dirname $(readlink -f "$0"))

function usage() {
     echo "Usage: $0 -p base-package(tgz) -o output-package(tgz)" 1>&2
     exit 1
}

while getopts "p:o:" flag
do
  case $flag in

    p)
      # Required. Package name
      package="$OPTARG"
      ;;

    o)
      # Required. Target for iso and source for components
      OUT="$OPTARG"
      ;;

    *)
    usage
    ;;
  esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "$package" -o -z "${OUT}" ]; then
    usage
fi

PKGDIR=$(mktemp -d)

# prep the build system
# Make sure we only try this as root
if [ "$(id -u)" == "0" ]; then
  apt-get update && apt-get -y install yum
else
  echo "Skipping apt-get - rerun as root if missing packages"
fi

# unpackage base package
mkdir -p ${PKGDIR} && tar -C ${PKGDIR} -xf $package
ROOTFS=${PKGDIR}/rootfs
BOOTFS=${PKGDIR}/bootfs

# Install VCH base packages
#
# List stable packages here
#   e2fsprogs # for mkfs.ext4
#   procps-ng # for ps
#   iputils   # for ping
#   iproute2  # for ip
#   openssh   # for ssh server
#
# Temporary packages list here
#   systemd   # for convenience only at this time
#   tndf      # so we can deploy other packages into the appliance live - MUST BE REMOVED FOR SHIPPING
#   vim       # basic editing function
/usr/bin/yum --installroot=${ROOTFS} install \
                                  systemd \
                                  openssh \
                                  e2fsprogs \
                                  procps-ng \
                                  iputils \
                                  iproute2 \
                                  tdnf \
                                  vim \
                              -y --nogpgcheck

# configure us for autologin of root
#COPY override.conf $ROOTFS/etc/systemd/system/getty@.service.d/
# HACK until the issues with override.conf above are dealt with
pwhash=$(openssl passwd -1 -salt vic password)
sed -i -e "s/^root:[^:]*:/root:${pwhash}:/" ${ROOTFS}/etc/shadow

# Allow root login via ssh
sed -i -e "s/\#*PermitRootLogin\s.*/PermitRootLogin yes/" ${ROOTFS}/etc/ssh/sshd_config

# package up the result
tar -C $PKGDIR -zcf $OUT rootfs bootfs && rm -fr PKGDIR
