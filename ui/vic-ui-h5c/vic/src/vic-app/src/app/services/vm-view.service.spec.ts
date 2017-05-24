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
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

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
        let queryStr: string = service.buildQueryString('test?some=1', {
            offset: 0,
            maxResultCount: 50,
            sorting: 'id,asc',
            filter: 'name=j'
        });

        expect(queryStr).toBe('test?some=1&offset=0' +
            '&maxResultCount=50&sorting=id,asc&filter=name=j');
    });

    it('should process a properly formatted response', async(() => {
        const mockResponse = getVchResponseStub();
        service.reloadVchs({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        service.vchs$.subscribe(vchs => {
            if (vchs.length) {
                expect(vchs[0]['name']).toBe('vice-9703');
                expect(vchs.length).toBe(31);
            }
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should return an error when processing a malformed vch response', async(() => {
        const mockResponse = getMalformedVchResponseStub();

        service.reloadVchs({
            offset: 0,
            maxResultCount: 10,
            sorting: 'id,asc'
        });

        service.vchs$.subscribe(vchs => { }, err => {
            console.info(err);
            expect(err).toBeTruthy();
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    // TODO: verify the link to vapp
});