/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone;

import static com.vmware.vsphere.client.automation.vm.clone.CloneVmTestConstants.CLONED_VM_TAG;
import static com.vmware.vsphere.client.automation.vm.clone.CloneVmTestConstants.ORIGINAL_VM_TAG;
import static com.vmware.vsphere.client.automation.vm.clone.CloneVmTestConstants.SOURCE_HOST_TAG;
import static com.vmware.vsphere.client.automation.vm.clone.CloneVmTestConstants.TARGET_HOST_TAG;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec.HddActionType;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.PowerOnVmByApiStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.edit.step.AddDiskStep;
import com.vmware.vsphere.client.automation.vm.edit.step.VerifyAddHddStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.spec.CloneVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.FinishCloneVmWizardPageStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.LaunchCloneVmStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.SelectCloneOptionsPageStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.SelectComputeResourcePageStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.SelectNameAndFolderPageStep;
import com.vmware.vsphere.client.automation.vm.lib.clone.step.SelectStoragePageStep;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;
import com.vmware.vsphere.client.automation.vm.lib.ops.step.VerifyVmPowerStateViaApiStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.MountDatastoreByApiStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.VerifyVmHostAtSummaryPageStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.VerifyVmHostByApiStep;

/**
 * Test for clone powered on VM Executes the following test work-flow:
 * 1. Open a browser
 * 2. Login as admin user
 * 3. Navigate to the powered on VM
 * 4. Clone the VM on another host using NFS storage
 * 5. Select clone Options > Power on VM after creation
 * 6. Select clone Options > Customize Hardware > Add Hard disk device
 * 6. Verify via UI that clone VM task completes successfully
 * 7. Verify via api that VM is powered on
 * 8. Verify via UI that the VM has been cloned and HDD is added
 */
public class CloneVmTest extends NGCTestWorkflow {

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBedProvider = testbedBridge
            .requestTestbed(CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec requestedVcSpec = testBedProvider
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the datacenter
      DatacenterSpec requestedDatacenterSpec = testBedProvider
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Spec for the source host
      HostSpec sourceHostSpec = testBedProvider
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_ENTITY);
      sourceHostSpec.tag.set(SOURCE_HOST_TAG);

      TestbedSpecConsumer hostProvider = testbedBridge
            .requestTestbed(HostProvider.class, false);

      // Spec for the target host
      HostSpec targetHostSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      targetHostSpec.parent.set(requestedDatacenterSpec);
      targetHostSpec.tag.set(TARGET_HOST_TAG);

      // Spec for the NFS datastore
      DatastoreSpec requestedDatastoreSpec = testBedProvider
            .getPublishedEntitySpec(
                  CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);
      requestedDatastoreSpec.parent.set(targetHostSpec);

