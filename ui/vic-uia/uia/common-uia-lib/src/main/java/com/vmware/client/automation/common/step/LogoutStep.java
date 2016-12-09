/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import com.vmware.client.automation.common.view.LoginView;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;

/**
 * General workflow step that logs out the user.
 */
public class LogoutStep extends CommonUIWorkflowStep {

   private static final LoginView _loginView = new LoginView();

   @Override
   public void execute() throws Exception {
      _loginView.logout();
   }
}
