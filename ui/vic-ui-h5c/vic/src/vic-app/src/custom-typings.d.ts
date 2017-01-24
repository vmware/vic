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

declare var ENV: string;
declare var HMR: boolean;
declare var WEB_PLATFORM: any;
declare var com_vmware_vic: any;
declare var customElements: any;

interface GlobalEnvironment {
    ENV;
    HMR;
}

interface WebpackModule {
    hot: {
        data?: any,
        idle: any,
        accept(dependencies?: string | string[], callback?: (updatedDependencies?: any) => void): void;
        decline(dependencies?: string | string[]): void;
        dispose(callback?: (data?: any) => void): void;
        addDisposeHandler(callback?: (data?: any) => void): void;
        removeDisposeHandler(callback?: (data?: any) => void): void;
        check(autoApply?: any, callback?: (err?: Error, outdatedModules?: any[]) => void): void;
        apply(options?: any, callback?: (err?: Error, outdatedModules?: any[]) => void): void;
        status(callback?: (status?: string) => void): void | string;
        removeStatusHandler(callback?: (status?: string) => void): void;
    };
}

interface WebpackRequire {
    context(file: string, flag?: boolean, exp?: RegExp): any;
}


interface ErrorStackTraceLimit {
    stackTraceLimit: number;
}



// Extend typings
interface NodeRequire extends WebpackRequire {}
interface ErrorConstructor extends ErrorStackTraceLimit {}
interface NodeModule extends WebpackModule {}
interface Global extends GlobalEnvironment {}
