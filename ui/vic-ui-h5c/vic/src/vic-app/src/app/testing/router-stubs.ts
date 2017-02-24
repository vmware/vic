import { NavigationExtras } from "@angular/router";

import { BehaviorSubject } from "rxjs/BehaviorSubject";

export class ActivatedRouteStub {

   // ActivatedRoute.params is Observable
   private subject = new BehaviorSubject(this.testParams);
   params = this.subject.asObservable();

   // Test parameters
   private _testParams: {};
   get testParams() {
      return this._testParams;
   }
   set testParams(params: {}) {
      this._testParams = params;
      this.subject.next(params);
   }
}

export class RouterStub {
   navigate(commands: any[], extras?: NavigationExtras) { }
}
