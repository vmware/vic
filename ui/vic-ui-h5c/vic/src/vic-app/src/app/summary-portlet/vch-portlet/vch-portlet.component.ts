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

import { Component, Input, OnInit } from '@angular/core';
import { VirtualMachine } from '../../vm.interface';

@Component({
    selector: 'vic-vch-portlet',
    styleUrls: ['./vch-portlet.scss'],
    template: `
    <table class="table table-vertical">
        <tbody>
            <tr>
                <th>{{ dockerApiEndpointLabel }}</th>
                <td>{{ activeVm.dockerEndpoint }}</td>
            </tr>
            <tr>
                <th>{{ vchAdminPortalLabel }}</th>
                <td [ngSwitch]="activeVm.powerState">
                    <span *ngSwitchCase="'poweredOn'"><a href="{{ activeVm.dockerLog }}" target="_blank">{{ activeVm.dockerLog }}</a></span>
                    <span *ngSwitchCase="'poweredOff'">{{ activeVm.dockerLog }}</span>
                </td>
            </tr>
        </tbody>
    </table>
    `
})
export class VchPortletComponent implements OnInit {
    @Input() activeVm: VirtualMachine;

    private dockerApiEndpointLabel: string;
    private vchAdminPortalLabel: string;

    constructor() {}

    ngOnInit() {
        this.dockerApiEndpointLabel = window.hasOwnProperty('com_vmware_vic')
            ? com_vmware_vic.getString('vch.dockerApiEndpoint.label') :
            'Docker API endpoint';
        this.vchAdminPortalLabel = window.hasOwnProperty('com_vmware_vic') ?
            com_vmware_vic.getString('vch.vchAdminPortal.label') :
            'VCH Admin portal';
    }
}
