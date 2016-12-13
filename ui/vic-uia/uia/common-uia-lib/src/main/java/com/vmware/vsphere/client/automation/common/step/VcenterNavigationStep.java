/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.DatacenterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VcenterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.NGCNavigationStep;

/**
 * Retrive the VcenterLocationSpec from the test spec and invoke the
 * NGCNavigationStep logic to navigate to the respective VcenterLocationSpec.
 * Use that step to navigate to vCenter related pages.
 */
public class VcenterNavigationStep extends NGCNavigationStep {

   @Override
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(VcenterLocationSpec.class);

      if(_locationSpec == null) {
         _logger.info("Prepare for navigation to the vCenter base page.");
         _locationSpec = new DatacenterLocationSpec();
      }
   }
}
