/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.List;

import com.vmware.client.automation.common.step.ReconfigureManagedEntityStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage;

import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.clickAddDevice;
import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.openHwDevicesMenu;
import static com.vmware.vsphere.client.automation.vm.lib.view.EditVmVirtualHardwarePage.selectHddAddDevice;

/**
 * Step that adds a new disk adapter to a VM
 */
public class AddDiskStep extends ReconfigureManagedEntityStep {

   private CustomizeHwVmSpec _vmOldSpec;
   private CustomizeHwVmSpec _vmNewSpec;
   private ReconfigureVmSpec _reconfigVmSpec;
   private int sizeInGb = 1;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      super.prepare(filteredWorkflowSpec);

      _vmNewSpec = (CustomizeHwVmSpec) super.getNewSpec(
            CustomizeHwVmSpec.class);
      _vmOldSpec = (CustomizeHwVmSpec) super.getOldSpec(
            CustomizeHwVmSpec.class);
      _reconfigVmSpec = filteredWorkflowSpec.get(ReconfigureVmSpec.class);

      ensureNotNull(_vmOldSpec, "VmOldSpec object is missing.");
      ensureNotNull(_vmNewSpec, "VmNewSpec object is missing.");
      ensureNotNull(_vmNewSpec.hddList, "Ensure there are hard disks added");
   }

   @Override
   public void execute() throws Exception {

      // Add disk
      EditVmVirtualHardwarePage hardwarePage = new EditVmVirtualHardwarePage();
      int disksNumber = _vmOldSpec.hddList.getAll().size();
      List<HddSpec> hddDevices = _vmNewSpec.hddList.getAll();

      for (HddSpec hddSpec : hddDevices) {
         if (hddSpec.vmDeviceAction.isAssigned() && hddSpec.vmDeviceAction.get()
               .value().equals(HddSpec.HddActionType.ADD.value())) {
            addHdd();
            disksNumber++;
            hardwarePage.setDiskSize(disksNumber - 1, sizeInGb);
         }
      }
   }

   /**
    * Add a HDD. The method invokes the hardware menu > Add New Hard Disk
    */
   public void addHdd() {
      openHwDevicesMenu();
      selectHddAddDevice();
      clickAddDevice();
   }

   @Override
   public void clean() throws Exception {
      if (VmSrvApi.getInstance().checkVmExists(_vmNewSpec)) {
         VmSrvApi.getInstance().reconfigureVm(_reconfigVmSpec);
      }
   }

}