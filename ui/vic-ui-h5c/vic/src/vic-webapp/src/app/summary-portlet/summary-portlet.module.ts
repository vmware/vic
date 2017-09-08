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
import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { VicSummaryPortletComponent } from './summary-portlet.component';
import { VchPortletComponent } from './vch-portlet/vch-portlet.component';
import { ContainerPortletComponent } from './container-portlet/container-portlet.component';
import { DataPropertyService } from '../services/data-property.service';

const routes: Routes = [
    { path: '', component: VicSummaryPortletComponent },
    { path: ':id', component: VicSummaryPortletComponent }
];

@NgModule({
    imports: [
        CommonModule,
        RouterModule.forChild(routes)
    ],
    declarations: [
        VicSummaryPortletComponent,
        VchPortletComponent,
        ContainerPortletComponent
    ],
    providers: [
        DataPropertyService
    ],
    exports: [
        VicSummaryPortletComponent,
        VchPortletComponent,
        ContainerPortletComponent
    ]
})
export class VicSummaryPortletModule { }
