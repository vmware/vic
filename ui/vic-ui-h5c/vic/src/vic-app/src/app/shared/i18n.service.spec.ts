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

import { I18nService } from "./i18n.service";
import { JASMINE_TIMEOUT } from "../testing/jasmine.constants";

let i18nService: I18nService;

// Simple service unit tests without assistance from Angular testing utilities
// Note: we don't really need to test initLocale() which is only for dev mode

describe("I18Service tests", () => {
   jasmine.DEFAULT_TIMEOUT_INTERVAL = JASMINE_TIMEOUT;
   beforeEach(() => {
      i18nService = new I18nService(null, null, null);
   });
   it ("interpolates messages correctly", () => {
      let msg1 = "message1 without params";
      let msg2 = "message2 with {0}";
      let msg4 = "message4 with {0} {1}";
      let result1 = i18nService.interpolate(msg1, null);
      let result2 = i18nService.interpolate(msg2, "param");
      let result3 = i18nService.interpolate(msg2, ["param"]);
      let result4 = i18nService.interpolate(msg4, ["param"]);


      expect(result1).toBe(msg1);
      expect(result2).toBe("message2 with param");
      expect(result3).toBe("message2 with param");
      expect(result4).toBe("message4 with param {1}");
   });
});
