# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Keywords ***
Trap Signal Command
    # Container command runs an infinite loop, trapping and logging the given signal name
    [Arguments]  ${sig}
    [Return]  ${busybox} sh -c "trap 'echo KillSignal${sig}' ${sig}; echo READY; while true; do date && sleep 1; done"

Remove All Containers
    Kill All Containers
    Run  docker ps -a | grep Exit | cut -d ' ' -f 1 | xargs docker rm
    Run  docker ps -a | grep Stopped | cut -d ' ' -f 1 | xargs docker rm

Stop All Containers
    Run  docker ps | grep 'Up ' | cut -d ' ' -f 1 | xargs docker stop

Kill All Containers
    Run  docker ps | grep 'Up ' | cut -d ' ' -f 1 | xargs docker kill
