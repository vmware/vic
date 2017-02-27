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
    VCH_VM_ENDPOINT_PORT,
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
        private gs: GlobalsService
    ) {
        // retrieve objectId from the frame's URL
        this._objectId = this.gs.getWebPlatform().getObjectId();
    }

    /**
     * Builds data URL for vSphere Client's REST API
     * @param   id
     * @param   props : properties to extract
     * @return  data URL
     */
    buildDataUrl(id: string = this._objectId, props: string[]): string {
        let url: string = window[APP_CONFIG.bundleName]
            .buildDataUrl(id, props);
        return url;
    }

    /**
     * Calls the vSphere Client's API endpoint to retrieve VM information,
     * and emits the results to Observable
     * @param props
     * @param stubVmType? : stub type (vch or container)
     */
    fetchVmInfo(props: string[], stubVmType?: string): void {
        if (!this.gs.isPluginMode()) {
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
        if (!this.gs.isPluginMode()) {
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
    }

    for (let i = 0; i < extConfig.length; i++) {
        if (extConfig[i].key === VCH_VM_CLIENT_IP_KEY) {
            const base64_decoded: string = atob(extConfig[i].value);
            const ipv4: string = base64_decoded.charCodeAt(0) + '.'
                + base64_decoded.charCodeAt(1) + '.'
                + base64_decoded.charCodeAt(2) + '.'
                + base64_decoded.charCodeAt(3);
            obj.dockerEndpoint = `DOCKER_HOST=tcp://${ipv4}${VCH_VM_ENDPOINT_PORT}`;
            obj.dockerLog = `https://${ipv4}${VCH_VM_LOG_PORT}`;
            break;
        }

        if (extConfig[i].key === CONTAINER_VM_IMAGE_NAME_KEY) {
            obj.image_name = extConfig[i].value;
            continue;
        }

        if (extConfig[i].key === CONTAINER_PRETTY_NAME_KEY) {
            obj.container_name = extConfig[i].value;
            continue;
        }

        if (extConfig[i].key === CONTAINER_VM_PORTMAPPING_KEY) {
            obj.portmapping = extConfig[i].value;
            continue;
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
