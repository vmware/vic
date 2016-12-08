/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Common workflow step for delete VM via the API.
 * This step works also for VM templates in folder.
 */
public class DeleteVmByApiStep extends BaseWorkflowStep {

   protected VmSpec _vmToDelete;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare() throws Exception {
      _vmToDelete = getSpec().get(VmSpec.class);
      if (_vmToDelete == null) {
         throw new IllegalArgumentException(
               "Delete VM step extpects one VmSpec to be present."
            );
      }
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      if (!VmSrvApi.getInstance().deleteVm(_vmToDelete)) {
         throw new Exception(
               String.format("Unable to delete VM with name '%s'", _vmToDelete.name.get())
            );
      }
   }
}
