#!/bin/bash
set -x
SSH="ssh -o StrictHostKeyChecking=no"
SCP="scp -o StrictHostKeyChecking=no"

function usage() {
    echo "Usage: $0 -h vch-address -p password -s" >&2
    exit 1
}

while getopts "h:p:s" flag
do
    case $flag in

        h)
            # Optional
            export VCH_HOST="$OPTARG"
            ;;

        p)
            export SSHPASS="$OPTARG"
            ;;

        s)
            COPY_SSH_KEY=true
            ;;
        *)
            usage
            ;;
    esac
done

shift $((OPTIND-1))

if [ -z "${VCH_HOST}" -o -z "${SSHPASS}" ]; then
    usage
fi

DLV_BIN="$GOPATH/bin/dlv"

if [ ! -f /usr/bin/sshpass ]; then
    echo sshpass must be installed. Run \"apt-get install sshpass\" 
    exit 1
fi

if [ -z "${GOPATH}" -o -z "${GOROOT}" ]; then
    echo GOROOT and GOPATH should be set to point to the current GOLANG enironment
    exit 1
fi

# copy dlv binary
if [ -f ${DLV_BIN} ]; then
    sshpass -e ${SCP} ${DLV_BIN} root@${VCH_HOST}:/usr/local/bin
else
    echo $DLV_BIN does not exist. Run \"go get github.com/derekparker/delve/cmd/dlv\"
    exit 1
fi

# copy GOROOT env
sshpass -e ${SSH} root@${VCH_HOST} "mkdir -p /usr/local/go"
sshpass -e ${SCP} -r ${GOROOT}/bin root@${VCH_HOST}:/usr/local/go
sshpass -e ${SCP} -r ${GOROOT}/api root@${VCH_HOST}:/usr/local/go
sshpass -e ${SCP} ${GOROOT}/VERSION root@${VCH_HOST}:/usr/local/go
sshpass -e ${SSH} root@${VCH_HOST} "ln -f -s /usr/local/go/bin/go /usr/local/bin/go"

# open IPTABLES
sshpass -e ${SSH} root@${VCH_HOST} "iptables -A INPUT -p tcp -m tcp --dport 2345:2349 -j ACCEPT"

# write remote dlv attach script
TEMPFILE=$(mktemp)
cat > ${TEMPFILE} <<EOF
#/bin/bash
if [ \$# != 2 ]; then
    echo "\$0 vic-init|vic-admin|docker-engine|port-layer|virtual-kubelet port"
    exit 1
fi

NAME=\$1
PORT=\$2

if [ -z "\${NAME}" -o -z "\${PORT}" ]; then
    echo "\$0 vic-init|vic-admin|docker-engine|port-layer|virtual-kubelet port"
    exit 1
fi

PID=\$(ps -e | grep \${NAME} | grep -v grep | tr -s ' ' | cut -d " " -f 2)

if [ -z "\${PID}" ]; then
    echo "\$0: cannot find process \${NAME}"
    exit 1
fi

dlv attach \${PID} --api-version 2 --headless --listen=:\${PORT} 
EOF

sshpass -e ${SCP} ${TEMPFILE} root@${VCH_HOST}:/usr/local/bin/dlv-attach-headless.sh

# write dlv detach script
cat > ${TEMPFILE} <<EOF
#/bin/bash
if [ \$# != 1 ]; then
    echo "\$0 port-number"
    exit 1
fi

PORT=\$1

if [ -z "\${PORT}" ]; then
    echo "\$0 port-number"
    exit 1
fi

# Find appropriate dlv instance
PID=\$(ps -ef | grep "dlv" |  grep "api-version" | grep \${PORT} | grep -v grep |  tr -s ' ' | cut -d " " -f 2)

if [ -z "\${PID}" ]; then
    echo "\$0: cannot find dlv listening on \${PORT}"
    exit 1
fi

kill \${PID}
EOF

sshpass -e ${SCP} ${TEMPFILE} root@${VCH_HOST}:/usr/local/bin/dlv-detach-headless.sh

sshpass -e ${SSH}  root@${VCH_HOST} 'chmod +x /usr/local/bin/*'

rm ${TEMPFILE}

if [ -n "${COPY_SSH_KEY}" ]; then 
    # setup authorized_keys
    sshpass -e ${SSH} root@${VCH_HOST} "mkdir -p .ssh"
    sshpass -e ${SCP} ${HOME}/.ssh/*.pub root@${VCH_HOST}:.ssh
    sshpass -e ${SSH} root@${VCH_HOST} 'cat ~/.ssh/*.pub > ~/.ssh/authorized_keys'
    sshpass -e ${SSH} root@${VCH_HOST} 'rm ~/.ssh/*.pub'
    sshpass -e ${SSH} root@${VCH_HOST} 'chmod 600 ~/.ssh/authorized_keys'
fi
