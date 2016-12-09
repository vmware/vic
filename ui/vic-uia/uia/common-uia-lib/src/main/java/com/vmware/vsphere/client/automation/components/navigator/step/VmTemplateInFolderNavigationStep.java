/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmTemplateInFolderLocationSpec;

/**
 * Retrieve the VmTemplateInFolderLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective
 * VmTemplateInFolderLocationSpec. Use that step to navigate to VM template's
 * related pages.
 */
public class VmTemplateInFolderNavigationStep extends NGCNavigationStep {
   protected static final Logger _logger = LoggerFactory
         .getLogger(VmTemplateInFolderNavigationStep.class);

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(VmTemplateInFolderLocationSpec.class);

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
      _locationSpec = filteredWorkflowSpec
            .get(VmTemplateInFolderLocationSpec.class);
      if (_locationSpec == null) {
         _logger.info("Prepare for navigation to the VM template's base page.");
         _locationSpec = new VmTemplateInFolderLocationSpec();
      }
   }
}
