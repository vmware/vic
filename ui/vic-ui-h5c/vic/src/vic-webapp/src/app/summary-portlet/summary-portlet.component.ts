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
    OnDestroy,
    ChangeDetectorRef
} from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
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

export class VicSummaryPortletComponent implements OnDestroy {
    public activeVm: VirtualMachine;
    private refreshSubscription: Subscription;
    private vmInfoSubscription: Subscription;
    private stubType = 'vch';

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
                this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
            }
        );
        this.vmInfoSubscription = this.service.vmInfo$.subscribe(
            (results: VirtualMachine) => {
                this.activeVm = results;
                this.cd.detectChanges();
            },
            (err) => {
                this.appAlertService.showError(err);
                console.error('data fetch failed!', err);
            }
        );

        // set up objectId in data property service
        const paramsIntervalTimer = setInterval(() => {
            if (this.route && this.route.params) {
                this.route.params.subscribe((params: any) => {
                    this.service.setObjectId(params.id);
                    this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT, this.stubType);
                });
                clearInterval(paramsIntervalTimer);
            }
        }, 5);
    }

    ngOnDestroy() {
        if (this.refreshSubscription) {
            this.refreshSubscription.unsubscribe();
        }

        if (this.vmInfoSubscription) {
            this.vmInfoSubscription.unsubscribe();
        }
    }

    toggleVmType() {
        this.stubType = this.stubType === 'vch' ? 'container' : 'vch';
        this.service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT, this.stubType);
    }
}
