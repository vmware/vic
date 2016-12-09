/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.HostLocationSpec;

/**
 * Retrieve the HostLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective HostLocationSpec.
 * Use that step to navigate to host's related pages.
 */
public class HostNavigationStep extends NGCNavigationStep {

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(HostLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required HostLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(HostLocationSpec.class);
      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the hosts base page.");
         _locationSpec = new HostLocationSpec();
      }
   }
}
