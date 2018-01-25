#!/bin/bash
set -e

for x in $(echo GOVC_USERNAME GOVC_PASSWORD VCH_NAME GOVC_URL GOPATH GOVC_INSECURE); do
    if [[ ! -v $x ]]; then
        echo "Insufficient inputs. Please set $x environment variable and re-execute this script.";
        exit 1
    fi;
done

function get-thumbprint () {
    openssl s_client -connect $GOVC_URL:443 </dev/null 2>/dev/null \
        | openssl x509 -fingerprint -noout \
        | cut -d= -f2
}

$GOPATH/src/github.com/vmware/vic/bin/vic-machine-linux debug --target=$GOVC_URL --name=$VCH_NAME --user=$GOVC_USERNAME --password=$GOVC_PASSWORD --thumbprint=$(get-thumbprint)

on-vch() {
    sshpass -ppassword ssh -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no root@$VCH_IP -C $@ 2>/dev/null
}

VCH_IP="$(govc vm.ip $VCH_NAME)"

for x in port-layer-server docker-engine-server vicadmin vic-init; do
    sshpass -p 'password' scp -oUserKnownHostsFile=/dev/null -oStrictHostKeyChecking=no $GOPATH/src/github.com/vmware/vic/bin/$x root@$VCH_IP:/tmp/$x
    on-vch mv /sbin/$x /tmp/old-$x
    on-vch mv /tmp/$x /sbin/$x
    on-vch chmod 755 /sbin/$x
    pid=$(on-vch ps -e --format='pid,args' | grep $x | grep -v grep | awk '{print $1}')
    on-vch kill -9 $pid
    on-vch rm -f /tmp/old-$x
done

on-vch vic-init &
