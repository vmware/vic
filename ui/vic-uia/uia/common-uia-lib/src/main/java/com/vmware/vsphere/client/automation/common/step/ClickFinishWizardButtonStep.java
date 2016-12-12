/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;

/**
 * Common workflow step to complete a wizard, by clicking the Finish button.
 *
 * Operations performed in this step:
 * - Click the Finish wizard button and wait for it to disappear, as indication that the wizard is closed
 * - Wait for the most recent task to complete
 */
public class ClickFinishWizardButtonStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.FULL, new WizardNavigator().finishWizard(),
            "The Finish wizard button is successfully clicked.");
      new BaseView().waitForRecentTaskCompletion();
   }
}
