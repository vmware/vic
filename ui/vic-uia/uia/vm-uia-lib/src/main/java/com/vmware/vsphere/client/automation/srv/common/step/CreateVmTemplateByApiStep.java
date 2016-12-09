package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Step for creating vm template in folder. Vm templates in folder are created by first
 * creating a VM in the standart way (as in CreateVmByApiStep) and then mark those VMs
 * as templates.
 */
public class CreateVmTemplateByApiStep extends CreateVmByApiStep {

   @Override
   public void execute() throws Exception {
      super.execute();
      for (VmSpec vm : _vmsToCreate) {
         VmSrvApi.getInstance().convertToTemplate(vm);
      }
   }

}
