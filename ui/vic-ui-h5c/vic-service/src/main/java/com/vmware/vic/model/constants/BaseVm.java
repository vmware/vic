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

public class BaseVm {
    public static final String ID = "id";
    public static final String VM_NAME = "name";
    public static final String VM_OVERALL_STATUS = "overallStatus";
    public static class Runtime {
        public static final String VM_POWERSTATE_FULLPATH = "runtime.powerState";
        public static final String VM_POWERSTATE_BASENAME = "powerState";
    }
    public static final String VM_SUMMARY = "summary";
    public static class Config {
        public static final String VM_GUESTFULLNAME = "config.guestFullName";
        public static final String VM_EXTRACONFIG = "config.extraConfig";
    }
    public static final String VM_RESOURCECONFIG = "resourceConfig";
    public static final String VM_RESOURCEPOOL = "resourcePool";
    public static final String VM_OVERALLCPUUSAGE = "overallCpuUsage";
    public static final String VM_GUESTMEMORYUSAGE = "guestMemoryUsage";
    public static final String VM_COMMITTEDSTORAGE = "committedStorage";
}
