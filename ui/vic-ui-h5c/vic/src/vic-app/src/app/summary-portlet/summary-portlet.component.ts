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
    AfterViewInit,
    OnDestroy,
    ChangeDetectorRef
} from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import 'rxjs/add/operator/switchMap';

import { DataPropertyService } from '../services/data-property.service';
import { GlobalsService, RefreshService, AppAlertService } from '../shared/index';
import { VirtualMachine } from '../vm.interface';
import { VM_PROPERTIES_TO_EXTRACT } from '../vm.constants';

@Component({
    selector: 'vic-summary-portlet',
    template: `
    <div *ngIf="!gs.isPluginMode()">
        <button (click)="toggleVmType()">Toggle VM Type</button>
    </div>
    <vic-vch-portlet [activeVm]="activeVm" *ngIf="activeVm && activeVm.isVCH"></vic-vch-portlet>
    <vic-container-portlet [activeVm]="activeVm" *ngIf="activeVm && activeVm.isContainer"></vic-container-portlet>
    `
})

export class VicSummaryPortletComponent implements
    AfterViewInit, OnDestroy {
    public activeVm: VirtualMachine;
    private refreshSubscription: Subscription;
    private vmInfoSubscription: Subscription;
    private stubType: string = 'vch';

    constructor(
        public gs: GlobalsService,
        private service: DataPropertyService,
        private refreshService: RefreshService,
        private appAlertService: AppAlertService,
        private route: ActivatedRoute,
        private router: Router,
        private cd: ChangeDetectorRef
    ) {
        this.refreshSubscription = refreshService.refreshObservable$.subscribe(
            () => {
                console.log('vm is being refreshed');
                this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
            }
        );
        this.vmInfoSubscription = this.service.vmInfo$.subscribe(
            (results: VirtualMachine) => {
                this.activeVm = results;
                this.cd.detectChanges();
                console.log('data fetched successfully', results);
            },
            (err) => {
                this.appAlertService.showError(err);
                console.error('data fetch failed!', err);
            }
        );
    }

    ngAfterViewInit() {
        setTimeout(() => {
            this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT, this.stubType);
        });
    }

    ngOnDestroy() {
        this.refreshSubscription.unsubscribe();
        this.vmInfoSubscription.unsubscribe();
    }

    toggleVmType() {
        this.stubType = this.stubType === 'vch' ? 'container' : 'vch';
        this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT, this.stubType);
    }
}
