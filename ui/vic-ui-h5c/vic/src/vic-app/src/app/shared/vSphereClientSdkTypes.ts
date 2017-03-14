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

/**
 * Interface with the vSphere Client SDK 6.5 Javascript APIs
 * Although not required we recommend to use this interface along with globals.services.ts
 * to take advantage of compiler type checking and code completion in your IDE
 */
export interface WebPlatform {
   callActionsController(url: string, jsonData: string, targets?: string): void;
   closeDialog(): void;
   getClientType(): string;
   getClientVersion(): string;
   getString(bundleName: string, key: string, params: any): string;
   getRootPath(): string;
   getUserSession(): UserSession;
   openModalDialog(title, url, width, height, objectId): void;
   sendModelChangeEvent(objectId, opType): void;
   sendNavigationRequest(targetViewId, objectId): void;
   setDialogSize(width, height): void;
   setDialogTitle(title): void;
   setGlobalRefreshHandler(callback, document): void;
}

export class UserSession {
   userName: string;
   clientId: string;
   locale: string;
   serversInfo: ServerInfo[];
}

export class ServerInfo {
   serviceGuid: string;
   serviceUrl: string;
   sessionId: string;
}
