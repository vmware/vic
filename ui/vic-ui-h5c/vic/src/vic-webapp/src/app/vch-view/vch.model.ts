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
    VchVmResponse
} from '../vm.interface';

const FORWARD_SLASH = '/';
const COLON = ':';

export class VirtualContainerHost implements VirtualMachine {
    private _parentObj: {
        value: string;
        type: string;
    };
    private _vchVmId: string;
    public readonly isVCH: boolean = true;
    public readonly isContainer: boolean = false;
    public isUsingTls = true;
    public name: string;
    public overallStatus: string;
    public powerState: string;
    public vchIp: string;
    public overallCpuUsage: number;
    public guestMemoryUsage: number;
    public committedStorage: number;
    public resourceConfig: any;

    constructor(data: VchVmResponse) {
        try {
            const splitVmId = data.id.split(FORWARD_SLASH);
            this._parentObj = data.resourcePool;
            this._vchVmId = `urn:vmomi:VirtualMachine:${splitVmId[1]}${COLON}${splitVmId[0]}`;
            this.name = data.name;
            this.overallStatus = data.overallStatus;
            this.powerState = data.powerState;
            this.vchIp = data.clientIp;
            this.isUsingTls = data.isUsingTls;
            this.overallCpuUsage = data.overallCpuUsage;
            this.guestMemoryUsage = data.guestMemoryUsage;
            this.committedStorage = data.committedStorage;
            this.resourceConfig = data.resourceConfig;
        } catch (e) {
            throw new Error('response does not fit into VchVmResponse! ' +
                e.toString());
        }
    }

    get parentType(): string {
        return this._parentObj.type;
    }

    get id(): string {
        return this._vchVmId;
    }

    get vchVmId(): string {
        return this._vchVmId;
    }
}
