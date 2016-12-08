/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate;

import static com.vmware.vsphere.client.automation.vm.migrate.MigrateVmTestConstants.MIGRATED_VM_TAG;
import static com.vmware.vsphere.client.automation.vm.migrate.MigrateVmTestConstants.ORIGINAL_VM_TAG;
import static com.vmware.vsphere.client.automation.vm.migrate.MigrateVmTestConstants.SOURCE_HOST_TAG;
import static com.vmware.vsphere.client.automation.vm.migrate.MigrateVmTestConstants.TARGET_HOST_TAG;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.step.VerifyTaskByUiStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec.AddressType;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.PowerOnVmByApiStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.migrate.spec.MigrateVmSpec;
import com.vmware.vsphere.client.automation.vm.migrate.step.FinishMigrateVmWizardPageStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.LaunchMigrateVmStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.MountDatastoreByApiStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.SelectComputeResourcePageStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.VerifyVmHostAtSummaryPageStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.VerifyVmHostByApiStep;

/**
 * Test for migrating powered on VM
 * Executes the following test work-flow:
 * 1. Open a browser
 * 2. Login as admin user
 * 3. Navigate to the powered on VM
 * 4. Migrate the VM on another host
 * 5. Verify via the API that the VM has been migrated
 * 6. Verify via UI that the relocate VM task completes successfully
 * 7. Verify via UI that the VM has been migrated
 */
public class VmotionVmTest extends NGCTestWorkflow {

   /**
    * {@inheritDoc}
    */
   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBedProvider = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

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

      TestbedSpecConsumer hostProvider = testbedBridge.requestTestbed(
            HostProvider.class, false);

      // Spec for the target host
      HostSpec targetHostSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      targetHostSpec.parent.set(requestedDatacenterSpec);
      targetHostSpec.tag.set(TARGET_HOST_TAG);

      // Spec for the NFS datastore
      DatastoreSpec requestedDatastoreSpec = testBedProvider
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);
      requestedDatastoreSpec.parent.set(targetHostSpec);

      // Source VM test spec
      VmSpec originalVmSpec = SpecFactory.getSpec(VmSpec.class, null, sourceHostSpec);
      originalVmSpec.datastore.set(requestedDatastoreSpec);
      originalVmSpec.tag.set(ORIGINAL_VM_TAG);

      // Migrate test spec
      MigrateVmSpec migrateVmSpec = new MigrateVmSpec();
      migrateVmSpec.targetEntity.set(targetHostSpec);
      migrateVmSpec.targetEntityType.set(MigrateVmSpec.TargetMigrationTypes.HOSTS);

      // Target VM test spec
      VmSpec postMigrateVmSpec = SpecFactory.getSpec(VmSpec.class,
            originalVmSpec.name.get(), targetHostSpec);
      postMigrateVmSpec.tag.set(MIGRATED_VM_TAG);

      // NIC spec
      NicSpec nic1 = SpecFactory.getSpec(NicSpec.class, originalVmSpec);
      nic1.addressType.set(AddressType.GENERATED);
      originalVmSpec.nicList.set(nic1);

      // Location spec for the VM
      VmLocationSpec vmLocationSpec = new VmLocationSpec(originalVmSpec,
            NGCNavigator.NID_VM_SUMMARY);

      // Spec for the relocate VM task
      TaskSpec relocateVmTaskSpec = new TaskSpec();
      relocateVmTaskSpec.name.set(VmUtil.getLocalizedString("task.relocateVm.name"));
      relocateVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      relocateVmTaskSpec.target.set(postMigrateVmSpec);

      testSpec.add(requestedVcSpec, sourceHostSpec, targetHostSpec, originalVmSpec,
            migrateVmSpec, postMigrateVmSpec, nic1, vmLocationSpec,
            requestedDatastoreSpec, relocateVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Add target host as standalone via API.", new AddHostStep(),
            new String[] { TARGET_HOST_TAG });

      flow.appendStep("Mount shared datastore to the target host via API.",
            new MountDatastoreByApiStep());

      flow.appendStep("Create VM via API", new CreateVmByApiStep(),
            new String[] { ORIGINAL_VM_TAG });

      flow.appendStep("Power on VM via API.", new PowerOnVmByApiStep(),
            new String[] { ORIGINAL_VM_TAG });
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigate to test VM.", new VmNavigationStep());

      flow.appendStep("Launch Migrate VM wizard.", new LaunchMigrateVmStep());

      flow.appendStep("Select Change Compute resource migration type.",
            new ClickNextWizardButtonStep());

      flow.appendStep("Select target host.", new SelectComputeResourcePageStep());

      flow.appendStep("Pass through Select network page.",
            new ClickNextWizardButtonStep());

      flow.appendStep("Pass through Select vMotion Priority page.",
            new ClickNextWizardButtonStep());

      flow.appendStep("Finish the wizard.", new FinishMigrateVmWizardPageStep());

      flow.appendStep("Verify relocate VM task via UI", new VerifyTaskByUiStep());

      flow.appendStep("Verify VM's host through API.", new VerifyVmHostByApiStep(),
            new String[] { MIGRATED_VM_TAG, TARGET_HOST_TAG });

      flow.appendStep("Verify new host appears in VM Summary.",
            new VerifyVmHostAtSummaryPageStep(), new String[] { TARGET_HOST_TAG });
   }

   /**
    * {@inheritDoc}
    */
   @Override
   @Test(description = "Create default VM", groups = { BAT, CAT })
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
