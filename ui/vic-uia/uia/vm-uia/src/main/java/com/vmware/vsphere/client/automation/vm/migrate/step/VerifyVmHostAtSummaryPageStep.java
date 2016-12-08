/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.vm.migrate.MigrateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.migrate.view.VmSummaryPageView;

/**
 * Verifies the host of the VM at VM summary page
 */
public class VerifyVmHostAtSummaryPageStep extends MigrateVmFlowStep {
   private HostSpec _targetHost;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _targetHost = filteredWorkflowSpec.get(HostSpec.class);

      ensureNotNull(_targetHost, "No HostSpec object was linked to the spec.");
   }

   @Override
   public void execute() throws Exception {
      // Verify the host label
      if (_targetHost != null) {
         verifyFatal(TestScope.FULL,
               new VmSummaryPageView().isHostFound(_targetHost.name.get()),
               "Verifying the VM host is: " + _targetHost.name.get());
      }

   }
}
