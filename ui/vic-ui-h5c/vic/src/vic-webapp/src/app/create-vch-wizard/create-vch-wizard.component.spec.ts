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
import { By } from '@angular/platform-browser';
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
      providers: [
        {
          provide: GlobalsService, useValue: {
            getWebPlatform: () => {
              return {
                closeDialog: () => { }
              };
            }
          }
        }],
      declarations: [CreateVchWizardComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateVchWizardComponent);
    component = fixture.componentInstance;
    spyOn(component, 'resizeToParentFrame').and.callThrough();
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have called wizard.open', async(() => {
    expect(component.resizeToParentFrame).toHaveBeenCalled();
  }));

  it('should navigate between pages successfully', async(() => {
    const nextBtn = fixture.debugElement.query(By.css('clr-wizard-button[ng-reflect-type="next"]'));
    const firstPage = fixture.debugElement.query(By.css('clr-wizard-page'));
    const lastPage = fixture.debugElement.query(By.css('clr-wizard-page:last-of-type'));
    expect(nextBtn).toBeTruthy();
    spyOn(component, 'onCancel').and.callThrough();
    spyOn(component, 'onCommit').and.callThrough();
    spyOn(component, 'onFinish').and.callThrough();
    spyOn(component, 'goBack').and.callThrough();

    // click on next button
    firstPage.triggerEventHandler('clrWizardPageOnCommit', null);
    expect(component.onCommit).toHaveBeenCalled();

    // click on finish button
    lastPage.triggerEventHandler('clrWizardPageOnCommit', null);
    expect(component.onFinish).toHaveBeenCalled();

    // click on previous button
    lastPage.triggerEventHandler('clrWizardPagePrevious', null);
    expect(component.goBack).toHaveBeenCalled();

    // click on cancel button
    lastPage.triggerEventHandler('clrWizardPageOnCancel', null);
    expect(component.onCancel).toHaveBeenCalled();
  }));
});
