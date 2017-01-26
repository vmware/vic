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
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable, Subject, BehaviorSubject } from 'rxjs/Rx';

import { VicSummaryPortletComponent } from './summary-portlet.component';
import { VchPortletComponent } from './vch-portlet/vch-portlet.component';
import { ContainerPortletComponent } from './container-portlet/container-portlet.component';
import { DataPropertyService } from '../data-property-service/data-property.service';
import { VirtualMachine, VM_PROPERTIES_TO_EXTRACT } from '../vm.interface';

describe('VCH Component Unit Test', () => {
    let fixture: ComponentFixture<VicSummaryPortletComponent>;
    let compInstance: VicSummaryPortletComponent;
    let svc: DataPropertyService;
    const defaultVmObj = {
        powerState: 'poweredOn'
    };

    const getVirtualMachineStub: Function = (obj: any): VirtualMachine => {
        let cloned = JSON.parse(JSON.stringify(defaultVmObj));
        for (let key in obj) {
            cloned[key] = obj[key];
        }
        return <VirtualMachine>cloned;
    };

    class DpServiceStub {
        private results: BehaviorSubject<VirtualMachine> = new BehaviorSubject<VirtualMachine>(null);

        getCurrent(): VirtualMachine {
            return this.results.getValue();
        }

        fetch(): Observable<VirtualMachine> {
            this.results.next(<VirtualMachine>defaultVmObj);

            return this.results;
        }
    };

    class ActivatedRouteStub {
        navigateByUrl(url: string) {
            return url;
        }
    }

    class RouterStub {
        navigateByUrl(url: string) {
            return url;
        }
    }

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                { provide: DataPropertyService, useClass: DpServiceStub },
                { provide: ActivatedRoute, useClass: ActivatedRouteStub },
                { provide: Router, useClass: RouterStub }
            ],
            declarations: [
                VicSummaryPortletComponent,
                VchPortletComponent,
                ContainerPortletComponent
            ],
            imports: [
                FormsModule
            ]
        })
        .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent<VicSummaryPortletComponent>(VicSummaryPortletComponent);
        compInstance = fixture.componentInstance;
        svc = fixture.debugElement.injector.get(DataPropertyService);
    });

    it('should have test service', () => {
        expect(svc).toBeDefined();
    });

    it('should display the container vm portlet - without portmapping info', () => {
        compInstance.activeVm = getVirtualMachineStub({
            id: '1',
            name: 'fake container vm',
            container_name: 'fcv',
            isVCH: false,
            isContainer: true,
            image_name: 'busybox'
        });

        fixture.detectChanges();
        let trs = fixture.debugElement.queryAll(By.css('tr'));
        expect(trs.length).toBe(2);
        expect(trs[0].nativeElement.querySelector('td').textContent).toBe('fcv');
        expect(trs[1].nativeElement.querySelector('td').textContent).toBe('busybox');
    });

    it('should display the container vm portlet - with portmapping info', () => {
        compInstance.activeVm = getVirtualMachineStub({
            id: '2',
            name: 'fake container vm 2',
            container_name: 'fcv2',
            isVCH: false,
            isContainer: true,
            image_name: 'nginx',
            portmapping: '8081:80/tcp'
        });

        fixture.detectChanges();
        let trs = fixture.debugElement.queryAll(By.css('tr'));
        expect(trs.length).toBe(3);
        expect(trs[0].nativeElement.querySelector('td').textContent).toBe('fcv2');
        expect(trs[1].nativeElement.querySelector('td').textContent).toBe('nginx');
        expect(trs[2].nativeElement.querySelector('td').textContent).toBe('8081:80/tcp');
    });

    it('should display the vch information when vm is on', () => {
        compInstance.activeVm = getVirtualMachineStub({
            id: '3',
            name: 'fake vch vm',
            isVCH: true,
            isContainer: false,
            dockerEndpoint: 'DOCKER_HOST=tcp://1.2.3.4:2376',
            dockerLog: 'https://1.2.3.4:2378'
        });

        fixture.detectChanges();
        let trs = fixture.debugElement.queryAll(By.css('tr'));
        expect(trs.length).toBe(2);
        expect(trs[0].nativeElement.querySelector('td').textContent.trim()).toBe('DOCKER_HOST=tcp://1.2.3.4:2376');
        expect(trs[1].nativeElement.querySelector('td').textContent.trim()).toBe('https://1.2.3.4:2378');
    });

    it('should display the placeholder values for vch when vm is off', () => {
        compInstance.activeVm = getVirtualMachineStub({
            id: '4',
            name: 'fake vch vm 2',
            isVCH: true,
            isContainer: false,
            dockerEndpoint: '-',
            dockerLog: '-',
            powerState: 'poweredOff'
        });

        fixture.detectChanges();
        let trs = fixture.debugElement.queryAll(By.css('tr'));
        expect(trs.length).toBe(2);
        expect(trs[0].nativeElement.querySelector('td').textContent.trim()).toBe('-');
        expect(trs[1].nativeElement.querySelector('td').textContent.trim()).toBe('-');
    });

    it('should not display portlets when activeVm contains corrupt data', () => {
        compInstance.activeVm = <VirtualMachine>{
            id: '5',
            name: 'hai'
        };

        fixture.detectChanges();
        let table = fixture.debugElement.queryAll(By.css('table'));
        expect(table.length).toBe(0);
    });

});