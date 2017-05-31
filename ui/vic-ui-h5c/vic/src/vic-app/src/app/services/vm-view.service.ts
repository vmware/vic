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
import { ContainerVm } from '../container-view/container.model';
import { GlobalsService } from '../shared/globals.service';
import { getVchResponseStub } from '../services/mocks/vch.response';
import { getContainerResponseStub } from '../services/mocks/container.response';

const ALL_VCHS_URL = '/ui/vic/rest/data/list/?targetType=' +
    'vic:VirtualContainerHostVm&properties=match,results';
const ALL_CONTAINERS_URL = '/ui/vic/rest/data/list/?targetType=' +
    'vic:ContainerVm&properties=match,results';

@Injectable()
export class VicVmViewService {
    private vchsSubj: Subject<VirtualContainerHost[]>;
    private containersSubj: Subject<ContainerVm[]>;
    public vchs$: Observable<VirtualContainerHost[]>;
    public containers$: Observable<ContainerVm[]>;

    constructor(
        private http: Http,
        private globalsService: GlobalsService
    ) {
        // sets up subjects and observables
        this.vchsSubj = new Subject<VirtualContainerHost[]>();
        this.containersSubj = new Subject<ContainerVm[]>();
        this.vchs$ = this.vchsSubj.asObservable();
        this.containers$ = this.containersSubj.asObservable();
    }

    /**
     * Queries vic-service and gets a list of VCHs
     * @param params : query parameters to query backend with
     */
    getVchsData(params: {
        offset?: number,
        maxResultCount?: number,
        sorting?: string,
        filter?: string
    }) {
        let vchsArray: VirtualContainerHost[] = [];
        let vchDataResponse;

        // uses mocked data if app is running in standalone (dev) mode
        if (!this.globalsService.isPluginMode()) {
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
     * Queries vic-service and gets a list of Container VMs
     * @param params : query parameters to query backend with
     */
    getContainersData(params: {
        offset?: number,
        maxResultCount?: number,
        sorting?: string,
        filter?: string
    }) {
        let containersArray: ContainerVm[] = [];
        let containersDataResponse;

        // uses mocked data if app is running in standalone (dev) mode
        if (!this.globalsService.isPluginMode()) {
            containersDataResponse = getContainerResponseStub().results;
            containersArray = [];

            for (let objectId in containersDataResponse) {
                if (containersDataResponse.hasOwnProperty(objectId)) {
                    containersArray.push(
                        new ContainerVm(containersDataResponse[objectId]));
                }
            }

            this.containersSubj.next(containersArray);
            return;
        }

        // makes a GET request to vic-service and emits data fetched
        this.http.get(this.buildQueryString(ALL_CONTAINERS_URL, params))
            .map(response => response.json())
            .subscribe(data => {
                try {
                    for (let objectId in data.data) {
                        if (data.data.hasOwnProperty(objectId)) {
                            containersArray.push(
                                new ContainerVm(data.data[objectId]));
                        }
                    }
                } catch (e) {
                    this.containersSubj.error(e);
                    return;
                }

                this.containersSubj.next(containersArray);
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
