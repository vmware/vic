/*
 Copyright 2017 VMware, Inc. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

export const VM_PROPERTIES_TO_EXTRACT = [
    'name',
    'isVCH',
    'isContainer',
    'summary.runtime.powerState',
    'config.extraConfig'
];
export const CONTAINER_VM_IMAGE_NAME_KEY = 'guestinfo.vice./repo';
export const CONTAINER_VM_PORTMAPPING_KEY = 'guestinfo.vice./networks|bridge/ports~';
export const CONTAINER_PRETTY_NAME_KEY = 'common/name';
export const VCH_VM_CLIENT_IP_KEY = 'guestinfo.vice..init.networks|client.assigned.IP';
export const DOCKER_PERSONALITY_ARGS_KEY = 'guestinfo.vice./init/sessions|docker-personality/cmd/Args~';
export const VCH_VM_LOG_PORT = ':2378';
