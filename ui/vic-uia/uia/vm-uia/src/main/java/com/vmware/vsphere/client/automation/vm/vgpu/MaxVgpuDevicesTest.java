/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.vgpu;

import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickCancelWizardButtonStep;
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.HostLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.HostNavigationStep;
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
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec.VmCreationType;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.LaunchNewVmStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectCreationTypePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectNameAndFolderPageStep;
import com.vmware.vsphere.client.automation.vm.vgpu.step.AddMaxVmVgpuDeviceStep;

/**
 * Test class for create vGPU VM and power it on in the NGC client. Executes the
 * following test work-flow:
 * 1. Open a browser
 * 2. Login as admin user
 * 3. Navigate to the vGPU host
 * 4. Launch Create new VM wizard > add Shared PCI Device
 * 5. Try to add second Shared PCI Device
 * 6. Verify message for max devices of this type appears
 */
public class MaxVgpuDevicesTest extends NGCTestWorkflow {

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
      sharedPciDeviceSpec.SharedPCIDeviceReserveMemory.set(false);
      sharedPciDeviceSpec.vmDeviceAction.set(SharedPCIDeviceActionType.ADD);

      CustomizeHwVmSpec customizeVmSpec = SpecFactory
            .getSpec(CustomizeHwVmSpec.class, vgpuHostSpec);
      customizeVmSpec.sharedPciDeviceList.set(sharedPciDeviceSpec);

      testSpec.add(requestedVcSpec, requestedDatacenterSpec, vgpuHostSpec,
            hostLocationSpec, createVmSpec, sharedPciDeviceSpec, customizeVmSpec);
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
            "Passed through Select VM Hardware page>Add Shared PCI Device>Try to add second Shared PCI Device device.",
            new AddMaxVmVgpuDeviceStep());

      flow.appendStep("Cancel the wizard.", new ClickCancelWizardButtonStep());
   }

   @Override
   @Test(description = "Verify only 1 Shared PCI Device / vGPU device can be added to a VM")
   @TestID(id = "624458")
   public void execute() throws Exception {
      super.execute();
   }
}
