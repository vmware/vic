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
import { State, Comparator } from 'clarity-angular';
import { ContainerVm } from './container.model';
import {
    GlobalsService,
    RefreshService,
    Vic18nService
} from '../shared';
import { VicVmViewService } from '../services/vm-view.service';
import {
    VM_COMMITTEDSTORAGE,
    VM_GUESTMEMORYUSAGE,
    VM_OVERALLCPUUSAGE,
    VSPHERE_VM_SUMMARY_KEY,
    WS_CONTAINER
} from '../shared/constants';

class GuestMemoryUsageComparator implements Comparator<any> {
    compare(a: any, b: any) {
        return a[VM_GUESTMEMORYUSAGE] - b[VM_GUESTMEMORYUSAGE];
    }

    toString(): string {
        return VM_GUESTMEMORYUSAGE;
    }
}

class OverallCpuUsageComparator implements Comparator<any> {
    compare(a: any, b: any) {
        return a[VM_OVERALLCPUUSAGE] - b[VM_OVERALLCPUUSAGE];
    }

    toString(): string {
        return VM_OVERALLCPUUSAGE;
    }
}

class CommittedStorageComparator implements Comparator<any> {
    compare(a: any, b: any) {
        return a[VM_COMMITTEDSTORAGE] - b[VM_COMMITTEDSTORAGE];
    }

    toString(): string {
        return VM_COMMITTEDSTORAGE;
    }
}

@Component({
    selector: 'vic-container-view',
    styleUrls: [],
    templateUrl: './container-view.template.html'
})
export class VicContainerViewComponent implements OnInit, OnDestroy {
    public readonly WS_CONTAINER_CONSTANTS = WS_CONTAINER;
    private refreshSubscription: Subscription;
    public guestMemoryUsageComparator = new GuestMemoryUsageComparator();
    public overallCpuUsageComparator = new OverallCpuUsageComparator();
    public committedStorageComparator = new CommittedStorageComparator();
    public isDgLoading: boolean = true;
    public containers: ContainerVm[] = [];
    public totalContainersCount: number = 0;
    public currentOffset: number = 0;
    public currentSort: string = 'id,asc';
    public currentFilter: string = '';
    public readonly maxResultCount: number = 10;
    public readonly MEGABYTE: number = Math.pow(1024, 2);
    public readonly GIGABYTE: number = Math.pow(1024, 3);

    constructor(
        private vmViewService: VicVmViewService,
        private refreshService: RefreshService,
        private gs: GlobalsService,
        public vicI18n: Vic18nService
    ) {
        // subscribes to the global refresh event
        this.refreshSubscription = this.refreshService
            .refreshObservable$.subscribe(() => {
                this.reloadContainers();
            });
    }

    ngOnInit() {
        // sets up a listener for updating this.containers
        this.vmViewService.containers$.subscribe(vms => {
            this.containers = vms;
            this.isDgLoading = false;
        }, err => {
            this.containers = [];
            console.error(err);
        });
    }

    ngOnDestroy() {
        this.refreshSubscription.unsubscribe();
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
            let sortBy = typeof state.sort.by === 'object' ?
                state.sort.by.toString() : state.sort.by;

            this.currentSort = `${sortBy},${state.sort.reverse ? 'desc' : 'asc'}`;
        }

        this.currentOffset = state.page.from;
        this.reloadContainers();
    }

    /**
     * Calls vm-view service to reload Containers
     */
    reloadContainers() {
        this.isDgLoading = true;
        this.vmViewService.reloadContainers({
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
        this.gs.getWebPlatform()
            .sendNavigationRequest(VSPHERE_VM_SUMMARY_KEY, objectId);
    }

    /**
     * Turns size in byte into a more legible unit
     * @param size 
     * @return pretty-formatted size
     */
    formatStorage(size: number): string {
        let results: string;
        if (size < this.MEGABYTE) {
            results = Math.round(size / 1024) + ' KB';
        } else if (size < this.GIGABYTE) {
            results = Math.round(size / this.MEGABYTE) + ' MB';
        } else {
            results = Math.round(size / this.GIGABYTE) + ' GB';
        }
        return results;
    }
}
