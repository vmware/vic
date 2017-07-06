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

import { Component, OnInit, Input } from '@angular/core';
import { VirtualMachine } from '../../vm.interface';

import {
    DEFAULT_CONTAINER_NAME_LABEL,
    DEFAULT_CONTAINER_IMAGE_LABEL,
    DEFAULT_CONTAINER_PORT_MAPPING_LABEL
} from '../../shared/constants/index';

@Component({
    selector: 'vic-container-portlet',
    styleUrls: ['container-portlet.scss'],
    template: `
    <table class="table table-vertical">
        <tbody>
            <tr>
                <th>{{ containerNameLabel }}</th>
                <td>{{ activeVm['container_name'] }}</td>
            </tr>
            <tr>
                <th>{{ containerImageLabel }}</th>
                <td>{{ activeVm['image_name'] }}</td>
            </tr>
            <tr *ngIf="activeVm.hasOwnProperty('portmapping')">
                <th>{{ containerPortmappingLabel }}</th>
                <td>{{ activeVm['portmapping'] }}</td>
            </tr>
        </tbody>
    </table>
    `
})
export class ContainerPortletComponent implements OnInit {
    @Input() activeVm: VirtualMachine;

    public containerNameLabel: string;
    public containerImageLabel: string;
    public containerPortmappingLabel: string;

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
        this.containerNameLabel = this.isPluginMode() ?
            this.i18n.translate('container.name.label') :
            DEFAULT_CONTAINER_NAME_LABEL;

        this.containerImageLabel = this.isPluginMode() ?
            this.i18n.translate('container.image.label') :
            DEFAULT_CONTAINER_IMAGE_LABEL;

        this.containerPortmappingLabel = this.isPluginMode() ?
            this.i18n.translate('container.portmapping.label') :
            DEFAULT_CONTAINER_PORT_MAPPING_LABEL;
    }

}
