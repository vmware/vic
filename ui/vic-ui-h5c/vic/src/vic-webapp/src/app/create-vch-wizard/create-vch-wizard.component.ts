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

import { Component, OnInit, ViewChild, ElementRef, Renderer } from '@angular/core';
import { Wizard } from 'clarity-angular';
import { GlobalsService } from 'app/shared';

@Component({
  selector: 'vic-create-vch-wizard',
  templateUrl: './create-vch-wizard.component.html',
  styleUrls: ['./create-vch-wizard.component.scss']
})
export class CreateVchWizardComponent implements OnInit {
  @ViewChild('wizardlg') wizard: Wizard;
  public loading = false;
  public errorFlag = false;
  public errorMsg: string;

  // TODO: remove the following
  public testVal = 0;

  constructor(
    private elRef: ElementRef,
    private renderer: Renderer,
    private globalsService: GlobalsService
  ) { }

  /**
   * Launch the wizard
   */
  ngOnInit() {
    this.wizard.open();
  }

  /**
   * Resize the parent modal where the inline wizard is instantiated such that
   * the wizard fits exactly in the modal. This method is called once the
   * clrWizardOpenChange event is fired. While inlined wizard is a pattern
   * not recommended by the Clarity team, this is the only way to interface
   * with the H5 Client properly through WEB_PLATFORM.openModalDialog()
   */
  resizeToParentFrame(p: Window = parent) {
    // "context error" warning shows up during unit tests (but they still pass).
    // this can be avoided by running the logic a tick later
    setTimeout(() => {
      const parentIframes = p.document.querySelectorAll('iframe');
      const targetIframeEl = parentIframes[parentIframes.length - 1];
      const activeModalContentEl = <HTMLElement>p.document.querySelector('clr-modal .modal-content');
      const activeModalHeaderEl = <HTMLElement>p.document.querySelector('clr-modal .modal-header');
      // resize only if the parent modal is there. this prevents the unit tests from failing
      if (activeModalContentEl !== null) {
        let targetIframeHeight = activeModalContentEl.offsetHeight - 2;
        if (activeModalHeaderEl !== null) {
          targetIframeHeight -= activeModalHeaderEl.offsetHeight;
          activeModalHeaderEl.remove();
        }

        this.renderer.setElementStyle(targetIframeEl, 'height', `${targetIframeHeight}px`);
        this.renderer.setElementStyle(
          this.elRef.nativeElement.querySelector('clr-wizard'),
          'height',
          `${targetIframeHeight}px`
        );
      }
    });
  }

  /**
   * Perform validation for the current WizardPage and proceed to the next page
   * if the data is valid. If not, display an error message
   * @param id : ID of WizardPage
   */
  onCommit(id: string) {
    this.loading = true;
    // TODO: validation model & logic
    this.testVal = Math.random();
    console.log('on commit', id);
    if (this.testVal >= 0.5) {
      this.wizard.forceNext();
    } else {
      this.errorFlag = true;
      this.errorMsg = `${this.testVal} is smaller than 0.5! Try again.`;
    }
  }

  /**
   * Go back to the previous step
   */
  goBack() {
    this.wizard.previous();
  }

  /**
   * Clear the error flag (this method might be removed)
   * @param id : ID of the current WizardPage
   */
  onPageLoad(id: string) {
    this.errorFlag = false;
  }

  /**
   * Perform the final data validation and send the data to the
   * OVA endpoint via a POST request
   */
  onFinish() {
    // TODO: send the results to the OVA endpoint via a POST request
    if (!this.loading && this.areAllUserInputsValid) {
      this.wizard.forceFinish();
      this.onCancel();
      return;
    }

    this.errorFlag = true;
    this.errorMsg = 'User inputs validation failed!';
  }

  /**
   * Close the H5 Client modal
   */
  onCancel() {
    const webPlatform = this.globalsService.getWebPlatform();
    webPlatform.closeDialog();
  }

  /**
   * Readonly attribute to determine whether the wizard
   * can be finished by validating user inputs
   */
  get areAllUserInputsValid(): boolean {
    // TODO: TBI
    return false;
  }
}
