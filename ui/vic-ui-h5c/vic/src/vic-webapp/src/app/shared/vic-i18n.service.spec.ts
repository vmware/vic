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

import { async, TestBed } from '@angular/core/testing';
import {
    Globals,
    GlobalsService,
    I18nService,
    Vic18nService
} from './index';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

class I18nServiceStub {
    public translate(key: string, params: string | string[] = null): string {
        // this.i18n.translate() is already tested by h5c team
        // so just assume the key is returned as-is in this case
        return key;
    }
}

describe('VicI18nService', () => {
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;
    let vicI18n: Vic18nService;
    const stub = {
        WS_SUMMARY: {
            keys: {
                NAME: 'a.b.c.d.e'
            },
            defaults: {
                'a.b.c.d.e': 'Default Value'
            }
        }
    };

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            providers: [
                { provide: I18nService, useClass: I18nServiceStub },
                Vic18nService
            ]
        });
        vicI18n = TestBed.get(Vic18nService);
    }));

    it('should be initialized', () => {
        expect(vicI18n).toBeTruthy();
    });

    it('should get a translated value for an existing key', () => {
        const v = vicI18n.translate(stub.WS_SUMMARY, 'NAME');
        expect(v).toBe('Default Value');
    });

    it('should display no value for a nonexistent key', () => {
        const v = vicI18n.translate(stub.WS_SUMMARY, 'A');
        expect(v).toBe('');
    });
});
