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

import { TestBed } from '@angular/core/testing';

import { GlobalsService } from '../shared/index';
import { webPlatformStub } from '../shared/dev/webPlatformStub';

export * from './router-stubs';

/**
 * Stub for testing in plugin more or dev mode
 */
export const globalStub = {
    pluginMode: true,
    webPlatform: webPlatformStub
};

/**
 * Initialization for unit tests
 */
export function initGlobalService(pluginMode: boolean): GlobalsService {

    globalStub.pluginMode = pluginMode;
    return TestBed.get(GlobalsService);
}

export const appErrorHandlerStub = {
    httpPromiseError(error: any): Promise<any> {
        return Promise.reject('error message from appErrorHandlerStub');
    }
};

// ---- Utilities copied from Angular2 doc ----

/**
 * Create custom DOM event the old fashioned way
 *
 * https://developer.mozilla.org/en-US/docs/Web/API/Event/initEvent
 * Although officially deprecated, some browsers (phantom) dont accept the preferred 'new Event(eventName)'
 */
export function newEvent(eventName: string, bubbles = false, cancelable = false) {
    const evt = document.createEvent('CustomEvent');
    evt.initCustomEvent(eventName, bubbles, cancelable, null);
    return evt;
}
