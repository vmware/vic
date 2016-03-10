#!/bin/bash
# Build the appliance filesystem ontop of the base

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
      set -x
fi

DIR=$(dirname $(readlink -f "$0"))

function usage() {
     echo "Usage: $0 -p staged-package(tgz) -b binary-dir" 1>&2
     exit 1
}

while getopts "p:b:" flag
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

# Select systemd for our init process
export INIT=/lib/systemd/systemd

# TEMP imagec wrapper
cp ${DIR}/appliance/imagec.sh ${ROOTFS}/sbin/imagec

# kick off our components at boot time
cp ${DIR}/appliance/launcher.sh ${ROOTFS}/bin/
cp ${DIR}/appliance/launcher.service ${ROOTFS}/etc/systemd/system/
ln -s /etc/systemd/system/launcher.service ${ROOTFS}/etc/systemd/system/multi-user.target.wants/launcher.service

cp ${BIN}/imagec ${ROOTFS}/sbin/imagec.bin
cp ${BIN}/{docker-engine-server,port-layer-server,rpctool,vicadmin} ${ROOTFS}/sbin/

# package up the result
rm -f $BIN/appliance.iso
$DIR/generate-iso.sh -p $PKGDIR -i $INIT $BIN/appliance.iso