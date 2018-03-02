#!/bin/bash

REMOTE_DLV_ATTACH=/usr/local/bin/dlv-attach-headless.sh
REMOTE_DLV_DETACH=/usr/local/bin/dlv-detach-headless.sh

function usage() {
    echo "Usage: $0 -h vch-address [-a/-d] -p password [attach/detach] target" >&2
    echo "Valid targets are: "
    echo "    vic-init"
    echo "    vic-admin"
    echo "    docker-engine"
    echo "    port-layer"
    echo "    virtual-kubelet"
    exit 1
}

while getopts "h:p:ad" flag
do
    case $flag in

        h)
            # Optional
            export VCH_HOST="$OPTARG"
            ;;

        p)
            export SSHPASS="$OPTARG"
            ;;

        a)
            export COMMAND="attach"
            ;;

        d)
            export COMMAND="detach"
            ;;

        *)
            usage
            ;;
    esac
done

echo $OPTIND
shift $((OPTIND-1))

if [[ -z "${COMMAND}" &&  $# != 2 ]]; then
    usage
elif [[ -n "${COMMAND}" && $# != 1 ]]; then
    usage
fi

if [ -z "${COMMAND}" ]; then
    COMMAND=$1
    TARGET=$2
else
    TARGET=$1
fi

case ${TARGET} in

    vic-init)
        PORT=2345
        ;;

    vic-admin)
        PORT=2346
        ;;

    docker-engine)
        PORT=2347
        ;;

    port-layer)
        PORT=2348
        ;;

    virtual-kubelet)
        PORT=2349
        ;;

    *)
        usage
        ;;
esac

if [ -z "${VCH_HOST}" ]; then
    usage
fi

if [ ${COMMAND} == "attach" ]; then
    if [ -n "${SSHPASS}" ]; then
        sshpass -e ssh root@${VCH_HOST} "nohup /usr/local/bin/dlv-attach-headless.sh $TARGET $PORT > /var/tmp/${TARGET}.log 2>&1 &"
    else
       ssh root@${VCH_HOST} "nohup /usr/local/bin/dlv-attach-headless.sh $TARGET $PORT >  /var/tmp/${TARGET}.log 2>&1 &"
    fi
elif [ ${COMMAND} == "detach" ]; then
    if [ -n "${SSHPASS}" ]; then
        sshpass -e ssh root@${VCH_HOST} "/usr/local/bin/dlv-detach-headless.sh $PORT"
    else
       ssh root@${VCH_HOST} "/usr/local/bin/dlv-detach-headless.sh $PORT"
    fi
else
    usage
fi
