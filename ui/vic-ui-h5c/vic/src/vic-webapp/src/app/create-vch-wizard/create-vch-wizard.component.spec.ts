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

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ClarityModule } from 'clarity-angular';
import { CreateVchWizardComponent } from './create-vch-wizard.component';
import { Globals, GlobalsService } from 'app/shared';

describe('CreateVchWizardComponent', () => {
  let component: CreateVchWizardComponent;
  let fixture: ComponentFixture<CreateVchWizardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ClarityModule, BrowserAnimationsModule],
      providers: [GlobalsService, Globals],
      declarations: [CreateVchWizardComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateVchWizardComponent);
    component = fixture.componentInstance;
    spyOn(component, 'resizeToParentFrame').and.callFake(() => {
      console.log('resizeToParentFrame called');
    });
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have called wizard.open', async(() => {
    expect(component.resizeToParentFrame).toHaveBeenCalled();
  }));
});
