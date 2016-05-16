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

load 'vendor/github.com/ztombol/bats-support/load'
load 'vendor/github.com/ztombol/bats-assert/load'

VIC_DIR=$(git rev-parse --show-toplevel)

# starts the port layer server in the background, waits for it to start, saves the pid to $port_layer_pid
# return value 0 on success, 1 on timeout
start_port_layer () {
    # FIXME: need to integrate with ESX_URL so disabling it now
    # https://github.com/vmware/vic/issues/304
    return

    [ "$1" = "" ] && port="8080" || port="$1"
    "$portlayer" --port="$port" --path=/tmp/portlayer 2>&1 &
    port_layer_pid="$!"

    retval=1
    for attempt in {1..10}; do
        curl localhost:"$port"/_ping > /dev/null 2>&1 && retval=0 && echo "Port library started" && break
        sleep 1
    done
    export port_layer_pid

    return $retval
}

# kills the port layer
kill_port_layer () {
    # FIXME: need to integrate with ESX_URL so disabling it now
    # https://github.com/vmware/vic/issues/304
    return

    kill "$port_layer_pid" > /dev/null 2>&1

    retval=1
    for attempt in {1..5}; do
        kill -0 "$port_layer_pid" || echo "Port library stopped" && retval=0 && break
        sleep 1
    done

    return $retval
}

# returns the IDs of each FS layer represented in the manifest
# in the order of appearance
get_ids() { # assumes cwd is basedir of image
    # find the id with jq
    cat manifest.json | jq -r ".history[].v1Compatibility|fromjson.id"
}

# usage:
# get_checksum INDEX
# returns the checksum for the INDEXth layer in the manifest in the cwd
get_checksum () { # assume cwd is basedir of image
    jq -r ".fsLayers[$1].blobSum" < manifest.json | cut -d: -f2
}

# usage:
# verify_checksums /path/to/basedir/containing/manifest
# returns an error code if calculated sha256sums of the FS tarballs
# do not match the hash recorded in the manifest
# otherwise terminates and returns nothing
verify_checksums () {
    # cd to the directory with a manifest file
    pushd "$1"

    index=0
    # get list of ids and iterate
    for id in $(get_ids); do
        manifest_checksum=$(get_checksum $index) # find manifest hash
        pushd $id # cd into fs basedir

        # check the manifest checksum against calculated hash from tarball
        [[ $manifest_checksum = $(sha256sum $id.tar | awk '{print $1}') ]] || exit 1

        popd
        index=$((index + 1))
    done

    popd
}

govc () {
    govc_bin=${GOVC:-$(which govc)}

    if [ -z "$govc_bin" ] ; then
        echo "govc: command not found" >&2
        exit 1
    fi

    $govc_bin "$@"
}
