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

import { async, TestBed, ComponentFixture } from '@angular/core/testing';
import { HttpModule } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { By } from '@angular/platform-browser';

import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';
import {
    RefreshService,
    GlobalsService,
    Globals,
    AppAlertService,
    Vic18nService,
    I18nService
} from '../shared';
import { VicVmViewService } from '../services/vm-view.service';
import { VicContainerViewComponent } from './container-view.component';
import { ContainerVm } from './container.model';
import { ClarityModule } from 'clarity-angular';
import {
    getContainerResponseStub,
    getMalformedContainerResponseStub
} from '../services/mocks/container.response';
import { WS_CONTAINER } from '../shared/constants';

let responseProperlyFormatted: boolean = true;

class VicVmViewServiceStub {
    private containersSubj: Subject<ContainerVm[]>;
    public containers$: Observable<ContainerVm[]>;
    private data: ContainerVm[] = [];

    constructor() {
        this.containersSubj = new Subject<ContainerVm[]>();
        this.containers$ = this.containersSubj.asObservable();
    }

    getContainersData() {
        // populates data with either correctly or incorrectly formatted data
        // based on the responseProperlyFormatted flag
        this.data = [];
        let cResponse = responseProperlyFormatted ?
            getContainerResponseStub().results :
            getMalformedContainerResponseStub().results;

        for (let objId in cResponse) {
            if (cResponse.hasOwnProperty(objId)) {
                this.data.push(new ContainerVm(cResponse[objId]));
            }
        }

        this.containersSubj.next(this.data);
    }
}

describe('VicContainerViewComponent', () => {
    let fixture: ComponentFixture<VicContainerViewComponent>;
    let vmViewService: VicVmViewService;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                { provide: VicVmViewService, useClass: VicVmViewServiceStub },
                GlobalsService,
                Globals,
                RefreshService,
                AppAlertService,
                I18nService,
                Vic18nService
            ],
            declarations: [
                VicContainerViewComponent
            ],
            imports: [
                ClarityModule.forRoot(),
                HttpModule
            ]
        }).compileComponents();
        fixture = TestBed.createComponent<VicContainerViewComponent>
            (VicContainerViewComponent);
    }));

    it('should have fixture', () => {
        expect(fixture).toBeTruthy();
    });

    it('should render the data grid with properly formatted data', async(() => {
        let compInstance = fixture.componentInstance;
        compInstance.ngOnInit();
        compInstance.reloadContainers();
        fixture.detectChanges();
        let rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
        let rowElementsLength = rowElements.length;
        expect(rowElementsLength).toBe(30);
    }));

    it('should render zero row for malformed data', async(() => {
        responseProperlyFormatted = false;
        let compInstance = fixture.componentInstance;
        compInstance.ngOnInit();
        compInstance.reloadContainers();
        fixture.detectChanges();
        let rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
        let rowElementsLength = rowElements.length;
        expect(rowElementsLength).toBe(0);
    }));

    it('should render default localized text for table headers', async(() => {
        let compInstance = fixture.componentInstance;
        compInstance.ngOnInit();
        compInstance.reloadContainers();
        fixture.detectChanges();

        // containerName column
        let containerEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="containerName"]'));
        expect(containerEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.CONTAINER_NAME
            ]);

        // powerState column
        let powerStateEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="powerState"]'));
        expect(powerStateEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.POWER_STATE
            ]);

        // guestMemoryUsage column
        let guestMemoryUsageEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-sort-by="guestMemoryUsage"]'));
        expect(guestMemoryUsageEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.MEMORY_USAGE
            ]);

        // overallCpuUsage column
        let overallCpuUsageEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-sort-by="overallCpuUsage"]'));
        expect(overallCpuUsageEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.CPU_USAGE
            ]);

        // committedStorage column
        let committedStorageEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-sort-by="committedStorage"]'));
        expect(committedStorageEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.STORAGE_USAGE
            ]);

        // portMapping column
        let portMappingEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="portMapping"]'));
        expect(portMappingEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.PORT_MAPPING
            ]);

        // name column
        let nameEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="name"]'));
        expect(nameEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.VM_NAME
            ]);

        // imageName column
        let imageNameEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="imageName"]'));
        expect(imageNameEl.nativeElement.textContent.trim()).toBe(
            WS_CONTAINER.DG.COLUMNS.defaults[
            WS_CONTAINER.DG.COLUMNS.keys.IMAGE_NAME
            ]);
    }));
});
