/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;

/**
 * Step for verifying Remove action in launched Edit VM Settings dialog > Shared
 * PCI Device
 */
public class DeleteSharedPciDeviceStep extends CommonUIWorkflowStep {

   private String expectedDeviceRemoveLabel = VmUtil
         .getLocalizedString("vgpu.device.remove.label");

   @Override
   public void execute() throws Exception {
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      // Device can be removed for powered off vm
      hardwarePage.clickRemoveSharedPciDevice();
      String actualRemoveText = hardwarePage.getVgpuDeviceRemoveLabel();

      verifySafely(actualRemoveText.equals(expectedDeviceRemoveLabel),
            String.format("Verify expected removal label %s appears actually in dialog: %s",
                  expectedDeviceRemoveLabel, actualRemoveText));
   }
}