/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.tree.spec.VcenterTreeLocationSpec;

/**
 * Retrieve the VcenterTreeLocationSpec from the test spec and invoke the
 * TreeNavigationStep logic to navigate to the respective
 * VcenterTreeLocationSpec. Use that step to navigate to VC's related pages.
 */
public class VcenterTreeNavigationStep extends TreeNavigationStep {
   protected static final Logger _logger = LoggerFactory
         .getLogger(VcenterTreeNavigationStep.class);

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(VcenterTreeLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required VcenterTreeLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(VcenterTreeLocationSpec.class);
      if (_locationSpec == null) {
         _logger.info("Prepare for navigation to the VC base page.");
         _locationSpec = new VcenterTreeLocationSpec();
      }
   }
}