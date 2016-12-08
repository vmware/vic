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
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.step.ClickOkSinglePageDialogStep;
import com.vmware.vsphere.client.automation.common.step.GlobalRefreshStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.HostLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.HostNavigationStep;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SharedPciDeviceSpec.SharedPCIDeviceActionType;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec.VmCreationType;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.FinishCreateVmWizardPageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.LaunchNewVmStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectCreationTypePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectNameAndFolderPageStep;
import com.vmware.vsphere.client.automation.vm.lib.step.LaunchEditSettingsStep;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;
import com.vmware.vsphere.client.automation.vm.ops.step.InvokeVmPowerOperationUiStep;
import com.vmware.vsphere.client.automation.vm.vgpu.step.SelectVmVgpuDevicePageStep;
import com.vmware.vsphere.client.automation.vm.vgpu.step.VerifyAddSharedPciDeviceStep;

/**
 * Test class for create vGPU VM and power it on in the NGC client. Executes the
 * following test work-flow:
 * 1. Open a browser
 * 2. Login as admin user
 * 3. Navigate to the vGPU host
 * 4. Create new VM > add Shared PCI Device
 * 5. Verify Reserve All Memory button appears
 * 6. Verify vGPU related warning and note appears
 * 7. Verify when Reserve All Memory button clicked it disappears as well as vGPU related warning
 * 8. Verify via UI that the create VM task completes successfully
 * 9. Verify via Edit VM Settings that device is added
 * 10. Verify via UI that the VM can be powered on
 */
public class CreateVgpuVmTest extends NGCTestWorkflow {
   private static final String TAG_POWERON = "PowerOn";
   private static final String TAG_CREATE_VM = "CreteVm";

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Add target host as standalone via API.",
            new AddHostStep());
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge
            .requestTestbed(CommonTestBedProvider.class, true);

      TestbedSpecConsumer hostProvider = testbedBridge
            .requestTestbed(HostProvider.class, false);

      // Spec for the VC
      VcSpec requestedVcSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the datacenter
      DatacenterSpec requestedDatacenterSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Spec for the host to be added
      HostSpec vgpuHostSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      vgpuHostSpec.parent.set(requestedDatacenterSpec);

      // Spec for the location of the host
      HostLocationSpec hostLocationSpec = new HostLocationSpec(vgpuHostSpec,
            NGCNavigator.NID_ENTITY_PRIMARY_TAB_SUMMARY);

      // Spec for the VM that will be created
      CreateVmSpec createVmSpec = SpecFactory.getSpec(CreateVmSpec.class,
            vgpuHostSpec);
      createVmSpec.creationType.set(VmCreationType.CREATE_NEW_VM);

      // Spec for the Shared PCI Device of the vm that to be added
      SharedPciDeviceSpec sharedPciDeviceSpec = SpecFactory
            .getSpec(SharedPciDeviceSpec.class, createVmSpec);
      sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.set(true);
      sharedPciDeviceSpec.vmDeviceAction.set(SharedPCIDeviceActionType.ADD);

      CustomizeHwVmSpec customizeVmSpec = SpecFactory
            .getSpec(CustomizeHwVmSpec.class, vgpuHostSpec);
      customizeVmSpec.sharedPciDeviceList.set(sharedPciDeviceSpec);

      // Spec for the location to the VM
      VmLocationSpec vmLocationSpec = new VmLocationSpec(createVmSpec);

      TaskSpec powerOnVmTaskSpec = new TaskSpec();
      powerOnVmTaskSpec.name
            .set(VmUtil.getLocalizedString("task.powerOnVm.name"));
      powerOnVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      powerOnVmTaskSpec.target.set(requestedDatacenterSpec);
      powerOnVmTaskSpec.tag.set(TAG_POWERON);

      VmPowerStateSpec vmPowerStateSpec = new VmPowerStateSpec();
      vmPowerStateSpec.vm.set(createVmSpec);
      vmPowerStateSpec.powerState.set(VmPowerState.POWER_ON);

      // Spec for the create VM task
      TaskSpec createVmTaskSpec = new TaskSpec();
      createVmTaskSpec.name
            .set(VmUtil.getLocalizedString("task.createVm.name"));
      createVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      createVmTaskSpec.target.set(requestedDatacenterSpec);
      createVmTaskSpec.tag.set(TAG_CREATE_VM);

      testSpec.add(requestedVcSpec, requestedDatacenterSpec, vgpuHostSpec,
            hostLocationSpec, createVmSpec, sharedPciDeviceSpec, customizeVmSpec, vmLocationSpec,
            createVmTaskSpec, vmPowerStateSpec, powerOnVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigated to vGPU host.", new HostNavigationStep());

      flow.appendStep("Launched Create VM wizard.", new LaunchNewVmStep());

      flow.appendStep("Selected Crerate new VM as a Creation type.",
            new SelectCreationTypePageStep());

      flow.appendStep("Set name for the VM.",
            new SelectNameAndFolderPageStep());

      flow.appendStep(
            "Passed through Select compute resource page with the default selected node.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Passed through Select storage page with the default selected node.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Passed through Select Compatibility page with default version.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Passed through Select Guest OS page with default settings.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Passed through Select VM Hardware page, with default settings and add Shared PCI Device.",
            new SelectVmVgpuDevicePageStep());
      
      flow.appendStep(
            "Passed through Select Guest OS page with default settings.",
            new ClickNextWizardButtonStep());

      flow.appendStep("Finish the wizard.", new FinishCreateVmWizardPageStep());

      flow.appendStep("Verifying create VM task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_CREATE_VM });

      flow.appendStep("Navigate to VM in object navigator",
            new VmNavigationStep());

      flow.appendStep("Launch Edit Settings of the VM",
            new LaunchEditSettingsStep());

      flow.appendStep(
            "Verify added Shared PCI Device in Edit Settings of the VM",
            new VerifyAddSharedPciDeviceStep());

      flow.appendStep("Click OK on the 'Edit Vm Settings' dialog",
            new ClickOkSinglePageDialogStep());

      // Due to timing issue Power On remains disabled for a few seconds
      flow.appendStep("Refresh the page", new GlobalRefreshStep());

      flow.appendStep("Power On VM", new InvokeVmPowerOperationUiStep());

      flow.appendStep("Verify Power On VM task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_POWERON });
   }

   @Override
   @Test(description = "Create VM with added Shared PCI Device and verify it can be powered on")
   @TestID(id = "616767")
   public void execute() throws Exception {
      super.execute();
   }
}
