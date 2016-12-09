/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;

/**
 * Step for verifying that in launched Edit VM Settings dialog > Shared PCI
 * Device listed
 */
public class VerifyAddSharedPciDeviceStep extends CommonUIWorkflowStep {

   private CustomizeHwVmSpec _customizeHwVmSpec;
   private String editVmSharedPciDeviceTitle = VmUtil
         .getLocalizedString("vm.summary.hardware.shared.pci.device");
   private String editVmSharedPciDeviceType = VmUtil
         .getLocalizedString("vgpu.vm.device.type");

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _customizeHwVmSpec = filteredWorkflowSpec.get(CustomizeHwVmSpec.class);

      ensureNotNull(_customizeHwVmSpec, "VmSpec object is missing.");
   }

   @Override
   public void execute() throws Exception {
      // Bug 1575876
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      verifySafely(
            hardwarePage.getSharedPciTitleLabel()
                  .equals(editVmSharedPciDeviceTitle),
            "Verify title label for added Shared PCI Device in Edit VM Settings dialog"
                  + "expected: " + editVmSharedPciDeviceTitle + "actual: "
                  + hardwarePage.getSharedPciTitleLabel());
      verifySafely(
            hardwarePage.getSharedPciDeviceType()
                  .equals(editVmSharedPciDeviceType),
            "Verify vGPU type for added Shared PCI Device in Edit VM Settings dialog"
                  + "expected: " + editVmSharedPciDeviceType + "actual: "
                  + hardwarePage.getSharedPciDeviceType());
   }
}