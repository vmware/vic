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
# Create a VCSA VM

usage() {
    echo "Usage: $0 [-n VM_NAME] [-i VCSA_OVA] [-a IP] ESX_URL" 1>&2
    exit 1
}

export GOVC_INSECURE=1

name=vcsa
# 6.5.0d - http://pubs.vmware.com/Release_Notes/en/vsphere/65/vsphere-vcenter-server-650d-release-notes.html
ova=VMware-vCenter-Server-Appliance-6.5.0.5500-5318154_OVF10.ova

while getopts a:i:n: flag
do
    case $flag in
        a)
            ip=$OPTARG
            ;;
        i)
            ova=$OPTARG
            ;;
        n)
            name=$OPTARG
            ;;
        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ $# -ne 1 ] ; then
    usage
fi

export GOVC_URL=$1

network=${GOVC_NETWORK:-$(basename "$(govc ls network)")}
datastore=${GOVC_DATASTORE:-$(basename "$(govc ls datastore)")}

if [ "$(uname -s)" = "Darwin" ]; then
    PATH="/Applications/VMware Fusion.app/Contents/Library/VMware OVF Tool:$PATH"
fi

if [ -z "$ip" ] ; then
    mode=dhcp
    ntp=time.vmware.com
else
    mode=static

    # Derive net config from the ESX server
    config=$(govc host.info -k -json | jq -r .HostSystems[].Config)
    gateway=$(jq -r .Network.IpRouteConfig.DefaultGateway <<<"$config")
    dns=$(jq -r .Network.DnsConfig.Address[0] <<<"$config")
    ntp=$(jq -r .DateTimeInfo.NtpConfig.Server[0] <<<"$config")
    route=$(jq -r ".Network.RouteTableInfo.IpRoute[] | select(.DeviceName == \"vmk0\") | select(.Gateway == \"0.0.0.0\")" <<<"$config")
    prefix=$(jq -r .PrefixLength <<<"$route")

    opts=(--prop:guestinfo.cis.appliance.net.addr=$ip
          --prop:guestinfo.cis.appliance.net.prefix=$prefix
          --prop:guestinfo.cis.appliance.net.dns.servers=$dns
          --prop:guestinfo.cis.appliance.net.gateway=$gateway)
fi

# Use the same password as GOVC_URL
password=$(govc env | grep GOVC_PASSWORD= | cut -d= -f 2-)

ovftool --acceptAllEulas --noSSLVerify --skipManifestCheck \
        --X:injectOvfEnv --allowExtraConfig --X:enableHiddenProperties \
        --deploymentOption=tiny --X:waitForIp --powerOn \
        "--name=$name" \
        "--net:Network 1=${network}" \
        "--datastore=${datastore}" \
        "--prop:guestinfo.cis.vmdir.password=$password" \
        "--prop:guestinfo.cis.appliance.root.passwd=$password" \
        --prop:guestinfo.cis.appliance.root.shell=/bin/bash \
        --prop:guestinfo.cis.deployment.node.type=embedded \
        --prop:guestinfo.cis.vmdir.domain-name=vsphere.local \
        --prop:guestinfo.cis.vmdir.site-name=VCSA \
        --prop:guestinfo.cis.appliance.net.addr.family=ipv4 \
        --prop:guestinfo.cis.appliance.net.mode=$mode \
        "${opts[@]}" \
        --prop:guestinfo.cis.appliance.ssh.enabled=True \
        --prop:guestinfo.cis.appliance.ntp.servers="$ntp" \
        --prop:guestinfo.cis.ceip_enabled=False \
        --prop:guestinfo.cis.deployment.autoconfig=True \
        "$ova" "vi://${GOVC_URL}"
