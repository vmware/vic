/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.StoragePolicyBasicSrvApi;

/**
 * Common test work-flow step that deletes storage policy by API.
 */
public class DeleteStoragePolicyByApiStep extends BaseWorkflowStep {

   private StoragePolicySpec _storagePolicySpec;

   @Override
   public void prepare() throws Exception {
      _storagePolicySpec = getSpec().links.get(StoragePolicySpec.class);

      if (_storagePolicySpec == null) {
         throw new IllegalArgumentException("No StoragePolicySpec found!");
      }
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(
            TestScope.FULL,
            StoragePolicyBasicSrvApi.getInstance().deleteStoragePolicy(_storagePolicySpec),
            "Delete a storage plolicy.");
   }
}
