/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.verification.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.common.verification.spec.WizardPageVerificationSpec;

/**
 * Test step that performs basic UI verifications.
 * Only these verifications are performed for which spec properties are set.
 * Mandatory verifications are - TIWO, help and cancel buttons availability.
 * This step expects that wizard step is opened and loaded.
 *
 * Performed verifications:
 * 1. Verifies left navigation page title
 * 2. Verifies page header
 * 3. Verifies page header description
 * 4. Verifies that TIWO minimize button is present and enabled
 * 5. Verifies that help button is present and enabled
 * 6. Verifies that back button is present and enabled/disabled
 * 7. Verifies that next button is present and enabled/disabled
 * 8. Verifies that finish button is present and enabled/disabled
 * 9. Verifies that cancel button is present and enabled
 */
public class WizardPageBasicUiVerificationStep extends BaseWorkflowStep {

   WizardPageVerificationSpec _verificationSpec;

   @Override
   public void prepare() throws Exception {
      _verificationSpec = getSpec().links.get(WizardPageVerificationSpec.class);

      if (_verificationSpec == null) {
         throw new IllegalArgumentException(
               "No verification spec is linked to the test spec!"
            );
      }
   }

   @Override
   public void execute() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();

      //TODO: implement left navigation title verification
      //TODO check the title, the step number, and the current state

      //verify page header
      if (_verificationSpec.pageHeader.isAssigned()) {
         verifyFatal(
               getTestScope(),
               _verificationSpec.pageHeader.get().equals(wizardNavigator.getPageTitle()),
               "Verifying that wizard page header matches the expected one."
               );
      }

      //verify page header description
      if (_verificationSpec.pageHeaderDesription.isAssigned()) {
         String pageHeaderDesription = wizardNavigator.getPageHeaderDescription();
         verifyFatal(
               getTestScope(),
               _verificationSpec.pageHeaderDesription.get().equals(pageHeaderDesription),
               "Verifying that wizard page header description matches the expected one."
               );
      }

      //verify that TIWO minimize button is enabled
      verifyFatal(
            getTestScope(),
            wizardNavigator.isMinimizeButonEnabled(),
            "Verifying that TIWO minimize button is enabled."
            );

      //verify that help button is enabled
      verifyFatal(
            getTestScope(),
            wizardNavigator.isHelpButtonEnabled(),
            "Verifying that help button is enabled."
            );

      //verify if back button is enabled/disabled
      if (_verificationSpec.backButtonEnabled.isAssigned()) {
         verifyFatal(
               getTestScope(),
               wizardNavigator.isPrevBtnEnabled() == _verificationSpec.backButtonEnabled.get(),
               String.format(
                     "Verifying if back button is %s.",
                     _verificationSpec.backButtonEnabled.get() ? "enabled" : "disbled"
                  )
               );
      }

      //verify if next button is enabled/disabled
      if (_verificationSpec.nextButtonEnabled.isAssigned()) {
         verifyFatal(
               getTestScope(),
               wizardNavigator.isNextBtnEnabled() == _verificationSpec.nextButtonEnabled.get(),
               String.format(
                     "Verifying if next button is %s.",
                     _verificationSpec.nextButtonEnabled.get() ? "enabled" : "disbled"
                     )
               );
      }

      //verify if finish button is enabled/disabled
      if (_verificationSpec.finishButtonEnabled.isAssigned()) {
         verifyFatal(
               getTestScope(),
               wizardNavigator.isFinishBtnEnabled() == _verificationSpec.finishButtonEnabled.get(),
               String.format(
                     "Verifying if finish button is %s.",
                     _verificationSpec.finishButtonEnabled.get() ? "enabled" : "disbled"
                     )
               );
      }

      //verify if cancel button is enabled
      verifyFatal(
            getTestScope(),
            wizardNavigator.isCancelBtnEnabled(),
            "Verifying if cancel button is enabled."
            );
   }
}
