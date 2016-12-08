/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotBlank;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow that creates a user.
 *
 * Steps performed by this step:
 * * Create SSO user
 */
public class CreateUserByApiStep extends BaseWorkflowStep {

   private UserSpec _user;

   @Override
   public void prepare() throws Exception {

      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create administrator user!");
      }

      _user = getSpec().links.get(UserSpec.class);
      ensureNotNull(_user, "The required UserSpec is not set.");
      ensureNotBlank(_user.username, "The user name is not set.");

   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create administrator user!");
      }

      _user = filteredWorkflowSpec.links.get(UserSpec.class);
      ensureNotNull(_user, "The required UserSpec is not set.");
      ensureNotBlank(_user.username, "The user name is not set.");

   }

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.FULL, UserBasicSrvApi.getInstance().createUser(_user), "Creating SSO user");
   }

   @Override
   public void clean() throws Exception {
      UserBasicSrvApi.getInstance().deleteUser(_user);
   }

}
