package com.vmware.vsphere.client.automation.storage.lib.core.steps.storagepolicy;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.StoragePolicyBasicSrvApi;

public class CreateStoragePolicyByApiStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   protected StoragePolicySpec storagePolicySpec;

   @Override
   public void execute() throws Exception {
      verifyFatal(
            StoragePolicyBasicSrvApi.getInstance().createStoragePolicy(
                  storagePolicySpec),
            "Verify that storage policy was created successfully");
      verifyFatal(StoragePolicyBasicSrvApi.getInstance()
            .isStoragePolicyPresent(storagePolicySpec),
            "Verify that storage policy is present");
   }

   @Override
   public void clean() throws Exception {
      verifySafely(StoragePolicyBasicSrvApi.getInstance()
            .deleteStoragePolicySafely(storagePolicySpec),
            "Verify that the storage policy is deleted successfully");
   }
}