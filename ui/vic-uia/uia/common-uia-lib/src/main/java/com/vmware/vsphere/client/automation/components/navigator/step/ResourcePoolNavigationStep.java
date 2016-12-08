/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.ResourcePoolLocationSpec;

/**
 * Retrieve the ResourcePoolLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective ResourcePoolLocationSpec.
 * Use that step to navigate to resource pool's related pages.
 */
public class ResourcePoolNavigationStep extends NGCNavigationStep {

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(ResourcePoolLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required ResourcePoolLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(ResourcePoolLocationSpec.class);
      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the resource pools base page.");
         _locationSpec = new ResourcePoolLocationSpec();
      }
   }
}
