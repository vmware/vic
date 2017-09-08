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
import { ExtendedUserSessionService } from '../services/extended-usersession.service';
import { VicVchViewComponent } from './vch-view.component';
import { VirtualContainerHost } from './vch.model';
import { ClarityModule } from 'clarity-angular';
import {
    getVchResponseStub,
    getMalformedVchResponseStub
} from '../services/mocks/vch.response';
import { WS_VCH } from '../shared/constants';
import { NO_ERRORS_SCHEMA } from '@angular/core';

let responseProperlyFormatted = true;

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
        const vchResponse = responseProperlyFormatted ?
            getVchResponseStub().results : getMalformedVchResponseStub().results;
        try {
            for (const objId in vchResponse) {
                if (vchResponse.hasOwnProperty(objId)) {
                    this.data.push(new VirtualContainerHost(vchResponse[objId]));
                }
            }
            this.vchsSubj.next(this.data);
        } catch (e) {
            this.vchsSubj.error('error');
        }
    }
}

describe('VicVchViewComponent', () => {
    let fixture: ComponentFixture<VicVchViewComponent>;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            schemas: [NO_ERRORS_SCHEMA],
            providers: [
                { provide: VicVmViewService, useClass: VicVmViewServiceStub },
                ExtendedUserSessionService,
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
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent<VicVchViewComponent>(VicVchViewComponent);
        responseProperlyFormatted = true;
    });

    it('should render the data grid with properly formatted data', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();
        const rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
        const rowElementsLength = rowElements.length;
        expect(rowElementsLength).toBe(31);
        if (rowElementsLength) {
            const cellElements = rowElements[0].queryAll(By.css('clr-dg-cell'));
            const cellElementsLength = cellElements.length;
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
            console.log('expected exception caught');
        } finally {
            const rowElements = fixture.debugElement.queryAll(By.css('clr-dg-row'));
            const rowElementsLength = rowElements.length;
            expect(rowElementsLength).toBe(0);
        }
    }));

    it('should render default localized text for table headers', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();

        // containerName column
        const nameEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="name"]'));
        expect(nameEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[WS_VCH.DG.COLUMNS.keys.NAME]);

        // overallStatus column
        const overallStatusEl = fixture.debugElement.query(
            By.css('clr-dg-column[ng-reflect-field="overallStatus"]'));
        expect(overallStatusEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[WS_VCH.DG.COLUMNS.keys.OVERALL_STATUS]);

        // Docker API Endpoint column
        const dockerApiEndpointEl = fixture.debugElement.queryAll(
            By.css('clr-dg-column[ng-reflect-field="vchIp"]'))[0];
        expect(dockerApiEndpointEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[
            WS_VCH.DG.COLUMNS.keys.DOCKER_API_ENDPOINT
            ]);

        // VCH Admin Portal column
        const vchAdminPortalEl = fixture.debugElement.queryAll(
            By.css('clr-dg-column[ng-reflect-field="vchIp"]'))[1];
        expect(vchAdminPortalEl.nativeElement.textContent.trim()).toBe(
            WS_VCH.DG.COLUMNS.defaults[
            WS_VCH.DG.COLUMNS.keys.VCH_ADMIN_PORTAL
            ]);
    }));

    it('should render the new VCH button for an admin user', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();
        fixture.componentInstance.isVsphereAdmin = true;
        fixture.detectChanges();
        const actionBarEl = fixture.debugElement.query(
            By.css('clr-dg-action-bar'));
        expect(actionBarEl).toBeTruthy();
    }));

    it('should not render the new VCH button for a non-admin user', async(() => {
        fixture.componentInstance.ngOnInit();
        fixture.componentInstance.reloadVchs();
        fixture.detectChanges();
        fixture.componentInstance.isVsphereAdmin = false;
        fixture.detectChanges();
        const actionBarEl = fixture.debugElement.query(
            By.css('clr-dg-action-bar'));
        expect(actionBarEl).toBeNull();
    }));
});
