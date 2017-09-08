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

function generateRandomVchs(length: number) {
    function getRandomIp() {
        return Math.round(Math.random() * 255) +
            '.' + Math.round(Math.random() * 255) +
            '.' + Math.round(Math.random() * 255) +
            '.' + Math.round(Math.random() * 255);
    }

    if (length < 1) {
        length = 1;
    }
    const results = {};
    for (let i = 0; i < length; i++) {
        results[`vm-${100 + i}`] = {
            clientIp: getRandomIp(),
            isUsingTls: Math.random() > 0.5,
            name: `VCH-VM-${i}`,
            powerState: 'POWERED_ON',
            overallStatus: 'GREEN',
            overallCpuUsage: Math.round(Math.random() * 100),
            guestMemoryUsage: Math.round(Math.random() * 2048),
            committedStorage: 123456789,
            id: `85421094-c58e-40f9-a42c-b624160d05f5/vm-${200 + i}`,
            type: 'vic:VirtualContainerHostVm',
            resourcePool: {
                value: `resgroup-v${i}`,
                type: Math.random() > 0.5 ? 'VirtualApp' : 'ResourcePool'
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

export const getVchResponseStub = () => {
    const resultsObj = {
        bla: {
            'clientIp': '10.17.109.80',
            'isUsingTls': true,
            'name': 'vice-9703',
            'powerState': 'POWERED_ON',
            'overallStatus': 'GREEN',
            'overallCpuUsage': 19,
            'guestMemoryUsage': 61,
            'committedStorage': 2263086006,
            'id': '85421094-c58e-40f9-a42c-b624160d05f5/vm-27',
            'type': 'vic:VirtualContainerHostVm',
            'resourcePool': {
                'value': `resgroup-v26`,
                'type': 'VirtualApp'
            },
            'resourceConfig': {
                'entity': {
                    'value': `vm-27`,
                    'type': 'VirtualMachine'
                },
                'changeVersion': null,
                'lastModified': null,
                'cpuAllocation': {
                    'reservation': 1,
                    'expandableReservation': true,
                    'limit': -1,
                    'shares': {
                        'shares': 4000,
                        'level': 'normal'
                    },
                },
                'memoryAllocation': {
                    'reservation': 1,
                    'expandableReservation': true,
                    'limit': -1,
                    'shares': {
                        'shares': 163840,
                        'level': 'normal'
                    },
                    'overheadLimit': 0
                }
            }
        }
    };
    const randomVchsObj = generateRandomVchs(30);
    for (const key in randomVchsObj) {
        if (randomVchsObj[key]) {
            resultsObj[key] = randomVchsObj[key];
        }
    }

    return {
        id: 'something',
        match: 1,
        results: resultsObj
    };
};

export const getMalformedVchResponseStub = () => {
    return {
        id: 'something-wrong',
        match: 2,
        results: {
            bla: {}
        }
    };
};
