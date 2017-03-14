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

import { Injectable, Inject, forwardRef } from '@angular/core';
import { Http, URLSearchParams, Response } from '@angular/http';
import { GlobalsService }     from '../shared/globals.service';
import { AppAlertService }    from '../shared/app-alert.service';
import { AppErrorHandler }    from '../shared/appErrorHandler';


// Store the ActionDevService instance in a global because we can't inject ActionDevService in webPlatformStub
export let actionDevService: ActionDevService;

/**
 * Dev mode service to call action controllers
 */
@Injectable()
export class ActionDevService {
   gs: GlobalsService;

   constructor(private http: Http,
               private appAlertService: AppAlertService,
               private errorHandler: AppErrorHandler,
               @Inject(forwardRef(() => GlobalsService)) gs: GlobalsService) {
      this.gs = gs;
      actionDevService = this;
   }

   /**
    * Dev implementation of ca
    * @param url
    * @param jsonData
    * @param targets
    */
   public callActionsController(url: string, jsonData: string, targets: string): void {
      if (this.gs.useLiveData()) {
         // Post to the Java service rest endpoint when testing live data
         let headers = this.gs.getHttpHeaders();
         let data = new URLSearchParams();
         data.append('targets', targets);
         data.append('json', jsonData);

         this.http.post(url, data, headers)
               .toPromise()
               .then((response: Response) => {
                  return (url + ' returned: ' + response.text());
               })
               .catch(error => this.errorHandler.httpPromiseError(error))
               .then(msg => this.appAlertService.showInfo(msg ? msg : 'no result'))
               .catch(errMsg => this.appAlertService.showError(errMsg));
      } else {
         // Just show that the action was called
         this.appAlertService.showInfo('URL: ' + url + ' called with json: ' + jsonData);
      }
   }

}
