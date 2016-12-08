/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Common workflow step that minimizes dialog to TIWO.
 */
public class MinimizeDialogToTiwoStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new WizardNavigator().minimize();
   }
}
