#!/bin/bash
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
# Install basic VCH appliance
#

if [ -n "$DEBUG" ]; then
      set -x
fi

THIS=$(readlink -f "$0")
DIR=$(dirname "$THIS")

# exit on failure
set -e

function usage() {
     echo "# Usage: $0 -t=target-url -p=compute-resource -i=image-datastore [-d=container-datastore] [-e=external-network] [-m=management-network] [-b=bridge-network] [-a=appliance-iso] [-c=bootstrap] [-g=stub] [-x=certificate-file] [-y=key-file] [-v:verbose] [-f] name" 2>&1
     echo "#   -g: generate the certificate and key files, using the value as a stub name"
     echo "#   -f: delete existing VM and image store if found"

     exit 1
}

# pick up GOVC_URL if present
targetURL=${GOVC_URL}
datastore=${GOVC_DATASTORE}

# defaults
externalNet="VM Network"
managementNet="VM Network"
clientNet="VM Network"

applianceIso="${DIR}/appliance.iso"
bootstrapIso="${DIR}/bootstrap.iso"


while getopts "fvt:g:p:i:d:e:m:b:a:c:x:y:" flag
do
  case $flag in
    v)
     # verbose debug
      set -x
      ;;

    t)
     # Required. Target URL - translated to GOVC_URL
      targetURL=${OPTARG}
      export GOVC_URL=${targetURL}
      ;;

    p)
     # Required. Compute resource path
     # This should fully specify the path to the compute resource, e.g.
     # /ha-datacenter/host/myCluster/Resources/myRP
      compute=${OPTARG}
      # primarily internal convenience for constructing other paths
      datacenter=$(echo $compute | cut -d '/' -f 2)
      ;;

    d)
      # Optional. Container datastore path - defaults to image datastore
      cdatastore=${OPTARG}
      ;;

    i)
      # Required. Image datastore path - sets cdatastore too for default
      idatastore=${OPTARG}
      cdatastore=${OPTARG}
      ;;

    e)
      # Optional. The external network (can see hub.docker.com)
      externalNet=${OPTARG}
      : ${GOVC_NETWORK:=$externalNet}
      ;;

    m)
      # Optional. The management network (can see target)
      managementNet=${OPTARG}
      # this is the one we *must* have
      export GOVC_NETWORK=${managementNet}
      ;;

    b)
      # Optional. The bridge network
      bridgeNet=$OPTARG
      ;;

    a)
      # Optional. The appliance iso
      applianceIso=${OPTARG}
      ;;

    c)
      # Optional. The bootstrap iso
      bootstrapIso=${OPTARG}
      ;;

    f)
      # Optional. Force the install, removing existing if present
      force=1
      ;;

    g)
      # Optional. Generate the cert and key and store them in $OPTARG-{cert,key}.pem
      keyf="${OPTARG}-key.pem"
      certf="${OPTARG}-cert.pem"
      # generate
      echo "# Generating certificate/key pair - private key in ${keyf}"
      openssl req -batch -nodes -new -x509  -keyout "${keyf}" -out "${certf}" > /dev/null 2>&1

      # load the bits
      key=$(cat "$keyf")
      certificate=$(cat "$certf")
      ;;

    x)
      # Optional. Certificate file
      certf="${OPTARG}"
      certificate=$(cat "$OPTARG")
      ;;

    y)
      # Optional. Key file
      keyf="${OPTARG}"
      key=$(cat "$OPTARG")
      ;;

    *)
    usage
    ;;
  esac
done

shift $((OPTIND-1))

vchName="$1"
if [ "$#" -ne 1 ]; then
   echo "Trailing arguments after name paramenter ($vchName)"
   usage
fi

if [ -z "$1" -o -z "$GOVC_URL" -o -z "$vchName" -o -z "$compute" -o -z "$idatastore" ]; then
     usage
fi


# if not set explicitly, bridge network has vch name
: ${bridgeNet:=$vchName}

# for now it's all insecure
export GOVC_INSECURE=true
export GOVC_PERSIST_SESSION=true

# login and persist
echo "# Logging into the target"
govc about > /dev/null

# delete target if present
if [ ! -z "${force}" ]; then
   echo "# Cleaning up prior VM if needed"
   govc vm.destroy "${vchName}" 2>/dev/null || echo "# Target VM does not need removing"
   govc datastore.rm -ds "${idatastore}" "${vchName}" 2>/dev/null|| echo "# Target VM does not need cleaning"

   echo "# Cleaning up image store if present"
   govc datastore.rm -ds "${idatastore}" "VIC" 2>/dev/null|| echo "# Target image store does not need cleaning"
fi

# upload the isos
echo "# Uploading ISOs"
appIsoPath="${vchName}/$(basename ${applianceIso})"
cIsoPath="${vchName}/$(basename ${bootstrapIso})"

govc datastore.mkdir "${vchName}"
govc datastore.upload "${applianceIso}" "${appIsoPath}"
if [ -e ${bootstrapIso} ]; then
    govc datastore.upload "${bootstrapIso}" "${cIsoPath}"
fi

# check the bridge network port group and create if needed
govc host.vswitch.info | grep -e "Portgroup:\\s*${bridgeNet}$" >/dev/null || {
   echo "# Creating vSwitch"
   govc host.vswitch.add "${bridgeNet}" || echo "# Switch already exists"
   echo "# Creating Portgroup"
   govc host.portgroup.add -vswitch="${bridgeNet}" "${bridgeNet}"
}


