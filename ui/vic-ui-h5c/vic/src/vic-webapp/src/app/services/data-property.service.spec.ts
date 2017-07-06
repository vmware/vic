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
import { GlobalsService, Globals } from '../shared/globals.service';
import { DataPropertyService } from './data-property.service';
import { VirtualMachine } from '../vm.interface';
import { VM_PROPERTIES_TO_EXTRACT } from '../vm.constants';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

describe('VicDataPropertyService', () => {
    let service: DataPropertyService;
    let backend: MockBackend;
    let connection: MockConnection;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(async(() => {
        TestBed.configureCompiler({
            providers: [
                Http,
                { provide: ConnectionBackend, useClass: MockBackend },
                { provide: RequestOptions, useClass: BaseRequestOptions },
                DataPropertyService,
                GlobalsService,
                Globals
            ]
        });
        service = TestBed.get(DataPropertyService);
        backend = TestBed.get(ConnectionBackend);
        backend.connections.subscribe((c: MockConnection) => connection = c);
    }));

    it('should process a VCH with TLS on and a client network configured to use DHCP', async(() => {
        const mockResponse = {
            'id': 'urn:vmomi:VirtualMachine:vm-76:85421094-c58e-40f9-a42c-b624160d05f5',
            'summary.runtime.powerState': 'poweredOn',
            'name': 'virtual-container-host',
            'isVCH': true,
            'config.extraConfig': [
                {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice./init/sessions|docker-personality/cmd/Args~',
                    'value': '/sbin/docker-engine-server|-port=2376|-port-layer-port=2377'
                }, {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice..init.networks|client.assigned.IP',
                    'value': 'ChFtcg=='
                }
            ],
            'isContainer': false
        };
        service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
        service.vmInfo$.subscribe((data: VirtualMachine) => {
            expect(data).toBeTruthy();
            expect(data.isVCH).toEqual(true);
            expect(data.isContainer).toEqual(false);
            expect(data['dockerLog']).toEqual('https://10.17.109.114:2378');
            expect(data['dockerEndpoint']).toEqual('DOCKER_HOST=tcp://10.17.109.114:2376');
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should process a VCH without TLS on and a client network configured to use DHCP', async(() => {
        const mockResponse = {
            'id': 'urn:vmomi:VirtualMachine:vm-76:85421094-c58e-40f9-a42c-b624160d05f5',
            'summary.runtime.powerState': 'poweredOn',
            'name': 'virtual-container-host',
            'isVCH': true,
            'config.extraConfig': [
                {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice./init/sessions|docker-personality/cmd/Args~',
                    'value': '/sbin/docker-engine-server|-port=2375|-port-layer-port=2377'
                }, {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice..init.networks|client.assigned.IP',
                    'value': 'ChFtcg=='
                }
            ],
            'isContainer': false
        };
        service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
        service.vmInfo$.subscribe((data: VirtualMachine) => {
            expect(data).toBeTruthy();
            expect(data.isVCH).toEqual(true);
            expect(data.isContainer).toEqual(false);
            expect(data['dockerLog']).toEqual('https://10.17.109.114:2378');
            expect(data['dockerEndpoint']).toEqual('DOCKER_HOST=tcp://10.17.109.114:2375');
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should process a VCH with TLS on and a client network configured to use a static IP', async(() => {
        const mockResponse = {
            'id': 'urn:vmomi:VirtualMachine:vm-76:85421094-c58e-40f9-a42c-b624160d05f5',
            'summary.runtime.powerState': 'poweredOn',
            'name': 'virtual-container-host',
            'isVCH': true,
            'config.extraConfig': [
                {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice./init/sessions|docker-personality/cmd/Args~',
                    'value': '/sbin/docker-engine-server|-port=2376|-port-layer-port=2377'
                }, {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice..init.networks|client.assigned.IP',
                    'value': 'AAAAAAAAAAAAAP//wKhkFg=='
                }
            ],
            'isContainer': false
        };
        service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
        service.vmInfo$.subscribe((data: VirtualMachine) => {
            expect(data).toBeTruthy();
            expect(data.isVCH).toEqual(true);
            expect(data.isContainer).toEqual(false);
            expect(data['dockerLog']).toEqual('https://192.168.100.22:2378');
            expect(data['dockerEndpoint']).toEqual('DOCKER_HOST=tcp://192.168.100.22:2376');
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));

    it('should process a VCH without TLS on and a client network configured to use a static IP', async(() => {
        const mockResponse = {
            'id': 'urn:vmomi:VirtualMachine:vm-76:85421094-c58e-40f9-a42c-b624160d05f5',
            'summary.runtime.powerState': 'poweredOn',
            'name': 'virtual-container-host',
            'isVCH': true,
            'config.extraConfig': [
                {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice./init/sessions|docker-personality/cmd/Args~',
                    'value': '/sbin/docker-engine-server|-port=2375|-port-layer-port=2377'
                }, {
                    'dynamicType': null,
                    'dynamicProperty': null,
                    'key': 'guestinfo.vice..init.networks|client.assigned.IP',
                    'value': 'AAAAAAAAAAAAAP//wKhkFg=='
                }
            ],
            'isContainer': false
        };
        service.fetchVmInfo(VM_PROPERTIES_TO_EXTRACT);
        service.vmInfo$.subscribe((data: VirtualMachine) => {
            expect(data).toBeTruthy();
            expect(data.isVCH).toEqual(true);
            expect(data.isContainer).toEqual(false);
            expect(data['dockerLog']).toEqual('https://192.168.100.22:2378');
            expect(data['dockerEndpoint']).toEqual('DOCKER_HOST=tcp://192.168.100.22:2375');
        });

        connection.mockRespond(new Response(new ResponseOptions({
            body: mockResponse
        })));
    }));
});
