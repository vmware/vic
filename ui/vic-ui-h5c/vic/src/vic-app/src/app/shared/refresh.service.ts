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
import { Subject }    from 'rxjs/Subject';

import { AppAlertService }   from './app-alert.service';

/**
 * Service used to send a 'refresh event' to any observer view
 */
@Injectable()
export class RefreshService {
   // Use an rxjs Subject to multicast to multiple observers.
   // See http://reactivex.io/rxjs/manual/overview.html#subject
   private refreshSource = new Subject();
   public refreshObservable$ = this.refreshSource.asObservable();

   constructor(private appAlertService: AppAlertService) {
   }

   public refreshView(): void {
      // Close any open alert box here before a view is refreshed
      this.appAlertService.closeAlert();

      // Propagate refresh event to subscribers
      this.refreshSource.next();
   }
}
