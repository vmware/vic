/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.PermissionSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.UserBasicSrvApi;

/**
 * Common workflow that sets permissions on root folder
 *
 * TODO: Ledda - check the usage of UserSrvApi.getInstance().getRootFolderPermissions();
 */
public class SetPermissionsOnRootFolderByApiStep extends BaseWorkflowStep {

   private List<PermissionSpec> _permissionSpecToSet;
   private PermissionSpec[] _permissionSpecToGet = new PermissionSpec[] {};

    @Override
    public void prepare() throws Exception {
        _permissionSpecToSet = getSpec().links.getAll(PermissionSpec.class);
        if (CollectionUtils.isEmpty(_permissionSpecToSet)) {
            throw new IllegalArgumentException("Set Permissions Step requires PermissionSpec.");
        }
    }

   @Override
   public void execute() throws Exception {
      PermissionSpec[] permissionSpecs = new PermissionSpec[_permissionSpecToSet.size()];
      //TODO Ledda - check the usage of UserSrvApi.getInstance().getRootFolderPermissions();
      _permissionSpecToGet = UserBasicSrvApi.getInstance().getRootFolderPermissions(null);
      verifyFatal(
            TestScope.FULL,
            UserBasicSrvApi.getInstance().setRootFolderPermissions(_permissionSpecToSet
                  .toArray(permissionSpecs)),
            "Setting root folder permissions through API.");
   }

   @Override
   public void clean() throws Exception {
      verifySafely(
            TestScope.FULL,
            UserBasicSrvApi.getInstance().setRootFolderPermissions(_permissionSpecToGet),
            "Setting original root folder permissions");
   }
}
