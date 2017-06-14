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
import { By } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { Location, LocationStrategy } from '@angular/common';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { AppRoutingComponent } from './app-routing.component';
import { Globals, GlobalsService } from './shared';
import { ClarityModule } from 'clarity-angular';
import { JASMINE_TIMEOUT } from './testing/jasmine.constants';
import { AppRoutingModule, routedComponents } from 'app/app-routing.module';

describe('VicAppRoutingComponent', () => {
    let fixture: ComponentFixture<AppRoutingComponent>;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            schemas: [NO_ERRORS_SCHEMA],
            imports: [
                ClarityModule,
                AppRoutingModule
            ],
            declarations: [
                AppRoutingComponent,
                routedComponents
            ],
            providers: [
                {
                    provide: ActivatedRoute, useValue: {
                        queryParams: Observable.of({
                            view: 'abc',
                            objectId: 'ddd',
                            actionUid: 'ddd',
                            targets: '',
                            locale: ''
                        })
                    }
                },
                {
                    provide: Router, useValue: {
                        navigate: (target) => {
                            return target;
                        }
                    }
                },
                Globals,
                GlobalsService,
                {
                    provide: Location, useValue: {
                        path: () => 'https://10.17.109.170/ui/' +
                            '#?extensionId=vsphere.core.vm.summary' +
                            '&objectId=urn:vmomi:VirtualMachine:vm-54:85421094-c58e-40f9-a42c-b624160d05f5' +
                            '&navigator=vsphere.core.viTree.hostsAndClustersView'
                    }
                },
            ]
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AppRoutingComponent);
        fixture.detectChanges();
    });

    it('should have component initializd', () => {
        expect(fixture.componentInstance).toBeTruthy();
    });
});
