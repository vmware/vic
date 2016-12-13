/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;

/**
 * Retrieve the VmLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective VmLocationSpec.
 * Use that step to navigate to VM's related pages.
 */
public class VmNavigationStep extends NGCNavigationStep {
   protected static final Logger _logger = LoggerFactory.getLogger(VmNavigationStep.class);

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(VmLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required VmLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   // TestWorkflowStep methods

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(VmLocationSpec.class);

      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the VM base page.");
         _locationSpec = new VmLocationSpec();
      }
   }
}
