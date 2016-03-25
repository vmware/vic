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

# args:
# component name
# arguments
launchComponent() {
    component=$(basename $1)
    args=${2}
    pidfile=/var/run/${component}.pid
    cmdfile=/var/run/${component}.cmd
    logfile="$logFileDir/${component}.log"

    # read existing file contents
    pid=$(cat $pidfile)
    cargs=$(cat $cmdfile)

    # if arguments have changed, kill it
    if [ -e "${pifdile}" -a "$cargs" != "$args" ]; then
        echo "Component $1 arguments have changed - relaunching" | tee -a $logfile

        for attempt in {1..5}; do
            kill "$pid" || { echo "Clean shutdown of component" && break; }
            sleep 1
        done
        # ensure it's gone
        kill -9 "$pid" > /dev/null

        rm -f "$pidfile"
    fi

    # if we expect the component to be there and it's not
    if [ -e "$pidfile" ] && ! kill -0 "$pid"; then
        echo "Component $1 has died - relaunching" | tee -a $logfile
        rm -f "$pidfile"
    fi

    # Launch it if it's not running
    if [ -z "$pid" ] || ! kill -0 "$pid"; then
        echo "Launching $1 $2" | tee -a $logfile
        echo "$args" > ${cmdfile}

        "${1}" ${args} >>$logfile 2>&1 &
        echo $! > $pidfile
    fi
}

watchdog() {
    echo "Launching watchdog loop"

    while true; do

        # reload the files in case they're dependencies of components
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

        # run through the components
        components=$(rpctool -get vch/components)
        for component in ${components}; do
            args="$(rpctool -get vch$component)"
            launchComponent "$component" "$args"
        done

        sleep 15
    done
}

nameInterfaces >>$logFileDir/watchdog.log
publishClientIP >>$logFileDir/watchdog.log

# systemd seems to want this to run in the backgroud, but that doesn't work
# for us currently so this blocks
watchdog >>$logFileDir/watchdog.log