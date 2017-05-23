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
package com.vmware.vic.model.constants;

public class Vch {
    public static final String VM_IS_USING_TLS = "isUsingTls";
    public static final String VM_CLIENT_IP = "clientIp";
    public static final String VM_VCH_IP = "vchIp";
    public static final String VM_EXTRACONFIG_VCH_KEY = "init/common/name";
    public static final String VM_EXTRACONFIG_CLIENT_IP_KEY =
            "guestinfo.vice..init.networks|client.assigned.IP";
    public static final String VM_EXTRACONFIG_DOCKER_PERSONALITY_ARGS_KEY =
            "guestinfo.vice./init/sessions|docker-personality/cmd/Args~";
}
