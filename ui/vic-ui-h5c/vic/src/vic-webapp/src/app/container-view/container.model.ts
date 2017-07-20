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

import {
    VirtualMachine,
    ContainerVmResponse
} from '../vm.interface';

const FORWARD_SLASH = '/';
const COLON = ':';

export class ContainerVm implements VirtualMachine {
    public parentObj: {
        id: string;
        value: string;
        type: string;
    };
    public vmId: string;
    public readonly isVCH: boolean = false;
    public readonly isContainer: boolean = true;
    public name: string;
    public overallStatus: string;
    public powerState: string;
    public containerName: string;
    public parentObjectName: string;
    public imageName: string;
    public portMapping: string;
    public overallCpuUsage: number;
    public guestMemoryUsage: number;
    public committedStorage: number;
    public resourceConfig: any;

    constructor(data: ContainerVmResponse) {
        try {
            // populate vm information
            const splitVmId = data.id.split(FORWARD_SLASH);
            this.parentObj = {
                id: `urn:vmomi:${data.resourcePool.type}:${data.resourcePool.value}${COLON}${splitVmId[0]}`,
                type: data.resourcePool.type,
                value: data.resourcePool.value
            };
            this.parentObjectName = data.parentObjectName;
            this.vmId = `urn:vmomi:VirtualMachine:${splitVmId[1]}${COLON}${splitVmId[0]}`;
            this.name = data.name;
            this.overallStatus = data.overallStatus;
            this.powerState = data.powerState;
            this.containerName = data.containerName;
            this.imageName = data.imageName;
            this.portMapping = data.portMapping;
            this.overallCpuUsage = data.overallCpuUsage;
            this.guestMemoryUsage = data.guestMemoryUsage;
            this.committedStorage = data.committedStorage;
            this.resourceConfig = data.resourceConfig;
        } catch (e) {
            throw new Error('response does not fit into the required type! ' +
                e.toString());
        }
    }

    get parentType() {
        return this.parentObj.type;
    }

    get id(): string {
        return this.vmId;
    }
}
