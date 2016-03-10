#!/bin/bash
# Build the base of a bootable ISO

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
      set -x
fi

DIR=$(dirname $(readlink -f "$0"))

REPO=https://dl.bintray.com/vmware/vic-repo/kernel/
KERNEL=linux-esx-4.2.0-10.x86_64.rpm

function usage() {
     echo "Usage: $0 -p package-name(tgz) [-k kernel-rpm] [-p photon-repo]" 1>&2
     exit 1
}

while getopts "p:r:k:" flag
do
  case $flag in

    k)
      # Optional.  Path to kernel rpm
      KERNEL="$OPTARG"
      ;;

    p)
      # Required. Package name
      package="$OPTARG"
      ;;

    r)
      # Optional. photon repo
      REPO="$OPTARG"
      ;;

    *)
    usage
    ;;
  esac
done

shift $((OPTIND-1))

# check there were no extra args and the required ones are set
if [ ! -z "$*" -o -z "$package"  ]; then
    usage
fi

# prep the build system
# Make sure we only try this as root
if [ "$(id -u)" == "0" ]; then
  apt-get update && apt-get -y install curl cpio rpm tar yum ca-certificates
else
  echo "Skipping apt-get - rerun as root if missing packages"
fi

KERNEL_TREE=$(mktemp -d)
RPM=$(mktemp)
PKGDIR=$(mktemp -d)

# base root filesystem & package management
ROOTFS=$PKGDIR/rootfs
BOOTFS=$PKGDIR/bootfs
mkdir -p $ROOTFS/{var/lib/rpm,etc/yum.repos.d} $BOOTFS/boot $KERNEL_TREE  
ln -s /lib $ROOTFS/lib64
rpm --root=$ROOTFS --initdb
cp $DIR/base/*.repo $ROOTFS/etc/yum.repos.d/
yum --installroot $ROOTFS install --nogpgcheck -y filesystem coreutils bash

# add kernel and modules
curl -L --insecure $REPO/$KERNEL -o $RPM && rpm2cpio $RPM | (cd $KERNEL_TREE && cpio -id) && rm -f $RPM
mv $KERNEL_TREE/boot/vmlinuz-esx-4.2.0 $BOOTFS/boot/vmlinuz64
rm -rf $ROOTFS/lib/modules/ $ROOTFS/lib/firmware/ && mv $KERNEL_TREE/lib/modules/ $ROOTFS/lib/
rm -fr $KERNEL_TREE 

# prep boot filesystem
cp -a $DIR/base/isolinux $BOOTFS/boot/isolinux

# package up the result
tar -C $PKGDIR -zcf $package rootfs bootfs && rm -fr PKGDIR
