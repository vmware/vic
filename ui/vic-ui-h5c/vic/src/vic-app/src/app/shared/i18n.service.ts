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
import { Http } from '@angular/http';
import { GlobalsService, APP_CONFIG } from './index';
import { AppAlertService } from './app-alert.service';

/**
 * Internationalization service for handling message translation within plugin views
 *
 * This implementation relies on the HTML SDK internationalization guidelines for
 * keeping text messages compatible between vSphere HTML client and Flex client.
 * For text that is directly displayed in an HTML view it may be possible to use
 * Angular i18n tools (see https://angular.io/docs/ts/latest/cookbook/i18n.html)
 */
@Injectable()
export class I18nService {
   bundle = {};
   bundleName = APP_CONFIG.bundleName;

   constructor(private gs: GlobalsService,
               private appAlertService: AppAlertService,
               private http: Http) {
   }

   /**
    * Initialization of i18n bundles for dev mode.
    */
   initLocale(locale): void {
      this.gs.locale = locale;

      if (this.gs.isPluginMode()) {
         // Local bundle is only used in dev mode
         return;
      }
      // Only handle 2 locales here
      let localeCode = 'en_US';
      if (locale.startsWith('fr')) {
         localeCode = 'fr_FR';
      }
      let jsonFile = this.bundleName + '_' + localeCode + '.json';

      let errorMsg = 'Cannot load /src/webapp/locales/' + jsonFile + '! ' +
            'Please check local file and start json-server with --static ./src/webapp';

      // This requires properties file to have been converted to .json ahead of time!
      this.http.get('http://localhost:3000/locales/' + jsonFile)
            .toPromise()
            .then(res => res.json())
            .catch(error => this.appAlertService.showError(errorMsg))
            .then(bundle => this.bundle = bundle);
   }

   /**
    * Get the translated message for the given key and optional parameters
    * @param key
    * @param params
    * @returns {any}
    */
   public translate(key: string, params: string|string[] = null): string {
      if (this.gs.isPluginMode() &&
          this.gs.getWebPlatform().getString) {
         // SDK's getString allows compatibility with vSphere Flex Client
         return this.gs.getWebPlatform().getString(this.bundleName, key, params);
      }
      if (this.bundle && this.bundle[key]) {
         return this.interpolate(this.bundle[key], params);
      }
      // Display non translated keys as is
      return key;
   }

   // Insert parameters in messages containing placeholders {0} {1} ...
   interpolate(message: string, params: string|string[]): string {
      if (params) {
         if (typeof params === 'string') {
            params = [params];
         }
         message = message.replace(/\{(\d+)\}/g, function (match, index) {
            if (index >= params.length) {
               // Less parameters than there are placeholders, so return the placeholder value.
               return match;
            }
            return params[ index ];
         });
      }
      return message;
   }
}
