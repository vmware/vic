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
import { Observable, BehaviorSubject } from 'rxjs/Rx';
import { Http, Response } from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

import {
    VirtualMachine,
    CONTAINER_VM_IMAGE_NAME_KEY,
    CONTAINER_VM_PORTMAPPING_KEY,
    CONTAINER_PRETTY_NAME_KEY,
    VCH_VM_CLIENT_IP_KEY,
    VCH_VM_ENDPOINT_PORT,
    VCH_VM_LOG_PORT
} from '../vm.interface';

@Injectable()
export class DataPropertyService {
    private _results: BehaviorSubject<VirtualMachine> = new BehaviorSubject<VirtualMachine>(null);

    constructor(private http: Http) {}

    buildDataUrl(props: string[]): string {
        let url: string = com_vmware_vic.buildDataUrl(WEB_PLATFORM.getObjectId(), props);
        return url;
    }

    getCurrent(): VirtualMachine {
        return this._results.getValue();
    }

    fetch(props: string[]): Observable<VirtualMachine> {
        this.http.get(this.buildDataUrl(props))
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
                    this._results.next(<VirtualMachine>res);
                }, err => {
                    this._results.error(err);
                }
            );

        return this._results.asObservable();
    }
}

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

function processPowerState(obj: any): any {
    obj.powerState = obj['summary.runtime.powerState'];
    delete obj['summary.runtime.powerState'];
    return obj;
}
