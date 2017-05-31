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

import { TestBed, async, ComponentFixture } from '@angular/core/testing';
import { HttpModule } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { By } from '@angular/platform-browser';

import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';
import {
    RefreshService,
    AppAlertService,
    I18nService,
    Vic18nService,
    Globals,
    GlobalsService,
} from '../shared';
import { VicVmViewService } from '../services/vm-view.service';
import { VicVchViewComponent } from './vch-view.component';
import { VirtualContainerHost } from './vch.model';
import { ClarityModule } from 'clarity-angular';
import {
    getVchResponseStub,
    getMalformedVchResponseStub
} from '../services/mocks/vch.response';
import { WS_VCH } from '../shared/constants';

let responseProperlyFormatted: boolean = true;

class VicVmViewServiceStub {
    public totalVchsLength$: Observable<number>;
    private vchsSubj: Subject<VirtualContainerHost[]>;
    public vchs$: Observable<VirtualContainerHost[]>;
    private data: VirtualContainerHost[] = [];

    constructor() {
        this.vchsSubj = new Subject<VirtualContainerHost[]>();
        this.vchs$ = this.vchsSubj.asObservable();
    }

    getVchsData() {
        // populates data with either correctly or incorrectly formatted data
        // based on the responseProperlyFormatted flag
        this.data = [];
        let vchResponse = responseProperlyFormatted ?
            getVchResponseStub().results : getMalformedVchResponseStub().results;
        for (let objId in vchResponse) {
            if (vchResponse.hasOwnProperty(objId)) {
                this.data.push(new VirtualContainerHost(vchResponse[objId]));
            }
        }
        this.vchsSubj.next(this.data);
    }
}

describe('VicVchViewComponent', () => {
    let fixture: ComponentFixture<VicVchViewComponent>;
    let vmViewservice: VicVmViewService;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                { provide: VicVmViewService, useClass: VicVmViewServiceStub },
                GlobalsService,
                Globals,
                RefreshService,
                AppAlertService,
                Vic18nService,
                I18nService
            ],
            declarations: [
                VicVchViewComponent
            ],
            imports: [
                ClarityModule.forRoot(),
                HttpModule
            ]
        }).compileComponents();
        fixture = TestBed.createComponent<VicVchViewComponent>(VicVchViewComponent);
        responseProperlyFormatted = true;
    }));

    it('should have fixture', () => {
        expect(fixture).toBeTruthy();
    });

    it('should render the data grid with properly formatted data', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();
        let rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
        let rowElementsLength = rowElements.length;
        expect(rowElementsLength).toBe(31);
        if (rowElementsLength) {
            let cellElements = rowElements[0].queryAll(By.css('clr-dg-cell'));
            let cellElementsLength = cellElements.length;
            expect(cellElementsLength).toBe(4);
            expect(cellElements[0].nativeElement.textContent).toContain('vice-9703');
            expect(cellElements[1].nativeElement.attributes['ng-reflect-ng-switch'].value).toBe('GREEN');
            expect(cellElements[2].nativeElement.textContent).toContain('10.17.109.80:2376');
            expect(cellElements[3].nativeElement.textContent).toContain('10.17.109.80:2378');
        }
    }));

    it('should render zero row for malformed data', async(() => {
        responseProperlyFormatted = false;
        try {
            fixture.componentInstance.ngOnInit();
            fixture.componentInstance.reloadVchs();
            fixture.detectChanges();
        } catch (e) {
            console.info('expected exception caught');
        } finally {
            let rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
            let rowElementsLength = rowElements.length;
            expect(rowElementsLength).toBe(0);
        }
    }));

    it('should render default localized text for table headers', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();

        // containerName column
        let nameEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="name"]'));
        expect(nameEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[WS_VCH.DG.COLUMNS.keys.NAME]);

        // overallStatus column
        let overallStatusEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="overallStatus"]'));
        expect(overallStatusEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[WS_VCH.DG.COLUMNS.keys.OVERALL_STATUS]);

        // Docker API Endpoint column
        let dockerApiEndpointEl = fixture.debugElement.queryAll(
            By.css('clr-dg-column[ng-reflect-field="vchIp"]'))[0];
        expect(dockerApiEndpointEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[
            WS_VCH.DG.COLUMNS.keys.DOCKER_API_ENDPOINT
            ]);

        // VCH Admin Portal column
        let vchAdminPortalEl = fixture.debugElement.queryAll(
            By.css('clr-dg-column[ng-reflect-field="vchIp"]'))[1];
        expect(vchAdminPortalEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[
            WS_VCH.DG.COLUMNS.keys.VCH_ADMIN_PORTAL
            ]);
    }));
});
