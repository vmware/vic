/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;

/**
 * Step for verifying if Shared PCI Device removed
 */
public class VerifySharedPciDeviceDeletedStep extends CommonUIWorkflowStep {
   private CustomizeHwVmSpec _customizeHwVmSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _customizeHwVmSpec = filteredWorkflowSpec.get(CustomizeHwVmSpec.class);

      ensureNotNull(_customizeHwVmSpec.sharedPciDeviceList,
            "_customizeHwVmSpec.sharedPciDeviceList object is missing.");
   }

   @Override
   public void execute() throws Exception {
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      // Verify Shared Pci Device is removed
      Boolean isDevicePresent = hardwarePage.isSharedPciDevicePresent();

      verifySafely(isDevicePresent.equals(false),
            "Edit Vm Settings dialog > Verify Shared Pci Device removed");
   }
}