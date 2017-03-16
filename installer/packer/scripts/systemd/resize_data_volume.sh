#!/usr/bin/bash
# Copyright 2017 VMware, Inc. All Rights Reserved.
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
set -euf -o pipefail

# Our data block device always sits at ID 1 on the first controller
block_device=/dev/disk/by-path/pci-0000:00:10.0-scsi-0:0:1:0
# Our data partition is always partition one on said block device
data_partition=${block_device}-part1

# Size threshold upon which we trigger automatic expansion of the partition,
# currently set at 10MiB
device_size_threshold=10485760

# Size threshold upon which we trigger automatic expansion of the filesystem,
# currently set at 4KiB
fs_size_threshold=4096

function check_integer {
  [[ $1 =~ ^-?[0-9]+$ ]] || ( echo "check did not return an integer, failing"; exit 1 )
}

function repartition {
  blkdev_size=`blockdev --getsize64 ${block_device}`
  check_integer $blkdev_size

  partition_size=`blockdev --getsize64 ${data_partition}`
  check_integer $partition_size

  device_size_difference=`expr ${blkdev_size} - ${partition_size}`
  check_integer $device_size_difference

  if [ $device_size_difference -gt $device_size_threshold ]; then
    # Resize Partition to use 100% of the block device
    echo "Repartitioning ${block_device}"
    parted -a minimal --script ${block_device} 'resizepart 1 100%'
    # Reload partition table of the device
    echo "Reloading partition table"
    partprobe ${block_device}
    sleep 3
  else
    echo "No repartition performed, size threshold not met"
  fi
}

function resize {
  fs_size=`dumpe2fs -h ${data_partition} |& gawk -F: '/Block count/{count=$2} /Block size/{size=$2} END{print count*size}'`
  check_integer $fs_size

  partition_size=`blockdev --getsize64 ${data_partition}`
  check_integer $partition_size

  fs_size_difference=`expr ${partition_size} - ${fs_size}`
  check_integer $fs_size_difference

  if [ $fs_size_difference -gt $fs_size_threshold ]; then
    # Force a filesystem check on the data partition
    echo "Force filesystem check on ${data_partition}"
    e2fsck -pf ${data_partition}
    # Resize the filesystem
    echo "Resize filesystem on ${data_partition}"
    resize2fs ${data_partition}
  else
    echo "No resize performed, size threshold not met"
  fi
}

function usage {
  echo $"Usage: $0 {repartition|resize}"
  exit 1
}

if [ $# -gt 0 ]; then
  case "$1" in
    repartition)
      repartition
      ;;
    resize)
      resize
      ;;         
    *)
      usage
  esac
else
  usage
fi