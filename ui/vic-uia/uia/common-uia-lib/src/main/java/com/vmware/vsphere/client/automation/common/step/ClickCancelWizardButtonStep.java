/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Common workflow step to cancel a wizard.
 *
 * Operations performed in this step:
 * - Click the Cancel wizard button
 * - Wait for Loading progress bar to disappear
 */
public class ClickCancelWizardButtonStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new WizardNavigator().cancel();
   }
}
