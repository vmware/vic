/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.ops.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

/**
 * Verify VM power state via API
 */
public class VerifyVmPowerStateViaApiStep extends BaseWorkflowStep {

   private VmPowerStateSpec _powerStateSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _powerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

      ensureNotNull(_powerStateSpec, "VmPowerStateSpec object is missing.");
      ensureAssigned(_powerStateSpec.vm, "VM is not assigned to the power spec.");
      ensureAssigned(
            _powerStateSpec.powerState,
            "Power state is not assigned to the power spec.");
   }

   @Override
   public void execute() throws Exception {
      switch (_powerStateSpec.powerState.get()) {
         case POWER_ON:
            verifyFatal(
                  TestScope.BAT,
                  VmSrvApi.getInstance().isVmPoweredOn(_powerStateSpec.vm.get()),
                  String.format(
                        "Verifying that VM '%s' is powered on!",
                        _powerStateSpec.vm.get().name.get()));
            break;
         case POWER_OFF:
            verifyFatal(
                  TestScope.BAT,
                  VmSrvApi.getInstance().isVmPoweredOff(_powerStateSpec.vm.get()),
                  String.format(
                        "Verifying that VM '%s' is powered off!",
                        _powerStateSpec.vm.get().name.get()));
            break;
      }
   }
}
