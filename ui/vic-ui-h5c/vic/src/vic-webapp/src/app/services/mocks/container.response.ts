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

function generateRandomContainers(length: number) {
    if (length < 1) {
        length = 1;
    }
    const results = {};
    for (let i = 0; i < length; i++) {
        results[`vm-${100 + i}`] = {
            name: `Container-VM-${i}-blablablabla`,
            containerName: `Container-VM-${i}`,
            portMapping: 'tcp/8081:80',
            imageName: 'busybox',
            powerState: 'POWERED_ON',
            overallStatus: 'GREEN',
            overallCpuUsage: Math.round(Math.random() * 100),
            guestMemoryUsage: Math.round(Math.random() * 2048),
            commitedStorage: 123456789,
            id: `85421094-c58e-40f9-a42c-b624160d05f5/vm-${200 + i}`,
            type: 'vic:ContainerVm',
            resourcePool: {
                value: `resgroup-v${100 + i}`,
                type: 'VirtualApp'
            },
            resourceConfig: {
                entity: {
                    value: `vm-${100 + i}`,
                    type: 'VirtualMachine'
                },
                changeVersion: null,
                lastModified: null,
                cpuAllocation: {
                    reservation: 0,
                    expandableReservation: false,
                    limit: -1,
                    shares: { shares: 1000, level: 'NORMAL' },
                    overheadLimit: null
                },
                memoryAllocation: {
                    reservation: 0,
                    expandableReservation: false,
                    limit: -1,
                    shares: { shares: 20480, level: 'NORMAL' },
                    overheadLimit: 57
                }
            }
        };
    }
    return results;
}

export const getContainerResponseStub = () => {
    const randomContainersObj = generateRandomContainers(30);

    return {
        id: 'something',
        match: 30,
        results: randomContainersObj
    };
};

export const getMalformedContainerResponseStub = () => {
    return {
        id: 'something-wrong',
        match: 0,
        results: {}
    };
};
