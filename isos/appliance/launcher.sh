#!/bin/bash
# Copyright 2016 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


logFileDir="/var/log/vic/"
mkdir -p "$logFileDir"

nameInterfaces() {
    echo "Waiting for network interface configuration"
    while true ; do
        networks=($(rpctool -get vch/networks))
        if [ ${#networks[@]} -ne 0 ] ; then
            break
        fi
        sleep 1
    done

    for net in "${networks[@]}" ; do
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

    echo "${ip} client.localhost" >> /etc/hosts
    if [ -d /sys/class/net/management ] ; then
        ip=$(ip addr show dev management | grep "inet " | tr -s ' ' | cut -d ' ' -f 3 | cut -d '/' -f 1)
    fi
    echo "${ip} management.localhost" >> /etc/hosts
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
    pid=$(cat $pidfile 2>/dev/null)
    cargs=$(cat $cmdfile 2>/dev/null)

    # if arguments have changed, kill it
    if [ -e "${pifdile}" -a "$cargs" != "$args" ]; then
        echo "Component $1 arguments have changed - relaunching" | tee -a $logfile

        for attempt in {1..5}; do
            kill "$pid" 2>/dev/null || { echo "Clean shutdown of component" && break; }
            sleep 1
        done
        # ensure it's gone
        kill -9 "$pid" >/dev/null 2>&1

        rm -f "$pidfile"
    fi

    # if we expect the component to be there and it's not
    if [ -e "$pidfile" ] && ! kill -0 "$pid" 2>/dev/null; then
        echo "Component $1 has died - relaunching" | tee -a $logfile
        rm -f "$pidfile"
    fi

    # Launch it if it's not running
    if [ -z "$pid" ] || ! kill -0 "$pid" 2>/dev/null; then
        echo "Launching $1 $2" | tee -a $logfile
        echo "$args" > ${cmdfile}

        # we change IFS and then replace
        # all " -" with ",-" so that bash can
        # escape from regular spaces in $args
        IFS=","
        "${1}" ${args//\ -/,-} >>$logfile 2>&1 &
        IFS="  "
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
