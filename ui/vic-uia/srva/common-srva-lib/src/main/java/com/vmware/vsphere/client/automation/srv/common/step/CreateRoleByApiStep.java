/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotEmpty;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.ArrayList;
import java.util.List;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.RoleSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow step for creation a custom Role
 */
public class CreateRoleByApiStep extends BaseWorkflowStep {

   private List<RoleSpec> _rolesToCreate;
   private List<RoleSpec> _rolesToClean;

   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create Roles");
      }
      _rolesToCreate = getSpec().links.getAll(RoleSpec.class);

      ensureNotNull(_rolesToCreate,
            "The spec has no links to 'RoleSpec' instances");
      ensureNotEmpty(_rolesToCreate, "The list with 'RoleSpec' is empty");

      _rolesToClean = new ArrayList<RoleSpec>();
   }

   public void prepare(WorkflowSpec filterWorkflowSpec) throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create Roles");
      }
      _rolesToCreate = filterWorkflowSpec.getAll(RoleSpec.class);

      ensureNotNull(_rolesToCreate,
            "The spec has no links to 'RoleSpec' instances");
      ensureNotEmpty(_rolesToCreate, "The list with 'RoleSpec' is empty");

      _rolesToClean = new ArrayList<RoleSpec>();
   }

   @Override
   public void execute() throws Exception {
      boolean result = false;
      for (RoleSpec roleSpec : _rolesToCreate) {
         if (roleSpec.roleId.isAssigned()) {
            // We don't need to create roles which are already present
            continue;
         }

         if (roleSpec.exceptPrivilegeIds.isAssigned()) {
            result = UserBasicSrvApi.getInstance()
                  .createRoleAllExceptSpecifiedPrivileges(
                        roleSpec,
                        roleSpec.exceptPrivilegeIds.getAll().toArray(
                              new String[] {}));
         } else {
            result = UserBasicSrvApi.getInstance().createRole(roleSpec);
         }

         if (result) {
            _rolesToClean.add(roleSpec);
         } else {
            throw new Exception(String.format("Unable to create role '%s'",
                  roleSpec.name.get()));
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (RoleSpec roleSpec : _rolesToClean) {
         UserBasicSrvApi.getInstance().deleteRole(roleSpec);
      }
   }

}
