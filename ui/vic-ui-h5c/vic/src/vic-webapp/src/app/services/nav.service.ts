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

import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { GlobalsService, APP_CONFIG } from '../shared/index';
import { extensionToRoutes } from '../app-routing.module';

export enum ObjectViewType {
    summary,
    monitor,
    manage
}

/**
 * Navigation service for jumping to another view
 */
@Injectable()
export class NavService {
    // Keep track of the selected view type.
    // Default value used the first time an object is selected.
    viewType: ObjectViewType = ObjectViewType.monitor;

    private navigate(extensionId: string, objectId: string = null): void {
        if (this.gs.isPluginMode()) {
            this.gs.getWebPlatform().sendNavigationRequest(extensionId, objectId);
        } else if (objectId) {
            this.router.navigate([extensionToRoutes[extensionId], objectId]);
        } else {
            this.router.navigate([extensionToRoutes[extensionId]]);
        }
    }

    constructor(private gs: GlobalsService,
        private router: Router) {
    }

    showMainView(): void {
        this.navigate(APP_CONFIG.packageName + '.mainView');
    }

    showSettingsView(): void {
        this.navigate(APP_CONFIG.packageName + '.settingsView');
    }

    /**
     * Navigate to the view of giving type for given object id
     * @param id
     * @param type (optional) an ObjectViewType enum, or the corresponding name.
     *             or re-use the current view type if no type is given.
     */
    showObjectView(id: string, type: ObjectViewType | string = this.viewType): void {
        this.setViewType(type);

        // The view extension name ends with .host.summaryView, .host.manageView, .host.monitorView
        const viewExtension = APP_CONFIG.packageName + '.host.' + this.getViewType() + 'View';
        this.navigate(viewExtension, id);
    }

    /**
     * Keep track of the selected view type, i.e. summary, monitor or manage
     * @param type an ObjectViewType enum, or the corresponding name
     */
    setViewType(type: ObjectViewType | string): void {
        if (typeof type === 'string') {
            if (typeof ObjectViewType[type] === 'undefined') {
                throw new Error('Invalid view type: ' + type);
            }
            this.viewType = ObjectViewType[type];
        } else {
            this.viewType = type;
        }
    }

    /**
     * @returns the name of the current view type
     */
    getViewType(): string {
        return ObjectViewType[this.viewType];
    }
}
