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
import { Headers } from '@angular/http';

import { APP_CONFIG } from './app-config';
import { WebPlatform } from './vSphereClientSdkTypes';
import { webPlatformStub } from './dev/webPlatformStub';

/**
 * Internal flags used by GlobalsService.  Do not import or modify!
 * Use the GlobalsService API below instead.
 * Use app-config.ts for application configuration
 */
export class Globals {
      readonly pluginMode: boolean = true;

      // The WEB_PLATFORM global object giving access to SDK 6.5  APIs
      readonly webPlatform: WebPlatform;

      constructor() {
            // pluginMode is set to true when the app runs inside an iFrame, i.e. inside the vSphere Client
            // (this assumes the app itself doesn't use iFrames, which is the normal case)
            this.pluginMode = (window.self !== window.parent);

            // WEB_PLATFORM global is defined on the vSphere Client window, except for older versions and dev mode
            this.webPlatform = window.parent['WEB_PLATFORM'] ||
                  (this.pluginMode ? this.getOldVersionPlatform() : webPlatformStub);

            // Patch for setGlobalRefreshHandler API added in 6.5
            if (this.webPlatform.getClientType() === 'flex') {
                  this.webPlatform.setGlobalRefreshHandler = function (handler) {
                        window.parent['WEB_PLATFORM.refresh' + window.name] = handler;
                  };
            }
      }

      // Case of Web Client version < 6.0u3 where the WEB_PLATFORM global didn't exist yet
      getOldVersionPlatform(): WebPlatform {
            let webPlatformElement: any = window.parent.document.getElementById('container_app');
            if (window['jasmine'] && !webPlatformElement) {
                  // Dummy element for jasmine tests
                  webPlatformElement = {};
            }
            if (window.navigator.userAgent.indexOf('Chrome/') >= 0) {
                  // Chrome browser needs a real JS object
                  webPlatformElement = Object.create(webPlatformElement);
            }
            webPlatformElement.getRootPath = function () { return '/vsphere-client'; };
            webPlatformElement.getClientType = function () { return 'flex'; };
            webPlatformElement.getClientVersion = function () { return '6.0'; };
            return webPlatformElement;
      }
}

/**
 *  A service dealing with application globals which facilitate the mixed-mode
 *  approach of plugin development: dev mode (standalone app) and plugin mode.
 */
@Injectable()
export class GlobalsService {
      private liveData = false;
      private sidenav = false;
      private viewInfo = true;
      private devUI = true;
      private clientId: string;
      private readonly webContextPath: string;

      // The locale id of the vSphere Client, or 'en' by default in standalone
      public locale = 'en';

      constructor(private globals: Globals) {
            if (this.isPluginMode()) {
                  this.liveData = true;
                  this.sidenav = false;
                  this.devUI = false;

            } else {
                  // Dev mode defaults: re-use what was last saved in browser local storage
                  this.liveData = (localStorage.getItem('liveData') === 'true' ? true : false);
                  this.devUI = (localStorage.getItem('devUI') === 'false' ? false : true);
                  this.sidenav = (localStorage.getItem('sidenav') === 'true' ? true : false);
                  this.viewInfo = (localStorage.getItem('viewInfo') === 'true' ? true : false);
                  this.clientId = localStorage.getItem('clientId');
            }
            this.webContextPath = this.getWebPlatform().getRootPath() + '/' + APP_CONFIG.pluginName;
      }

      /**
       * @return true when the app runs as a plugin inside vSphere Client,
       *       false when it runs standalone during development.
       */
      public isPluginMode(): boolean {
            return this.globals.pluginMode;
      }

      /**
       * @return the WEB_PLAFORM object containing SDK 6.5 APIs
       */
      public getWebPlatform(): WebPlatform {
            return this.globals.webPlatform;
      }

      /**
       * @return the context path for this plugin object.
       */
      public getWebContextPath(): string {
            return this.webContextPath;
      }

      /**
       * @return true when pluginMode is true, or when testing live data in dev mode
       */
      public useLiveData(): boolean {
            return this.liveData;
      }

      public toggleLiveData(): void {
            if (!this.isPluginMode()) {
                  this.liveData = !this.liveData;
                  localStorage.setItem('liveData', '' + this.liveData);
            }
      }

      /**
       * @return true to show all dev UI elements, false to hide them.
       *       Always return false in plugin mode.
       */
      public showDevUI(): boolean {
            return this.devUI;
      }

      public toggleDevUI(): void {
            if (!this.isPluginMode()) {
                  this.devUI = !this.devUI;
                  localStorage.setItem('devUI', '' + this.devUI);
            }
      }

      /**
       * @return true to show the sidenav component in dev mode.
       *          Always return false in plugin mode, or when showDevUI is false.
       */
      public showSidenav(): boolean {
            return this.sidenav && this.devUI;
      }

      public toggleSidenav(): void {
            if (!this.isPluginMode()) {
                  this.sidenav = !this.sidenav;
                  localStorage.setItem('sidenav', '' + this.sidenav);
            }
      }

      /**
       * @return true to show the viewInfo component in standalone mode.
       */
      public showViewInfo(): boolean {
            return this.viewInfo;
      }

      public toggleViewInfo(): void {
            if (!this.isPluginMode()) {
                  this.viewInfo = !this.viewInfo;
                  localStorage.setItem('viewInfo', '' + this.viewInfo);
            }
      }

      /**
       *  Sets the vSphere Client session id used for accessing live data in dev mode
       */
      public setClientId(id): void {
            if (!this.isPluginMode()) {
                  this.clientId = id;
                  localStorage.setItem('clientId', this.clientId);
            }
      }
      public getClientId(): string {
            return this.clientId;
      }

      /**
       * Flag for checking both the dev mode or a local environment
       *
       * @returns true if the app is running standalone, or as a plugin in a local environment
       */
      isLocalhostDevMode(): boolean {
            return !this.isPluginMode() ||
                  window['self'].parent.location.href.indexOf('https://localhost:9443/') === 0;
      }

      /**
       * Get the necessary headers in case of dev mode + live data mode
       *
       * @returns an empty object if live data is turned off, or in plugin mode. Else this returns
       * http headers with the current clientId in order to authorize requests to the local Virgo server.
       * This is a simple way to access live data in dev mode which relies on a valid vSphere client session
       * (i.e. vSphere object data for instance)
       *
       * - clientId must have been initialized beforehand (see clientid.component)
       * - this function is harmless in pluginMode because it returns an empty header.
       */
      getHttpHeaders(): any {
            if (!this.useLiveData() || this.isPluginMode()) {
                  return {};
            }
            return { headers: new Headers({ 'webClientSessionId': this.getClientId() }) };
      }
}
