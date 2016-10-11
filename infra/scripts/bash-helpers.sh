#!/bin/bash

tls () {
        unset TLS_OPTS
}

no-tls () {
        export TLS_OPTS="--no-tls"
}

unset-vic () {
	unset MAPPED_NETWORKS NETWORKS IMAGE_STORE DATASTORE COMPUTE VOLUME_STORES IPADDR GOVC_INSECURE TLS
	unset DOCKER_CERT_PATH DOCKER_TLS_VERIFY
	unalias docker 2>/dev/null
}

vic-path () {
	echo "${GOPATH}/src/github.com/vmware/vic"
}

vic-deploy () {
        pushd $(vic-path)/bin/

        $(vic-path)/bin/vic-machine-linux create -target="$GOVC_URL" -image-store="$IMAGE_STORE" -compute-resource="$COMPUTE" ${TLS} ${TLS_OPTS} --name=${VIC_NAME:-${USER}test} ${MAPPED_NETWORKS} ${VOLUME_STORES} ${NETWORKS} ${IPADDR} ${TIMEOUT} $*

	envfile=${VIC_NAME:-${USER}test}/${VIC_NAME:-${USER}test}.env
	if [ -f "$envfile" ]; then
		set -a
		source $envfile
		set +a
	fi

	if [ -z ${DOCKER_TLS_VERIFY+x} ]; then
		alias docker='docker --tls'
	fi

        popd
}

vic-rm () {
        $(vic-path)/bin/vic-machine-linux delete -target="$GOVC_URL" -compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} --force $*
}

vic-inspect () {
       $(vic-path)/bin/vic-machine-linux inspect -target="$GOVC_URL" -compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} $*
}

vic-ls () {
	$(vic-path)/bin/vic-machine-linux ls -target="$GOVC_URL" $*
}

vic-ssh () {
	unset keyarg
	if [ -e $HOME/.ssh/authorized_keys ]; then
		keyarg="--authorized-key=$HOME/.ssh/authorized_keys"
	fi

        out=$($(vic-path)/bin/vic-machine-linux debug -target="$GOVC_URL" -compute-resource="$COMPUTE" --name=${VIC_NAME:-${USER}test} --enable-ssh $keyarg --rootpw=password $*)
	host=$(echo $out | grep DOCKER_HOST | sed -n 's/.*DOCKER_HOST=\([^i:]*\).*/\1/p')

	echo "SSH to ${host}"
	sshpass -ppassword ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@${host}
}

addr-from-dockerhost () {
	echo $DOCKER_HOST | sed -e 's/:[0-9]*$//'
}

# import the custom sites
# exmaple entry, actived by typing "example"
#example () {
#	target='https://user:password@host.domain.com/datacenter'
#	unset-vic
#
#	export GOVC_URL=$target
#	export COMPUTE=cluster/pool
#	export DATASTORE=datastore1
#	export IMAGE_STORE=$DATASTORE/image/path
#	export NETWORKS="--bridge-network=private-dpg-vlan --external-network=extern-dpg"
#	export TIMEOUT="--timeout=10m"
#	export IPADDR="--client-network-ip=vch-hostname.domain.com --client-network-gateway=x.x.x.x/22 --dns-server=y.y.y.y --dns-server=z.z.z.z"
#	export TLS="--tls-cname=vch-hostname.domain.com --organisation=MyCompany"
#	export VIC_NAME="MyVCH"
#}

. ~/.vic
