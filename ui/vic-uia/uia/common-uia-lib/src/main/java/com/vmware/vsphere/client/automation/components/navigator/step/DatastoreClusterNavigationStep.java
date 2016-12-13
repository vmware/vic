/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.DatastoreClusterLocationSpec;

/**
 * Retrieve the DatastoreClusterLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective DatastoreClusterLocationSpec.
 * Use that step to navigate to datastore cluster's related pages.
 */
public class DatastoreClusterNavigationStep extends NGCNavigationStep {

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(DatastoreClusterLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required DatastoreClusterLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(DatastoreClusterLocationSpec.class);
      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the datastore cluster base page.");
         _locationSpec = new DatastoreClusterLocationSpec();
      }
   }

}
