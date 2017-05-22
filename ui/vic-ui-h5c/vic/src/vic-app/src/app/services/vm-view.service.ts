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
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { Http } from '@angular/http';

import { VirtualContainerHost } from '../vch-view/vch.model';
import { GlobalsService } from '../shared/globals.service';
import { getVchResponseStub } from '../services/mocks/vch.response';

const ALL_VCHS_URL = '/ui/vic/rest/data/list/?targetType=vic:VirtualContainerHostVm' +
    '&properties=match,results';

@Injectable()
export class VicVmViewService {
    private vchsSubj: Subject<VirtualContainerHost[]>;
    public vchs$: Observable<VirtualContainerHost[]>;

    constructor(
        private http: Http,
        private gs: GlobalsService
    ) {
        // sets up subject and observable
        this.vchsSubj = new Subject<VirtualContainerHost[]>();
        this.vchs$ = this.vchsSubj.asObservable();
    }

    /**
     * Queries vic-service and gets a list of VCHs
     * @param params : query parameters to query backend with
     */
    reloadVchs(params: {
        offset?: number,
        maxResultCount?: number,
        sorting?: string,
        filter?: string
    }) {
        let vchsArray: VirtualContainerHost[] = [];
        let vchDataResponse;

        // uses mocked data if app is running in standalone (dev) mode
        if (!this.gs.isPluginMode()) {
            vchDataResponse = getVchResponseStub().results;
            vchsArray = [];

            for (let objectId in vchDataResponse) {
                if (vchDataResponse.hasOwnProperty(objectId)) {
                    vchsArray.push(new VirtualContainerHost(vchDataResponse[objectId]));
                }
            }

            this.vchsSubj.next(vchsArray);
            return;
        }

        // makes a GET request to vic-service and emits data fetched
        this.http.get(this.buildQueryString(ALL_VCHS_URL, params))
            .map(response => response.json())
            .subscribe(data => {
                try {
                    for (let objectId in data.data) {
                        if (data.data.hasOwnProperty(objectId)) {
                            vchsArray.push(new VirtualContainerHost(data.data[objectId]));
                        }
                    }
                } catch (e) {
                    this.vchsSubj.error(e);
                    return;
                }

                this.vchsSubj.next(vchsArray);
            });
    }

    /**
     * Builds a query string to query vic-service
     * @param baseUrl : endpoint as defind by vic-service
     * @param params : query parameters to query backend with
     */
    buildQueryString(baseUrl: string, params: {
        offset?: number,
        maxResultCount?: number,
        sorting?: string,
        filter?: string
    }) {
        let queryString = `${baseUrl}&`;
        let paramKeys: string[] = Object.keys(params);
        let paramsLength: number = paramKeys.length;

        paramKeys.forEach((key, idx) => {
            queryString += `${key}=${params[key]}`;
            if (idx < paramsLength - 1) {
                queryString += '&';
            }
        });

        return queryString;
    }
}
