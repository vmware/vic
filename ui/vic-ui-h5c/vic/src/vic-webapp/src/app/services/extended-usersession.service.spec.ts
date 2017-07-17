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
import { GlobalsService } from '../shared';
import { ExtendedUserSessionService } from './extended-usersession.service';
import { IExtendedServerInfo } from './extended-serverinfo.interface';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

describe('ExtendedUserSessionService', () => {
    let service: ExtendedUserSessionService;
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    const getMockedGlobalsService = (groups: string, serversInfo: IExtendedServerInfo[]) => {
        return {
            isPluginMode: () => {
                return true;
            },
            getWebPlatform: () => {
                return {
                    getUserSession() {
                        return {
                            groups: groups,
                            locale: 'en_US',
                            serversInfo: serversInfo,
                            samlTokenXml: 'loremipsum'
                        };
                    }
                };
            }
        };
    };

    it('should return true for isVsphereAdmin for an admin user ', async(() => {
        TestBed.configureTestingModule({
            providers: [
                {
                    provide: GlobalsService,
                    useValue: getMockedGlobalsService('vsphere.local\\Administrators', [])
                },
                ExtendedUserSessionService
            ]
        }).compileComponents();
        service = TestBed.get(ExtendedUserSessionService);
        expect(service.isVsphereAdmin).toBeTruthy();
    }));

    it('should return false for isVsphereAdmin for a non-admin user', async(() => {
        TestBed.configureTestingModule({
            providers: [
                {
                    provide: GlobalsService,
                    useValue: getMockedGlobalsService('vsphere.local\\Everyones', [])
                },
                ExtendedUserSessionService
            ]
        }).compileComponents();
        service = TestBed.get(ExtendedUserSessionService);
        expect(service.isVsphereAdmin).toBeFalsy();
    }));

    it('should return locale, samlTokenXml and ServerInfos', async(() => {
        TestBed.configureTestingModule({
            providers: [
                {
                    provide: GlobalsService,
                    useValue: getMockedGlobalsService(
                        'vsphere.local\\Administrators',
                        [{
                            name: 'test-name',
                            serverGuid: 'test-guid',
                            sessionCookie: 'test-cookie',
                            thumbprint: 'test-thumbprint',
                            version: 'test-version'
                        }])
                },
                ExtendedUserSessionService
            ]
        }).compileComponents();
        service = TestBed.get(ExtendedUserSessionService);
        expect(service.locale).toBe('en_US');
        expect(service.samlTokenXml).toBe('loremipsum');
        const serverInfos = service.getVcenterServersInfo();
        expect(serverInfos.length).toBe(1);
        expect(serverInfos[0]).toEqual({
            name: 'test-name',
            serverGuid: 'test-guid',
            sessionCookie: 'test-cookie',
            thumbprint: 'test-thumbprint',
            version: 'test-version'
        });
    }));
});
