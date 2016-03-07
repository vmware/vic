#!/bin/bash -e

# Deploy Vagrant box to esx
set -e

if [ "$(uname -s)" = "Darwin" ]; then
  PATH="/Applications/VMware Fusion.app/Contents/Library:$PATH"
fi

export GOVC_URL=${GOVC_URL-"root:vagrant@localhost:18443"}
export GOVC_DATASTORE=${GOVC_DATASTORE-"datastore1"}
export GOVC_NETWORK=${GOVC_NETWORK-"VM Network"}
export GOVC_INSECURE=1

echo "deploying to $(awk -F@ '{print $2}' <<<"$GOVC_URL"):"
govc about

config="$(git rev-parse --show-toplevel)/Vagrantfile"
box=$(grep config.vm.box "$config" | awk -F'"' '{print $2}')
provider=$(dirname "$box")
name=$(basename "$box")
disk="${name}.vmdk"

pushd "$(dirname "$0")" >/dev/null

if ! govc datastore.ls "${name}/${disk}" 1>/dev/null 2>&1 ; then
    if [ ! -e "$disk" ] ; then
        src=$(echo ~/.vagrant.d/boxes/"${provider}"-*-"${name}"/*.*.*/vmware_desktop/disk.vmdk)

        if [ ! -e "$src" ] ; then
            echo "box not found, install via: vagrant box add --provider vmware_desktop $box"
            exit 1
        fi

        echo "converting vagrant box for ESX..."
        vmware-vdiskmanager -r "$src" -t 0 "$disk"
    fi

    echo "importing vmdk to datastore ${GOVC_DATASTORE}..."
    govc import.vmdk "$disk" "$name"
fi

vm_name=${VM_NAME-"${USER}-${name}"}
vm_memory=${VM_MEMORY-$(grep memsize "$config" | awk -F'"' '{print $4}')}

if [ -z "$(govc ls "vm/$vm_name")" ] ; then
    echo "creating VM ${vm_name}..."

    govc vm.create -m "$vm_memory" -c 2 -g ubuntu64Guest -on=false "$vm_name"

    govc vm.disk.attach -vm "$vm_name" -link=true -disk "$name/$disk"

    govc device.cdrom.add -vm "$vm_name"

    govc vm.power -on "$vm_name"
fi

# An ipv6 is reported by tools when the machine is first booted.
# Wait until we get an ipv4 address from tools.
while true
do
    ip=$(govc vm.ip "$vm_name")
    ipv4="${ip//[^.]}"
    if [ "${#ipv4}" -eq 3 ]
    then
        break
    fi
    sleep 1
done

echo "VM ip=$ip"

for script in provision*.sh ; do
    echo "Applying $script..."
    ssh -i ~/.vagrant.d/insecure_private_key "vagrant@$ip" sudo bash -s - < "$script"
done