# determine how many network interfaces we need. External has a default value
declare -A networks
# these are ordered such that the last value set against a key will be the interface name
networks[${externalNet}]=external
networks[${managementNet}]=management
networks[${clientNet}]=client
networks[${bridgeNet}]=bridge

# create the VM
echo "# Creating the Virtual Container Host appliance"
govc vm.create -iso="${appIsoPath}" -net="${bridgeNet}" -disk.controller=pvscsi -net.adapter=vmxnet3 -m 2048 -on=false -g other3xLinux64Guest "${vchName}"

# short-hand the vm
uuid=$(govc vm.info ${vchName}| grep -e "^\\s*UUID:" | awk '{print$2}')
vmpath=$(govc vm.info ${vchName}| grep -e "^\\s*Path:" | awk '{print$2}')

echo "# Adding network interfaces"
for net in "${!networks[@]}"; do
   # squash the bridgeNet as we already added that in create
   if [ "$net" != "$bridgeNet" ]; then
      govc vm.network.add -net="${net}" -net.adapter=vmxnet3 -vm.uuid=${uuid}
   fi
done


# ensure we have defaults for datastore and compute if needed
if [ -z "${datastore}" ]; then
   datastore=$(govc datastore.info | grep -e "^Name:" | awk '{print$2}')
fi
if [ -z "${compute}" ]; then
   compute=$(govc pool.info */Resources| grep -e "^\\s*Path:" | awk '{print$2}')
fi


echo "# Setting component configuration"
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/components="/sbin/docker-engine-server /sbin/port-layer-server /sbin/vicadmin"
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/sbin/imagec="-debug -logfile=/var/log/vic/imagec.log -insecure"
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/sbin/port-layer-server="--host=localhost --port=8080 --insecure --sdk=${targetURL} --datacenter=${datacenter} --cluster=${compute} --datastore=/${datacenter}/datastore/${idatastore} --network=/ha-datacenter/network/${externalNet} --vch=${vchName}"
files="/var/tmp/images/ /var/log/vic/"

# now we see if we configure TLS
if [ -n "${certificate}" -a -n "${key}" ] ; then
   echo "# Configuring TLS server"
   govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/etc/pki/tls/certs/vic-host-cert.pem="${certificate}"
   govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/etc/pki/tls/certs/vic-host-key.pem="${key}"
   files="${files} /etc/pki/tls/certs/vic-host-cert.pem /etc/pki/tls/certs/vic-host-key.pem"
   dockertlsargs="-TLS -tls-certificate=/etc/pki/tls/certs/vic-host-cert.pem -tls-key=/etc/pki/tls/certs/vic-host-key.pem"
   vicadmintlsargs=" -hostcert=/etc/pki/tls/certs/vic-host-cert.pem -hostkey=/etc/pki/tls/certs/vic-host-key.pem"
   port=2376
else
   port=2375
fi

# and finalize the config (this is the components that have frontend TLS considerations)
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/files="${files}"
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/sbin/docker-engine-server="-serveraddr=0.0.0.0 -port=${port} -port-layer-port=8080 ${dockertlsargs}"
govc vm.change -vm.uuid="${uuid}" -e guestinfo.vch/sbin/vicadmin="-docker-host=unix:///var/run/docker.sock -insecure -sdk=${targetURL} -ds=/${datacenter}/datastore/${idatastore} -vm-path=${vmpath} ${vicadmintlsargs}"


# power on the appliance
echo "# Powering on the Virtual Container Host"
powerstatus=$(mktemp)
govc vm.power -on=true -vm.uuid="${uuid}" > ${powerstatus} 2>&1 || {
   cat ${powerstatus} && rm -f ${powerstatus} && exit 1
}

echo "# Setting network identities"
# Generate the network interface identiy map - MACs aren't generated until power on
for net in "${!networks[@]}"; do
  mac=$(govc device.ls -vm.uuid $uuid | grep "${net}" | grep ethernet | awk '{print $1}' | xargs govc device.info -vm.uuid $uuid | grep MAC | awk '{print$3}')
  govc vm.change -e guestinfo.vch/networks/${networks[$net]}=$mac -vm.uuid="${uuid}"
done
# jump through some hoops to deal with quote behaviour for arrays
govc vm.change -e guestinfo.vch/networks="$(echo "${networks[@]}")" -vm.uuid="${uuid}"

echo "# Waiting for IP information"
while [ -z "$host" ]; do
   host=$(govc vm.info -e -vm.uuid="${uuid}" | grep guestinfo.vch.clientip | awk '{print$2}')
   sleep 1
done

echo "# "
echo "# SSH to appliance (default=root:password)"
echo "# root@${host}"
echo "# "
echo "# Log server:"
echo "# https://${host}:2378"
echo "# "
if [ -n "${dockertlsargs}" ]; then
   echo "# Connect to docker:"
   echo "docker -H ${host}:${port} --tls --tlscert='${certf}' --tlskey='${keyf}'"
   echo
   echo "DOCKER_OPTS=\"--tls --tlscert='${certf}' --tlskey='${keyf}'\""
fi
echo "DOCKER_HOST=${host}:${port}"
