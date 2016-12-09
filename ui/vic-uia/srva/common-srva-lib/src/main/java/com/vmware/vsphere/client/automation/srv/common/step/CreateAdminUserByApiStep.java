/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.Collections;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow that created administrative user.
 *
 * Steps performed by this step:
 * * Create SSO user
 * * Add the user to Administrators group
 */
public class CreateAdminUserByApiStep extends BaseWorkflowStep {

   private static final String ADMINISTRATORS_GROUP = "Administrators";

   private UserSpec _userSpec;

   @Override
   public void prepare() throws Exception {
      _userSpec = getSpec().get(UserSpec.class);
      if (_userSpec == null) {
         throw new IllegalArgumentException("Create admin user step requires UserSpec.");
      }
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.BAT, UserBasicSrvApi.getInstance().createUser(_userSpec), "Creating SSO user");
      verifyFatal(
            TestScope.BAT,
            UserBasicSrvApi.getInstance().addUserToGroups(
                  _userSpec,
                  Collections.singletonList(ADMINISTRATORS_GROUP)),
            "Adding user to administrators group");
   }

   @Override
   public void clean() throws Exception {
      UserBasicSrvApi.getInstance().deleteUser(_userSpec);
   }


   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create administrator user!");
      }

      _userSpec = filteredWorkflowSpec.links.get(UserSpec.class);

      if (_userSpec == null) {
         throw new IllegalArgumentException("The required UserSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_userSpec.username.get())) {
         throw new IllegalArgumentException("The user name is not set.");
      }

   }

}
