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

import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { Subscription } from 'rxjs';

import { GlobalsService } from './shared/index';

@Component({
    template: `
    <div *ngIf="!gs.isPluginMode()" class="floating-left">
        <a (click)="gs.toggleDevUI()" class="tooltip tooltip-sm tooltip-bottom-right"
           role="tooltip" aria-haspopup="true" href="javascript://">
            <clr-icon [attr.shape]="gs.showDevUI() ? 'remove' : 'plus-circle'" size="16"
                      [attr.class]="gs.showDevUI() ? 'is-inverse' : ''"></clr-icon>
            <span class="tooltip-content">
                {{gs.showDevUI() ? "Remove dev UI" : "Show dev UI"}}
            </span>
        </a>
    </div>
    <router-outlet></router-outlet>
    `
})
export class AppRoutingComponent implements OnInit, OnDestroy {
    private subscription: Subscription;

    constructor(
        public gs: GlobalsService,
        private router: Router,
        private location: Location,
        private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        let path = this.location.path();
        console.log('app.component path = ' + path);

        const FORWARD_SLASH_ENCODED2 = '%252F';
        const FORWARD_SLASH_ENCODED = /%2F/g;

        // Extract query parameters and navigate to view
        this.subscription = this.route.queryParams.subscribe(
            (param: any) => {
                let view = param['view'];
                let objectId = param['objectId'];
                let actionUid = param['actionUid'];
                let targets = param['targets'];
                let locale = param['locale'];
                let params = {};

                if (!view) {
                    throw new Error('Missing view parameter! path = ' + path);
                }
                if (objectId) {
                    objectId = objectId.replace(FORWARD_SLASH_ENCODED, FORWARD_SLASH_ENCODED2);
                }
                if (actionUid) {
                    params['actionUid'] = actionUid;
                    if (targets) {
                        objectId = targets.replace(FORWARD_SLASH_ENCODED, FORWARD_SLASH_ENCODED2);
                    } else {
                        objectId = 'undefined';
                    }
                }
                if (locale) {
                    this.gs.locale = locale;
                }
                let commands: [any] = ['/' + view];
                if (objectId) {
                    commands[1] = objectId;
                }
                if (actionUid) {
                    commands[2] = actionUid;
                }

                setTimeout(() => this.router.navigate(commands), 0);
            }
        );
    }

    ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }
}
