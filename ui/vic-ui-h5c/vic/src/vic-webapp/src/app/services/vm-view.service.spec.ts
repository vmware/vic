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

import { TestBed, async } from '@angular/core/testing';
import {
    BaseRequestOptions,
    ConnectionBackend,
    RequestOptions,
    HttpModule,
    Http,
    Response,
    ResponseOptions,
    ResponseType,
    XHRBackend,
    Headers
} from '@angular/http';
import { MockBackend, MockConnection } from '@angular/http/testing';
import { VicVmViewService } from './vm-view.service';
import { GlobalsService, Globals } from '../shared/globals.service';
import { Observable } from 'rxjs/Observable';
import {
    getVchResponseStub,
    getMalformedVchResponseStub
} from './mocks/vch.response';
import {
    getContainerResponseStub,
    getMalformedContainerResponseStub
} from './mocks/container.response';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';
import { VirtualContainerHost } from '../vch-view/vch.model';
import { ContainerVm } from '../container-view/container.model';
import { ContainerVmResponse, VchVmResponse } from '../vm.interface';

describe('VicVchViewService', () => {
    let service: VicVmViewService;
    let backend: MockBackend;
    let connection: MockConnection;
    let http: Http;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                Http,
                { provide: ConnectionBackend, useClass: MockBackend },
                { provide: RequestOptions, useClass: BaseRequestOptions },
                VicVmViewService,
                GlobalsService,
                Globals
            ]
        });
        service = TestBed.get(VicVmViewService);
        backend = TestBed.get(ConnectionBackend);
        http = TestBed.get(Http);

        backend.connections.subscribe((c: MockConnection) => connection = c);
    }));

    it('should have service instance', () => {
        expect(service).toBeTruthy();
    });

    it('should generate a correct query string', () => {
        const queryStr: string = service.buildQueryString('test?some=1', {
            offset: 0,
            maxResultCount: 50,
            sorting: 'id,asc',
            filter: 'name=j'
        });

        expect(queryStr).toBe('test?some=1&offset=0' +
            '&maxResultCount=50&sorting=id,asc&filter=name=j');
    });

    it('should process a properly formatted response for vch', async(() => {
        const mockResponse = getVchResponseStub();
        service.getVchsData({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        const vch = new VirtualContainerHost(<VchVmResponse>mockResponse.results.bla);
        expect(vch.parentType).toBe('VirtualApp');
        expect(vch.id).toBe('urn:vmomi:VirtualMachine:vm-27:85421094-c58e-40f9-a42c-b624160d05f5');
        expect(vch.vchVmId).toBe('urn:vmomi:VirtualMachine:vm-27:85421094-c58e-40f9-a42c-b624160d05f5');

        service.vchs$.subscribe(vchs => {
            if (vchs.length) {
                expect(vchs[0].name).toBe('vice-9703');
                expect(vchs.length).toBe(31);
            }
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should return an error when processing a malformed vch response', async(() => {
        const mockResponse = getMalformedVchResponseStub();

        service.getVchsData({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        service.vchs$.subscribe(vchs => { }, err => {
            expect(err).toBeTruthy();
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should return a properly formatted response for container', async(() => {
        const mockResponse = getContainerResponseStub();
        service.getContainersData({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        const cvm = new ContainerVm(<ContainerVmResponse>mockResponse.results['vm-100']);
        expect(cvm.parentType).toBe('VirtualApp');
        expect(cvm.id).toBe('urn:vmomi:VirtualMachine:vm-200:85421094-c58e-40f9-a42c-b624160d05f5');

        service.containers$.subscribe(cvms => {
            if (cvms.length) {
                expect(cvms[0].name).toBe('Container-VM-0-blablablabla');
                expect(cvms.length).toBe(30);
            }
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should return an error when processing a malformed container response', async(() => {
        const mockResponse = getMalformedContainerResponseStub();

        service.getContainersData({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        service.containers$.subscribe(cvms => { }, err => {
            expect(err).toBeTruthy();
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));
});
