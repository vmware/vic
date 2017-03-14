import { TestBed } from "@angular/core/testing";

import { GlobalsService } from "../shared/index";
import { webPlatformStub } from "../shared/dev/webPlatformStub";

export * from "./router-stubs";

/**
 * Stub for testing in plugin more or dev mode
 */
export const globalStub = {
   pluginMode: true,
   webPlatform: webPlatformStub
};

/**
 * Initialization for unit tests
 */
export function initGlobalService(pluginMode: boolean): GlobalsService {

   globalStub.pluginMode = pluginMode;
   return TestBed.get(GlobalsService);
}

export const appErrorHandlerStub = {
   httpPromiseError(error: any): Promise<any> {
      return Promise.reject("error message from appErrorHandlerStub");
   }
};

// ---- Utilities copied from Angular2 doc ----

/**
 * Create custom DOM event the old fashioned way
 *
 * https://developer.mozilla.org/en-US/docs/Web/API/Event/initEvent
 * Although officially deprecated, some browsers (phantom) dont accept the preferred "new Event(eventName)"
 */
export function newEvent(eventName: string, bubbles = false, cancelable = false) {
   let evt = document.createEvent("CustomEvent");
   evt.initCustomEvent(eventName, bubbles, cancelable, null);
   return evt;
}
