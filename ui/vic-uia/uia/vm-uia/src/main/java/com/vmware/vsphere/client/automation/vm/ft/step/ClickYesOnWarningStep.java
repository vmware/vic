/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.vsphere.client.automation.vm.ft.view.FtWarningPage;

/**
 * Click Yes in the warning dialog before the Fault Tolerance wizard
 */
public class ClickYesOnWarningStep extends CommonUIWorkflowStep {

   @Override
   public void execute() throws Exception {
      final FtWarningPage warningPage = new FtWarningPage();

      warningPage.clickYes();
   }
}