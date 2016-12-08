/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.step.VerifyTaskByUiStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickOkSinglePageDialogStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.VerifyVmExistenceByApiStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.step.LaunchEditSettingsStep;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;
import com.vmware.vsphere.client.automation.vm.ops.step.InvokeVmPowerOperationUiStep;
import com.vmware.vsphere.client.automation.vm.vgpu.step.AddSharedPciDeviceStep;
import com.vmware.vsphere.client.automation.vm.vgpu.step.VerifyAddSharedPciDeviceStep;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec.SharedPCIDeviceActionType;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Test class for add vGPU device to VM version 11. The workflow is:
 * Pre-requisite steps:
 * 1. Create new VM via API
 * 2. Verify via the API that the VM has been created
 * Test steps:
 * 1. Open a browser
 * 2. Login as Administrator user
 * 3. Navigate to the vGPU host > select created vm via api
 * 4. Edit VM - add new Shared PCI Device
 * 5. Verify warning, note labels and Reserve All Memory buttons appear
 * 6. Verify via UI that the Shared PCI Device has been added
 * 7. Power on VM > verify task succeeds
 */
public class VgpuDeviceVmCompatibilityTest extends NGCTestWorkflow {

   private static final String TAG_INIT_CONFIG = "oldConfig";
   private static final String TAG_NEW_CONFIG = "newConfig";
   private static final String TAG_POWERON_TASK = "powerOn";
   private static final String TAG_RECONFIG_TASK = "reconfigVm";
   private static final String VM_DISK_VERSION11 = "vmx-11";

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Add target host as standalone via API.",
            new AddHostStep());

      flow.appendStep("Create new test VM through the API",
            new CreateVmByApiStep(), new String[] { TAG_INIT_CONFIG });

      flow.appendStep("Verified that VM exists through API.",
            new VerifyVmExistenceByApiStep(), new String[] { TAG_INIT_CONFIG });
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge
            .requestTestbed(CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec requestedVcSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the datacenter
      DatacenterSpec requestedDatacenterSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      TestbedSpecConsumer hostProvider = testbedBridge
            .requestTestbed(HostProvider.class, false);

      // Spec for the host to be added
      HostSpec vgpuHostSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      vgpuHostSpec.parent.set(requestedDatacenterSpec);

      // Spec for the vm that will be edit (init configuration)
      CustomizeHwVmSpec vmSpec = SpecFactory.getSpec(CustomizeHwVmSpec.class,
            vgpuHostSpec);
      vmSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.OLD);
      vmSpec.tag.set(TAG_INIT_CONFIG);
      vmSpec.hardwareVersion.set(VM_DISK_VERSION11);

      // Spec for the Shared PCI Device of the vm that to be added
      SharedPciDeviceSpec sharedPciDeviceSpec = SpecFactory
            .getSpec(SharedPciDeviceSpec.class, vmSpec);
      sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.set(true);
      sharedPciDeviceSpec.vmDeviceAction.set(SharedPCIDeviceActionType.ADD);

      // Spec for the vm that will be edit (target configuration)
      CustomizeHwVmSpec editVmSpec = SpecFactory
            .getSpec(CustomizeHwVmSpec.class, vgpuHostSpec);
      editVmSpec.name.set(vmSpec.name.get());
      editVmSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.NEW);
      editVmSpec.tag.set(TAG_NEW_CONFIG);
      editVmSpec.sharedPciDeviceList.set(sharedPciDeviceSpec);

      // Spec for the reconfigurable VMs
      ReconfigureVmSpec reconfigVmSpec = new ReconfigureVmSpec();
      reconfigVmSpec.newVmConfigs.set(vmSpec);
      reconfigVmSpec.targetVm.set(editVmSpec);

      VmLocationSpec vmLocationSpec = new VmLocationSpec(vmSpec,
            NGCNavigator.NID_VM_SUMMARY);
      vmLocationSpec.tag.set(TAG_INIT_CONFIG);

      // Spec for the edit VM task
      TaskSpec editVmTaskSpec = new TaskSpec();
      editVmTaskSpec.name
            .set(VmUtil.getLocalizedString("task.reconfigureVm.name"));
      editVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      editVmTaskSpec.target.set(vmSpec);
      editVmTaskSpec.tag.set(TAG_RECONFIG_TASK);

      VmLocationSpec vmEditSummaryLocationSpec = new VmLocationSpec(editVmSpec,
            NGCNavigator.NID_VM_SUMMARY);
      vmEditSummaryLocationSpec.tag.set(TAG_NEW_CONFIG);

      TaskSpec powerOnVmTaskSpec = new TaskSpec();
      powerOnVmTaskSpec.name
            .set(VmUtil.getLocalizedString("task.powerOnVm.name"));
      powerOnVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      powerOnVmTaskSpec.target.set(requestedDatacenterSpec);
      powerOnVmTaskSpec.tag.set(TAG_POWERON_TASK);

      VmPowerStateSpec vmPowerStateSpec = new VmPowerStateSpec();
      vmPowerStateSpec.vm.set(vmSpec);
      vmPowerStateSpec.powerState.set(VmPowerState.POWER_ON);

      testSpec.add(requestedVcSpec, requestedDatacenterSpec, vgpuHostSpec,
            vmSpec, sharedPciDeviceSpec, editVmSpec, vmLocationSpec,
            editVmTaskSpec, vmEditSummaryLocationSpec, reconfigVmSpec,
            vmPowerStateSpec, powerOnVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigate to VM in object navigator",
            new VmNavigationStep(), new String[] { TAG_INIT_CONFIG });

      flow.appendStep("Launch Edit Settings of the VM",
            new LaunchEditSettingsStep());

      flow.appendStep("Add new Shared PCI Device for the VM via UI",
            new AddSharedPciDeviceStep());

      flow.appendStep("Click OK on the 'Edit Vm Settings' dialog",
            new ClickOkSinglePageDialogStep());

      flow.appendStep("Verify Reconfigure VM task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_RECONFIG_TASK });

      flow.appendStep("Launch Edit Settings of the VM",
            new LaunchEditSettingsStep());

      flow.appendStep("Verify the new Shared PCI Device in Hardware portlet",
            new VerifyAddSharedPciDeviceStep(), new String[] { TAG_NEW_CONFIG });

      flow.appendStep("Click OK on the 'Edit Vm Settings' dialog",
            new ClickOkSinglePageDialogStep());

      flow.appendStep("Power On VM", new InvokeVmPowerOperationUiStep());

      flow.appendStep("Verify Power On VM task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_POWERON_TASK });
   }

   @Override
   @Test(description = "Edit VM version 11 > Add Shared PCI Device and power on vm.")
   @TestID(id = "618553")
   public void execute() throws Exception {
      super.execute();
   }
}
