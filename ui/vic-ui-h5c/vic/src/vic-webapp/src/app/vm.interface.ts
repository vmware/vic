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

export interface VirtualMachine {
    id: string;
    name: string;
    powerState: string;
    isVCH: boolean;
    isContainer: boolean;
}

export interface VchVmResponse {
    id: string;
    type: string;
    name: string;
    clientIp: string;
    isUsingTls: boolean;
    powerState: string;
    overallStatus: string;
    overallCpuUsage: number;
    guestMemoryUsage: number;
    committedStorage: number;
    resourcePool: {
        value: string;
        type: string;
    };
    resourceConfig: {
        entity: {
            value: string;
            type: string;
        };
        changeVersion: any | null;
        lastModified: string | null;
        cpuAllocation: {
            reservation: number;
            expandableReservation: boolean;
            limit: number;
            shares: {
                shares: number;
                level: string;
            };
            overheadLimit: number | null;
        };
        memoryAllocation: {
            reservation: number;
            expandableReservation: boolean;
            limit: number;
            shares: {
                shares: number;
                level: string;
            };
            overheadLimit: number | null;
        }
    };
}

export interface ContainerVmResponse {
    id: string;
    type: string;
    containerName: string;
    parentObjectName: string;
    imageName: string;
    portMapping?: string | null;
    name: string;
    powerState: string;
    overallStatus: string;
    overallCpuUsage: number;
    guestMemoryUsage: number;
    committedStorage: number;
    resourcePool: {
        value: string;
        type: string;
    };
    resourceConfig: {
        entity: {
            value: string;
            type: string;
        };
        changeVersion: any | null;
        lastModified: string | null;
        cpuAllocation: {
            reservation: number;
            expandableReservation: boolean;
            limit: number;
            shares: {
                shares: number;
                level: string;
            };
            overheadLimit: number | null;
        };
        memoryAllocation: {
            reservation: number;
            expandableReservation: boolean;
            limit: number;
            shares: {
                shares: number;
                level: string;
            };
            overheadLimit: number | null;
        }
    };
}
