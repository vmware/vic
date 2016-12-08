/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

/**
 * Step for verifying Remove button enabled status in launched Edit VM Settings
 * dialog > Shared PCI Device
 */
public class VerifyDeleteBtnForSharedPciDeviceStep
      extends CommonUIWorkflowStep {

   private CustomizeHwVmSpec _customizeHwVmSpec;
   private VmPowerStateSpec _vmPowerStateSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _customizeHwVmSpec = filteredWorkflowSpec.get(CustomizeHwVmSpec.class);
      _vmPowerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

      ensureNotNull(_customizeHwVmSpec, "VmSpec object is missing.");
      ensureNotNull(_vmPowerStateSpec, "VmPowerStateSpec object is missing.");
   }

   @Override
   public void execute() throws Exception {
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      // Device can be removed only on powered off vm
      Boolean enabledState = _vmPowerStateSpec.powerState.get()
            .equals(VmPowerState.POWER_OFF) ? true : false;
      Boolean deleteBtEnabled = hardwarePage.getVgpuDeviceRemoveBtnEnabled();
      verifySafely(deleteBtEnabled.equals(enabledState),
            "Verify Remove button enabled state" + " expected: " + enabledState
                  + " actual: " + deleteBtEnabled);
   }
}