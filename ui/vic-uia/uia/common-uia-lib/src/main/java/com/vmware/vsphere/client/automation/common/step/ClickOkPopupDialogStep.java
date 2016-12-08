/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.PopupDialogNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Clicks the OK button of a popup dialog.
 */
public class ClickOkPopupDialogStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new PopupDialogNavigator().clickOk();
      new BaseView().waitForRecentTaskCompletion();
   }
}
