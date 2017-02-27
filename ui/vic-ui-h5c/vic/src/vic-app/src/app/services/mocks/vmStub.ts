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

export const getVmStubData: Function = (type: string): any => {
    if (type === 'vch') {
        return {
            id: '1',
            name: 'fake vch vm',
            isVCH: true,
            isContainer: false,
            dockerEndpoint: '1.2.3.4:2376',
            dockerLog: 'https://1.2.3.4:2378',
            powerState: 'poweredOff'
        };
    } else {
        return {
            id: '2',
            name: 'fake container vm',
            container_name: 'container name',
            isVCH: false,
            isContainer: true,
            image_name: 'nginx',
            portmapping: '8081:80/tcp',
            powerState: 'poweredOn'
        };
    }
};
