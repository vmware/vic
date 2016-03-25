#!/bin/bash

logFileDir="/var/log/vic/"
mkdir -p "$logFileDir"

nameInterfaces() {
    for net in $(rpctool -get vch/networks); do
        mac=$(rpctool -get vch/networks/$net)
        
        # interface name for mac, with trailing :
        dev=$(ip link | grep $mac -B 1 | head -n 1 | cut -d ' ' -f 2 | tr -d ':')
        ip link set dev $dev down
        echo "Renaming $dev to $net"
        ip link set dev $dev name $net
        ip link set dev $net up
    done
}

publishClientIP() {
    echo "Waiting for IP assigning on client interface"
    while [ -z "$ip" ]; do    
        ip=$(ip addr show dev client | grep "inet " | tr -s ' ' | cut -d ' ' -f 3 | cut -d '/' -f 1)
        sleep 1
    done
    echo "Publishing client IP: $ip"    
    rpctool -set IpAddress ${ip}
    rpctool -set vch.clientip ${ip}
}

watchdog() {
    echo "Launching watchdog loop"

    while true; do
        # do we need to launch any components?
        components=$(rpctool -get vch/components)
        for component in ${components}; do
            pidfile=/var/run/$(basename $component).pid

            if [ -e "$pidfile" ]; then
                pid=$(cat $pidfile)
                if ! kill -0 $pid; then
                    echo "Component $component has died - relaunching all services."
                    launch="${components}"
                    break
                fi
            else
                echo "Component $component has no pid file - relaunching all services."
                launch="${components}"
                break
            fi
        done
        
        # if we found a need to relaunch
        if [ ! -z "$launch" ]; then
            # if we relaunching anything, then reload, kill everything, and relaunch
            files=$(rpctool -get vch/files)
            for file in ${files}; do
                if [ "/" == "${file:(-1)}" ]; then
                    # it's a directory
                    mkdir -p $file
                else
                    mkdir -p $(dirname "$file")

                    data=$(rpctool -get vch$file)
                    if [ ! -z "$data" ]; then
                        echo "Generating file: $file"
                        echo "$data" > "${file}"
                    fi
                fi
            done

            for component in ${components}; do
                logfile="$logFileDir/$(basename ${component}).log"
                pidfile=/var/run/$(basename $component).pid
                if [ -e "$pidfile" ]; then
                    pid=$(cat $pidfile)
                    echo "Killing $component via pidfile"
                    kill $pid && echo "Restarting all services" >> $logfile
                    rm -f $pidfile
                fi
            
                echo "Launching component: $component" | tee $logfile
                args=$(rpctool -get vch$component)
                "${component}" ${args} >>$logfile 2>&1 &
                echo $! > /var/run/$(basename $component).pid
            done
        fi

        unset launch
        sleep 15
    done
}

nameInterfaces >>$logFileDir/watchdog.log
publishClientIP >>$logFileDir/watchdog.log

# systemd seems to want this to run in the backgroud, but that doesn't work
# for us currently so this blocks
watchdog >>$logFileDir/watchdog.log