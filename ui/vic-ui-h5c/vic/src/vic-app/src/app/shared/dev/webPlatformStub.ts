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

import { WebPlatform, UserSession }  from '../vSphereClientSdkTypes';

import { actionDevService } from '../../services/action-dev.service';

/**
 * The WebPlatform stub used for dev mode and unit testing
 */
export const webPlatformStub: WebPlatform = {

   callActionsController(url: string, jsonData: string, targets?: string): void {
      console.log('callActionsController called: url = ' + url + ', jsonData = ' + jsonData);
      actionDevService.callActionsController(url, jsonData, targets);
   },
   closeDialog(): void {
      console.log('closeDialog called');
   },
   getClientVersion(): string { return '6.5.0'; },
   getClientType(): string { return 'html'; },

   getRootPath(): string {
      // FIXME: change it to 3000 later
      return 'http://localhost:4201/ui';
   },

   getString(bundleName: string, key: string, params: any): string { return ''; },
   getUserSession(): UserSession { return null; },
   openModalDialog(): void {
      console.log('openModalDialog called');
   },
   sendModelChangeEvent(): void {
      console.log('sendModelChangeEvent called');
   },
   sendNavigationRequest(): void {
      console.log('sendNavigationRequest called');
   },
   setDialogSize(width: string, height: string): void {
      console.log('setDialogSize called: width = ' + width + 'height = ' + height);
   },
   setDialogTitle(title: string): void {
      console.log('setDialogTitle called: title = ' + title);
   },
   setGlobalRefreshHandler(callback, document): void {
      console.log('setGlobalRefreshHandler called');
   }
};



