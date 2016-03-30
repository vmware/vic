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
echo "Usage: $0 -c yum-cache(tgz) -p base-package(tgz) -o output-package(tgz) -d <activates debug when set>" 1>&2
exit 1
}

while getopts "c:p:o:d:" flag
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

        d)
            # Optional. directs script to make a debug iso instead of a production iso.
            debug='$OPTARG'
            ;;

        c)
            # Optional. Offline cache of yum packages
            cache="$OPTARG"
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

#################################################################
# Above: arg parsing and setup
# Below: the image authoring
#################################################################

PKGDIR=$(mktemp -d)

unpack $package $PKGDIR

if [ -v debug ]; then
    # These are the packages we install to create an interactive bootstrapVM
    # Install bootstrap base packages
    #
    # packages list here
    #   tndf      # allows package install during debugging.
    #   vim       # basic editing function for debugging.
    yum_cached -c $cache -u -p $PKGDIR install \
        bash \
        shadow \
        tdnf \
        vim \
        -y --nogpgcheck

    # HACK until the issues with override.conf above are dealt with
    pwhash=$(openssl passwd -1 -salt vic password)
    sed -i -e "s/^root:[^:]*:/root:${pwhash}:/" $(rootfs_dir $PKGDIR)/etc/shadow
fi

# Install bootstrap base packages
#
# List stable packages here
#   iproute2  # for ip
#   libtirpc  # due to a previous package reliance on rpc
#
yum_cached -c $cache -u -p $PKGDIR install \
    iproute2 \
    libtirpc \
    systemd \
    -y --nogpgcheck

# https://www.freedesktop.org/wiki/Software/systemd/InitrdInterface/
touch $(rootfs_dir $PKGDIR)/etc/initrd-release

# ensure we're not including a cache in the staging bundle
# but don't update the cache bundle we're using to install
yum_cached -p $PKGDIR clean all

# package up the result
pack $PKGDIR $OUT
