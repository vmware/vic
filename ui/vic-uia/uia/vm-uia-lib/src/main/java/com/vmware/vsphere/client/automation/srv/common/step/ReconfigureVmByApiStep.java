/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Common workflow step for reconfiguring VM via API.
 * This step works also for VM templates in folder.
 */
public class ReconfigureVmByApiStep extends BaseWorkflowStep {

   private ReconfigureVmSpec _reconfigureVmSpec;

   @Override
   public void prepare() throws Exception {
      _reconfigureVmSpec = getSpec().get(ReconfigureVmSpec.class);
      if (_reconfigureVmSpec == null) {
         throw new IllegalArgumentException(
               "Reconfigur VM step expects ReconfigureVmSpec!"
            );
      }
   }

   @Override
   public void execute() throws Exception {
      if (!VmSrvApi.getInstance().reconfigureVm(_reconfigureVmSpec)) {
         throw new Exception(
               String.format(
                     "Unable to reconfigure VM: '%s'",
                     _reconfigureVmSpec.targetVm.get().name.get()
                  )
            );
      }
   }

   @Override
   public void clean() throws Exception {
      VmSrvApi.getInstance().deleteVmSafely(_reconfigureVmSpec.newVmConfigs.get());
   }
}
