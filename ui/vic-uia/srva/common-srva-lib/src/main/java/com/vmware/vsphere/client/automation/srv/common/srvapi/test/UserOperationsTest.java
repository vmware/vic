/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import java.util.List;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.model.PrivilegesCommonConstants;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.PermissionSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.RoleSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Test the User API operations exposed in UserSrvApi class. The test scenario
 * is: Create user, create role and assign permissions on datacenter
 */
public class UserOperationsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      getSpec().links.remove(UserSpec.class);
      UserSpec userSpec = SpecFactory.getUserSpec();
      getSpec().links.add(userSpec);

      RoleSpec roleSpec = new RoleSpec();
      roleSpec.name.set("UserOperationsTestRole_" + RandomStringUtils.randomAlphanumeric(5));
      roleSpec.privilegeIds.set(PrivilegesCommonConstants.Privileges.SYSTEM_READ.getId());

      // this is here before admin user spec is removed in order to load ssoUser and ssoPwd
      VcServiceSpec serviceSpec = new VcServiceSpec();
      serviceSpec.endpoint.set(testBed.getVc());
      serviceSpec.username.set(testBed.getAdminUser().username.get());
      serviceSpec.password.set(testBed.getAdminUser().password.get());

      PermissionSpec permissionsSpec = new PermissionSpec();
      permissionsSpec.name.set(userSpec.username.get());
      permissionsSpec.group.set(false);
      permissionsSpec.role.set(roleSpec);
      permissionsSpec.propagate.set(true);
      permissionsSpec.links.add(serviceSpec);

      // entity on which to set permissions
      DatacenterSpec dcSpec = SpecFactory.getSpec(DatacenterSpec.class,
            testBed.getCommonDatacenterName(), null);

      getSpec().links.add(roleSpec, permissionsSpec, dcSpec);

   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Nothing to do
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate create user operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            UserSpec user = getSpec().get(UserSpec.class);

            verifyFatal(TestScope.BAT, UserBasicSrvApi.getInstance().createUser(user), "Verifying the Create User Operation");
         }
      }, "Validating create user operation.");

      // Validate create role operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            RoleSpec roleSpec = getSpec().links.get(RoleSpec.class);

            verifyFatal(TestScope.BAT, UserBasicSrvApi.getInstance().createRole(roleSpec),
                  "Verify the create role operation");
         }
      }, "Validating create role operation.");

      // Validate set permissions operation

      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            List<PermissionSpec>  permissionsSpec = getSpec().links.getAll(PermissionSpec.class);
            ManagedEntitySpec entitySpec = getSpec().links.get(DatacenterSpec.class);

            PermissionSpec[] permissionSpecs = new PermissionSpec[permissionsSpec.size()];
            verifyFatal(
                  TestScope.FULL,
                  UserBasicSrvApi.getInstance().setEntityPermissions(entitySpec, permissionsSpec.toArray(permissionSpecs)),
                  "Setting entity permissions through API.");
         }
      }, "Validating set permissions on entity operation.");

      // Validate set global permissions operation

      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            List<PermissionSpec>  permissionsSpec = getSpec().links.getAll(PermissionSpec.class);

            PermissionSpec[] permissionSpecs = new PermissionSpec[permissionsSpec.size()];
            verifyFatal(
                  TestScope.FULL,
                  UserBasicSrvApi.getInstance().setGlobalPermissions( permissionsSpec.toArray(permissionSpecs)),
                  "Setting global permissions through API.");
         }
      }, "Validating set global permissions operation.");

      // Validate remove global permissions operation

      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            List<PermissionSpec>  permissionsSpec = getSpec().links.getAll(PermissionSpec.class);

            PermissionSpec[] permissionSpecs = new PermissionSpec[permissionsSpec.size()];
            UserBasicSrvApi.getInstance().deleteGlobalPermissions(permissionsSpec.toArray(permissionSpecs));
         }
      }, "Validating delete global permissions operation.");
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }

}
