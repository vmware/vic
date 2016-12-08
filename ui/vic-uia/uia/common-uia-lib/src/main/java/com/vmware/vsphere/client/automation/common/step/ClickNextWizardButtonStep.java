/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;

/**
 * Common workflow step to navigate to the next page of a wizard, by clicking
 * the Next button.
 *
 * Operations performed in this step:
 * - Click the Next wizard button
 * - Wait for Loading progress bar to disappear
 */
public class ClickNextWizardButtonStep extends CommonUIWorkflowStep {

   private static final WizardNavigator _wizardNavigator = new WizardNavigator();

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.FULL, _wizardNavigator.gotoNextPage(),
            "The Next wizard button is successfully clicked.");
      _wizardNavigator.waitForDialogToLoad();
   }
}
