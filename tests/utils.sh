#!/bin/bash
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

function filter_url_array {
    for url in $@; do
     if ! grep -q 'VCH-[0-9]\{4\}-[0-9]\{4\}' <<< $(GOVC_URL=$url GOVC_USERNAME=${TEST_USERNAME} GOVC_PASSWORD=${TEST_PASSWORD} GOVC_INSECURE=1 govc ls /ha-datacenter/vm/); then
         echo $url
     fi;
    done | paste -sd " " - | tr -d '\n'
}

function count_url_in_array {
     echo $#
}

# Note: this function directly works on TEST_URL_ARRAY and reset it before return
function wait_for_idle_server {
    filtered_urls=$(filter_url_array $TEST_URL_ARRAY)
    number_of_current_available_servers=$(count_url_in_array ${filtered_urls})

    while [[ ${number_of_current_available_servers} -le 0 ]]; do
     echo "Waiting 5 minutes for idle build server";
     sleep 300;
     filtered_urls=$(filter_url_array $TEST_URL_ARRAY)
     number_of_current_available_servers=$(count_url_in_array ${filtered_urls})
    done

    TEST_URL_ARRAY=${filtered_urls}
}
