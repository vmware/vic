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

import { Component, OnInit, Input } from '@angular/core';
import { VirtualMachine } from '../../vm.interface';

@Component({
    selector: 'vic-container-portlet',
    styleUrls: ['./container-portlet.scss'],
    template: `
    <table class="table table-vertical">
        <tbody>
            <tr>
                <th>{{ containerNameLabel }}</th>
                <td>{{ activeVm.container_name }}</td>
            </tr>
            <tr>
                <th>{{ containerImageLabel }}</th>
                <td>{{ activeVm.image_name }}</td>
            </tr>
            <tr *ngIf="activeVm.portmapping">
                <th>{{ containerPortmappingLabel }}</th>
                <td>{{ activeVm.portmapping }}</td>
            </tr>
        </tbody>
    </table>
    `
})
export class ContainerPortletComponent implements OnInit {
    @Input() activeVm: VirtualMachine;

    private containerNameLabel: string;
    private containerImageLabel: string;
    private containerPortmappingLabel: string;

    constructor() {}

    ngOnInit() {
        this.containerNameLabel = window.hasOwnProperty('com_vmware_vic') ?
            com_vmware_vic.getString('container.name.label') :
            'Container';
        this.containerImageLabel = window.hasOwnProperty('com_vmware_vic') ?
            com_vmware_vic.getString('container.image.label') :
            'Image';
        this.containerPortmappingLabel = window.hasOwnProperty('com_vmware_vic') ?
            com_vmware_vic.getString('container.portmapping.label') :
            'Port Mapping';
    }

}
