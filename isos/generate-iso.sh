#!/bin/bash
# Build the base of a bootable ISO

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
      set -x
fi

DIR=$(dirname $(readlink -f "$0"))

INIT=/bin/bash

function usage() {
     echo "Usage: $0 -p package-path [-i init-path] output-name" 1>&2
     exit 1
}

while getopts "p:i:" flag
do
  case $flag in

    i)
      # Optional.  Path to kernel rpm
      INIT="$OPTARG"
      ;;

    p)
      # Required. Package name
      package="$OPTARG"
      ;;

    *)
    usage
    ;;
  esac
done

shift $((OPTIND-1))

ISOOUT="$*"


# check there were no extra args and the required ones are set
if [ -z "$ISOOUT" -o -z "$package"  ]; then
    usage
fi

# prep the build system
# Make sure we only try this as root
if [ "$(id -u)" == "0" ]; then
  apt-get update && apt-get -y install cpio xorriso
else
  echo "Skipping apt-get - rerun as root if missing packages"
fi

# ISO to stdout
#${ISOOUT:="stdio:/dev/fd/1"}

ROOTFS=$package/rootfs
BOOTFS=$package/bootfs

: ${BOOTFS:?BOOTFS must be set to the path containing boot/isolinux} || exit 1
test -e $BOOTFS/boot/isolinux/isolinux.bin -a \
     -e $BOOTFS/boot/isolinux/isolinux.cfg || exit 2

# ensure the target init exists
test -e ${ROOTFS:?ROOTFS must be set}/${INIT} || exit 3
# set the init binary in isolinux.cfg
sed -i -e "s|^#\(\s*append rdinit\)=_INIT_BINARY_|\1=$INIT|" $BOOTFS/boot/isolinux/isolinux.cfg || exit 4

# create the initramfs archive - subshell to avoid changing directory
( cd $ROOTFS && find | cpio -o -H newc | gzip --fast > $BOOTFS/boot/core.gz ) || exit 5

# generate the ISO and write it to $ISOOUT
xorriso -publisher 'VMware Inc.' -dev "$ISOOUT" -blank as_needed -map "$BOOTFS" / -boot_image isolinux dir=/boot/isolinux || exit 6
