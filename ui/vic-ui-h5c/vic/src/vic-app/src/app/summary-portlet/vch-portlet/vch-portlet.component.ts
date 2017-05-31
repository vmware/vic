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
import { I18nService } from '../../shared/i18n.service';
import { GlobalsService } from '../../shared/globals.service';

import { Component, Input, OnInit } from '@angular/core';
import { VirtualMachine } from '../../vm.interface';

import {
    DEFAULT_VCH_DOCKER_ENDPOINT_LABEL,
    DEFAULT_VCH_VCH_ADMIN_PORTAL_LABEL
} from '../../shared/constants/index';

@Component({
    selector: 'vic-vch-portlet',
    styleUrls: ['vch-portlet.scss'],
    template: `
    <table class="table table-vertical">
        <tbody>
            <tr>
                <th>{{ dockerApiEndpointLabel }}</th>
                <td>{{ activeVm['dockerEndpoint'] }}</td>
            </tr>
            <tr>
                <th>{{ vchAdminPortalLabel }}</th>
                <td [ngSwitch]="activeVm['powerState']">
                    <span *ngSwitchCase="'poweredOn'">
                        <a href="{{ activeVm['dockerLog'] }}" target="_blank">
                        {{ activeVm['dockerLog'] }}
                        </a>
                    </span>
                    <span *ngSwitchCase="'poweredOff'">{{ activeVm['dockerLog'] }}</span>
                </td>
            </tr>
        </tbody>
    </table>
    `
})
export class VchPortletComponent implements OnInit {
    @Input() activeVm: VirtualMachine;

    public dockerApiEndpointLabel: string;
    public vchAdminPortalLabel: string;

    constructor(
        private i18n: I18nService,
        private globalsService: GlobalsService
    ) { }

    /**
     * Evaluates if the application is running in the vSphere Client environment,
     * which is called the "Plugin Mode"
     * @returns true if plugin mode. false if not
     */
    isPluginMode() {
        return this.globalsService.isPluginMode() &&
            !(<any>this.activeVm)['is_testing'];
    }

    ngOnInit() {
        this.dockerApiEndpointLabel = this.isPluginMode() ?
            this.i18n.translate('vch.dockerApiEndpoint.label') :
            DEFAULT_VCH_DOCKER_ENDPOINT_LABEL;
        this.vchAdminPortalLabel = this.isPluginMode() ?
            this.i18n.translate('vch.vchAdminPortal.label') :
            DEFAULT_VCH_VCH_ADMIN_PORTAL_LABEL;
    }
}
