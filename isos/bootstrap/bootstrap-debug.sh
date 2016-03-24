#!/bin/bash
# Build the appliance filesystem ontop of the base

# exit on failure
set -x

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
export INIT=/sbin/init

#cleanup bootfs, we need a different configuration for the bootstrap vms

# copy in our components
cp ${BIN}/tether-linux ${ROOTFS}/bin/tether &&
    ln -sf ../bin/tether ${ROOTFS}/sbin/init

cp ${BIN}/rpctool ${ROOTFS}/sbin/

# package up the result
rm -f $BIN/bootstrap-debug.iso
$DIR/../generate-iso.sh -p $PKGDIR -i $INIT $BIN/bootstrap-debug.iso
