/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotEmpty;

import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.PermissionSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow that sets permissions on provided entity
 */
public class SetPermissionsOnEntityByApiStep extends BaseWorkflowStep {

   private List<PermissionSpec> _permissionSpecToSet;
   private List<ManagedEntitySpec> _entitySpecs;

   @Override
   public void prepare(WorkflowSpec filterWorkflowSpec) throws Exception {
      _permissionSpecToSet = filterWorkflowSpec.getAll(PermissionSpec.class);
      ensureNotEmpty(_permissionSpecToSet,
            "Set Permissions Step requires PermissionSpec.");

      _entitySpecs = filterWorkflowSpec.getAll(ManagedEntitySpec.class);
      ensureNotEmpty(_entitySpecs, "Set Permissions Step requires EntitySpec.");
   }

   @Override
   public void execute() throws Exception {
      PermissionSpec[] permissionSpecs = new PermissionSpec[_permissionSpecToSet
            .size()];
      for (ManagedEntitySpec entitySpec : _entitySpecs) {
         if (!(entitySpec instanceof PermissionSpec)) {
            verifyFatal(
                  TestScope.FULL,
                  UserBasicSrvApi.getInstance().setEntityPermissions(
                        entitySpec,
                        _permissionSpecToSet.toArray(permissionSpecs)),
                  "Setting entity permissions through API.");
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (ManagedEntitySpec entitySpec : _entitySpecs) {
         if (!(entitySpec instanceof PermissionSpec)) {
            for (PermissionSpec permission : _permissionSpecToSet) {
               verifySafely(TestScope.FULL, UserBasicSrvApi.getInstance()
                     .removeEntityPermission(entitySpec, permission),
                     "Setting original entity permissions");
            }
         }
      }
   }
}
