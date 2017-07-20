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

// Application level constants can be kept here
export interface AppConfig {
   pluginName: string;
   bundleName: string;
   packageName: string;
}

export const APP_CONFIG: AppConfig = {
   // Names used during plugin generation, can be changed once here
   // and the rest of the generated .ts and .html code will use them.
   pluginName: 'vic',
   bundleName: 'com_vmware_vic',
   packageName: 'com.vmware.vic'
};
