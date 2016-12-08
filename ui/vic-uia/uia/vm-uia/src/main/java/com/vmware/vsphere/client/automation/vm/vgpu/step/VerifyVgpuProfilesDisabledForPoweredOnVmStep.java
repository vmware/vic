/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

/**
 * Step for verifying vGPU profiles disabled status when VM powered on in
 * launched Edit VM Settings dialog > Shared PCI Device
 */
public class VerifyVgpuProfilesDisabledForPoweredOnVmStep
      extends CommonUIWorkflowStep {

   private VmPowerStateSpec _vmPowerStateSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmPowerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

      ensureNotNull(_vmPowerStateSpec, "VmPowerStateSpec object is missing.");
   }

   @Override
   public void execute() throws Exception {
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      // Vgpu Profiles can be edited only on powered off vm
      boolean expectedVgpuProfileEnabled = _vmPowerStateSpec.powerState.get()
            .equals(VmPowerState.POWER_OFF);

      hardwarePage.expandSharedPciDeviceStackblock();
      boolean actualVgpuState = hardwarePage.getVgpuProfileEnabled();

      verifySafely(actualVgpuState == expectedVgpuProfileEnabled,
            String.format(
                  "Verify vGPU profiles combo box disabled state when vm powered off, expected %s, actual: %s",
                  expectedVgpuProfileEnabled, actualVgpuState));
   }
}