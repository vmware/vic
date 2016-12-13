/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.StoragePolicyBasicSrvApi;

/**
 * Common new provisioning step that creates storage policy by API.
 */
public class CreateStoragePolicyByApiStep extends BaseWorkflowStep {

   private StoragePolicySpec _storagePolicySpec;
   private BackingCategorySpec _backingCategory;
   private List<BackingTagSpec> _backingTags;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare() {
      _storagePolicySpec = getSpec().links.get(StoragePolicySpec.class);
      _backingCategory = getSpec().links.get(BackingCategorySpec.class);
      _backingTags = getSpec().links.getAll(BackingTagSpec.class);

      if (_storagePolicySpec == null) {
         throw new IllegalArgumentException("No StoragePolicySpec found!");
      }

      if (_backingCategory == null) {
         throw new IllegalArgumentException("No BackingCategorySpec found!");
      }

      if (CollectionUtils.isEmpty(_backingTags)) {
         throw new IllegalArgumentException("No BackingTagSpec found!");
      }
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) {
      _storagePolicySpec = filteredWorkflowSpec.get(StoragePolicySpec.class);
      _backingCategory = filteredWorkflowSpec.get(BackingCategorySpec.class);
      _backingTags = filteredWorkflowSpec.getAll(BackingTagSpec.class);
  
      if (_storagePolicySpec == null) {
         throw new IllegalArgumentException("No StoragePolicySpec found!");
      }

      if (_backingCategory == null) {
         throw new IllegalArgumentException("No BackingCategorySpec found!");
      }

      if (_backingTags == null || _backingTags.size() == 0) {
         throw new IllegalArgumentException("No BackingTagSpec found!");
      }
   }   
   
   
   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      // Create storage policy
      verifyFatal(
            TestScope.FULL,
            StoragePolicyBasicSrvApi.getInstance().createStoragePolicy(
                  _storagePolicySpec,
                  _backingCategory,
                  _backingTags
                  ),
                  "Verify that storage policy was created successfully."
            );

      // check if storage policy is present
      verifyFatal(
            TestScope.FULL,
            StoragePolicyBasicSrvApi.getInstance().isStoragePolicyPresent(_storagePolicySpec),
            "Verify that storage plolicy is present."
            );
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void clean() throws Exception {
      StoragePolicyBasicSrvApi.getInstance().deleteStoragePolicySafely(_storagePolicySpec);
   }
}
