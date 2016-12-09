/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.createvm.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.vm.createvm.view.VmListOnDatastoreClusterView;

/**
 * Verifies that a VM exists via UI, by checking that it appears in the Virtual
 * Machines tab of the parent datastore cluster.
 *
 * Prerequisite for this step is to navigate to VMs > Virtual Machines tab of
 * the datastore cluster beforehand.
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.step.VerifyVmExistsInDsClusterVmsListStep}
 */
@Deprecated
public class VerifyVmExistsInDsClusterVmsListStep extends CommonUIWorkflowStep {

   private String _vmName;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      VmSpec vmSpec = filteredWorkflowSpec.get(VmSpec.class);
      ensureNotNull(vmSpec, "VmSpec object is missing.");
      _vmName = vmSpec.name.get();
   }

   @Override
   public void execute() throws Exception {
      verifySafely(
            getTestScope(),
            new VmListOnDatastoreClusterView().isFoundInGrid(_vmName),
            String.format(
                  "VM '%s' appears in the datastore cluster's Virtual Machines list.",
                  _vmName));
   }
}
