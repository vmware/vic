/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage.openHwDevicesMenu;
import static com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage.selectSharedPciAddDevice;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;
import static com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage.clickAddDevice;

import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.CreateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.messages.VmHardwareMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Create VM>Customize Hardware>Add Shared PCI Device>Try to add second Shared
 * PCI Device
 *
 */
public class AddMaxVmVgpuDeviceStep extends CreateVmFlowStep {
   private SharedPciDeviceSpec _sharedPciDeviceSpec;
   private String vgpuMaxDevicesMsg = VmUtil
         .getLocalizedString("max.devices.message");

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _sharedPciDeviceSpec = filteredWorkflowSpec
            .get(SharedPciDeviceSpec.class);

      ensureNotNull(_sharedPciDeviceSpec,
            "SharedPciDeviceSpec object is missing.");
   }

   @Override
   public void execute() throws Exception {
      CustomizeHardwarePage customizePage = new CustomizeHardwarePage();

      customizePage.waitForLoadingProgressBar();
      customizePage.selectCustomizeHardwareTab(
            I18n.get(VmHardwareMessages.class).vmVirtualHardwareTab());
      customizePage.waitForLoadingProgressBar();
      // Add Shared PCI Device
      if (_sharedPciDeviceSpec.vmDeviceAction.isAssigned()
            && _sharedPciDeviceSpec.vmDeviceAction.get().value().equals(
                  SharedPciDeviceSpec.SharedPCIDeviceActionType.ADD.value())) {
         addSharedPciDevice();
         // Try to add second Shared PCI Device
         addSharedPciDevice();
         // Verify max devices warning appears
         verifySafely(
               vgpuMaxDevicesMsg.equals(customizePage.getMaxDevicesMsg()),
               "Verify max devices message appears when second Shared PCI Device Added");
      }
   }

   /**
    * Add a Shared PCI Device. The method invokes the hardware menu > Add New
    * Shared PCI Device
    */
   public void addSharedPciDevice() {
      openHwDevicesMenu();
      selectSharedPciAddDevice();
      clickAddDevice();
   }
}
