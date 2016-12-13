/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Refreshing the UI.
 */
public class RefreshUIStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new BaseView().refreshPage();
   }

}
