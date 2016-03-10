#!/bin/bash

# This generates an ISO using the isolinux boot image in $BOOTFS and the
# and initramfs generated from the contents of $ROOTFS
: ${ROOTFS:?ROOTFS must be set} || exit 1

: ${BOOTFS:?BOOTFS must be set to the path containing boot/isolinux} || exit 2
test -e $BOOTFS/boot/isolinux/isolinux.bin -a \
     -e $BOOTFS/boot/isolinux/boot.cat -a \
     -e $BOOTFS/boot/isolinux/isolinux.cfg || exit 2

if [ "$INIT" != "" ]; then
    # ensure the target init exists
    test -e ${ROOTFS:?ROOTFS must be set}/${INIT} || exit 3
    # set the init binary in isolinux.cfg
    sed -i -e "s|^#\(\s*append rdinit\)=_INIT_BINARY_|\1=$INIT|" $BOOTFS/boot/isolinux/isolinux.cfg || exit 4
fi

# create the initramfs archive
cd $ROOTFS && find | cpio -o -H newc | gzip --fast > $BOOTFS/boot/core.gz || exit 5

# generate the ISO with a label of $ISOLABEL and write it to $ISOOUT
xorriso -publisher 'VMware Inc.' -as mkisofs -V ${ISOLABEL:-boot-iso} \
	    -l -J -R -no-emul-boot -boot-load-size 4 -boot-info-table \
	    -b boot/isolinux/isolinux.bin -c boot/isolinux/boot.cat \
	    -o ${ISOOUT} ${BOOTFS} || exit 6
