/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Clicks the OK button of a single page dialog.
 */
public class ClickOkSinglePageDialogStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new SinglePageDialogNavigator().clickOk();
      new BaseView().waitForRecentTaskCompletion();
   }
}
