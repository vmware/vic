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

import { NgModule, ErrorHandler } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpModule } from '@angular/http';

import { ClarityModule } from 'clarity-angular';

import {
    Globals,
    GlobalsService,
    I18nService,
    RefreshService,
    AppAlertComponent,
    AppAlertService
} from './shared/index';
import { ActionDevService } from './services/action-dev.service';
import { DataPropertyService } from './services/data-property.service';
import { AppErrorHandler } from './shared/appErrorHandler';

import { AppRoutingModule, routedComponents } from './app-routing.module';
import { VchPortletComponent } from './summary-portlet/vch-portlet/vch-portlet.component';
import { ContainerPortletComponent } from './summary-portlet/container-portlet/container-portlet.component';
import { AppComponent } from './app.component';

@NgModule({
    imports: [
        BrowserModule,
        AppRoutingModule,
        HttpModule,
        ClarityModule.forRoot()
    ],
    declarations: [
        AppComponent,
        AppAlertComponent,
        routedComponents,
        VchPortletComponent,
        ContainerPortletComponent
    ],
    providers: [
        ActionDevService,
        AppAlertService,
        AppErrorHandler,
        Globals,
        GlobalsService,
        I18nService,
        RefreshService,
        DataPropertyService,
        { provide: ErrorHandler, useClass: AppErrorHandler }
    ],
    bootstrap: [AppComponent]
})
export class AppModule {

}
