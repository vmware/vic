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

import { Component, OnInit, OnDestroy, NgZone } from '@angular/core';
import { GlobalsService, RefreshService } from '../shared/index';
import { Vic18nService } from '../shared/vic-i18n.service';
import { DataPropertyService } from '../services/data-property.service';
import { Subscription } from 'rxjs/Rx';

import {
    VIC_LOGO_100X100,
    WS_SUMMARY
} from '../shared/constants/index';

@Component({
    selector: 'vic-summary-view',
    styleUrls: ['vic-summary-view.scss'],
    template: `
    <div class="summary-container">
        <div class="col-xs-12">
            <img class="vic-summary-view-logo" src="{{vicLogoPath}}">
            <ul class="summary-items-list">
                <li id="vendor">
                    <span class="summary-label">
                        {{vicI18n.translate(
                            this.WS_SUMMARY_CONSTANTS,
                            'VENDOR')}}
                    </span>
                    <span class="summary-value">
                        VMware
                    </span>
                </li>
                <li id="version">
                    <span class="summary-label">
                        {{vicI18n.translate(
                            this.WS_SUMMARY_CONSTANTS,
                            'VERSION')}}
                    </span>
                    <span class="summary-value">
                        {{ pluginVersion }}
                    </span>
                </li>
                <li id="vch_len">
                    <span class="summary-label">
                        {{vicI18n.translate(
                            this.WS_SUMMARY_CONSTANTS,
                            'VCH')}}
                    </span>
                    <span class="summary-value">
                        {{ vchVmsLen }}
                    </span>
                </li>
            </ul>
        </div>
    </div>
    `
})

export class VicSummaryViewComponent implements OnInit, OnDestroy {
    public vicLogoPath: string;
    public pluginVersion: string;
    public vchVmsLen: number;
    public readonly WS_SUMMARY_CONSTANTS = WS_SUMMARY;
    private rootInfoSubscription: Subscription;
    private refreshSubscription: Subscription;

    constructor(
        private zone: NgZone,
        private globalsService: GlobalsService,
        public vicI18n: Vic18nService,
        private refreshService: RefreshService,
        private dataPropertyService: DataPropertyService
    ) {
        this.vicLogoPath = this.isPluginMode() ?
            this.globalsService.getWebContextPath() + VIC_LOGO_100X100 :
            VIC_LOGO_100X100;

        this.rootInfoSubscription = this.dataPropertyService
            .vicObject$.subscribe(res => {
                this.pluginVersion = res.uiVersion;
                this.vchVmsLen = res.vchVmsLen;
            }, err => {
                console.error(err);
            });

        this.refreshSubscription = this.refreshService
            .refreshObservable$.subscribe(() => {
                this.zone.run(() => {
                    this.fetchRootInfo();
                });
            });
    }

    /**
     * Evaluates if the application is running in the vSphere Client environment,
     * which is called the "Plugin Mode"
     * @returns true if plugin mode. false if not
     */
    isPluginMode() {
        return this.globalsService.isPluginMode();
    }

    fetchRootInfo(): void {
        this.dataPropertyService.fetchRootInfo(
            ['vchVmsLen', 'containerVmsLen', 'uiVersion']
        );
    }

    ngOnInit() {
        this.fetchRootInfo();
    }

    ngOnDestroy() {
        if (this.rootInfoSubscription) {
            this.rootInfoSubscription.unsubscribe();
        }

        if (this.refreshSubscription) {
            this.refreshSubscription.unsubscribe();
        }
    }
}
