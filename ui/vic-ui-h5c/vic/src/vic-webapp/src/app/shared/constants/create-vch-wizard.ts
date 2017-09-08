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
import { environment } from '../../../environments/environment';

export const VIC_ROOT_OBJECT_ID_WITH_NAME = 'urn:vic:vic:Root:vic%25252Fvic-root?properties=name';
export const CREATE_VCH_WIZARD_URL =
    `/ui/vic/resources/${environment.production ? 'dist' : 'build-dev'}/index.html?view=create-vch`;
export const WIZARD_MODAL_WIDTH = 920;
export const WIZARD_MODAL_HEIGHT = 600;
