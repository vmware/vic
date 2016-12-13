/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Common workflow step that refreshes the general NGC screen.
 * This step will fail if the refresh does not finish on time(configurable).
 */
public class GlobalRefreshStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      BaseView view = new BaseView();

      view.refreshPage();

      //TODO: tmp fix for Work in Progress list refresh - RP 1216314
      view.refreshPage();
   }
}
