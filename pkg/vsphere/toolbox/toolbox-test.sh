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
# Create (or reuse) a VM to run toolbox and/or toolbox.test
# Requires ESX to be configured with:
# govc host.esxcli system settings advanced set -o /Net/GuestIPHack -i 1

set -o pipefail

vm="toolbox-test-$(uuidgen)"
destroy=true

while getopts n:stv flag
do
    case $flag in
        n)
            vm=$OPTARG
            unset destroy
            ;;
        s)
            start=true
            ;;
        t)
            test=true
            ;;
        *)
            echo "unknown option" 1>&2
            exit 1
            ;;
    esac
done

echo "Building toolbox binaries..."
pushd "$(git rev-parse --show-toplevel)" >/dev/null
go install -v ./cmd/toolbox
go test -i -c ./pkg/vsphere/toolbox -o "$GOPATH/bin/toolbox.test"
popd >/dev/null

iso=coreos_production_iso_image.iso

govc datastore.mkdir -p images

if ! govc datastore.ls images | grep -q $iso ; then
    echo "Downloading ${iso}..."
    if [ ! -e $iso ] ; then
        wget http://beta.release.core-os.net/amd64-usr/current/$iso
    fi

    echo "Uploading ${iso}..."
    govc datastore.upload $iso images/$iso
fi

if [ ! -e config.iso ] ; then
    echo "Generating config.iso..."
    keys=$(cat ~/.ssh/id_[rd]sa.pub)

    dir=$(mktemp -d toolbox.XXXXXX)
    pushd "${dir}" >/dev/null

    mkdir -p drive/openstack/latest

    cat > drive/openstack/latest/user_data <<EOF
#!/bin/bash

# Add ${USER}'s public key(s) to .ssh/authorized_keys
echo "$keys" | update-ssh-keys -u core -a coreos-cloudinit
EOF

    genisoimage -R -V config-2 -o config.iso ./drive

    popd >/dev/null

    mv -f "$dir/config.iso" .
    rm -rf "$dir"
fi

destroy() {
    echo "Destroying VM ${vm}..."
    govc vm.destroy "$vm"
    govc datastore.rm -f "$vm"
}

govc datastore.mkdir -p "$vm"

if ! govc datastore.ls "$vm" | grep -q "${vm}.vmx" ; then
    echo "Creating VM ${vm}..."
    govc vm.create -g otherGuest64 -m 1024 -on=false "$vm"

    if [ -n "$destroy" ] ; then
        trap destroy EXIT
    fi

    device=$(govc device.cdrom.add -vm "$vm")
    govc device.cdrom.insert -vm "$vm" -device "$device" images/$iso

    govc datastore.upload config.iso "$vm/config.iso" >/dev/null
    device=$(govc device.cdrom.add -vm "$vm")
    govc device.cdrom.insert -vm "$vm" -device "$device" "$vm/config.iso"
fi

state=$(govc vm.info -json "$vm" | jq -r .VirtualMachines[].Runtime.PowerState)

if [ "$state" != "poweredOn" ] ; then
    govc vm.power -on "$vm"
fi

echo -n "Waiting for ${vm} ip..."
ip=$(govc vm.ip -esxcli "$vm")

opts=(-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o LogLevel=error -o BatchMode=yes)

scp "${opts[@]}" "$GOPATH"/bin/toolbox{,.test} "core@${ip}:"

if [ -n "$test" ] ; then
    echo "Running toolbox tests..."
    ssh "${opts[@]}" "core@${ip}" ./toolbox.test -test.v -test.run TestServiceRunESX -toolbox.testesx \
        -toolbox.testpid="$$" -toolbox.powerState="$state" &

    echo "Waiting for VM ip from toolbox..."
    ip=$(govc vm.ip "$vm")
    echo "toolbox vm.ip=$ip"

    echo "Testing guest operations via govc..."
    out=$(govc guest.start -vm "$vm" -l user:pass /bin/date)

    if [ "$out" != "$$" ] ; then
        echo "'$out' != '$$'" 1>&2
    fi

    echo "Waiting for tests to complete..."
    wait
fi

if [ -n "$start" ] ; then
    echo "Starting toolbox..."
    ssh "${opts[@]}" "core@${ip}" ./toolbox -toolbox.trace
fi
