#!/bin/bash
# Build the bootstrap filesystem ontop of the base

# exit on failure
set -e

if [ -n "$DEBUG" ]; then
    set -x
fi

DIR=$(dirname $(readlink -f "$0"))
. $DIR/base/utils.sh

function usage() {
echo "Usage: $0 -p staged-package(tgz) -b binary-dir -d <activates debug when set>" 1>&2
exit 1
}

while getopts "p:b:d:" flag
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
        d)
            # Optional. directs script to make a debug iso instead of a production iso.
            debug="$OPTARG"
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

#################################################################
# Above: arg parsing and setup
# Below: the image authoring
#################################################################

PKGDIR=$(mktemp -d)

#selecting the init script as our entry point.
if [ -v debug ]; then
    export INIT=/bin/bash
    export ISONAME="bootstrap-debug.iso"
else
    export INIT=/sbin/init
    export ISONAME="bootstrap.iso"
fi

unpack $package $PKGDIR

# copy in our components
cp ${BIN}/tether-linux $(rootfs_dir $PKGDIR)/bin/tether
cp ${BIN}/rpctool $(rootfs_dir $PKGDIR)/sbin/

cp ${DIR}/bootstrap/rc.local $(rootfs_dir $PKGDIR)/etc/rc.d/rc.local

generate_iso $PKGDIR $BIN/$ISONAME $INIT
