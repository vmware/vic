/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * General workflow step that logs out the user.
 */
public class LogoutWithWaitStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      new BaseView().logoutWithWait();
   }
}
