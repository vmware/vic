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

import { Subject } from 'rxjs/Subject';

/**
 * Service used to display top-level alerts, see app-alert.component.
 */
export class AppAlertService {
    // Observable sources:
    // alertMessageSource array contains the message to display and the alert type (see Clarity doc)
    // closeAlertSource is for closing the alert component
    private alertMessageSource = new Subject<[string, string]>();
    private closeAlertSource = new Subject();

    // Observable streams
    alertMessage$ = this.alertMessageSource.asObservable();
    closeAlert$ = this.closeAlertSource.asObservable();

    showError(message: string) {
        this.alertMessageSource.next([message, 'alert-danger']);
    }

    showInfo(message: string) {
        this.alertMessageSource.next([message, 'alert-info']);
    }

    showWarning(message: string) {
        this.alertMessageSource.next([message, 'alert-warning']);
    }

    showSuccess(message: string) {
        this.alertMessageSource.next([message, 'alert-success']);
    }

    closeAlert() {
        this.closeAlertSource.next();
    }
}
