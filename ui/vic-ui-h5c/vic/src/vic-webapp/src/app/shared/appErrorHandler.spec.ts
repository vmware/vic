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

// External imports
import { Response, ResponseOptions } from '@angular/http';

// Internal imports
import { AppErrorHandler, liveDataHelp, jsonServerHelp } from '../shared/appErrorHandler';
import { JASMINE_TIMEOUT } from '../testing/jasmine.constants';

// Simple service unit tests without assistance from Angular testing utilities

describe('AppErrorHandler tests', () => {
    let appErrorHandler: AppErrorHandler;
    const errorFromServer = 'some error message\nsome more lines';
    const errorToDisplay = 'some error message';
    jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;

    beforeEach(() => {
        appErrorHandler = new AppErrorHandler(null);
    });

    it('formats server errors correctly - with body object', () => {
        const resOptions: ResponseOptions = {
            status: 500,
            body: { message: errorFromServer },
            headers: null,
            url: 'some url',
            merge: null
        };
        const error = new Response(resOptions);
        const statusText = 'Server error';
        error.statusText = statusText;

        const errMsg = appErrorHandler.formatHttpError(error);
        expect(errMsg).toBe(statusText + ' => ' + errorToDisplay);
    });

    it('formats server errors correctly - with body string', () => {
        const resOptions: ResponseOptions = {
            status: 500,
            body: errorFromServer,
            headers: null,
            url: 'some url',
            merge: null
        };
        const error = new Response(resOptions);
        const statusText = 'Server error';
        error.statusText = statusText;

        const errMsg = appErrorHandler.formatHttpError(error);
        expect(errMsg).toBe(statusText + ' => ' + errorToDisplay);
    });


    it('formats http errors correctly for live data', () => {
        const resOptions: ResponseOptions = {
            status: 401,
            body: '',
            headers: null,
            url: 'some url',
            merge: null
        };
        const error = new Response(resOptions);
        const statusText = 'some status';
        error.statusText = statusText;

        let pluginMode = true;
        let errMsg = appErrorHandler.formatHttpError(error, pluginMode);
        expect(errMsg).toBe('Http error: 401, ' + statusText);

        pluginMode = false;
        errMsg = appErrorHandler.formatHttpError(error, pluginMode);
        expect(errMsg).toBe('Http error: 401, ' + statusText + liveDataHelp);

        pluginMode = true;
        error.status = 404;
        errMsg = appErrorHandler.formatHttpError(error, pluginMode);
        expect(errMsg).toBe('Http error: 404, ' + statusText + ' at URL: ' + resOptions.url);
    });

    it('returns json-server help in dev mode', () => {
        const resOptions: ResponseOptions = {
            status: 0,
            body: '',
            headers: null,
            url: 'some url',
            merge: null
        };
        const error = new Response(resOptions);

        const pluginMode = false;
        const useLiveData = false;
        const errMsg = appErrorHandler.formatHttpError(error, pluginMode, useLiveData);
        expect(errMsg).toBe(jsonServerHelp);
    });
});
