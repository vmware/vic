#!/bin/bash
# Build the appliance filesystem ontop of the base

# exit on failure
set -e
echo "BEGINNING OF BOOTSTRAP_STAGING_"
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
#   iproute2  # for ip 
#   openssh   # for ssh server
#
# Temporary packages list here
#   shadow    # chpasswd hack
#   tndf      # so we can deploy other packages into the appliance live - MUST BE REMOVED FOR SHIPPING
#   vim       # basic editing function
/usr/bin/yum --installroot=${ROOTFS} install \
                                  bash \
                                  iproute2 \
                                  libtirpc \
                                  tdnf \
                                  vim \
                              -y --nogpgcheck 

# package up the result
tar -C $PKGDIR -zcf $OUT rootfs bootfs && rm -fr PKGDIR
