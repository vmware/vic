package com.vmware.vsphere.client.automation.storage.lib.core.steps.storagepolicy;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicyComponentSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.StoragePolicyBasicSrvApi;

public class CreateStoragePolicyComponentByApiStep extends
      EnhancedBaseWorkflowStep {

   @UsesSpec()
   protected StoragePolicyComponentSpec storagePolicyComponentSpec;

   @Override
   public void execute() throws Exception {
      // Create storage policy component
      verifyFatal(StoragePolicyBasicSrvApi.getInstance()
            .createStoragePolicyComponent(storagePolicyComponentSpec),
            "Verify that storage policy component was created successfully");

      // check if storage policy component is present
      verifyFatal(StoragePolicyBasicSrvApi.getInstance()
            .isStoragePolicyComponentPresent(storagePolicyComponentSpec),
            "Verify that storage policy component is present");
   }

   @Override
   public void clean() throws Exception {
      verifySafely(StoragePolicyBasicSrvApi.getInstance()
            .deleteStoragePolicyComponentSafely(storagePolicyComponentSpec),
            "Verify that storage policy component is deleted successfully");
   }
}
