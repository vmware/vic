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

data_mount=/data/harbor
cfg=${data_mount}/harbor.cfg

function harborDataSanityCheck {
  harbor_dirs=( 
    cert
    database
    job_logs
    registry
  )

  for harbor_dir in "${harbor_dirs[@]}"
  do
    if [ ! -d "$1"/"$harbor_dir" ]; then
      echo "Harbor directory $1/${harbor_dir} not found"
      return 1
    fi
  done

}

function cleanupProcedure {
  # Umount old disks
  umount $system_tmp_mount
  umount $data_tmp_mount
}

#Configure attr in harbor.cfg
function configureHarborCfg {
  cfg_key=$1
  cfg_value=$2

  basedir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

  if [ -n "$cfg_key" ]; then
    cfg_value=$(echo $cfg_value | sed -r -e 's%[\/&%]%\\&%g')
    sed -i -r "s%#?$cfg_key\s*=\s*.*%$cfg_key = $cfg_value%" $cfg
  fi
}

function upgradeHarborConfiguration {
  local source_configuration=$1
  XIFS=$IFS
  IFS=$'\r\n'
  for config in $(grep = ${source_configuration} | sed 's/ *= */=/g' | awk -F= '{ printf "%s=\"%s\"\n", $1, $2 }'); do
    local key=$(echo ${config} | cut -d= -f1)
    local value=$(echo ${config} | cut -d= -f2 | sed -e 's/^"//' -e 's/"$//')
    if [[ x$value == "x" ]]; then 
      continue
    fi
    configureHarborCfg $key $value
  done
  IFS=$XIFS
}

# Before trying anything on the upgrade side, let's check if there is any data
# from harbor already in the new data folder.
if harborDataSanityCheck $data_mount; then 
  echo "Harbor Data is already present in ${data_mount}, this script might corrupt/delete/overwrite your data."
  read -r -p "Do you want to continue? [y/N] " response
  if [[ ! "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
  then
      exit 1
  else
    read -r -p "This script will overwrite your data, are you sure you want to continue? [y/N] " response
    if [[ ! "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    then
        exit 1
    fi
  fi
fi

echo "Performing sanity check..."

old_data_disk=$(pvs ---noheadings | grep data1_vg | gawk '{ print $1 }')

# If old harbor disk is not available, exit.
if [[ x$old_data_disk == "x" ]]; then
  echo "Old Harbor Data disk not available... exiting..." && exit 0
fi

data_tmp_mount=$(mktemp -d)

mount /dev/data1_vg/data $data_tmp_mount

# Perform sanity check on data volume
if ! harborDataSanityCheck $data_tmp_mount; then
  echo "Harbor Data is not present in ${data_tmp_mount}, can't continue with upgrade operation" 
  cleanupProcedure
  exit 0
fi

# Bash black magic to extract the third partition on an unmounted volume that
# is formatted ext3, which should be our old system disk coming from the
# to-be-upgraded harbor instance.
old_system_disk=/dev/$(lsblk -f --noheadings --raw -o NAME,FSTYPE,MOUNTPOINT | awk '$1~/s.*3/ && $2~/ext3/ && $3==""' | awk '{ print $1 }')

# If old harbor disk is not available, exit.
if [[ x$old_system_disk == "x/dev/" ]]; then
  echo "Old Harbor System disk not available... exiting..." && exit 0
fi

system_tmp_mount=$(mktemp -d)

mount $old_system_disk $system_tmp_mount

if [ ! -d "$system_tmp_mount"/harbor ]; then
  echo "Harbor system directory not found, exiting..."
  cleanupProcedure
  exit 0
fi

# Start migration
echo "Starting migration"

echo "[=] shutting down harbor"
systemctl stop harbor_startup.service
systemctl stop harbor.service

echo "[=] copying data"
rsync -av --info=progress $data_tmp_mount/ $data_mount/

echo "[=] migrating harbor configuration"
upgradeHarborConfiguration "$system_tmp_mount"/harbor/harbor/harbor.cfg
echo "[=] cleaning up"
cleanupProcedure

echo "Upgrade procedure complete."
echo "You can now shutdown the appliance, detach the old Harbors disk and start it again."