      // Source VM test spec
      CustomizeHwVmSpec originalVmSpec = SpecFactory
            .getSpec(CustomizeHwVmSpec.class, null, sourceHostSpec);
      originalVmSpec.datastore.set(requestedDatastoreSpec);
      originalVmSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.OLD);
      originalVmSpec.tag.set(ORIGINAL_VM_TAG);

      // Cloned new VM
      CloneVmSpec clonedVmSpec = SpecFactory.getSpec(CloneVmSpec.class,
            targetHostSpec);
      clonedVmSpec.name
            .set("ClonedVmNew" + RandomStringUtils.randomAlphanumeric(10));
      clonedVmSpec.datastore.set(requestedDatastoreSpec);
      clonedVmSpec.targetComputeResource.set(targetHostSpec);
      clonedVmSpec.customizeHw.set(true);
      clonedVmSpec.powerOnVm.set(true);
      clonedVmSpec.tag.set(CLONED_VM_TAG);

      // Spec for the HDD of the vm that to be edit
      HddSpec hddSpec = SpecFactory.getSpec(HddSpec.class, clonedVmSpec);
      hddSpec.name
            .set(VmUtil.getLocalizedString("vm.summary.hardware.harddisk"));
      hddSpec.vmDeviceAction.set(HddActionType.ADD);
      hddSpec.hddCapacity.set("1.00");
      hddSpec.hddCapacityType.set("GB");

      // Spec for customized hardware during Clone VM
      CustomizeHwVmSpec customizeHwSpec = SpecFactory
            .getSpec(CustomizeHwVmSpec.class, targetHostSpec);
      customizeHwSpec.name.set(clonedVmSpec.name.get());
      customizeHwSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.NEW);
      customizeHwSpec.hddList.set(hddSpec);
      customizeHwSpec.datastore.set(requestedDatastoreSpec);
      customizeHwSpec.tag.set(CLONED_VM_TAG);

      // Spec for the reconfigurable VMs
      ReconfigureVmSpec reconfigVmSpec = new ReconfigureVmSpec();
      reconfigVmSpec.newVmConfigs.set(clonedVmSpec);
      reconfigVmSpec.targetVm.set(customizeHwSpec);

      // Location spec for the VM
      VmLocationSpec vmLocationSpec = new VmLocationSpec(originalVmSpec);
      vmLocationSpec.tag.set(ORIGINAL_VM_TAG);

      // Location spec for the coned VM
      VmLocationSpec clonedVmLocationSpec = new VmLocationSpec(clonedVmSpec,
            NGCNavigator.NID_VM_SUMMARY);
      clonedVmLocationSpec.tag.set(CLONED_VM_TAG);

      // Spec for the required VM power state
      VmPowerStateSpec vmPowerStateSpec = new VmPowerStateSpec();
      vmPowerStateSpec.vm.set(clonedVmSpec);
      vmPowerStateSpec.powerState.set(VmPowerState.POWER_ON);

      testSpec.add(requestedVcSpec, sourceHostSpec, targetHostSpec,
            originalVmSpec, clonedVmSpec, hddSpec, customizeHwSpec,
            reconfigVmSpec, vmLocationSpec, clonedVmLocationSpec,
            requestedDatastoreSpec, vmPowerStateSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Add target host as standalone via API.",
            new AddHostStep(), new String[] { TARGET_HOST_TAG });

      flow.appendStep("Mount shared datastore to the target host via API.",
            new MountDatastoreByApiStep());

      flow.appendStep("Create VM via API", new CreateVmByApiStep(),
            new String[] { ORIGINAL_VM_TAG });

      flow.appendStep("Power on VM via API.", new PowerOnVmByApiStep(),
            new String[] { ORIGINAL_VM_TAG });
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigate to test VM.", new VmNavigationStep(),
            new String[] { ORIGINAL_VM_TAG });

      flow.appendStep("Launch Clone VM wizard.", new LaunchCloneVmStep());

      flow.appendStep("Set a name and folder for the VM.",
            new SelectNameAndFolderPageStep());

      flow.appendStep("Select a compute resource - the target host.",
            new SelectComputeResourcePageStep());

      flow.appendStep("Select storage - created NFS.",
            new SelectStoragePageStep());

      flow.appendStep("Select Clone options.",
            new SelectCloneOptionsPageStep());

      flow.appendStep("Add HDD in Customize Hw page", new AddDiskStep());

      flow.appendStep("Click Next", new ClickNextWizardButtonStep());

      flow.appendStep("Finish the wizard.", new FinishCloneVmWizardPageStep());

      flow.appendStep("Verify via API that the VM is powered on",
            new VerifyVmPowerStateViaApiStep());

      flow.appendStep("Verify VM's host through API.",
            new VerifyVmHostByApiStep(),
            new String[] { CLONED_VM_TAG, TARGET_HOST_TAG });

      flow.appendStep("Navigate to cloned VM > Summary tab",
            new VmNavigationStep(), new String[] { CLONED_VM_TAG });

      flow.appendStep("Verify new host appears in VM Summary.",
            new VerifyVmHostAtSummaryPageStep(),
            new String[] { CLONED_VM_TAG, TARGET_HOST_TAG });

      flow.appendStep("Verify the new hdd in Hardware portlet",
            new VerifyAddHddStep(), new String[] { CLONED_VM_TAG });
   }

   @Override
   @Test(description = "Clone powered on VM > Customize HW and power it on after creation", groups = { BAT, CAT })
   @TestID(id = "621167")
   public void execute() throws Exception {
      super.execute();
   }
}
