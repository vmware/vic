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

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { APP_CONFIG } from './shared/index';
import { AppRoutingComponent } from './app-routing.component';

const appRoutes: Routes = [
    { path: 'index.html', component: AppRoutingComponent },
    { path: 'portlet', loadChildren: './summary-portlet/summary-portlet.module#VicSummaryPortletModule' },
    { path: 'summary-view', loadChildren: './summary-view/summary-view.module#VicSummaryViewModule' },
    { path: 'vch-view', loadChildren: './vch-view/vch-view.module#VicVchViewModule' },
    { path: 'container-view', loadChildren: './container-view/container-view.module#VicContainerViewModule' },
    { path: 'create-vch', loadChildren: './create-vch-wizard/create-vch-wizard.module#CreateVchWizardModule' }
];

export const extensionToRoutes = {};
extensionToRoutes[APP_CONFIG.packageName + '.objectView.summaryView'] = '/objectViewSummary';
extensionToRoutes[APP_CONFIG.packageName + '.objectView.monitorView'] = '/objectViewMonitor';
extensionToRoutes[APP_CONFIG.packageName + '.objectView.vchView'] = '/objectViewVch';
extensionToRoutes[APP_CONFIG.packageName + '.objectView.containerView'] = '/objectViewContainer';

export const routedComponents = [AppRoutingComponent];

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule { }
