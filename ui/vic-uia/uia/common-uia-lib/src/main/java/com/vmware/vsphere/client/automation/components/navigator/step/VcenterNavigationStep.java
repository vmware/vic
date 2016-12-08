/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.components.navigator.step;

import org.testng.util.Strings;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VcenterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.NGCNavigationStep;

/**
 * Retrieves the <code>VcenterLocationSpec</code> from the test spec and invokes
 * a navigation to the specified location.
 */
public class VcenterNavigationStep extends NGCNavigationStep {

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(VcenterLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required VcenterLocationSpec is missing");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("Path to navigation destination is missing");
      }
   }

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(VcenterLocationSpec.class);
      if (_locationSpec == null) {
         _logger.info("Prepare a navigation to the VC.");
         _locationSpec = new VcenterLocationSpec();
      }
   }
}
