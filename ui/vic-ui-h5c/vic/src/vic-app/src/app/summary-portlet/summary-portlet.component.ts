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

import { Component, AfterViewInit, ChangeDetectorRef } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import 'rxjs/add/operator/switchMap';

import { DataPropertyService } from '../data-property-service/data-property.service';
import { VirtualMachine, VM_PROPERTIES_TO_EXTRACT } from '../vm.interface';

@Component({
    selector: 'vic-summary-portlet',
    styleUrls: [],
    template: `
    <vic-vch-portlet [activeVm]="activeVm" *ngIf="activeVm && activeVm.isVCH"></vic-vch-portlet>
    <vic-container-portlet [activeVm]="activeVm" *ngIf="activeVm && activeVm.isContainer"></vic-container-portlet>
    `
})

export class VicSummaryPortletComponent implements AfterViewInit {
    public activeVm: VirtualMachine;

    constructor(
        private service: DataPropertyService,
        private route: ActivatedRoute,
        private router: Router,
        private cd: ChangeDetectorRef
    ) {}

    ngAfterViewInit() {
        setTimeout(() => {
            this.service.fetch(VM_PROPERTIES_TO_EXTRACT).subscribe(
                (results: VirtualMachine) => {
                    this.activeVm = results;
                    this.cd.detectChanges();
                    console.log('data fetched successfully', results);
                },
                (err) => {
                    console.error('data fetch failed!', err);
                }
            );

            if (window.hasOwnProperty('WEB_PLATFORM')) {
                WEB_PLATFORM.setGlobalRefreshHandler(() => {
                    this.service.fetch(VM_PROPERTIES_TO_EXTRACT);
                }, document);
            }
        });
    }
}
