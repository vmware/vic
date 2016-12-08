/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.view.LoginView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;

/**
 * Common workflow step for logging in the client login screen.
 */
public class LoginStep extends CommonUIWorkflowStep {

   private static final LoginView _loginView = new LoginView();
   private UserSpec _userSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Login Step");
      }

      _userSpec = getSpec().links.get(UserSpec.class);

      if (_userSpec == null) {
         throw new IllegalArgumentException("The required UserSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_userSpec.username.get())) {
         throw new IllegalArgumentException("The user name is not set.");
      }
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      // TODO: The view doesn't have to be aware of the UserSpec.
      // Remove the dependency.
      _loginView.login(_userSpec);
   }

   /**
    * {@inheritDoc}
    * @throws Exception
    */
   @Override
   public void clean() throws Exception {
      _loginView.logout();
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Login Step");
      }

      _userSpec = filteredWorkflowSpec.links.get(UserSpec.class);

      if (_userSpec == null) {
         throw new IllegalArgumentException("The required UserSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_userSpec.username.get())) {
         throw new IllegalArgumentException("The user name is not set.");
      }

      if (Strings.isNullOrEmpty(_userSpec.password.get())) {
         throw new IllegalArgumentException("The user password is not set.");
      }
   }
}
