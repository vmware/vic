/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.step.ReconfigureManagedEntityStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;

import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.clickAddDevice;
import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.openHwDevicesMenu;
import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.selectSharedPciAddDevice;

/**
 * Step that adds a new Shared PCI Device to a VM
 */
public class AddSharedPciDeviceStep extends ReconfigureManagedEntityStep {

   private CustomizeHwVmSpec _vmOldSpec;
   private CustomizeHwVmSpec _vmNewSpec;
   private ReconfigureVmSpec _reconfigVmSpec;
   private SharedPciDeviceSpec _sharedPciDeviceSpec;
   private String vgpuWarningLabel = VmUtil
         .getLocalizedString("vgpu.vm.poweron.warning");
   private String vgpuNoteLabel = VmUtil.getLocalizedString("vgpu.vm.note");

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      super.prepare(filteredWorkflowSpec);

      _vmNewSpec = (CustomizeHwVmSpec) super.getNewSpec(
            CustomizeHwVmSpec.class);
      _vmOldSpec = (CustomizeHwVmSpec) super.getOldSpec(
            CustomizeHwVmSpec.class);
      _reconfigVmSpec = filteredWorkflowSpec.get(ReconfigureVmSpec.class);
      _sharedPciDeviceSpec = filteredWorkflowSpec
            .get(SharedPciDeviceSpec.class);

      ensureNotNull(_sharedPciDeviceSpec,
            "sharedPciDeviceSpec object is missing.");
      ensureNotNull(_vmOldSpec, "VmOldSpec object is missing.");
      ensureNotNull(_vmNewSpec, "VmNewSpec object is missing.");
      ensureNotNull(_sharedPciDeviceSpec.vmDeviceAction, "_sharedPciDeviceSpec.vmDeviceAction object is missing.");
      ensureNotNull(_sharedPciDeviceSpec.SharedPCIDeviceReserveMemory, "_sharedPciDeviceSpec.SharedPCIDeviceReserveMemory object is missing.");
   }

   @Override
   public void execute() throws Exception {

      // Add disk
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();

      if (_sharedPciDeviceSpec.vmDeviceAction.isAssigned()
            && _sharedPciDeviceSpec.vmDeviceAction.get().value().equals(
                  SharedPciDeviceSpec.SharedPCIDeviceActionType.ADD.value())) {
         addSharedPciDevice();
         verifySafely(hardwarePage.getVgpuWarning().equals(vgpuWarningLabel),
               "Verify Warning appears when Shared PCI Device Added");

         // Click Reserve All Memory Button if SharedPCIDeviceReserveMemory true
         if (_sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.isAssigned()
               && _sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.get()
                     .booleanValue()) {
            hardwarePage.clickReserveAllMemoryBtn();
            verifySafely(!hardwarePage.isVgpuWarningVisible(),
                  "Verify Warning disappears when Reserve Memory button clicked");
            verifySafely(!hardwarePage.isReserveAllMemoryBtnVisible(),
                  "Verify Reserve All Memory button disappears when clicked");
         }
         verifySafely(hardwarePage.getVgpuNoteLabel().equals(vgpuNoteLabel),
               "Verify note appears when Shared PCI Device Added");
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

   @Override
   public void clean() throws Exception {
      if (VmSrvApi.getInstance().checkVmExists(_vmNewSpec)) {
         VmSrvApi.getInstance().reconfigureVm(_reconfigVmSpec);
      }
   }
}