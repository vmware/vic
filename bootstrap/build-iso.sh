#!/bin/bash

#
# Create booststrap iso bootable by VIC.
#

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
      set -x
fi

# Make sure only root can run our script
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root (or fakeroot)." 1>&2
   exit 1
fi

function usage() {
     echo "Usage: $0 [-i isolinux.bin] [-l ldlinux.c32] -k kernel.rpm -c corepure64.gz -d "tinycore dependencies.gz" output.iso" 1>&2
     exit 1
}

isolinux=/usr/lib/isolinux/isolinux.bin
ldlinux=/usr/lib/syslinux/modules/bios/ldlinux.c32

while getopts "i:l:k:c:d:" flag
do
  case $flag in

    i)
     # Path to isolinux.bin.
      isolinux="$OPTARG"
      ;;

    l)
     # Optional.  isolinux.bin may require ldlinux32.
      ldlinux="$OPTARG"
      ;;

    k)
      # Required.  Path to kernel rpm
      photon_kernel_rpm="$OPTARG"
      ;;

    c)
      # tiny core root gz
      tcl_core="$OPTARG"
      ;;

    d)
      # tiny core dependency gzip files
      tcz_deps="$OPTARG"
      ;;

    *)
    usage
    ;;
  esac
done

shift $((OPTIND-1))

isofile=$*

if [ -z "$GOPATH" ]; then
     echo "$GOPATH must be set"
     exit 1
fi

if [ -z "$1" ]; then
     usage
fi

# the root of our source relative to GOPATH
vic=$GOPATH/src/github.com/vmware/vic

# Where we stick the binaries and build artifact
BINARY=${BINARY:-$vic/binary}

if [ ! -d "$BINARY" ]; then
     echo "Can't find ${BINARY} directory" 1>&2
     exit 1
fi

# Prepare kernel

# Extract the kernel rpm contents here
kernel_tree="$BINARY"/kernel
rm -rf "$kernel_tree"
mkdir -p "$kernel_tree"

# GNU cpio doesn't allow you to archive/extract from/to a specific dir.  Use subshell instead.
(cd "$kernel_tree" && rpm2cpio ../"$photon_kernel_rpm" | cpio -id)

# Root of the initrd filesystem
INITRDFS=${INITRDFS:-$BINARY/bootfs}
rm -rf "$INITRDFS"
mkdir -p "$INITRDFS"/boot

# Copy the kernel binary to the boot folder
mv "$kernel_tree"/boot/vmlinuz-esx-4.2.0 "$INITRDFS"/boot/vmlinuz64

# The root fs of the iso.  We build a root filesystem here.
ISOFS=${ISOFS:-${BINARY}/iso}
rm -rf "$ISOFS"
mkdir -p "$ISOFS"

(cd "$ISOFS" && gunzip -c ../"$tcl_core" | cpio -id)

# Install the TCZ dependencies
for dep in $tcz_deps; do
     unsquashfs -f -d "$ISOFS" "$dep"
done

# Create a default ssh key pair
ssh-keygen -t rsa -b 1024 -f "$ISOFS"/id_rsa -N ""

# Replace the modules with what came from the "photon kernel" rpm.
rm -rf "$ISOFS"/lib/modules/ "$ISOFS"/lib/firmware/
mkdir -p "$ISOFS"/lib/modules/
mv "$kernel_tree"/lib/modules/ "$ISOFS"/lib/

# Copy static busybox binary into ISOFS
cp /bin/busybox "$ISOFS"/bin/busybox.static

bootstrap="$vic"/bootstrap/targets/linux

install -D "$bootstrap"/base/default.script "$ISOFS"/etc/udhcpc/default.script

patch --unified -d "$ISOFS" -p 0 --input="$bootstrap"/60-persistent-storage.rules.diff

install -D "$bootstrap"/99-vmware-memhotplug.rules "$ISOFS"/etc/udev/rules.d/99-vmware-memhotplug.rules
install -D "$bootstrap"/99-vmware-cpuhotplug.rules "$ISOFS"/etc/udev/rules.d/99-vmware-cpuhotplug.rules

cp -rv "$bootstrap"/isolinux "$INITRDFS"/boot/isolinux
cp "$isolinux" "$INITRDFS"/boot/isolinux/isolinux.bin

if [ -f "$ldlinux" ]; then
     cp "$ldlinux" "$INITRDFS"/boot/isolinux/ldlinux.c32
fi

cp "$bootstrap"/vmfork.sh "$ISOFS"/

cp "$bootstrap"/init "$ISOFS"/init

# verify the binaries are statically linked.  ld will return 0 for dynamically linked files
! /lib64/ld-linux-x86-64.so.2 --verify "$BINARY"/tether-linux
! /lib64/ld-linux-x86-64.so.2 --verify "$BINARY"/rpctool

# Copy the binaries to their relevant locations on the ISOFS
cp -f "$BINARY"/tether-linux "$ISOFS"/bin/tether && \
     ln -sf ../bin/tether "$ISOFS"/sbin/init

cp -f "$BINARY"/rpctool "$ISOFS"/sbin/rpctool

sed -i -e "s!BUILD_ID!$BUILD_ID!" "$INITRDFS"/boot/isolinux/boot.msg

(cd "${ISOFS}" && find | cpio -o -H newc | gzip --fast > "$INITRDFS"/boot/core.gz)

# No idea why that extra "/" is needed, but it is.  ¯\_(ツ)_/¯
xorriso -publisher 'VMware Inc.' -outdev "${isofile}" -blank as_needed -map "$INITRDFS" / -boot_image isolinux dir=/boot/isolinux

# Print the md5s of all of the files in case we need to log lineage.
md5sum "$BINARY"/tether-linux
md5sum "$BINARY"/rpctool
md5sum "$isofile"

