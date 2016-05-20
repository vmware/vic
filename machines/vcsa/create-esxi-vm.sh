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

usage() {
    echo "Usage: $0 [-d DISK_GB] [-m MEM_GB] [-i ESX_ISO] [-s] ESX_URL VM_NAME" 1>&2
    exit 1
}

disk=48
mem=16
iso=VMware-VMvisor-6.0.0-3634798.x86_64.iso

while getopts d:i:m:s: flag
do
    case $flag in
        d)
            disk=$OPTARG
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
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ $# -ne 2 ] ; then
    usage
fi

export GOVC_INSECURE=1
export GOVC_URL=$1
shift

name=$1
shift

boot=$(basename "$iso")
if ! govc datastore.ls "$boot" > /dev/null 2>&1 ; then
    govc datastore.upload "$iso" "$boot"
fi

echo "Creating vm ${name}..."
govc vm.create -on=false -net "VM Network" -m $((mem*1024)) -c 2 -g "vmkernel6Guest" -net.adapter=e1000e "$name"

echo "Adding a second nic for ${name}..."
govc vm.network.add -net "VM Network" -net.adapter=e1000e -vm "$name"

echo "Enabling nested hv for ${name}..."
govc vm.change -vm "$name" -nested-hv-enabled

echo "Adding cdrom device to ${name}..."
id=$(govc device.cdrom.add -vm "$name")

echo "Inserting $boot in $name cdrom device..."
govc device.cdrom.insert -vm "$name" -device "$id" "$boot"

if [ -n "$standalone" ] ; then
    echo "Creating $name disk for use by ESXi..."
    govc vm.disk.create -vm "$name" -name "$name"/disk1 -size "${disk}G"
else
    echo "Creating $name disks for use by vSAN..."
    govc vm.disk.create -vm "$name" -name "$name"/vsan-cache -size "$((disk/2))G"
    govc vm.disk.create -vm "$name" -name "$name"/vsan-store -size "${disk}G"

    govc vm.change -e scsi0:0.virtualSSD=1 -e scsi0:1.virtualSSD=0 -vm "$name"
fi

echo "Powering on $name VM..."
govc vm.power -on "$name"

echo "Waiting for $name ESXi IP..."
vm_ip=$(govc vm.ip "$name")

! govc events -n 100 "vm/$name" | egrep 'warning|error'

# extract password from $GOVC_URL
password=$(govc env | grep GOVC_PASSWORD= | cut -d= -f 2-)

GOVC_URL="root:@${vm_ip}"
echo "Waiting for $name hostd (via GOVC_URL=$GOVC_URL)..."
while true; do
    if govc about 2>/dev/null; then
        break
    fi

    printf "."
    sleep 1
done

if [ -z "$standalone" ] ; then
    # Stateless esx will claim disks for its own use,
    # vSAN cannot autoclaim mounted disks, so unmount.
    echo "Unmounting datastores for use with vSAN..."

    govc ls datastore | xargs -n1 -I% govc datastore.remove -ds % '*'
fi

echo "Installing host client..."
govc host.esxcli -- software vib install -v http://download3.vmware.com/software/vmw-tools/esxui/esxui-signed-3843236.vib

echo "Enabling MOB..."
govc host.option.set Config.HostAgent.plugins.solo.enableMob true

echo "Propagating \$GOVC_URL password to $name host root account..."
govc host.account.update -id root -password "$password"

echo "Done."
