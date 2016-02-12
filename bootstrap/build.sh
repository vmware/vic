#!/bin/bash

#
# Build a booststrap iso bootable by VIC.
#

[ -n "$DEBUG" ] && set -x

set -e

function usage () {
     echo "$0 { docker | local }"
     exit 1
}

if [ $# -lt 1 ]; then
     usage
fi

BINARY=${BINARY:-../binary}
name=bootstrap

git_args=(--git-dir=../.git --work-tree=../)
branch_name=$(git "${git_args[@]}" symbolic-ref HEAD 2>/dev/null)
branch=${branch_name##refs/heads/}
sha=$(git "${git_args[@]}" rev-parse --short HEAD)
date=$(date -u +%Y/%m/%d_@_%H:%M:%S)
BUILD_ID="$date@$branch:$sha"

# The iso being created
isofile=${ISOFILE:-../binary/bootstrap-${sha}.iso}

# The vendor directory.
vendor=${VENDOR:-../vendor}

# Grab the kernel
PHOTON_RPM_REPO=${PHOTON_RPM_REPO:-http://bonneville.eng.vmware.com:8080/job/bonneville-kernel/lastSuccessfulBuild/artifact/binary/}
photon_kernel_rpm=linux-esx-4.2.0-10.x86_64.rpm

if [ ! -d "$BINARY" ]; then
     mkdir "$BINARY"
fi

# we may already have the kernel rpm
if [ ! -f "$BINARY"/${photon_kernel_rpm} ]; then
     curl -L "$PHOTON_RPM_REPO"/${photon_kernel_rpm} -o "$BINARY"/${photon_kernel_rpm}
fi

TCL_BASE=${TCL_BASE:-http://tinycorelinux.net/6.x}
tcl_core_path=$TCL_BASE/x86_64/release/distribution_files/
tcl_core=corepure64.gz
tcl_tcz=$TCL_BASE/x86_64/tcz
tcz_deps=(iproute2.tcz libtirpc.tcz)

# Download the tcl core root fs
if [ ! -f "$BINARY"/${tcl_core} ]; then
     curl -L "${tcl_core_path}"/${tcl_core} -o "$BINARY"/${tcl_core}
fi

# Download tinycore dependencies
for dep in "${tcz_deps[@]}"; do
     if [ ! -f "$BINARY"/$dep ]; then
          echo "Downloading ${tcl_tcz}/$dep"
          curl -L "${tcl_tcz}"/$dep -o "$BINARY"/$dep
     fi
done

case $1 in
     "docker")
          (cd ../ && docker build \
               -f=./Dockerfile.bootstrap-linux \
               -t ${name}-"${sha}" \
               --build-arg BUILD_ID="$BUILD_ID" \
               --build-arg KERNEL_RPM="$BINARY"/${photon_kernel_rpm} \
               --build-arg TINYCORE="$BINARY"/${tcl_core} \
               --build-arg TINYCORE_DEPS="$BINARY/${tcz_deps[0]} $BINARY/${tcz_deps[1]}" \
               .)
          docker run --rm ${name}-"$sha":latest > "$isofile"
     ;;

     "local")
          if [ "$(id -u)" != "0" ]; then
               fakeroot=fakeroot
          fi
          make -C ../ bootstrap
          BUILD_ID=$BUILD_ID $fakeroot ./build-iso.sh \
               -i /usr/lib/syslinux/isolinux.bin \
               -k "$BINARY"/$photon_kernel_rpm \
               -c "$BINARY"/${tcl_core} \
               -d "$BINARY/${tcz_deps[0]} $BINARY/${tcz_deps[1]}" \
               "$isofile"
     ;;

     *)
     usage
     ;;
esac

