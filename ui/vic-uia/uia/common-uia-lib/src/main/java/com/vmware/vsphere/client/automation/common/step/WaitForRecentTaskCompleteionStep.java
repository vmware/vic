/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Waits for a recent task's completion status
 */
public class WaitForRecentTaskCompleteionStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new BaseView().waitForRecentTaskCompletion();
   }
}
