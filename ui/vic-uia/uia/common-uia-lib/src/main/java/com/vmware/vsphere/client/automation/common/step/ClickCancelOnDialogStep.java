/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.BaseDialogNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Common workflow step that click cancel on a dialog.
 */
public class ClickCancelOnDialogStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new BaseView().waitForPageToRefresh();
      new BaseDialogNavigator().cancel();
      new BaseView().waitForPageToRefresh();
   }
}
