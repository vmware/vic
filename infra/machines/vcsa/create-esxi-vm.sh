#!/bin/bash -e
# Copyright 2016 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Create a VM and boot stateless ESXi via cdrom/iso

set -o pipefail

usage() {
    cat <<'EOF'
Usage: $0 [-d DISK_GB] [-m MEM_GB] [-i ESX_ISO] [-s] ESX_URL VM_NAME

GOVC_* environment variables also apply, see https://github.com/vmware/govmomi/tree/master/govc#usage
If GOVC_USERNAME is set, it is used to create an account on the ESX vm.  Default is to use the existing root account.
If GOVC_PASSWORD is set, the account password will be set to this value.  Default is to use the given ESX_URL password.
EOF
}

disk=48
mem=16
iso=VMware-VMvisor-6.0.0-3620759.x86_64.iso # 6.0u2
vib=http://download3.vmware.com/software/vmw-tools/esxui/esxui-signed-3976049.vib

while getopts c:d:hi:m:s flag
do
    case $flag in
        c)
            vib=$OPTARG
            ;;
        d)
            disk=$OPTARG
            ;;
        h)
            usage
            exit
            ;;
        i)
            iso=$OPTARG
            ;;
        m)
            mem=$OPTARG
            ;;
        s)
            standalone=true
            ;;
        *)
            usage 1>&2
            exit 1
            ;;
    esac
done

shift $((OPTIND-1))

if [ $# -ne 2 ] ; then
    usage
fi

if [[ "$iso" == *"-Installer-"* ]] ; then
    echo "Invalid iso name (need stateless, not installer): $iso" 1>&2
    exit 1
fi

echo -n "Checking govc version..."
govc version -require 0.8.0

username=$GOVC_USERNAME
password=$GOVC_PASSWORD
unset GOVC_USERNAME GOVC_PASSWORD

if [ -z "$password" ] ; then
    # extract password from $GOVC_URL
    password=$(govc env | grep GOVC_PASSWORD= | cut -d= -f 2-)
fi

export GOVC_INSECURE=1
export GOVC_URL=$1
export GOVC_DATASTORE=${GOVC_DATASTORE:-$(basename "$(govc ls datastore)")}
network=${GOVC_NETWORK:-"VM Network"}
shift

name=$1
shift

boot=$(basename "$iso")
if ! govc datastore.ls "$boot" > /dev/null 2>&1 ; then
    govc datastore.upload "$iso" "$boot"
fi

echo "Creating vm ${name}..."
govc vm.create -on=false -net "$network" -m $((mem*1024)) -c 2 -g "vmkernel6Guest" -net.adapter=e1000e "$name"

echo "Adding a second nic for ${name}..."
govc vm.network.add -net "$network" -net.adapter=e1000e -vm "$name"

echo "Enabling nested hv for ${name}..."
govc vm.change -vm "$name" -nested-hv-enabled

echo "Enabling Mac Learning dvFilter for ${name}..."
seq 0 1 | xargs -I% govc vm.change -vm "$name" \
                -e ethernet%.filter4.name=dvfilter-maclearn \
                -e ethernet%.filter4.onFailure=failOpen

echo "Adding cdrom device to ${name}..."
id=$(govc device.cdrom.add -vm "$name")

echo "Inserting $boot into $name cdrom device..."
govc device.cdrom.insert -vm "$name" -device "$id" "$boot"

echo "Powering on $name VM..."
govc vm.power -on "$name"

echo "Waiting for $name ESXi IP..."
vm_ip=$(govc vm.ip "$name")

! govc events -n 100 "vm/$name" | egrep 'warning|error'

esx_url="root:@${vm_ip}"
echo "Waiting for $name hostd (via GOVC_URL=$esx_url)..."
while true; do
    if govc about -u "$esx_url" 2>/dev/null; then
        break
    fi

    printf "."
    sleep 1
done

if [ -n "$standalone" ] ; then
    echo "Creating $name disk for use by ESXi..."
    govc vm.disk.create -vm "$name" -name "$name"/disk1 -size "${disk}G"
else
    echo "Creating $name disks for use by vSAN..."
    govc vm.disk.create -vm "$name" -name "$name"/vsan-cache -size "$((disk/2))G"
    govc vm.disk.create -vm "$name" -name "$name"/vsan-store -size "${disk}G"
fi

# Set target to the ESXi VM
GOVC_URL="$esx_url"

if [ -n "$standalone" ] ; then
    disk=$(govc host.storage.info -rescan | grep /vmfs/devices/disks | awk '{print $1}' | xargs basename)
    echo "Creating datastore on disk ${disk}..."
    govc datastore.create -type vmfs -name datastore1 -disk="$disk" '*'
else
    echo "Rescanning HBA for new devices..."
    disk=($(govc host.storage.info -rescan | grep /vmfs/devices/disks | awk '{print $1}' | sort))

    echo "Marking disk ${disk[0]} as SSD..."
    govc host.storage.mark -ssd "${disk[0]}"

    echo "Marking disk ${disk[1]} as HDD..."
    govc host.storage.mark -ssd=false "${disk[1]}"
fi

echo "Enabling MOB..."
govc host.option.set Config.HostAgent.plugins.solo.enableMob true

echo "Enabling ESXi Shell and SSH..."
for id in TSM TSM-SSH ; do
    govc host.service enable $id
    govc host.service start $id
done

if [ -z "$username" ] ; then
    username=root
    action="update"
else
    action="create"
fi

echo "ESX host account $action for user $username..."
govc host.account.$action -id $username -password "$password"

echo "Granting Admin permissions for user $username..."
govc permissions.set -principal $username -role Admin

echo "Enabling guest ARP inspection to get vm IPs without vmtools..."
govc host.esxcli system settings advanced set -o /Net/GuestIPHack -i 1

if [ -n "$vib" ] ; then
    echo -n "Installing host client ($(basename "$vib"))..."

    if govc host.esxcli -- software vib install -v "$vib" > /dev/null 2>&1 ; then
        echo "OK"
    else
        echo "Failed"
    fi
fi

echo "Done."
