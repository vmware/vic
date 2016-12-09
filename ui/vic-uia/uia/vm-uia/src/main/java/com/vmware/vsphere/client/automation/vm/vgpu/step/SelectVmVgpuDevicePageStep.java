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
 * Create VM>Customize Hardware>Add Shared Passthrough Device
 *
 */
public class SelectVmVgpuDevicePageStep extends CreateVmFlowStep {
   private SharedPciDeviceSpec _sharedPciDeviceSpec;
   String vgpuWarningLabel = VmUtil
         .getLocalizedString("vgpu.vm.poweron.warning");
   String vgpuNoteLabel = VmUtil.getLocalizedString("vgpu.vm.note");

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _sharedPciDeviceSpec = filteredWorkflowSpec
            .get(SharedPciDeviceSpec.class);

      ensureNotNull(_sharedPciDeviceSpec,
            "SharedPciDeviceSpec object is missing.");
      ensureNotNull(_sharedPciDeviceSpec.SharedPCIDeviceReserveMemory,
            "SharedPCIDeviceReserveMemory is missing");
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

         verifySafely(customizePage.getVgpuWarning().equals(vgpuWarningLabel),
               "Verify Warning appears when Shared PCI Device Added");
         // Click Reserve All Memory Button if SharedPCIDeviceReserveMemory true
         if (_sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.isAssigned()
               && _sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.get()
                     .booleanValue()) {
            customizePage.clickReserveAllMemoryBtn();
            verifySafely(!customizePage.isVgpuWarningVisible(),
                  "Verify Warning disappears when Reserve Memory button clicked");
            verifySafely(!customizePage.isReserveAllMemoryBtnVisible(),
                  "Verify Reserve All Memory button disappears when clicked");
         }

         verifySafely(customizePage.getVgpuNoteLabel().equals(vgpuNoteLabel),
               "Verify note appears when Shared PCI Device Added");
      }

      customizePage.waitForLoadingProgressBar();
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
