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

import { Injectable } from '@angular/core';
import { I18nService } from '../shared/i18n.service';
@Injectable()
export class Vic18nService {
    constructor(
        private i18n: I18nService
    ) { }

    /**
     * Returns localized text for a given key
     * @param key : key defined in com_vmware_vic.properties
     * @returns localized text
     */
    translate(ns: any, alias: string) {
        const key = ns['keys'][alias];
        const results = this.i18n.translate(key);
        if (results === key) {
            // when unit testing or run in a standalone mode,
            // key is returned as-is. for this case, look for its default
            // value defined in constants/resources.path.ts
            return ns['defaults'][key] || '';
        }
        return results;
    }
}
