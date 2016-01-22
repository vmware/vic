#!/bin/bash

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/.tether/

vmfork() {
    id=`cat /.tether-init/docker-id`
    newid=$id
    until [ "$id" != "$newid" ];do
        #sleep 5;
        /.tether/modprobe -r vmxnet3
        /.tether/rpctool -fork
        newid=`/.tether/rpctool -get docker_id`
    done

    echo "Updating docker id in forked container"
    echo -n $newid > /.tether-init/docker-id

    hwclock --hctosys &
    #echo "- - -" > /sys/class/scsi_host/host0/scan
    /.tether/modprobe vmxnet3

    #NOTE: consider putting hotadd memory online here if udev isn't running

    hostname ${newid:0:12}
    echo ${newid:0:12} >/etc/hostname
}

vmfork
# restart the tether
echo "Reinitializing the tether"
kill -HUP 1
