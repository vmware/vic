/* tslint:disable:no-unused-variable */
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
import { By } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';
import { Observable, Subject } from 'rxjs/Rx';

import {
    AppAlertService,
    Globals,
    GlobalsService,
    I18nService,
    Vic18nService,
    RefreshService,
} from '../shared/index';
import { DataPropertyService } from '../services/data-property.service';
import { AppErrorHandler } from '../shared/appErrorHandler';
import { VicSummaryViewComponent } from './summary-view.component';
import { ClarityModule } from 'clarity-angular';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

describe('VIC object view: Summary tab', () => {
    let fixture: ComponentFixture<VicSummaryViewComponent>;
    let compInstance: VicSummaryViewComponent;
    let dpService: DataPropertyService;
    let globalService: GlobalsService;
    let refreshService: RefreshService;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    class DpServiceStub {
        private vicObjectSource = new Subject<any>();
        public vicObject$ = this.vicObjectSource.asObservable();
        fetchRootInfo(props: string[]): void {
            this.vicObjectSource.next({
                uiVersion: '1.1',
                vchVmsLen: 1,
                containerVmsLen: 1
            });
        }
    }

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                { provide: DataPropertyService, useClass: DpServiceStub },
                AppAlertService,
                Globals,
                GlobalsService,
                I18nService,
                Vic18nService,
                RefreshService
            ],
            imports: [
                HttpModule
            ],
            declarations: [
                VicSummaryViewComponent
            ]
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(VicSummaryViewComponent);
        compInstance = fixture.componentInstance;
        dpService = fixture.debugElement.injector.get(DataPropertyService);
        globalService = fixture.debugElement.injector.get(GlobalsService);
        refreshService = fixture.debugElement.injector.get(RefreshService);
    });

    it('should display version and vch length correctly', () => {
        fixture.detectChanges();
        const ul = fixture.debugElement.query(By.css('ul.summary-items-list'));
        const span_uiVersion = ul
            .query(By.css('li#version'))
            .query(By.css('span:nth-child(2)'));
        const span_vchVmsLen = ul
            .query(By.css('li#vch_len'))
            .query(By.css('span:nth-child(2)'));

        expect(span_uiVersion.nativeElement.textContent.trim()).toBe('1.1', 'should be 1.1');
        expect(span_vchVmsLen.nativeElement.textContent.trim()).toBe('1', 'should be 1');
    });

    // TODO: i18n test
});
