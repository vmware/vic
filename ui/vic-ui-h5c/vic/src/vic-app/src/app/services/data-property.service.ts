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

import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs/Rx';
import { Http, Response } from '@angular/http';

import { GlobalsService } from '../shared/globals.service';
import { APP_CONFIG } from '../shared/app-config';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

import {
    CONTAINER_VM_IMAGE_NAME_KEY,
    CONTAINER_VM_PORTMAPPING_KEY,
    CONTAINER_PRETTY_NAME_KEY,
    VCH_VM_CLIENT_IP_KEY,
    DOCKER_PERSONALITY_ARGS_KEY,
    VCH_VM_LOG_PORT
} from '../vm.constants';
import { VirtualMachine } from '../vm.interface';
import { getVmStubData } from './mocks/vmStub';

@Injectable()
export class DataPropertyService {
    private _objectId: string;
    private vmInfoSource: Subject<VirtualMachine> = new Subject<VirtualMachine>();
    private vicObjectSource: Subject<any> = new Subject<any>();
    public vmInfo$: Observable<VirtualMachine> = this.vmInfoSource.asObservable();
    public vicObject$: Observable<any> = this.vicObjectSource.asObservable();

    constructor(
        private http: Http,
        private globalsService: GlobalsService
    ) { }

    setObjectId(id: string) {
        this._objectId = id;
    }

    /**
     * Builds data URL for vSphere Client's REST API
     * @param   id
     * @param   props : properties to extract
     * @return  data URL
     */
    buildDataUrl(id: string = this._objectId, props: string[]): string {
        let namespace = window[APP_CONFIG.bundleName];
        if (namespace) {
            return namespace.buildDataUrl(id, props);
        }
        return null;
    }

    /**
     * Calls the vSphere Client's API endpoint to retrieve VM information,
     * and emits the results to Observable
     * @param props
     * @param stubVmType? : stub type (vch or container)
     */
    fetchVmInfo(props: string[], stubVmType?: string): void {
        if (!this.globalsService.isPluginMode()) {
            this.vmInfoSource.next(<VirtualMachine>getVmStubData(stubVmType));
        }

        this.http.get(this.buildDataUrl(this._objectId, props))
            .map(res => {
                let parsed = res.json();
                return parsed;
            })
            .map(processPowerState)
            .map(processVmType)
            .catch((err: Response | any) => {
                return Observable.throw(err);
            })
            .subscribe(
            res => {
                this.vmInfoSource.next(<VirtualMachine>res);
            }, err => {
                this.vmInfoSource.error(err);
            }
            );
    }

    /**
     * Calls the vSphere Client's API endpoint to retrieve VIC Root information,
     * and emits the results to Observable
     * @param props
     */
    fetchRootInfo(props: string[]): void {
        if (!this.globalsService.isPluginMode()) {
            this.vicObjectSource.next({
                uiVersion: '3.14159265',
                vchVmsLen: 1000,
                containerVmsLen: 50000
            });
            return;
        }
        this.http.get(
            this.buildDataUrl(
                'urn:vic:vic:Root:vic%252Fvic-root',
                props)
        )
            .map(res => res.json())
            .catch((err: Response | any) => {
                return Observable.throw(err);
            })
            .subscribe(
            res => {
                this.vicObjectSource.next(res);
            }, err => {
                this.vicObjectSource.next(err);
            }
            );
    }
}

/**
 * Process raw response from vSphere Client data service
 * to process properties specific to VCH or Container VM
 * @param obj : raw json object
 */
function processVmType(obj: any): any {
    if (!obj) {
        return {};
    }

    let extConfig: any[] = obj['config.extraConfig'];
    delete obj['config.extraConfig'];

    // initial values for dockerLog and dockerEndpoint
    if (obj.isVCH) {
        obj.dockerEndpoint = obj.dockerLog = '-';
        if (obj.powerState === 'poweredOff') {
            return obj;
        }

        let isUsingTls = true;

        for (let { key, value } of extConfig) {
            if (key === VCH_VM_CLIENT_IP_KEY) {
                const base64Decoded: string = atob(value);
                const decIpLength = base64Decoded.length;
                // if the ip is in ipv6 format, the decoded string is
                // 16 bytes long
                const decIpIdx = decIpLength === 16 ? decIpLength - 4 : 0;
                const ipv4: string = base64Decoded.charCodeAt(decIpIdx) + '.'
                    + base64Decoded.charCodeAt(decIpIdx + 1) + '.'
                    + base64Decoded.charCodeAt(decIpIdx + 2) + '.'
                    + base64Decoded.charCodeAt(decIpIdx + 3);
                obj.dockerEndpoint = `DOCKER_HOST=tcp://${ipv4}`;
                obj.dockerLog = `https://${ipv4}${VCH_VM_LOG_PORT}`;
                continue;
            }

            if (key === DOCKER_PERSONALITY_ARGS_KEY) {
                isUsingTls = value.indexOf('2376') > -1;
                continue;
            }
        }

        // since the order in which list items are processed is not much guaranteed,
        // we set the port for Docker API endpoint at the end of the loop
        obj.dockerEndpoint += isUsingTls ? ':2376' : ':2375';

    } else {
        for (let { key, value } of extConfig) {
            if (key === CONTAINER_VM_IMAGE_NAME_KEY) {
                obj.image_name = value;
                continue;
            }

            if (key === CONTAINER_PRETTY_NAME_KEY) {
                obj.container_name = value;
                continue;
            }

            if (key === CONTAINER_VM_PORTMAPPING_KEY) {
                obj.portmapping = value;
                continue;
            }
        }
    }

    return obj;
}

/**
 * Process raw response from vSphere Client data service
 * to get powerState
 * @param obj : raw json object
 */
function processPowerState(obj: any): any {
    obj.powerState = obj['summary.runtime.powerState'];
    delete obj['summary.runtime.powerState'];
    return obj;
}
