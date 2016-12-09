/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.step.VerifyEntityNameByUiStep;
import com.vmware.client.automation.common.step.VerifyTaskByUiStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.spec.ClusterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.ClusterNavigationStep;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.step.VerifyVmExistenceByApiStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec.VmCreationType;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.FinishCreateVmWizardPageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.LaunchNewVmStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectCreationTypePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectNameAndFolderPageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectStoragePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectVmHardwareWithChangedHddSizePageStep;

/**
 * Test class for create VM in the NGC client.
 * Executes the following test work-flow:
 *  1. Open a browser
 *  2. Login as admin user
 *  3. Navigate to the cluster
 *  4. Create new VM
 *  5. Verify via the API that the VM has been created
 *  6. Verify via UI that the create VM task completes successfully
 *  7. Verify via UI that the VM can be reached
 */
public class CreateDefaultVmTest extends NGCTestWorkflow {

   /**
    * {@inheritDoc}
    */
   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec requestedVcSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the datacenter
      DatacenterSpec requestedDatacenterSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Spec for the cluster
      ClusterSpec requestedClusterSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_ENTITY);

      // Spec for the cluster
      DatastoreSpec requestedDatastoreSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);

      // Spec for the VM that will be created
      CreateVmSpec createVmSpec = SpecFactory.getSpec(CreateVmSpec.class,
            requestedClusterSpec);
      createVmSpec.creationType.set(VmCreationType.CREATE_NEW_VM);
      createVmSpec.datastore.set(requestedDatastoreSpec);

      // Spec for the location to the Cluster
      ClusterLocationSpec clusterLocationSpec = new ClusterLocationSpec(
            requestedClusterSpec);

      // Spec for the location to the VM
      VmLocationSpec vmLocationSpec = new VmLocationSpec(createVmSpec);

      // Spec for the create VM task
      TaskSpec createVmTaskSpec = new TaskSpec();
      createVmTaskSpec.name.set(VmUtil.getLocalizedString("task.createVm.name"));
      createVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      createVmTaskSpec.target.set(requestedDatacenterSpec);

      testSpec.add(requestedVcSpec, requestedClusterSpec, createVmSpec,
            clusterLocationSpec, vmLocationSpec, createVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigated to DRS cluster.", new ClusterNavigationStep());

      flow.appendStep("Launched Create VM wizard.", new LaunchNewVmStep());

      flow.appendStep("Selected Crerate new VM as a Creation type.",
            new SelectCreationTypePageStep());

      flow.appendStep("Set name for the VM.", new SelectNameAndFolderPageStep());

      flow.appendStep(
            "Passed through Select compute resource page with the default selected node - cluster.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Select storage to be used for the vm.",
            new SelectStoragePageStep());

      flow.appendStep("Passed through Select Compatibility page with default version.",
            new ClickNextWizardButtonStep());

      flow.appendStep("Passed through Select Guest OS  page with default settings.",
            new ClickNextWizardButtonStep());

      flow.appendStep(
            "Passed through Select VM Hardware page, with default settings, except HDD size - set to 1 GB.",
            new SelectVmHardwareWithChangedHddSizePageStep());

      flow.appendStep("Finish the wizard.", new FinishCreateVmWizardPageStep());

      flow.appendStep("Verified that VM exists through API.",
            new VerifyVmExistenceByApiStep());

      flow.appendStep("Verifying create VM task via UI", new VerifyTaskByUiStep());

      flow.appendStep("Navigating to VM", new VmNavigationStep());

      flow.appendStep("Verifying VM existence via UI",
            new VerifyEntityNameByUiStep(CreateVmSpec.class));
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
