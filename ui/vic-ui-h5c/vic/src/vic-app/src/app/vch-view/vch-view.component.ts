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
    OnDestroy
} from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { State } from 'clarity-angular';
import { VirtualContainerHost } from './vch.model';
import { GlobalsService } from '../shared/globals.service';
import { VicVmViewService } from '../services/vm-view.service';
import { RefreshService } from '../shared/refresh.service';

@Component({
    selector: 'vic-vch-view',
    styleUrls: [],
    templateUrl: './vch-view.template.html'
})
export class VicVchViewComponent implements OnInit, OnDestroy {
    private refreshSubscription: Subscription;
    public isDgLoading: boolean = true;
    public vchs: VirtualContainerHost[] = [];
    public totalVchsCount: number = 0;
    public currentOffset: number = 0;
    public currentSort: string = 'id,asc';
    public currentFilter: string = '';
    public readonly maxResultCount: number = 10;

    constructor(
        private vmViewService: VicVmViewService,
        private refreshService: RefreshService,
        private gs: GlobalsService
    ) {
        // subscribes to the global refresh event
        this.refreshSubscription = this.refreshService
            .refreshObservable$.subscribe(() => {
                this.reloadVchs();
            });
    }

    ngOnInit() {
        // sets up a listener for updating this.vchs
        this.vmViewService.vchs$.subscribe(vchs => {
            this.vchs = vchs;
            this.isDgLoading = false;
        }, err => {
            this.vchs = [];
            console.error(err);
        });
    }

    ngOnDestroy() {
        this.refreshSubscription.unsubscribe();
    }

    /**
     * Builds and returns Docker API endpoint string
     * @param item : VirtualContainerHost instance
     * @return DOCKER_HOST environment variable
     */
    getDockerEndpointString(item: VirtualContainerHost): string {
        return `DOCKER_HOST=${item.vchIp}:${item.isUsingTls ? '2376' : '2375'}`;
    }

    /**
     * Queries vic-service with the current Datagrid state
     * @param state current Datagrid state
     */
    refreshGrid(state: State) {
        this.currentFilter = state.filters ? state.filters
            .map(item => item['property'] + '=' + item['value'])
            .join(',') : '';

        if (state.sort) {
            this.currentSort = `${state.sort.by},${state.sort.reverse ? 'desc' : 'asc'}`;
        }

        this.currentOffset = state.page.from;
        this.reloadVchs();
    }

    /**
     * Calls vm-view service to reload VCHs
     */
    reloadVchs() {
        this.isDgLoading = true;
        this.vmViewService.reloadVchs({
            offset: this.currentOffset,
            maxResultCount: this.maxResultCount,
            sorting: this.currentSort,
            filter: this.currentFilter
        });
    }

    /**
     * Navigates to an object specified by objectId
     * @param objectId Full vSphere objectId which starts with urn:
     */
    navigateToObject(objectId: string) {
        if (objectId.indexOf('VirtualMachine') > -1) {
            this.gs.getWebPlatform().sendNavigationRequest('vsphere.core.vm.summary', objectId);
        } else {
            window.parent.location.href = '/ui/#?extensionId=' +
                'vsphere.core.inventory.serverObjectViewsExtension&' +
                'objectId=' + objectId + '&' +
                'navigator=vsphere.core.viTree.hostsAndClustersView';
        }
    }
}
