/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.StorageLocationSpec;

/**
 * Retrieve the StorageLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective StorageLocationSpec.
 * Use that step to navigate to storage related pages.
 */
public class StorageNavigationStep extends NGCNavigationStep {

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(StorageLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required StorageLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(StorageLocationSpec.class);
      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the storage base page.");
         _locationSpec = new StorageLocationSpec();
      }
   }
}
