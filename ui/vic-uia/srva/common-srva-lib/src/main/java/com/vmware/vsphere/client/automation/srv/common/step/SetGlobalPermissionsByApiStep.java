/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotEmpty;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.PermissionSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow that sets global permissions
 */
public class SetGlobalPermissionsByApiStep extends BaseWorkflowStep {

   private List<PermissionSpec> _permissionSpecToSet;

   @Override
   public void prepare() throws Exception {
      _permissionSpecToSet = getSpec().links.getAll(PermissionSpec.class);
      if (CollectionUtils.isEmpty(_permissionSpecToSet)) {
         throw new IllegalArgumentException(
               "Set Permissions Step requires PermissionSpec.");
      }
   }
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create Global Permissions");
      }
      _permissionSpecToSet = filteredWorkflowSpec.getAll(PermissionSpec.class);

      ensureNotNull(_permissionSpecToSet,
            "The spec has no links to 'PermissionSpec' instances");
      ensureNotEmpty(_permissionSpecToSet, "The list with 'PermissionSpec' is empty");
   }

   @Override
   public void execute() throws Exception {
      PermissionSpec[] permissionSpecs = new PermissionSpec[_permissionSpecToSet.size()];
      verifyFatal(TestScope.FULL, UserBasicSrvApi.getInstance().setGlobalPermissions(_permissionSpecToSet.toArray(permissionSpecs)),
            "Setting Global Permissions");
   }

   @Override
   public void clean() throws Exception {
      PermissionSpec[] permissionSpecs = new PermissionSpec[_permissionSpecToSet.size()];
      UserBasicSrvApi.getInstance().deleteGlobalPermissions(_permissionSpecToSet.toArray(permissionSpecs));
   }
}
