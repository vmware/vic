/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.MultiPageDialogNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Clicks the OK button of a multi page dialog and waits for the started task to complete
 */
public class ClickOkMultiPageDialogStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new MultiPageDialogNavigator().clickOk();
      new BaseView().waitForRecentTaskCompletion();
   }
}
