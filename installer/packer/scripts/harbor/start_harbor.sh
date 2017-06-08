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

function getRegistryVersion {
FILE="$1" python - <<END
import yaml, os
file = os.environ['FILE']
f = open(file, "r")
dataMap = yaml.safe_load(f)
print dataMap["services"]["registry"]["image"]
f.close()
END
}

function gc {
  echo "======================= $(date)====================="

  registry_image=$(getRegistryVersion /etc/vmware/harbor/docker-compose.yml)

  local data_dir=/data/harbor/registry
  local config_dir=/etc/vmware/harbor/common/config/registry

  /usr/bin/docker run --name gc --rm --volume $data_dir:/storage \
    --volume $config_dir/:/etc/registry/ \
    $registry_image garbage-collect /etc/registry/config.yml
  
  echo "===================================================="
}

gc_enabled=$(ovfenv --key registry.gc_enabled)

if [ ${gc_enabled,,} == "true" ]; then
  gc 2>&1 >> /var/log/harbor/gc.log || true
fi

/usr/local/bin/docker-compose -f /etc/vmware/harbor/docker-compose.yml -f /etc/vmware/harbor/docker-compose.notary.yml up
