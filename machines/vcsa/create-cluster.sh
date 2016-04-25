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
# Configure a vCenter cluster with vSAN datastore, DVS and DVPGs

export GOVC_INSECURE=1
export GOVC_USERNAME=${GOVC_USERNAME:-"Administrator@vsphere.local"}
if [ -z "$GOVC_PASSWORD" ] ; then
    # extract password from $GOVC_URL
    eval "$(govc env | grep GOVC_PASSWORD=)"
    export GOVC_PASSWORD
fi

usage() {
    echo "Usage: $0 [-d DATACENTER] [-c CLUSTER] VCSA_IP ESX_IP1 ESX_IP2 ESX_IP3..." 1>&2
    exit 1
}

# Defaults
dc_name="dc1"
cluster_name="cluster1"
vsan_vnic="vmk0"

while getopts c:d: flag
do
    case $flag in
        c)
            cluster_name=$OPTARG
            ;;
        d)
            dc_name=$OPTARG
            ;;
        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ $# -lt 4 ] ; then
    usage
fi

vc_ip=$1
shift

export GOVC_URL="${GOVC_USERNAME}:${GOVC_PASSWORD}@${vc_ip}"

cluster_path="/$dc_name/host/$cluster_name"
dvs_path="/$dc_name/network/DSwitch"
external_network="/$dc_name/network/ExternalNetwork"
internal_network="/$dc_name/network/InternalNetwork"

if [ -z "$(govc ls "/$dc_name")" ] ; then
    echo "Creating datacenter ${dc_name}..."
    govc datacenter.create "$dc_name"
fi

if [ -z "$(govc ls "$cluster_path")" ] ; then
    echo "Creating cluster ${cluster_path}..."
    govc cluster.create "$cluster_name"
fi

if [ -z "$(govc ls "$dvs_path")" ] ; then
    echo "Creating dvs ${dvs_path}..."
    govc dvs.create -parent "$(dirname "$dvs_path")" "$(basename "$dvs_path")"
fi

if [ -z "$(govc ls "$external_network")" ] ; then
    govc dvs.portgroup.add -dvs "$dvs_path" -type earlyBinding -nports 16 "$(basename "$external_network")"
fi

if [ -z "$(govc ls "$internal_network")" ] ; then
    govc dvs.portgroup.add -dvs "$dvs_path" -type ephemeral "$(basename "$internal_network")"
fi

hosts=()

for host_ip in "$@" ; do
    host_path="$cluster_path/$host_ip"
    hosts+=($host_path)

    if [ -z "$(govc ls "$host_path")" ] ; then
        echo "Adding host ($host_ip) to cluster $cluster_name"
        govc cluster.add -cluster "$cluster_path" -noverify -force \
             -hostname "$host_ip" -username root -password "$GOVC_PASSWORD"
    fi

    echo "Enabling vSAN traffic on ${vsan_vnic} for ${host_path}..."
    govc host.vnic.service -host "$host_path" -enable vsan "$vsan_vnic"

    echo "Opening firewall for serial port traffic for ${host_path}..."
    govc host.esxcli -host "$host_path" -- network firewall ruleset set -r remoteSerialPort -e true
done

govc dvs.add -dvs "$dvs_path" -pnic vmnic1 "${hosts[@]}"

echo "Enabling DRS and vSAN for ${cluster_path}..."
govc cluster.change -drs-enabled -vsan-enabled -vsan-autoclaim "$cluster_path"

echo "Granting Admin permissions for user root..."
govc permissions.set -principal root -role Admin

echo "Done."
