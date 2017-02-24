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

import { Component, Injector } from '@angular/core';
import { GlobalsService, RefreshService, I18nService } from './shared/index';
import { ActionDevService } from './services/action-dev.service';

@Component({
    selector: 'vic-app',
    template: `
    <div *ngIf="!gs.isPluginMode()" class="floating-left">
        <a (click)="gs.toggleDevUI()" class="tooltip tooltip-sm tooltip-bottom-right"
           role="tooltip" aria-haspopup="true" href="javascript://">
           <clr-icon [attr.shape]="gs.showDevUI() ? 'remove' : 'plus-circle'" size="16"
                     [attr.class]="gs.showDevUI() ? 'is-inverse' : ''"></clr-icon>
           <span class="tooltip-content">{{gs.showDevUI() ? "Remove dev UI" : "Show dev UI"}}</span>
        </a>
    </div>
    <router-outlet></router-outlet>
    `
})

export class AppComponent {

    constructor(
        public gs: GlobalsService,
        private injector: Injector,
        private refreshService: RefreshService,
        private i18nService: I18nService
    ) {
        // Refresh handler to be used in plugin mode
        this.gs.getWebPlatform().setGlobalRefreshHandler(
            this.refresh.bind(this), document
        );

        // Manual injection of ActionDevService, used in webPlatformStub
        if (!this.gs.isPluginMode()) {
            this.injector.get(ActionDevService);
        }

        // Start the app in english by default (dev mode)
        // In plugin mode the current locale is passed as parameter
        this.i18nService.initLocale('en');
    }

    refresh(): void {
        this.refreshService.refreshView();
    }
}
