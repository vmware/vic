/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.ClusterLocationSpec;

/**
 * Retrieve the ClusterLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective ClusterLocationSpec.
 * Use that step to navigate to cluster related pages.
 */
public class ClusterNavigationStep extends NGCNavigationStep {

   @Override
   /**
    * {@inheritDoc}
    */
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(ClusterLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required ClusterLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(ClusterLocationSpec.class);
      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the cluster base page.");
         _locationSpec = new ClusterLocationSpec();
      }
   }

}
