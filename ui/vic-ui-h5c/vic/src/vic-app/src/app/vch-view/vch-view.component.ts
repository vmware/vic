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
    Component,
    OnInit,
    OnDestroy,
    NgZone
} from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { State } from 'clarity-angular';
import { VirtualContainerHost } from './vch.model';
import {
    GlobalsService,
    RefreshService,
    Vic18nService
} from '../shared';
import { VicVmViewService } from '../services/vm-view.service';
import {
    VSPHERE_VM_SUMMARY_KEY,
    DOCKER_ENGINE_PORT_NOTLS,
    DOCKER_ENGINE_PORT_TLS,
    VSPHERE_SERVEROBJ_VIEWEXT_KEY,
    VSPHERE_VITREE_HOSTCLUSTERVIEW_KEY,
    WS_VCH
} from '../shared/constants';

@Component({
    selector: 'vic-vch-view',
    styleUrls: [],
    templateUrl: './vch-view.template.html'
})
export class VicVchViewComponent implements OnInit, OnDestroy {
    public readonly WS_VCH_CONSTANTS = WS_VCH;
    private refreshSubscription: Subscription;
    public isDgLoading: boolean = true;
    public vchs: VirtualContainerHost[] = [];
    public totalVchsCount: number = 0;
    public currentState: {
        offset: number;
        sorting: string;
        filter: string;
    } = { offset: 0, sorting: 'id,asc', filter: '' };
    public readonly maxResultCount: number = 10;

    constructor(
        private zone: NgZone,
        private vmViewService: VicVmViewService,
        private refreshService: RefreshService,
        private globalsService: GlobalsService,
        public vicI18n: Vic18nService
    ) { }

    ngOnInit() {
        // subscribes to the global refresh event and calls the
        // reloadVchs() method to query the server for new data
        this.refreshSubscription = this.refreshService
            .refreshObservable$.subscribe(() => {
                this.zone.run(() => {
                    this.reloadVchs();
                });
            });

        // listens to an observable that gets the updated vchs data
        // from the server, and updates this.vchs
        this.vmViewService.vchs$.subscribe(vchs => {
            this.vchs = vchs;
            this.isDgLoading = false;
        }, err => {
            this.vchs = [];
            console.error(err);
        });
    }

    ngOnDestroy() {
        if (this.refreshSubscription) {
            this.refreshSubscription.unsubscribe();
        }
    }

    /**
     * Builds and returns Docker API endpoint string
     * @param item : VirtualContainerHost instance
     * @return DOCKER_HOST environment variable
     */
    getDockerEndpointString(item: VirtualContainerHost): string {
        return `DOCKER_HOST=${item.vchIp}:${item.isUsingTls ?
            DOCKER_ENGINE_PORT_TLS : DOCKER_ENGINE_PORT_NOTLS}`;
    }

    /**
     * Queries vic-service with the current Datagrid state
     * @param state current Datagrid state
     */
    refreshGrid(state: State) {
        this.currentState.filter = state.filters ? state.filters
            .map(item => item['property'] + '=' + item['value'])
            .join(',') : '';

        if (state.sort) {
            this.currentState.sorting = `${state.sort.by},${state.sort.reverse ? 'desc' : 'asc'}`;
        }

        this.currentState.offset = state.page.from;
        this.reloadVchs();
    }

    /**
     * Calls vm-view service to reload VCHs
     */
    reloadVchs() {
        this.isDgLoading = true;
        this.vmViewService.getVchsData({
            offset: this.currentState.offset,
            maxResultCount: this.maxResultCount,
            sorting: this.currentState.sorting,
            filter: this.currentState.filter
        });
    }

    /**
     * Navigates to an object specified by objectId
     * @param objectId Full vSphere objectId which starts with urn:
     */
    navigateToObject(objectId: string) {
        if (objectId.indexOf('VirtualMachine') > -1) {
            this.globalsService.getWebPlatform().sendNavigationRequest(
                VSPHERE_VM_SUMMARY_KEY, objectId);
        } else {
            window.parent.location.href = '/ui/#?extensionId=' +
                VSPHERE_SERVEROBJ_VIEWEXT_KEY + '&' +
                'objectId=' + objectId + '&' +
                'navigator=' + VSPHERE_VITREE_HOSTCLUSTERVIEW_KEY;
        }
    }
}
